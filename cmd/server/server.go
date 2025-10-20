package main

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/http/routes"
	"github.com/shunwuse/go-hris/internal/infra"
)

type Server struct {
	env               infra.Env
	router            infra.RequestHandler
	routes            routes.Routes
	commonMiddlewares middlewares.CommonMiddlewares
	logger            infra.Logger
	database          infra.Database
}

func NewServer(
	env infra.Env,
	router infra.RequestHandler,
	routes routes.Routes,
	commonMiddlewares middlewares.CommonMiddlewares,
	logger infra.Logger,
	database infra.Database,
) Server {
	return Server{
		env:               env,
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

	port := server.env.ServerPort

	if port == "" {
		port = "8080" // default port
	}

	server.logger.Info("Running server on :" + port)
	if err := http.ListenAndServe(":"+port, server.router.Router); err != nil {
		server.logger.Fatalf("Error running server: %v", err)
	}
}
