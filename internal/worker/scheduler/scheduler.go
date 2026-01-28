package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"github.com/shunwuse/go-hris/internal/utils"
	"github.com/shunwuse/go-hris/internal/worker/scheduler/jobs"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger *logger.Logger
	cron   *cron.Cron
	jobs   jobs.CronJobs
}

func NewScheduler(
	log *logger.Logger,
	jobs jobs.CronJobs,
) *Scheduler {
	return &Scheduler{
		logger: log,
		cron:   cron.New(),
		jobs:   jobs,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.logger.Info("starting scheduler")

	// Register cron jobs.
	s.registerJobs()

	// Start Cron Scheduler.
	s.cron.Start()

	// Block until context is done.
	<-ctx.Done()
	s.logger.Info("stopping scheduler")
}

func (s *Scheduler) Stop(ctx context.Context) {
	stopCtx := s.cron.Stop()

	select {
	case <-stopCtx.Done():
		s.logger.Info("scheduler stopped successfully")
	case <-ctx.Done():
		s.logger.Warn("scheduler stop timed out")
	}
}

func (s *Scheduler) registerJobs() {
	for _, job := range s.jobs {
		_, err := s.cron.AddFunc(job.Schedule(), func() {
			// Generate trace ID for the job execution.
			traceID := utils.NewTraceID()
			ctx := contextx.WithTraceID(context.Background(), traceID)

			// Add system identity for background jobs.
			ctx = contextx.WithSystemIdentity(ctx)

			log := s.logger.WithContext(ctx).With(zap.String("job_name", job.Name()))

			log.Info("running job")

			if err := job.Run(ctx); err != nil {
				log.Error("job failed", zap.Error(err))
				return
			}

			log.Info("job completed successfully")
		})

		if err != nil {
			s.logger.Fatal("failed to register job", zap.Error(err))
		}
	}
}
