//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/worker/consumer"
	"github.com/shunwuse/go-hris/internal/worker/scheduler"
)

func InitializeWorker() *Worker {
	wire.Build(
		infra.ProviderSet,
		// repositories.ProviderSet,
		// services.ProviderSet,
		scheduler.ProviderSet,
		consumer.ProviderSet,
		NewWorker,
	)

	return &Worker{}
}
