package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra/app"
	"github.com/shunwuse/go-hris/internal/pkg/cryptox"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/infra"
	"github.com/shunwuse/go-hris/internal/ports/query"
	"github.com/shunwuse/go-hris/internal/ports/repository"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type authService struct {
	logger       *logger.Logger
	tokenService infra.TokenService
	transactor   infra.Transactor

	refreshTokenRepository repository.AuthRepository

	reader query.UserIdentityReader

	jwtRefreshExpire time.Duration
}

func NewAuthService(
	cfg *app.AuthConfig,
	log *logger.Logger,
	transactor infra.Transactor,
	tokenService infra.TokenService,
	refreshTokenRepository repository.AuthRepository,
	reader query.UserIdentityReader,
) service.AuthService {
	// Set refresh expire time with default value.
	refreshExpireDuration := time.Duration(cfg.JWTRefreshExpireMinutes) * time.Minute
	if refreshExpireDuration <= 0 {
		refreshExpireDuration = 24 * time.Hour // default to 24 hours
	}

	log.Info("Auth service initialized",
		zap.Duration("refresh_token_expire", refreshExpireDuration),
	)

	return &authService{
		logger:                 log,
		tokenService:           tokenService,
		transactor:             transactor,
		refreshTokenRepository: refreshTokenRepository,
		reader:                 reader,
		jwtRefreshExpire:       refreshExpireDuration,
	}
}

func (s *authService) Login(ctx context.Context, username string, password string) (*domains.LoginResult, error) {
	// Convert username to lowercase.
	username = strings.ToLower(username)

	user, err := s.reader.GetUserWithPermissionsByUsername(ctx, username)
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to get user", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	// Verify password.
	if user.PasswordHash == "" || !cryptox.CheckPasswordHash(password, user.PasswordHash) {
		s.logger.WithContext(ctx).Error("password verification failed")
		return nil, errors.ErrInvalidCredentials
	}

	// Generate access token.
	accessToken, err := s.GenerateAccessToken(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token.
	refreshToken, err := s.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}

	roles := user.Roles

	return &domains.LoginResult{
		Username:     user.Username,
		Roles:        roles,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) GenerateAccessToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
	return s.tokenService.GenerateAccessToken(ctx, user)
}

func (s *authService) ValidateAccessToken(ctx context.Context, tokenString string) (*domains.Claims, error) {
	claims, err := s.tokenService.ValidateAccessToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	// Fetch user data.
	if claims.Identity.UserID != 0 {
		user, err := s.reader.GetUserWithPermissionsByID(ctx, claims.Identity.UserID)
		if err != nil {
			s.logger.WithContext(ctx).Error("failed to get user data", zap.Uint("userID", claims.Identity.UserID), zap.Error(err))
			return nil, errors.ErrInternalError
		}

		// Set user info.
		claims.Identity.Username = user.Username
		claims.Identity.CreatedAt = user.CreatedAt

		// Set roles.
		claims.Identity.Roles = user.Roles

		// Set permissions.
		claims.Identity.Permissions = user.Permissions
	}

	return claims, nil
}

func (s *authService) GenerateRefreshToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
	// Generate random token (opaque token, not JWT).
	tokenBytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(tokenBytes); err != nil {
		s.logger.WithContext(ctx).Error("failed to generate random token", zap.Error(err))
		return "", errors.ErrInternalError
	}
	rawToken := base64.URLEncoding.EncodeToString(tokenBytes)

	tokenHash := cryptox.SHA256Hex(rawToken)

	expiresAt := time.Now().Add(s.jwtRefreshExpire)

	_, err := s.refreshTokenRepository.CreateRefreshToken(
		ctx,
		tokenHash,
		user.ID,
		expiresAt,
	)
	if err != nil {
		return "", err
	}

	return rawToken, nil
}

