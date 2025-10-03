package domains

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/shunwuse/go-hris/constants"
)

type TokenPayload struct {
	UserID      uint                  `json:"user_id"`
	Username    string                `json:"username"`
	CreatedAt   time.Time             `json:"created_at"`
	Roles       []constants.Role      `json:"roles"`
	Permissions constants.Permissions `json:"permissions"`
}

type Claims struct {
	jwt.StandardClaims

	TokenPayload
}
