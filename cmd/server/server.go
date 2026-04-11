package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shunwuse/go-hris/internal/http/router"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"go.uber.org/zap"
)

type Server struct {
	config   *Config
	logger   *logger.Logger
	database *database.Database
	cache    *cache.Cache
	router   *router.Router
}

func NewServer(
	cfg *Config,
	log *logger.Logger,
	db *database.Database,
	cache *cache.Cache,
	router *router.Router,
) *Server {
	return &Server{
		config:   cfg,
		logger:   log,
		database: db,
		cache:    cache,
		router:   router,
	}
}

func (server *Server) Run() {
	server.logger.Info("starting server initialization")

	port := server.config.Service.Port

	if port == "" {
		port = "8080" // default port
	}

	// Create HTTP server with timeout configurations.
	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           server.router,
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

		// Shutdown HTTP server.
		if err := httpServer.Shutdown(ctx); err != nil {
			server.logger.Error("graceful shutdown failed", zap.Error(err))

			// if shutdown fails, we should try to close other resources.
			if closeErr := httpServer.Close(); closeErr != nil {
				server.logger.Error("forced closure failed", zap.Error(closeErr))
			}

		}

		// Close cache connection.
		if err := server.cache.Close(); err != nil {
			server.logger.Error("failed to close cache connection", zap.Error(err))
		}

		// Close database connection.
		if err := server.database.Close(); err != nil {
			server.logger.Error("failed to close database connection", zap.Error(err))
		}

		server.logger.Info("server stopped gracefully")

		_ = server.logger.Sync()
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
