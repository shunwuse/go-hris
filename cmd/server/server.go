package main

import (
	"net/http"
	"time"

	"github.com/shunwuse/go-hris/internal/http/middlewares"
	"github.com/shunwuse/go-hris/internal/http/routes"
	"github.com/shunwuse/go-hris/internal/infra"
	"go.uber.org/zap"
)

type Server struct {
	config            infra.Config
	handler           *infra.RequestHandler
	routes            routes.Routes
	commonMiddlewares middlewares.CommonMiddlewares
	logger            *infra.Logger
	database          *infra.Database
}

func NewServer(
	config infra.Config,
	handler *infra.RequestHandler,
	routes routes.Routes,
	commonMiddlewares middlewares.CommonMiddlewares,
	logger *infra.Logger,
	database *infra.Database,
) *Server {
	return &Server{
		config:            config,
		handler:           handler,
		routes:            routes,
		commonMiddlewares: commonMiddlewares,
		logger:            logger,
		database:          database,
	}
}

func (server *Server) Run() {
	server.logger.Info("starting server initialization")

	// Setup common middlewares.
	server.commonMiddlewares.Setup(server.handler.Router)

	// Setup routes.
	server.routes.Setup(server.handler.Router)

	port := server.config.ServerPort

	if port == "" {
		port = "8080" // default port
	}

	// Create HTTP server with timeout configurations.
	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           server.handler.Router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Start server.
	server.logger.Info("starting HTTP server", zap.String("port", port))
	if err := httpServer.ListenAndServe(); err != nil {
		server.logger.Fatal("failed to start HTTP server", zap.Error(err))
	}
}
