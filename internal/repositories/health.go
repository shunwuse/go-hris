package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type HealthRepository struct {
	logger *infra.Logger
	*infra.Database
	*infra.Cache
}

func NewHealthRepository(
	logger *infra.Logger,
	db *infra.Database,
	cache *infra.Cache,
) *HealthRepository {
	return &HealthRepository{
		logger:   logger,
		Database: db,
		Cache:    cache,
	}
}

func (r *HealthRepository) CheckDatabase(ctx context.Context) bool {
	// Ping the database.
	if err := r.GetRawDB(ctx).PingContext(ctx); err != nil {
		r.logger.WithContext(ctx).Error("database health check failed", zap.Error(err))
		return false
	}

	return true
}

func (r *HealthRepository) CheckRedis(ctx context.Context) bool {
	// Ping the redis.
	if err := r.Cache.Client.Ping(ctx).Err(); err != nil {
		r.logger.WithContext(ctx).Error("redis health check failed", zap.Error(err))
		return false
	}

	return true
}
