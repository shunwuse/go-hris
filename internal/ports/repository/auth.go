package repository

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/ent/entgen"
)

type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, tokenHash string, userID uint, expiresAt time.Time) (*entgen.RefreshToken, error)
	FindRefreshTokenByTokenHash(ctx context.Context, tokenHash string) (*entgen.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllRefreshTokensForUser(ctx context.Context, userID uint) error
	DeleteExpiredTokens(ctx context.Context) (int, error)
}
