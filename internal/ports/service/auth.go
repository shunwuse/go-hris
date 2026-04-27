package service

import (
	"context"

	"github.com/shunwuse/go-hris/internal/domains"
)

type AuthService interface {
	Login(ctx context.Context, username string, password string) (*domains.LoginResult, error)
	ValidateAccessToken(ctx context.Context, tokenString string) (*domains.Claims, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*domains.TokenPair, error)
	Logout(ctx context.Context, refreshToken string, claims *domains.Claims) error
	LogoutAll(ctx context.Context, claims *domains.Claims) error
	CleanupExpiredTokens(ctx context.Context) (int, error)
}
