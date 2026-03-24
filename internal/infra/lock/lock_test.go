package lock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bsm/redislock"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocker(t *testing.T) {
	cfg := &cache.Config{
		UseMiniredis: true,
	}
	log := logger.New()

	c := cache.New(cfg, log)
	defer c.Close() //nolint:errcheck

	locker := New(c, log)
	ctx := context.Background()
	key := "test-lock"

	t.Run("Obtain and Release", func(t *testing.T) {
		lock, err := locker.Obtain(ctx, key, 100*time.Millisecond, nil)
		require.NoError(t, err)
		assert.NotNil(t, lock)

		err = lock.Release(ctx)
		assert.NoError(t, err)
	})

	t.Run("Obtain with too short ttl", func(t *testing.T) {
		lock, err := locker.Obtain(ctx, key, 50*time.Millisecond, nil)
		assert.ErrorIs(t, err, ErrInvalidLockTTL)
		assert.Nil(t, lock)
	})

	t.Run("Mutual Exclusion", func(t *testing.T) {
		lock1, err := locker.Obtain(ctx, key, 1*time.Second, nil)
		require.NoError(t, err)
		assert.NotNil(t, lock1)
		defer lock1.Release(ctx) //nolint:errcheck

		lock2, err := locker.Obtain(ctx, key, 1*time.Second, nil)
		assert.ErrorIs(t, err, ErrLockNotObtained)
		assert.Nil(t, lock2)
	})

	t.Run("Refresh", func(t *testing.T) {
		lock, err := locker.Obtain(ctx, key, 1*time.Second, nil)
		require.NoError(t, err)
		defer lock.Release(ctx) //nolint:errcheck

		err = lock.Refresh(ctx, 2*time.Second, nil)
		assert.NoError(t, err)
	})

	t.Run("AutoRefresh", func(t *testing.T) {
		// Acquire a 300ms lock
		lock, err := locker.Obtain(ctx, key, 300*time.Millisecond, nil)
		require.NoError(t, err)
		defer lock.Release(ctx) //nolint:errcheck

		// Start auto-refresh: will automatically use interval 100ms (300/3)
		lock.AutoRefresh(ctx)

		// Wait for 500ms. Without auto-refresh, the lock would have expired.
		time.Sleep(500 * time.Millisecond)

		// Try to obtain the same lock - it should still be held by us (fail to obtain)
		lock2, err := locker.Obtain(ctx, key, 100*time.Millisecond, nil)
		assert.ErrorIs(t, err, ErrLockNotObtained)
		assert.Nil(t, lock2)
	})

	t.Run("ObtainAutoRefresh keeps lock alive", func(t *testing.T) {
		autoCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		lock, err := locker.ObtainAutoRefresh(autoCtx, key, 300*time.Millisecond, nil)
		require.NoError(t, err)
		require.NotNil(t, lock)

		time.Sleep(500 * time.Millisecond)

		lock2, err := locker.Obtain(ctx, key, 100*time.Millisecond, nil)
		assert.ErrorIs(t, err, ErrLockNotObtained)
		assert.Nil(t, lock2)

		cancel()
		time.Sleep(50 * time.Millisecond)

		err = lock.Release(ctx)
		if err != nil {
			assert.True(t, errors.Is(err, redislock.ErrLockNotHeld))
		}
	})
}
