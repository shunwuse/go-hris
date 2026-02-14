package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/worker/consumer"
	"github.com/shunwuse/go-hris/internal/worker/scheduler"
	"go.uber.org/zap"
)

type Worker struct {
	logger    *logger.Logger
	database  *database.Database
	cache     *cache.Cache
	scheduler *scheduler.Scheduler
	consumer  *consumer.Consumer
}

func NewWorker(
	log *logger.Logger,
	db *database.Database,
	cache *cache.Cache,
	scheduler *scheduler.Scheduler,
	consumer *consumer.Consumer,
) *Worker {
	return &Worker{
		logger:    log,
		database:  db,
		cache:     cache,
		scheduler: scheduler,
		consumer:  consumer,
	}
}

func (worker *Worker) Run() {
	worker.logger.Info("starting worker initialization")

	ctx, cancel := context.WithCancel(context.Background())

	// Start Scheduler.
	go worker.scheduler.Start(ctx)

	// Start Consumer.
	go worker.consumer.Start(ctx)

	worker.logger.Info("worker is running and waiting for tasks")

	// Channel to listen for interrupt signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	gracefulShutdown := func() {
		worker.logger.Info("initiating graceful shutdown", zap.Duration("timeout", 30*time.Second))

		// Stop components.
		cancel()

		// Give components some time to shut down.
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// Wait for components to stop.
		worker.scheduler.Stop(shutdownCtx)
		worker.consumer.Stop(shutdownCtx)

		// Close cache connection.
		if err := worker.cache.Close(); err != nil {
			worker.logger.Error("failed to close cache connection", zap.Error(err))
		}

		// Close database connection.
		if err := worker.database.Close(); err != nil {
			worker.logger.Error("failed to close database connection", zap.Error(err))
		}

		worker.logger.Info("worker stopped gracefully")

		_ = worker.logger.Sync()
	}

	sig := <-quit
	worker.logger.Info("shutdown signal received", zap.String("signal", sig.String()))
	gracefulShutdown()
}
