package scheduler

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/robfig/cron/v3"
	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/pkg/random"
	"github.com/shunwuse/go-hris/internal/ports/infra"
	"github.com/shunwuse/go-hris/internal/worker/scheduler/jobs"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger  *logger.Logger
	alerter infra.Alerter
	cron    *cron.Cron
	jobs    jobs.CronJobs
}

func NewScheduler(
	log *logger.Logger,
	alerter infra.Alerter,
	jobs jobs.CronJobs,
) *Scheduler {
	return &Scheduler{
		logger:  log,
		alerter: alerter,
		cron:    cron.New(),
		jobs:    jobs,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.logger.Info("starting scheduler")

	// Register cron jobs.
	s.registerJobs(ctx)

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

func (s *Scheduler) registerJobs(ctx context.Context) {
	for _, job := range s.jobs {
		cronJob := job
		_, err := s.cron.AddFunc(cronJob.Schedule(), func() {
			s.runJob(ctx, cronJob)
		})

		if err != nil {
			s.logger.Fatal("failed to register job", zap.Error(err))
		}
	}
}

func (s *Scheduler) runJob(baseCtx context.Context, job jobs.ICronJob) {
	executionCtx := baseCtx
	if executionCtx == nil {
		executionCtx = context.Background()
	}
	executionCtx = contextx.WithSystemIdentity(executionCtx)

	alertCtx := contextx.WithSystemIdentity(context.Background())
	traceID := random.NewTraceID()
	executionCtx = contextx.WithTraceID(executionCtx, traceID)
	alertCtx = contextx.WithTraceID(alertCtx, traceID)

	jobName := "unknown"
	jobSchedule := "unknown"
	log := s.logger.WithContext(executionCtx)

	defer func() {
		if recovered := recover(); recovered != nil {
			stack := string(debug.Stack())
			fields := []zap.Field{
				zap.String("job_name", jobName),
				zap.String("job_schedule", jobSchedule),
				zap.Any("panic", recovered),
				zap.String("stack", stack),
				zap.String("trace_id", traceID),
			}

			log.Error("job panic recovered", fields...)

			if s.alerter == nil {
				return
			}

			if err := s.alerter.Send(alertCtx, infra.Message{
				Level:      infra.LevelCritical,
				TraceID:    traceID,
				Title:      "Scheduler Job Panic Recovered",
				Content:    fmt.Sprintf("Job: %s, Schedule: %s, Panic: %v", jobName, jobSchedule, recovered),
				StackTrace: stack,
			}); err != nil {
				log.Warn("failed to send scheduler panic alert", zap.Error(err))
			}
		}
	}()

	jobName = job.Name()
	jobSchedule = job.Schedule()

	log = log.With(
		zap.String("job_name", jobName),
		zap.String("job_schedule", jobSchedule),
		zap.String("trace_id", traceID),
	)

	log.Info("running job")

	if err := job.Run(executionCtx); err != nil {
		log.Error("job failed", zap.Error(err))
		return
	}

	log.Info("job completed successfully")
}
