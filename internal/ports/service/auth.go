package service

import (
	"context"

	"github.com/shunwuse/go-hris/internal/domains"
)

type AuthService interface {
	GenerateAccessToken(ctx context.Context, user *domains.UserWithPermissions) (string, error)
	ValidateAccessToken(ctx context.Context, tokenString string) (*domains.Claims, error)
	GenerateRefreshToken(ctx context.Context, user *domains.UserWithPermissions) (string, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*domains.TokenPair, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
	RevokeAllUserTokens(ctx context.Context, userID uint) error
}
