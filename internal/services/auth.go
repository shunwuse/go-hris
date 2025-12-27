package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/oklog/ulid/v2"
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
	cache  *infra.Cache

	refreshTokenRepository *repositories.AuthRepository

	userService service.UserService

	jwtKey           jwk.Key
	jwtExpire        time.Duration
	jwtRefreshExpire time.Duration
}

func NewAuthService(
	config infra.Config,
	logger *infra.Logger,
	cache *infra.Cache,
	userService service.UserService,
	refreshTokenRepository *repositories.AuthRepository,
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
		logger:                 logger,
		cache:                  cache,
		refreshTokenRepository: refreshTokenRepository,
		userService:            userService,
		jwtKey:                 key,
		jwtExpire:              expireDuration,
		jwtRefreshExpire:       refreshExpireDuration,
	}
}

func (s *authService) Login(ctx context.Context, username string, password string) (*domains.LoginResult, error) {
	// Convert username to lowercase.
	username = strings.ToLower(username)

	user, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to get user", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	// Verify password.
	if !utils.CheckPasswordHash(password, user.Edges.Password.Hash) {
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

	roles := make([]constants.Role, len(user.Edges.Roles))
	for idx, role := range user.Edges.Roles {
		roles[idx] = role.Name
	}

	return &domains.LoginResult{
		Username:     user.Username,
		Roles:        roles,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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
		// Issuer("go-hris").
		IssuedAt(now).
		Expiration(expiration).
		JwtID(ulid.Make().String()).
		Subject(strconv.FormatUint(uint64(user.ID), 10)).
		Claim(constants.ClaimUsername, user.Username).
		Claim(constants.ClaimCreatedAt, user.CreatedAt).
		Claim(constants.ClaimRoles, roles).
		Claim(constants.ClaimPermissions, user.Permissions).
		Claim(constants.ClaimType, constants.TokenTypeAccess).
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
		// jwt.WithIssuer("go-hris"),
	)
	if err != nil {
		if strings.Contains(err.Error(), "exp") {
			return nil, errors.ErrTokenExpired
		}
		s.logger.WithContext(ctx).Error("failed to parse JWT token", zap.Error(err))
		return nil, errors.ErrTokenInvalid
	}

	// Verify token type.
	tokenType, ok := token.Get(constants.ClaimType)
	if !ok || tokenType != string(constants.TokenTypeAccess) {
		s.logger.WithContext(ctx).Error("invalid token type", zap.Any("type", tokenType))
		return nil, errors.ErrTokenInvalid
	}

	// Extract claims.
	claims := &domains.Claims{}

	claims.JTI = token.JwtID()
	claims.ExpiresAt = token.Expiration()

	// Check if token is blacklisted.
	if claims.JTI != "" {
		blacklisted, err := s.cache.Client.Exists(ctx, constants.GetBlacklistKey(claims.JTI)).Result()
		if err == nil && blacklisted > 0 {
			return nil, errors.ErrTokenInvalid
		}
	}

	// Parse Subject (User ID) from string back to uint.
	if sub := token.Subject(); sub != "" {
		if id, err := strconv.ParseUint(sub, 10, 64); err == nil {
			claims.UserID = uint(id)
		}
	}

	if username, ok := token.Get(constants.ClaimUsername); ok {
		claims.Username, _ = username.(string)
	}

	if createdAt, ok := token.Get(constants.ClaimCreatedAt); ok {
		switch v := createdAt.(type) {
		case time.Time:
			claims.CreatedAt = v
		case string:
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				claims.CreatedAt = t
			}
		case float64:
			claims.CreatedAt = time.Unix(int64(v), 0)
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
	tokenHash := utils.SHA256Hex(refreshToken)

	var tokenPair *domains.TokenPair
	// Use WithTx to ensure token rotation is atomic (revoke old + create new).
	err := s.refreshTokenRepository.WithTx(ctx, func(txCtx context.Context) error {
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

		user, err := s.userService.GetUserByID(txCtx, tokenRecord.UserID)
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

func (s *authService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	tokenHash := utils.SHA256Hex(refreshToken)
	return s.refreshTokenRepository.RevokeRefreshToken(ctx, tokenHash)
}

func (s *authService) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	return s.refreshTokenRepository.RevokeAllRefreshTokensForUser(ctx, userID)
}

func (s *authService) BlacklistToken(ctx context.Context, jti string, expiration time.Duration) error {
	if jti == "" {
		return nil
	}

	return s.cache.Client.Set(ctx, constants.GetBlacklistKey(jti), "1", expiration).Err()
}
