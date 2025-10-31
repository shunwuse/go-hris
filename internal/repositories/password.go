package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/ent/entgen"
	"github.com/shunwuse/go-hris/internal/errors"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type PasswordRepository struct {
	logger *infra.Logger
	*infra.Database
}

func NewPasswordRepository(
	logger *infra.Logger,
	db *infra.Database,
) *PasswordRepository {
	return &PasswordRepository{
		logger:   logger,
		Database: db,
	}
}

func (r *PasswordRepository) Create(ctx context.Context, hash string, owner *entgen.User) (*entgen.Password, error) {
	password, err := r.Client.Password.Create().
		SetHash(hash).
		SetOwner(owner).
		Save(ctx)
	if err != nil {
		r.logger.WithContext(ctx).Error("failed to create password", zap.Error(err))
		return nil, errors.ErrDatabaseError
	}

	return password, nil
}
