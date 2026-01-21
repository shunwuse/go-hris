package infra

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRoutineGroup_RunAndWait(t *testing.T) {
	logger := GetLogger()
	g := NewRoutineGroup(logger)

	taskFinished := make(chan bool, 1)

	g.Run("test_task", func(ctx context.Context) {
		<-ctx.Done() // Wait for stop signal
		taskFinished <- true
	})

	waitCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := g.Wait(waitCtx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	select {
	case finished := <-taskFinished:
		if !finished {
			t.Error("expected task to finish successfully")
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("task did not complete in time")
	}
}

func TestRoutineGroup_PanicRecovery(t *testing.T) {
	logger := GetLogger()
	g := NewRoutineGroup(logger)

	// This task will panic; should be recovered and not crash the test
	g.Run("panic_task", func(ctx context.Context) {
		panic("intentional panic")
	})

	waitCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := g.Wait(waitCtx)
	if err != nil {
		t.Errorf("expected no error after panic recovery, got %v", err)
	}
}

func TestRoutineGroup_Timeout(t *testing.T) {
	logger := GetLogger()
	g := NewRoutineGroup(logger)

	// Hanging task
	g.Run("hanging_task", func(ctx context.Context) {
		select {} // Never finish
	})

	// Very short timeout
	waitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := g.Wait(waitCtx)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("expected deadline exceeded error, got %v", err)
	}
}
