package handlers

import (
	"context"
)

type Handlers []Handler

type Handler interface {
	Topic() string
	Handle(ctx context.Context, payload []byte) error
}

func NewHandlers(
	cleanupTokensHandler *CleanupTokensHandler,
) Handlers {
	return []Handler{
		cleanupTokensHandler,
	}
}
