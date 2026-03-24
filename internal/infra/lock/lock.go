package lock

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"go.uber.org/zap"
)

// Lock represents an obtained, distributed lock.
type Lock struct {
	lock   *redislock.Lock
	ttl    time.Duration
	key    string
	logger *logger.Logger
}

// Release releases the lock.
func (l *Lock) Release(ctx context.Context) error {
	return l.lock.Release(ctx)
}

// Refresh extends the lock with a new TTL.
func (l *Lock) Refresh(ctx context.Context, ttl time.Duration, options *redislock.Options) error {
	return l.lock.Refresh(ctx, ttl, options)
}

// TTL returns the remaining time-to-live.
func (l *Lock) TTL(ctx context.Context) (time.Duration, error) {
	return l.lock.TTL(ctx)
}

// Metadata returns the metadata of the lock.
func (l *Lock) Metadata() string {
	return l.lock.Metadata()
}

// Token returns the unique token value of the lock.
func (l *Lock) Token() string {
	return l.lock.Token()
}

// AutoRefresh starts a background goroutine to refresh the lock at 1/3 of its TTL interval.
// It stops when the provided context is canceled or if the refresh fails.
func (l *Lock) AutoRefresh(ctx context.Context) {
	l.autoRefresh(ctx, false)
}

func (l *Lock) autoRefresh(ctx context.Context, immediateRefresh bool) {
	log := l.logger.WithContext(ctx)

	if immediateRefresh {
		if err := l.Refresh(ctx, l.ttl, nil); err != nil {
			log.Warn("Failed to refresh lock immediately.",
				zap.String("key", l.key),
				zap.String("token", l.Token()),
				zap.Error(err),
			)
			return
		}
	}

	// Formula: interval = TTL / 3.
	// This ensures we have 3 attempts to refresh before the lock actually expires.
	interval := l.ttl / 3

	// Safety: minimum interval of 50ms to prevent infinite tight loops if TTL is very small.
	if interval < 50*time.Millisecond {
		interval = l.ttl / 2
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := l.Refresh(ctx, l.ttl, nil); err != nil {
					log.Warn("Failed to refresh lock.",
						zap.String("key", l.key),
						zap.String("token", l.Token()),
						zap.Error(err),
					)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
