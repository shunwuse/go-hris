package repository

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/internal/domains"
)

type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, tokenHash string, userID uint, expiresAt time.Time) (*domains.RefreshToken, error)
	FindRefreshTokenByTokenHash(ctx context.Context, tokenHash string) (*domains.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllRefreshTokensForUser(ctx context.Context, userID uint) error
	DeleteExpiredTokens(ctx context.Context) (int, error)
}
