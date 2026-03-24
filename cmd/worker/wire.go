//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/queries"
	"github.com/shunwuse/go-hris/internal/repositories"
	"github.com/shunwuse/go-hris/internal/services"
	"github.com/shunwuse/go-hris/internal/worker"
)

func InitializeWorker(
	cfg *Config,
	log *logger.Logger,
) *Worker {
	wire.Build(
		wire.FieldsOf(new(*Config),
			"Database",
			"Cache",
			"Auth",
		),
		infra.ProviderSet,
		repositories.ProviderSet,
		queries.ProviderSet,
		services.ProviderSet,
		worker.ProviderSet,
		NewWorker,
	)

	return &Worker{}
}
