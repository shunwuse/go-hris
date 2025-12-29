package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/utils"
	"github.com/shunwuse/go-hris/internal/worker/scheduler/jobs"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger *infra.Logger
	cron   *cron.Cron
	jobs   jobs.CronJobs
}

func NewScheduler(
	logger *infra.Logger,
	jobs jobs.CronJobs,
) *Scheduler {
	return &Scheduler{
		logger: logger,
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
	s.cron.Stop()
}

func (s *Scheduler) registerJobs() {
	for _, job := range s.jobs {
		_, err := s.cron.AddFunc(job.Schedule(), func() {
			// Generate trace ID for the job execution.
			traceID := utils.NewTraceID()
			ctx := context.WithValue(context.Background(), constants.TraceID, traceID)

			logger := s.logger.WithContext(ctx).With(zap.String("job_name", job.Name()))

			logger.Info("running job")

			if err := job.Run(ctx); err != nil {
				logger.Error("job failed", zap.Error(err))
				return
			}

			logger.Info("job completed successfully")
		})

		if err != nil {
			s.logger.Fatal("failed to register job", zap.Error(err))
		}
	}
}
