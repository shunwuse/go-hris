package lifecycle

import (
	"context"
	"errors"
	"slices"
	"sync"
)

// Hook represents a set of callbacks for the start and stop phases of the application.
type Hook struct {
	OnStart func(ctx context.Context) error
	OnStop  func(ctx context.Context) error
}

// Lifecycle manage the starting and stopping of application components.
// It follows LIFO (Last-In-First-Out) order for stopping.
type Lifecycle struct {
	mu    sync.Mutex
	hooks []Hook
}

// New creates a new Lifecycle manager.
func New() *Lifecycle {
	return &Lifecycle{
		hooks: make([]Hook, 0),
	}
}

// Append adds a new hook to the lifecycle.
func (l *Lifecycle) Append(hook Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks = append(l.hooks, hook)
}

// Start executes all registered OnStart hooks in order.
func (l *Lifecycle) Start(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, hook := range l.hooks {
		if hook.OnStart != nil {
			if err := hook.OnStart(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

// Stop executes all registered OnStop hooks in reverse order (LIFO).
func (l *Lifecycle) Stop(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var errs []error
	// LIFO Shutdown.
	for _, hook := range slices.Backward(l.hooks) {
		if hook.OnStop != nil {
			if err := hook.OnStop(ctx); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}
