package routine

import (
	"context"
	"testing"
	"time"

	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutineGroup_RunAndWait(t *testing.T) {
	log := logger.GetLogger()
	g := NewRoutineGroup(log)

	taskFinished := make(chan bool, 1)

	g.Run("test_task", func(ctx context.Context) {
		<-ctx.Done() // Wait for stop signal
		taskFinished <- true
	})

	waitCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := g.Wait(waitCtx)
	require.NoError(t, err)

	select {
	case finished := <-taskFinished:
		assert.True(t, finished)
	case <-time.After(500 * time.Millisecond):
		t.Error("task did not complete in time")
	}
}

func TestRoutineGroup_PanicRecovery(t *testing.T) {
	log := logger.GetLogger()
	g := NewRoutineGroup(log)

	// This task will panic; should be recovered and not crash the test
	g.Run("panic_task", func(ctx context.Context) {
		panic("intentional panic")
	})

	waitCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := g.Wait(waitCtx)
	require.NoError(t, err)
}

func TestRoutineGroup_Timeout(t *testing.T) {
	log := logger.GetLogger()
	g := NewRoutineGroup(log)

	// Hanging task
	g.Run("hanging_task", func(ctx context.Context) {
		select {} // Never finish
	})

	// Very short timeout
	waitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := g.Wait(waitCtx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}
