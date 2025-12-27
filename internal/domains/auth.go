package domains

import (
	"time"

	"github.com/shunwuse/go-hris/internal/constants"
)

type TokenPayload struct {
	JTI         string                `json:"jti"`
	UserID      uint                  `json:"user_id"`
	Username    string                `json:"username"`
	CreatedAt   time.Time             `json:"created_at"`
	ExpiresAt   time.Time             `json:"expires_at"`
	Roles       []constants.Role      `json:"roles"`
	Permissions constants.Permissions `json:"permissions"`
}

type Claims struct {
	TokenPayload
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResult struct {
	Username     string
	Roles        []constants.Role
	AccessToken  string
	RefreshToken string
}
