package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/worker/consumer"
	"github.com/shunwuse/go-hris/internal/worker/scheduler"
	"go.uber.org/zap"
)

type Worker struct {
	logger    *infra.Logger
	scheduler *scheduler.Scheduler
	consumer  *consumer.Consumer
}

func NewWorker(
	logger *infra.Logger,
	scheduler *scheduler.Scheduler,
	consumer *consumer.Consumer,
) *Worker {
	return &Worker{
		logger:    logger,
		scheduler: scheduler,
		consumer:  consumer,
	}
}

func (w *Worker) Run() {
	w.logger.Info("starting worker initialization")

	ctx, cancel := context.WithCancel(context.Background())

	// Start Scheduler.
	go w.scheduler.Start(ctx)

	// Start Consumer.
	go w.consumer.Start(ctx)

	w.logger.Info("worker is running and waiting for tasks")

	// Block until we receive an interrupt signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	w.logger.Info("shutdown signal received", zap.String("signal", sig.String()))

	// Initiate graceful shutdown.
	cancel()

	// Give components some time to shut down.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	w.logger.Info("initiating graceful shutdown", zap.Duration("timeout", 30*time.Second))

	// Wait for components to stop.
	w.scheduler.Stop(shutdownCtx)
	w.consumer.Stop(shutdownCtx)

	w.logger.Info("worker stopped gracefully")
}
