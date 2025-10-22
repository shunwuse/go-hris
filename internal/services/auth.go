package services

import (
	"context"
	"encoding/json"

	"github.com/golang-jwt/jwt"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type authService struct {
	logger infra.Logger

	secreteKey string
}

func NewAuthService(
	config infra.Config,
	logger infra.Logger,
) service.AuthService {
	return authService{
		logger: logger,

		secreteKey: config.JWTSecret,
	}
}

func (s authService) GenerateToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
	roles := make([]constants.Role, 0)
	for _, role := range user.Edges.Roles {
		roles = append(roles, constants.Role(role.Name))
	}

	payload := domains.TokenPayload{
		UserID:      user.ID,
		Username:    user.Username,
		CreatedAt:   user.CreatedAt,
		Roles:       roles,
		Permissions: user.Permissions,
	}

	// Convert payload to JSON for JWT claims.
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to marshal token payload", zap.Error(err))
		return "", err
	}

	var claims jwt.MapClaims
	// Unmarshal JSON payload into JWT claims.
	err = json.Unmarshal(payloadJson, &claims)
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to unmarshal token payload into claims", zap.Error(err))
		return "", err
	}

	// Generate JWT token with HS256 signing method.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secreteKey))
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to sign JWT token", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}

func (s authService) AuthenticateToken(ctx context.Context, tokenString string) (*domains.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domains.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secreteKey), nil
	})
	if err != nil {
		s.logger.WithContext(ctx).Error("failed to parse JWT token", zap.Error(err))
		return nil, err
	}

	claims, ok := token.Claims.(*domains.Claims)
	if !ok || !token.Valid {
		s.logger.WithContext(ctx).Error("invalid JWT token")
		return nil, err
	}

	return claims, nil
}
