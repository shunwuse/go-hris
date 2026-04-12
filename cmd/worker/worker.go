package main

import (
	"context"
	"fmt"
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
	"golang.org/x/sync/errgroup"
)

type runtimeComponent interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
}

type closable interface {
	Close() error
}

type Worker struct {
	logger    *logger.Logger
	database  closable
	cache     closable
	scheduler runtimeComponent
	consumer  runtimeComponent
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

func (worker *Worker) Run() error {
	worker.logger.Info("starting worker initialization")

	// Handle OS signals for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	return worker.run(quit)
}

func (worker *Worker) run(quit <-chan os.Signal) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		// Start Scheduler.
		return worker.runComponent(ctx, "scheduler", worker.scheduler.Start)
	})

	group.Go(func() error {
		// Start Consumer.
		return worker.runComponent(ctx, "consumer", worker.consumer.Start)
	})

	group.Go(func() error {
		// Wait for either a shutdown signal or component cancellation.
		select {
		case sig := <-quit:
			if sig != nil {
				worker.logger.Info("shutdown signal received", zap.String("signal", sig.String()))
			}
			// Stop components.
			cancel()
		case <-ctx.Done():
		}

		return nil
	})

	worker.logger.Info("worker is running and waiting for tasks")

	runErr := group.Wait()
	worker.shutdown()

	return runErr
}

func (worker *Worker) shutdown() {
	worker.logger.Info("initiating graceful shutdown", zap.Duration("timeout", 30*time.Second))

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

	worker.logger.Info("worker shutdown complete")

	_ = worker.logger.Sync()
}

func (worker *Worker) runComponent(ctx context.Context, name string, start func(context.Context) error) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("%s panicked: %v", name, recovered)
		}
	}()

	if err := start(ctx); err != nil {
		return fmt.Errorf("%s stopped with error: %w", name, err)
	}

	return nil
}
