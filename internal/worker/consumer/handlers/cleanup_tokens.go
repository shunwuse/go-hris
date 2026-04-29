package handlers

import (
	"context"

	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type CleanupTokensHandler struct {
	logger      *logger.Logger
	authService service.AuthService
}

func NewCleanupTokensHandler(
	log *logger.Logger,
	authService service.AuthService,
) *CleanupTokensHandler {
	return &CleanupTokensHandler{
		logger:      log,
		authService: authService,
	}
}

func (h *CleanupTokensHandler) Topic() string {
	return "auth.token.cleanup"
}

func (h *CleanupTokensHandler) Handle(ctx context.Context, payload []byte) error {
	count, err := h.authService.CleanupExpiredTokens(ctx)
	if err != nil {
		h.logger.Error("failed to cleanup expired tokens", zap.Error(err))
		return err
	}

	if count > 0 {
		h.logger.Info("cleaned up expired tokens", zap.Int("count", count))
	}

	return nil
}
