//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/http/routes"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/repositories"
	"github.com/shunwuse/go-hris/internal/services"
)

func InitializeServer() *Server {
	wire.Build(
		infra.ProvideSet,
		repositories.ProvideSet,
		services.ProvideSet,
		middlewares.ProvideSet,
		controllers.ProvideSet,
		routes.ProvideSet,
		NewServer,
	)

	return &Server{}
}
