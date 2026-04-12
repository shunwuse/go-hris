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

func (server *Server) Run() error {
	server.logger.Info("starting server initialization")

	// Channel to listen for interrupt signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	return server.run(quit)
}

func (server *Server) run(quit <-chan os.Signal) error {
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

	// Capture HTTP server startup and runtime errors.
	serverErrors := make(chan error, 1)

	go func() {
		server.logger.Info("starting HTTP server", zap.String("port", port))
		serverErrors <- httpServer.ListenAndServe()
	}()
	defer server.shutdown(httpServer)

	// Block until we receive an error or interrupt signal.
	select {
	case err := <-serverErrors:
		if err == http.ErrServerClosed {
			return nil
		}

		server.logger.Error("server exited unexpectedly", zap.Error(err))

		return err

	case sig := <-quit:
		if sig != nil {
			server.logger.Info("shutdown signal received", zap.String("signal", sig.String()))
		}

		return nil
	}
}

func (server *Server) shutdown(httpServer *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.logger.Info("initiating graceful shutdown", zap.Duration("timeout", 30*time.Second))

	// Shutdown HTTP server.
	if err := httpServer.Shutdown(ctx); err != nil {
		server.logger.Error("graceful shutdown failed", zap.Error(err))

		// If shutdown fails, we should try to close other resources.
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

	server.logger.Info("server shutdown complete")

	_ = server.logger.Sync()
}
