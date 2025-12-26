package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type HealthRepository struct {
	logger *infra.Logger
	*infra.Database
}

func NewHealthRepository(
	logger *infra.Logger,
	db *infra.Database,
) *HealthRepository {
	return &HealthRepository{
		logger:   logger,
		Database: db,
	}
}

func (r *HealthRepository) Check(ctx context.Context) bool {
	// Ping the database.
	if err := r.GetRawDB(ctx).PingContext(ctx); err != nil {
		r.logger.WithContext(ctx).Error("database health check failed", zap.Error(err))
		return false
	}

	return true
}
