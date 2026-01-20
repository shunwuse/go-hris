package domains

import (
	"slices"
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

func (i *Identity) Can(p constants.Permission) bool {
	return i.Permissions.Contains(p)
}

func (i *Identity) CanAll(ps ...constants.Permission) bool {
	return i.Permissions.ContainsAll(constants.Permissions(ps))
}

func (i *Identity) HasRole(r constants.Role) bool {
	return slices.Contains(i.Roles, r)
}

type Claims struct {
	JTI       string    `json:"jti"`
	ExpiresAt time.Time `json:"expires_at"`

	Identity Identity
}

func (c *Claims) ExpiresIn() time.Duration {
	return time.Until(c.ExpiresAt)
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
