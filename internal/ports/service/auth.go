package service

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/internal/domains"
)

type AuthService interface {
	Login(ctx context.Context, username string, password string) (*domains.LoginResult, error)
	GenerateAccessToken(ctx context.Context, user *domains.UserWithPermissions) (string, error)
	ValidateAccessToken(ctx context.Context, tokenString string) (*domains.Claims, error)
	GenerateRefreshToken(ctx context.Context, user *domains.UserWithPermissions) (string, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*domains.TokenPair, error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
	RevokeAllUserTokens(ctx context.Context, userID uint) error
	BlacklistToken(ctx context.Context, jti string, expiration time.Duration) error
	CleanupExpiredTokens(ctx context.Context) (int, error)
}
