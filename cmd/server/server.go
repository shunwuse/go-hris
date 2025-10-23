package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// Channel to capture server startup errors.
	serverErrors := make(chan error, 1)

	go func() {
		server.logger.Info("starting HTTP server", zap.String("port", port))
		serverErrors <- httpServer.ListenAndServe()
	}()

	// Channel to listen for interrupt signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	gracefulShutdown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.logger.Info("initiating graceful shutdown", zap.Duration("timeout", 30*time.Second))

		if err := httpServer.Shutdown(ctx); err != nil {
			server.logger.Error("graceful shutdown failed", zap.Error(err))

			if closeErr := httpServer.Close(); closeErr != nil {
				server.logger.Error("forced closure failed", zap.Error(closeErr))
			}

			return
		}

		server.logger.Info("server stopped gracefully")
	}

	// Block until we receive an error or interrupt signal.
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			server.logger.Fatal("server startup failed", zap.Error(err))
		}

	case sig := <-quit:
		server.logger.Info("shutdown signal received", zap.String("signal", sig.String()))
		gracefulShutdown()
	}
}
