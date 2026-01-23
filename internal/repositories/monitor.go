package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/shunwuse/go-hris/internal/ports/repository"
	"go.uber.org/zap"
)

type MonitorRepository struct {
	logger *logger.Logger
	*database.Database
	*cache.Cache
}

func NewMonitorRepository(
	log *logger.Logger,
	db *database.Database,
	c *cache.Cache,
) repository.MonitorRepository {
	return &MonitorRepository{
		logger:   log,
		Database: db,
		Cache:    c,
	}
}

func (r *MonitorRepository) CheckDatabase(ctx context.Context) bool {
	// Ping the database.
	if err := r.GetRawDB(ctx).PingContext(ctx); err != nil {
		r.logger.WithContext(ctx).Error("database health check failed", zap.Error(err))
		return false
	}

	return true
}

func (r *MonitorRepository) CheckRedis(ctx context.Context) bool {
	// Ping the redis.
	if err := r.Cache.Client.Ping(ctx).Err(); err != nil {
		r.logger.WithContext(ctx).Error("redis health check failed", zap.Error(err))
		return false
	}

	return true
}
