package service

import (
	"context"

	"github.com/shunwuse/go-hris/domains"
)

type AuthService interface {
	GenerateToken(ctx context.Context, user *domains.UserWithPermissions) (string, error)
	AuthenticateToken(ctx context.Context, tokenString string) (*domains.Claims, error)
}
