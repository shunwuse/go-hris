package http

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/http/controllers"
	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/http/routes"
)

var ProviderSet = wire.NewSet(
	middlewares.ProviderSet,
	controllers.ProviderSet,
	routes.ProviderSet,
)
