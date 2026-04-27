package infra

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/internal/domains"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, user *domains.UserWithPermissions) (string, error)
	ValidateAccessToken(ctx context.Context, tokenString string) (*domains.Claims, error)
	BlacklistToken(ctx context.Context, jti string, expiration time.Duration) error
}
