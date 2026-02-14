//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/infra/config"
	"github.com/shunwuse/go-hris/internal/repositories"
	"github.com/shunwuse/go-hris/internal/services"
	"github.com/shunwuse/go-hris/internal/worker"
)

func InitializeWorker(
	cfg *config.Config,
) *Worker {
	wire.Build(
		infra.ProviderSet,
		repositories.ProviderSet,
		services.ProviderSet,
		worker.ProviderSet,
		NewWorker,
	)

	return &Worker{}
}
