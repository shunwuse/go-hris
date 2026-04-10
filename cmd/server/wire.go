//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/http"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/queries"
	"github.com/shunwuse/go-hris/internal/repositories"
	"github.com/shunwuse/go-hris/internal/services"
)

func InitializeServer(
	cfg *Config,
	log *logger.Logger,
) *Server {
	wire.Build(
		wire.FieldsOf(new(*Config),
			"Database",
			"Cache",
			"Service",
			"Auth",
		),
		infra.ProviderSet,
		repositories.ProviderSet,
		queries.ProviderSet,
		services.ProviderSet,
		http.ProviderSet,
		NewServer,
	)

	return &Server{}
}
