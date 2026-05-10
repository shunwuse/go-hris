package token

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	stdErrors "errors"
	"strconv"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/oklog/ulid/v2"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra/app"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	portinfra "github.com/shunwuse/go-hris/internal/ports/infra"
	"go.uber.org/zap"
)

type JWTService struct {
	logger *logger.Logger
	cache  *cache.Cache

	jwkKey         jwk.Key
	accessTokenTTL time.Duration
}

func New(cfg *app.AuthConfig, log *logger.Logger, cache *cache.Cache) portinfra.TokenService {
	key, err := jwk.FromRaw([]byte(cfg.JWTSecret))
	if err != nil {
		log.Fatal("failed to create JWK from secret", zap.Error(err))
	}

	_ = key.Set(jwk.AlgorithmKey, jwa.HS256)

	hash := sha256.Sum256([]byte(cfg.JWTSecret))
	kid := hex.EncodeToString(hash[:])[:8]
	_ = key.Set(jwk.KeyIDKey, kid)

	// Set expire time with default value.
	expireDuration := time.Duration(cfg.JWTExpireMinutes) * time.Minute
	if expireDuration <= 0 {
		expireDuration = 1 * time.Hour
	}

	log.Info("Token service initialized",
		zap.String("kid", kid),
		zap.Duration("access_token_expire", expireDuration),
	)

	return &JWTService{
		logger:         log,
		cache:          cache,
		jwkKey:         key,
		accessTokenTTL: expireDuration,
	}
}

func (s *JWTService) GenerateAccessToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
	now := time.Now()
	expiration := now.Add(s.accessTokenTTL)

	// Build JWT token with jwx.
	token, err := jwt.NewBuilder().
		// Issuer("go-hris").
		IssuedAt(now).
		Expiration(expiration).
		JwtID(ulid.Make().String()).
		Subject(strconv.FormatUint(uint64(user.ID), 10)).
		Claim(constants.ClaimType, constants.TokenTypeAccess).
		Build()
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to build JWT token", zap.Error(err))
		return "", errors.ErrInternalError
	}

	// Sign token with HS256 using JWK.
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, s.jwkKey))
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to sign JWT token", zap.Error(err))
		return "", errors.ErrInternalError
	}

	return string(signed), nil
}

func (s *JWTService) ValidateAccessToken(ctx context.Context, tokenString string) (*domains.Claims, error) {
	// Parse and verify token using JWK.
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKey(jwa.HS256, s.jwkKey),
		jwt.WithValidate(true),
	)
	if err != nil {
		if stdErrors.Is(err, jwt.ErrTokenExpired()) {
			return nil, errors.ErrTokenExpired
		}

		if jwt.IsValidationError(err) {
			s.logger.WithContext(ctx).Error("failed to validate JWT token", zap.Error(err))
			return nil, errors.ErrTokenInvalid
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
		if err != nil {
			s.logger.WithContext(ctx).Error("failed to check token blacklist",
				zap.String("jti", claims.JTI),
				zap.Error(err),
			)
			return nil, errors.ErrServiceUnavailable
		}

		if blacklisted > 0 {
			return nil, errors.ErrTokenInvalid
		}
	}

	// Parse Subject (User ID) from string back to uint.
	if sub := token.Subject(); sub != "" {
		if id, err := strconv.ParseUint(sub, 10, 64); err == nil {
			claims.Identity.UserID = uint(id)
		}
	}

	return claims, nil
}

func (s *JWTService) BlacklistToken(ctx context.Context, jti string, expiration time.Duration) error {
	if jti == "" {
		return nil
	}

	if err := s.cache.Client.Set(ctx, constants.GetBlacklistKey(jti), "1", expiration).Err(); err != nil {
		s.logger.WithContext(ctx).Error("failed to blacklist token",
			zap.String("jti", jti),
			zap.Error(err),
		)
		return errors.ErrServiceUnavailable
	}

	return nil
}