func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (*domains.TokenPair, error) {
	tokenHash := cryptox.SHA256Hex(refreshToken)

	var tokenPair *domains.TokenPair
	// Use WithTx to ensure token rotation is atomic (revoke old + create new).
	err := s.transactor.WithTx(ctx, func(txCtx context.Context) error {
		tokenRecord, err := s.refreshTokenRepository.FindRefreshTokenByTokenHash(txCtx, tokenHash)
		if err != nil {
			return err
		}

		// Reuse detection.
		if tokenRecord.Revoked {
			s.logger.WithContext(txCtx).Warn("Refresh token reuse detected! Revoking all tokens for user.",
				zap.Uint("user_id", tokenRecord.UserID))
			// Revoke all tokens for this user as a security measure.
			_ = s.RevokeAllUserTokens(txCtx, tokenRecord.UserID)
			return errors.ErrTokenInvalid
		}

		// Check expiration.
		if tokenRecord.ExpiresAt.Before(time.Now()) {
			return errors.ErrTokenExpired
		}

		user, err := s.reader.GetUserWithPermissionsByID(txCtx, tokenRecord.UserID)
		if err != nil {
			s.logger.WithContext(txCtx).Error("failed to get user for refresh", zap.Uint("user_id", tokenRecord.UserID), zap.Error(err))
			return errors.ErrTokenInvalid
		}

		// Generate new access token.
		accessToken, err := s.GenerateAccessToken(txCtx, user)
		if err != nil {
			return err
		}

		// Token Rotation: revoke old refresh token.
		if err := s.refreshTokenRepository.RevokeRefreshToken(txCtx, tokenHash); err != nil {
			s.logger.WithContext(txCtx).Error("failed to revoke old refresh token", zap.Error(err))
			return err
		}

		// Generate new refresh token.
		newRefreshToken, err := s.GenerateRefreshToken(txCtx, user)
		if err != nil {
			s.logger.WithContext(txCtx).Error("failed to create new refresh token", zap.Error(err))
			return err
		}

		tokenPair = &domains.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string, claims *domains.Claims) error {
	if refreshToken != "" {
		// Revoke refresh token.
		if err := s.RevokeRefreshToken(ctx, refreshToken); err != nil {
			s.logger.WithContext(ctx).Error("failed to revoke refresh token", zap.Error(err))
			return err
		}
	}

	// Blacklist current access token.
	if claims != nil && claims.JTI != "" {
		expiration := claims.ExpiresIn()
		if expiration > 0 {
			if err := s.tokenService.BlacklistToken(ctx, claims.JTI, expiration); err != nil {
				s.logger.WithContext(ctx).Error("failed to blacklist token", zap.Error(err))
				return err
			}
		}
	}

	return nil
}

func (s *authService) LogoutAll(ctx context.Context, claims *domains.Claims) error {
	if claims == nil {
		return errors.ErrUnauthorized
	}

	// Revoke all refresh tokens for this user.
	if err := s.RevokeAllUserTokens(ctx, claims.Identity.UserID); err != nil {
		s.logger.WithContext(ctx).Error("failed to revoke all tokens", zap.Uint("user_id", claims.Identity.UserID), zap.Error(err))
		return err
	}

	// Blacklist current access token.
	if claims.JTI != "" {
		expiration := claims.ExpiresIn()
		if expiration > 0 {
			if err := s.tokenService.BlacklistToken(ctx, claims.JTI, expiration); err != nil {
				s.logger.WithContext(ctx).Error("failed to blacklist token", zap.Error(err))
				return err
			}
		}
	}

	return nil
}

func (s *authService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	tokenHash := cryptox.SHA256Hex(refreshToken)
	return s.refreshTokenRepository.RevokeRefreshToken(ctx, tokenHash)
}

func (s *authService) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	return s.refreshTokenRepository.RevokeAllRefreshTokensForUser(ctx, userID)
}

func (s *authService) CleanupExpiredTokens(ctx context.Context) (int, error) {
	return s.refreshTokenRepository.DeleteExpiredTokens(ctx)
}
