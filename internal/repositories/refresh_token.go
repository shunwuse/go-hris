package repositories

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type RefreshTokenRepository struct {
	logger *infra.Logger
	*infra.Database
}

func NewRefreshTokenRepository(
	logger *infra.Logger,
	db *infra.Database,
) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		logger:   logger,
		Database: db,
	}
}

func (r *RefreshTokenRepository) Create(
	ctx context.Context,
	tokenHash string,
	userID uint,
	expiresAt time.Time,
) (*entgen.RefreshToken, error) {
	query := r.Client.RefreshToken.Create().
		SetTokenHash(tokenHash).
		SetUserID(userID).
		SetExpiresAt(expiresAt)

	token, err := query.Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create refresh token", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return token, nil
}
