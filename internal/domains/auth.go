package domains

import (
	"time"

	"github.com/shunwuse/go-hris/internal/constants"
)

type Identity struct {
	UserID      uint                  `json:"user_id"`
	Username    string                `json:"username"`
	CreatedAt   time.Time             `json:"created_at"`
	Roles       []constants.Role      `json:"roles"`
	Permissions constants.Permissions `json:"permissions"`
}

type Claims struct {
	JTI       string    `json:"jti"`
	ExpiresAt time.Time `json:"expires_at"`

	Identity Identity
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
