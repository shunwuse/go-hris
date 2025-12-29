package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/worker/consumer"
	"github.com/shunwuse/go-hris/internal/worker/scheduler"
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
	defer cancel()

	// Start Scheduler.
	go w.scheduler.Start(ctx)

	// Start Consumer.
	go w.consumer.Start(ctx)

	w.logger.Info("worker is running and waiting for tasks")

	// Block until we receive an interrupt signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	w.logger.Info("shutting down worker...")
}
