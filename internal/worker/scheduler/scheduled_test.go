package scheduler

import (
	"context"
	"errors"
	"testing"

	"github.com/shunwuse/go-hris/internal/pkg/contextx"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/infra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type panicJob struct {
	name     string
	schedule string
}

func (j panicJob) Name() string {
	return j.name
}

func (j panicJob) Schedule() string {
	return j.schedule
}

func (j panicJob) Run(context.Context) error {
	panic("boom")
}

type captureAlerter struct {
	messages []infra.Message
	contexts []context.Context
}

func (a *captureAlerter) Send(ctx context.Context, msg infra.Message) error {
	a.messages = append(a.messages, msg)
	a.contexts = append(a.contexts, ctx)
	return nil
}

type errorAlerter struct{}

func (errorAlerter) Send(context.Context, infra.Message) error {
	return errors.New("alert failed")
}

type okJob struct {
	name     string
	schedule string
	called   bool
	ctx      context.Context
}

func (j *okJob) Name() string {
	return j.name
}

func (j *okJob) Schedule() string {
	return j.schedule
}

func (j *okJob) Run(ctx context.Context) error {
	j.called = true
	j.ctx = ctx
	return nil
}

func TestScheduledJob_Run_RecoversPanicAndAlerts(t *testing.T) {
	alerter := &captureAlerter{}
	scheduler := &Scheduler{logger: logger.NewNopLogger(), alerter: alerter}

	require.NotPanics(t, func() {
		scheduler.runJob(context.Background(), panicJob{name: "cleanup_expired_tokens", schedule: "0 2 * * *"})
	})
	require.Len(t, alerter.messages, 1)
	require.Len(t, alerter.contexts, 1)

	msg := alerter.messages[0]
	assert.Equal(t, infra.LevelCritical, msg.Level)
	assert.Equal(t, "Scheduler Job Panic Recovered", msg.Title)
	assert.Equal(t, contextx.GetTraceID(alerter.contexts[0]), msg.TraceID)
	assert.NotEmpty(t, msg.TraceID)
	assert.Contains(t, msg.Content, "cleanup_expired_tokens")
	assert.Contains(t, msg.Content, "0 2 * * *")
	assert.Contains(t, msg.Content, "boom")
	assert.NotEmpty(t, msg.StackTrace)
}

func TestScheduledJob_Run_DoesNotPanicWhenAlertingFails(t *testing.T) {
	scheduler := &Scheduler{logger: logger.NewNopLogger(), alerter: errorAlerter{}}

	require.NotPanics(t, func() {
		scheduler.runJob(context.Background(), panicJob{name: "cleanup_expired_tokens", schedule: "0 2 * * *"})
	})
}

func TestScheduledJob_Run_UsesExecutionContext(t *testing.T) {
	alerter := &captureAlerter{}
	job := &okJob{name: "cleanup_expired_tokens", schedule: "0 2 * * *"}
	scheduler := &Scheduler{logger: logger.NewNopLogger(), alerter: alerter}

	require.NotPanics(t, func() {
		scheduler.runJob(context.Background(), job)
	})
	assert.True(t, job.called)
	assert.NotNil(t, job.ctx)
	assert.NotEmpty(t, contextx.GetTraceID(job.ctx))
	assert.Empty(t, alerter.messages)
}
