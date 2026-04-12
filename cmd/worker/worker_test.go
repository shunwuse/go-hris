package main

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeLifecycle struct {
	startCalls int
	stopCalls  int
	startFn    func(context.Context) error
}

func (f *fakeLifecycle) Start(ctx context.Context) error {
	f.startCalls++

	if f.startFn != nil {
		return f.startFn(ctx)
	}

	<-ctx.Done()
	return nil
}

func (f *fakeLifecycle) Stop(context.Context) {
	f.stopCalls++
}

type fakeCloser struct {
	closeCalls int
	closeErr   error
}

func (f *fakeCloser) Close() error {
	f.closeCalls++
	return f.closeErr
}

func TestWorkerRun_ShutsDownOnSignal(t *testing.T) {
	scheduler := &fakeLifecycle{}
	consumer := &fakeLifecycle{}
	cacheCloser := &fakeCloser{}
	databaseCloser := &fakeCloser{}

	worker := &Worker{
		logger:    logger.NewNopLogger(),
		database:  databaseCloser,
		cache:     cacheCloser,
		scheduler: scheduler,
		consumer:  consumer,
	}

	quit := make(chan os.Signal, 1)
	done := make(chan error, 1)

	go func() {
		done <- worker.run(quit)
	}()

	quit <- syscall.SIGTERM

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("worker run timed out")
	}

	assert.Equal(t, 1, scheduler.startCalls)
	assert.Equal(t, 1, scheduler.stopCalls)
	assert.Equal(t, 1, consumer.startCalls)
	assert.Equal(t, 1, consumer.stopCalls)
	assert.Equal(t, 1, cacheCloser.closeCalls)
	assert.Equal(t, 1, databaseCloser.closeCalls)
}

func TestWorkerRun_ReturnsComponentError(t *testing.T) {
	schedulerErr := errors.New("scheduler failed")
	scheduler := &fakeLifecycle{
		startFn: func(context.Context) error {
			return schedulerErr
		},
	}
	consumer := &fakeLifecycle{}
	cacheCloser := &fakeCloser{}
	databaseCloser := &fakeCloser{}

	worker := &Worker{
		logger:    logger.NewNopLogger(),
		database:  databaseCloser,
		cache:     cacheCloser,
		scheduler: scheduler,
		consumer:  consumer,
	}

	quit := make(chan os.Signal, 1)
	done := make(chan error, 1)

	go func() {
		done <- worker.run(quit)
	}()

	select {
	case err := <-done:
		require.Error(t, err)
		assert.ErrorIs(t, err, schedulerErr)
		assert.Contains(t, err.Error(), "scheduler stopped with error")
	case <-time.After(5 * time.Second):
		t.Fatal("worker run timed out")
	}

	assert.Equal(t, 1, scheduler.startCalls)
	assert.Equal(t, 1, scheduler.stopCalls)
	assert.Equal(t, 1, consumer.startCalls)
	assert.Equal(t, 1, consumer.stopCalls)
	assert.Equal(t, 1, cacheCloser.closeCalls)
	assert.Equal(t, 1, databaseCloser.closeCalls)
}
