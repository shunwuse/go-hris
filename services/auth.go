package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shunwuse/go-hris/constants"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
)

type AuthService struct {
	logger lib.Logger

	secreteKey string
}

func NewAuthService(
	env lib.Env,
	logger lib.Logger,
) AuthService {
	return AuthService{
		logger: logger,

		secreteKey: env.JWTSecret,
	}
}

type TokenPayload struct {
	UserID      uint                  `json:"user_id"`
	Username    string                `json:"username"`
	CreatedAt   time.Time             `json:"created_at"`
	Roles       []constants.Role      `json:"roles"`
	Permissions constants.Permissions `json:"permissions"`
}

func (s AuthService) GenerateToken(ctx context.Context, user *models.User) (string, error) {
	roles := make([]constants.Role, 0)
	for _, role := range user.Roles {
		roles = append(roles, constants.Role(role.Name))
	}

	payload := TokenPayload{
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

type claims struct {
	jwt.StandardClaims
	TokenPayload
}

func (s AuthService) AuthenticateToken(ctx context.Context, tokenString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secreteKey), nil
	})
	if err != nil {
		s.logger.Errorf("parsing token failed: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		s.logger.Errorf("invalid token")
		return nil, err
	}

	return claims, nil
}
