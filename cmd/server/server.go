package main

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/http/routes"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type Server struct {
	config            infra.Config
	router            infra.RequestHandler
	routes            routes.Routes
	commonMiddlewares middlewares.CommonMiddlewares
	logger            infra.Logger
	database          infra.Database
}

func NewServer(
	config infra.Config,
	router infra.RequestHandler,
	routes routes.Routes,
	commonMiddlewares middlewares.CommonMiddlewares,
	logger infra.Logger,
	database infra.Database,
) Server {
	return Server{
		config:            config,
		router:            router,
		routes:            routes,
		commonMiddlewares: commonMiddlewares,
		logger:            logger,
		database:          database,
	}
}

func (server *Server) Run() {
	server.logger.Info("Starting to run server...")

	// setup common middlewares
	server.commonMiddlewares.Setup(server.router.Router)

	// setup routes
	server.routes.Setup(server.router.Router)

	port := server.config.ServerPort

	if port == "" {
		port = "8080" // default port
	}

	server.logger.Info("Running server on :" + port)
	if err := http.ListenAndServe(":"+port, server.router.Router); err != nil {
		server.logger.Fatal("Error running server", zap.Error(err))
	}
}
