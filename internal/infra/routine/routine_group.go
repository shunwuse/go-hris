package routine

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/shunwuse/go-hris/internal/infra/logger"
	"go.uber.org/zap"
)

// RoutineGroup manages the lifecycle of background goroutines.
// It provides structured concurrency, panic recovery, and graceful shutdown.
type RoutineGroup struct {
	logger *logger.Logger
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewRoutineGroup creates a new RoutineGroup.
func NewRoutineGroup(log *logger.Logger) *RoutineGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &RoutineGroup{
		logger: log,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run starts a new goroutine managed by the group.
// It includes panic recovery to prevent the entire application from crashing.
func (g *RoutineGroup) Run(name string, fn func(ctx context.Context)) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				g.logger.Error("panic recovered in routine",
					zap.String("routine_name", name),
					zap.Any("panic", r),
					zap.String("stack", string(debug.Stack())),
				)
			}
		}()

		g.logger.Debug("starting routine", zap.String("name", name))

		fn(g.ctx)

		g.logger.Debug("routine finished", zap.String("name", name))
	}()
}

// Context returns the context associated with the group.
// This context is canceled when Wait is called.
func (g *RoutineGroup) Context() context.Context {
	return g.ctx
}

// Wait signals all routines to stop and waits for them to finish.
// It supports a timeout via the provided context.
func (g *RoutineGroup) Wait(ctx context.Context) error {
	g.logger.Info("waiting for routines to finish")
	g.cancel() // Signal all routines to stop

	done := make(chan struct{})
	go func() {
		g.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		g.logger.Info("all routines finished gracefully")
		return nil
	case <-ctx.Done():
		g.logger.Warn("timeout reached while waiting for routines")
		return fmt.Errorf("routine group wait timeout: %w", ctx.Err())
	}
}
