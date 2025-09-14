//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/controllers"
	"github.com/shunwuse/go-hris/lib"
	"github.com/shunwuse/go-hris/middlewares"
	"github.com/shunwuse/go-hris/repositories"
	"github.com/shunwuse/go-hris/routes"
	"github.com/shunwuse/go-hris/services"
)

func InitializeServer() Server {
	wire.Build(
		lib.ProvideSet,
		repositories.ProvideSet,
		services.ProvideSet,
		middlewares.ProvideSet,
		controllers.ProvideSet,
		routes.ProvideSet,
		NewServer,
	)

	return Server{}
}
