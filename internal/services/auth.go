package services

import (
	"context"
	"crypto/sha256"
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
	"go.uber.org/zap"
)

type authService struct {
	logger *infra.Logger

	jwtKey jwk.Key
}

func NewAuthService(
	config infra.Config,
	logger *infra.Logger,
) service.AuthService {
	key, err := jwk.FromRaw([]byte(config.JWTSecret))
	if err != nil {
		logger.Fatal("failed to create JWK from secret", zap.Error(err))
	}

	_ = key.Set(jwk.AlgorithmKey, jwa.HS256)

	hash := sha256.Sum256([]byte(config.JWTSecret))
	kid := hex.EncodeToString(hash[:])[:8]
	_ = key.Set(jwk.KeyIDKey, kid)

	logger.Info("JWT signing key initialized", zap.String("kid", kid))

	return &authService{
		logger: logger,

		jwtKey: key,
	}
}

func (s *authService) GenerateToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
	roles := make([]constants.Role, 0)
	for _, role := range user.Edges.Roles {
		roles = append(roles, constants.Role(role.Name))
	}

	// Build JWT token with jwx.
	token, err := jwt.NewBuilder().
		IssuedAt(time.Now()).
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

func (s *authService) AuthenticateToken(ctx context.Context, tokenString string) (*domains.Claims, error) {
	// Parse and verify token using JWK.
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKey(jwa.HS256, s.jwtKey),
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
