package jobs

import (
	"context"

	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"go.uber.org/zap"
)

type CleanupTokensJob struct {
	logger      *infra.Logger
	authService service.AuthService
}

func NewCleanupTokensJob(
	logger *infra.Logger,
	authService service.AuthService,
) *CleanupTokensJob {
	return &CleanupTokensJob{
		logger:      logger,
		authService: authService,
	}
}

func (j *CleanupTokensJob) Name() string {
	return "cleanup_expired_tokens"
}

func (j *CleanupTokensJob) Schedule() string {
	// Every day at 02:00.
	return "0 2 * * *"
}

func (j *CleanupTokensJob) Run(ctx context.Context) error {
	count, err := j.authService.CleanupExpiredTokens(ctx)
	if err != nil {
		j.logger.Error("failed to cleanup expired tokens", zap.Error(err))
		return err
	}

	if count > 0 {
		j.logger.Info("cleaned up expired tokens", zap.Int("count", count))
	}

	return nil
}
