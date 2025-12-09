package service

import (
	"context"

	"github.com/shunwuse/go-hris/internal/domains"
)

type AuthService interface {
	GenerateAccessToken(ctx context.Context, user *domains.UserWithPermissions) (string, error)
	ValidateAccessToken(ctx context.Context, tokenString string) (*domains.Claims, error)
	GenerateRefreshToken(ctx context.Context, user *domains.UserWithPermissions) (string, error)
}
