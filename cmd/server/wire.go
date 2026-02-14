//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/http"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/pkg/config"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/repositories"
	"github.com/shunwuse/go-hris/internal/services"
)

func InitializeServer(
	cfg *config.Config,
	log *logger.Logger,
) *Server {
	wire.Build(
		infra.ProviderSet,
		repositories.ProviderSet,
		services.ProviderSet,
		http.ProviderSet,
		NewServer,
	)

	return &Server{}
}
