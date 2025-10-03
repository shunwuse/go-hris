package services

import (
	"context"
	"encoding/json"

	"github.com/golang-jwt/jwt"
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/domains"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
	"github.com/shunwuse/go-hris/ports/service"
)

type authService struct {
	logger lib.Logger

	secreteKey string
}

func NewAuthService(
	env lib.Env,
	logger lib.Logger,
) service.AuthService {
	return authService{
		logger: logger,

		secreteKey: env.JWTSecret,
	}
}

func (s authService) GenerateToken(ctx context.Context, user *models.User) (string, error) {
	roles := make([]constants.Role, 0)
	for _, role := range user.Roles {
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
		s.logger.Errorf("marshalling payload failed: %v", err)
		return "", err
	}

	var claims jwt.MapClaims
	// unmarshal json payload
	err = json.Unmarshal(payloadJson, &claims)
	if err != nil {
		s.logger.Errorf("unmarshalling payload failed: %v", err)
		return "", err
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secreteKey))
	if err != nil {
		s.logger.Errorf("signing token failed: %v", err)
		return "", err
	}

	return tokenString, nil
}

func (s authService) AuthenticateToken(ctx context.Context, tokenString string) (*domains.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domains.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secreteKey), nil
	})
	if err != nil {
		s.logger.Errorf("parsing token failed: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(*domains.Claims)
	if !ok || !token.Valid {
		s.logger.Errorf("invalid token")
		return nil, err
	}

	return claims, nil
}
