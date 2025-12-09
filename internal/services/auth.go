package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/repositories"
	"github.com/shunwuse/go-hris/internal/utils"
	"go.uber.org/zap"
)

type authService struct {
	logger *infra.Logger

	jwtKey           jwk.Key
	jwtExpire        time.Duration
	jwtRefreshExpire time.Duration

	refreshTokenRepository *repositories.RefreshTokenRepository
	userRepository         *repositories.UserRepository
}

func NewAuthService(
	config infra.Config,
	logger *infra.Logger,
	refreshTokenRepository *repositories.RefreshTokenRepository,
	userRepository *repositories.UserRepository,
) service.AuthService {
	key, err := jwk.FromRaw([]byte(config.JWTSecret))
	if err != nil {
		logger.Fatal("failed to create JWK from secret", zap.Error(err))
	}

	_ = key.Set(jwk.AlgorithmKey, jwa.HS256)

	hash := sha256.Sum256([]byte(config.JWTSecret))
	kid := hex.EncodeToString(hash[:])[:8]
	_ = key.Set(jwk.KeyIDKey, kid)

	// Set expire hours with default value.
	expireDuration := time.Duration(config.JWTExpireMinutes) * time.Minute
	if expireDuration <= 0 {
		expireDuration = 1 * time.Hour // default to 1 hour
	}

	// Set refresh expire hours with default value.
	refreshExpireDuration := time.Duration(config.JWTRefreshExpireMinutes) * time.Minute
	if refreshExpireDuration <= 0 {
		refreshExpireDuration = 24 * time.Hour // default to 24 hours
	}

	logger.Info("Auth service initialized",
		zap.String("kid", kid),
		zap.Duration("access_token_expire", expireDuration),
		zap.Duration("refresh_token_expire", refreshExpireDuration),
	)

	return &authService{
		logger: logger,

		jwtKey:           key,
		jwtExpire:        expireDuration,
		jwtRefreshExpire: refreshExpireDuration,

		refreshTokenRepository: refreshTokenRepository,
		userRepository:         userRepository,
	}
}

func (s *authService) GenerateAccessToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
	roles := make([]constants.Role, 0)
	for _, role := range user.Edges.Roles {
		roles = append(roles, constants.Role(role.Name))
	}

	now := time.Now()
	expiration := now.Add(s.jwtExpire)

	// Build JWT token with jwx.
	token, err := jwt.NewBuilder().
		IssuedAt(now).
		Expiration(expiration).
		Claim(constants.ClaimUserID, user.ID).
		Claim(constants.ClaimUsername, user.Username).
		Claim(constants.ClaimCreatedAt, user.CreatedAt).
		Claim(constants.ClaimRoles, roles).
		Claim(constants.ClaimPermissions, user.Permissions).
		Build()
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to build JWT token", zap.Error(err))
		return "", errors.ErrInternalError
	}

	// Sign token with HS256 using JWK.
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, s.jwtKey))
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to sign JWT token", zap.Error(err))
		return "", errors.ErrInternalError
	}

	return string(signed), nil
}

func (s *authService) ValidateAccessToken(ctx context.Context, tokenString string) (*domains.Claims, error) {
	// Parse and verify token using JWK.
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKey(jwa.HS256, s.jwtKey),
		jwt.WithValidate(true), // validate claims like exp, nbf, iat
	)
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to parse JWT token", zap.Error(err))
		return nil, errors.ErrTokenInvalid
	}

	// Extract claims.
	claims := &domains.Claims{}

	if userID, ok := token.Get(constants.ClaimUserID); ok {
		if id, ok := userID.(float64); ok {
			claims.UserID = uint(id)
		}
	}

	if username, ok := token.Get(constants.ClaimUsername); ok {
		claims.Username, _ = username.(string)
	}

	if createdAt, ok := token.Get(constants.ClaimCreatedAt); ok {
		if t, ok := createdAt.(time.Time); ok {
			claims.CreatedAt = t
		}
	}

	if rolesRaw, ok := token.Get(constants.ClaimRoles); ok {
		if roles, ok := rolesRaw.([]any); ok {
			for _, r := range roles {
				if roleStr, ok := r.(string); ok {
					claims.Roles = append(claims.Roles, constants.Role(roleStr))
				}
			}
		}
	}

	if permsRaw, ok := token.Get(constants.ClaimPermissions); ok {
		if perms, ok := permsRaw.([]any); ok {
			for _, p := range perms {
				if permStr, ok := p.(string); ok {
					claims.Permissions = append(claims.Permissions, constants.Permission(permStr))
				}
			}
		}
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

	tokenHash := utils.SHA256Hex(rawToken)

	expiresAt := time.Now().Add(s.jwtRefreshExpire)

	_, err := s.refreshTokenRepository.Create(
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
	tokenHash := utils.SHA256Hex(refreshToken)

	tokenRecord, err := s.refreshTokenRepository.FindValidByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindByID(ctx, tokenRecord.UserID)
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to get user for refresh", zap.Uint("user_id", tokenRecord.UserID), zap.Error(err))
		return nil, errors.ErrTokenInvalid
	}

	// Generate new access token.
	accessToken, err := s.GenerateAccessToken(ctx, user)
	if err != nil {
		return nil, err
	}

	// Token Rotation: revoke old refresh token, issue a new one.
	if err := s.refreshTokenRepository.Revoke(ctx, tokenHash); err != nil {
		s.logger.WithContext(ctx).Warn("failed to revoke old refresh token", zap.Error(err))
		// Don't return error.
	}

	// Generate new refresh token.
	newRefreshToken, err := s.GenerateRefreshToken(ctx, user)
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to create new refresh token", zap.Error(err))
		// Only return access token.
		return &domains.TokenPair{
			AccessToken: accessToken,
		}, nil
	}

	return &domains.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	tokenHash := utils.SHA256Hex(refreshToken)
	return s.refreshTokenRepository.Revoke(ctx, tokenHash)
}
