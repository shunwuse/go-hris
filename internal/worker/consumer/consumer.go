package consumer

import (
	"context"

	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/worker/consumer/handlers"
	"go.uber.org/zap"
)

type Consumer struct {
	logger   *infra.Logger
	handlers handlers.Handlers
}

func NewConsumer(
	logger *infra.Logger,
	handlers handlers.Handlers,
) *Consumer {
	return &Consumer{
		logger:   logger,
		handlers: handlers,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	c.logger.Info("starting consumer")

	// register handlers.
	for _, handler := range c.handlers {
		c.logger.Info("registered handler", zap.String("topic", handler.Topic()))
	}

	// Block until context is done.
	<-ctx.Done()
	c.logger.Info("stopping consumer")
}

func (c *Consumer) Stop(ctx context.Context) {
	// Implement consumer cleanup here if needed.

	c.logger.Info("consumer stopped successfully")
}
