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

	// convert payload to json
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error("marshalling payload failed", zap.Error(err))
		return "", err
	}

	var claims jwt.MapClaims
	// unmarshal json payload
	err = json.Unmarshal(payloadJson, &claims)
	if err != nil {
		s.logger.Error("unmarshalling payload failed", zap.Error(err))
		return "", err
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secreteKey))
	if err != nil {
		s.logger.Error("signing token failed", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}

func (s authService) AuthenticateToken(ctx context.Context, tokenString string) (*domains.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domains.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secreteKey), nil
	})
	if err != nil {
		s.logger.Error("parsing token failed", zap.Error(err))
		return nil, err
	}

	claims, ok := token.Claims.(*domains.Claims)
	if !ok || !token.Valid {
		s.logger.Error("invalid token")
		return nil, err
	}

	return claims, nil
}
