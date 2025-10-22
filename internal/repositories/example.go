package repositories

import (
	"context"

	"github.com/shunwuse/go-hris/internal/infra"
)

type ExampleRepository struct {
	logger infra.Logger
	infra.Database
}

func NewExampleRepository(
	logger infra.Logger,
	db infra.Database,
) ExampleRepository {
	return ExampleRepository{
		logger:   logger,
		Database: db,
	}
}

func (r ExampleRepository) Ping(ctx context.Context) string {
	r.logger.WithContext(ctx).Info("ping repository invoked")

	return "pong"
}
