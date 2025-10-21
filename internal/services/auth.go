package services

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/golang-jwt/jwt"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
)

type authService struct {
	secreteKey string
}

func NewAuthService(
	env infra.Env,
) service.AuthService {
	return authService{
		secreteKey: env.JWTSecret,
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
		slog.Error("marshalling payload failed", "error", err)
		return "", err
	}

	var claims jwt.MapClaims
	// unmarshal json payload
	err = json.Unmarshal(payloadJson, &claims)
	if err != nil {
		slog.Error("unmarshalling payload failed", "error", err)
		return "", err
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secreteKey))
	if err != nil {
		slog.Error("signing token failed", "error", err)
		return "", err
	}

	return tokenString, nil
}

func (s authService) AuthenticateToken(ctx context.Context, tokenString string) (*domains.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domains.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secreteKey), nil
	})
	if err != nil {
		slog.Error("parsing token failed", "error", err)
		return nil, err
	}

	claims, ok := token.Claims.(*domains.Claims)
	if !ok || !token.Valid {
		slog.Error("invalid token")
		return nil, err
	}

	return claims, nil
}
