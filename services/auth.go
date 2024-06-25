package services

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/models"
)

type AuthService struct {
	logger lib.Logger

	secreteKey string
}

func NewAuthService() AuthService {
	env := lib.NewEnv()

	logger := lib.GetLogger()

	return AuthService{
		logger: logger,

		secreteKey: env.JWTSecret,
	}
}

type TokenPayload struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (s AuthService) GenerateToken(user *models.User) (string, error) {
	payload := TokenPayload{
		UserID:    user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
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
