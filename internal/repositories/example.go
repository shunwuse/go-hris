package repositories

import (
	"context"
	"log/slog"

	"github.com/shunwuse/go-hris/internal/infra"
)

type ExampleRepository struct {
	infra.Database
}

func NewExampleRepository(
	db infra.Database,
) ExampleRepository {
	return ExampleRepository{
		Database: db,
	}
}

func (r ExampleRepository) Ping(ctx context.Context) string {
	slog.Info("Ping repository")

	return "pong"
}
