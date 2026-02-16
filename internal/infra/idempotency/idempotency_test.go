package idempotency

import (
	"context"
	"testing"
	"time"

	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdempotencyManager(t *testing.T) {
	// Setup: Use Miniredis for testing.
	cfg := &cache.Config{
		UseMiniredis: true,
	}
	log := logger.New()

	// We need to access the internal miniRedis instance or manually manage time.
	c := cache.New(cfg, log)
	defer c.Close()

	// Since miniRedis field in Cache struct is private and not exported,
	// we use direct deletion to simulate expiration in this test case.
	mgr := New(log, c)
	ctx := context.Background()
	key := "test-idempotency-key"

	t.Run("Record Not Found", func(t *testing.T) {
		record, err := mgr.Get(ctx, "non-existent")
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, record)
	})

	t.Run("Set and Get Success", func(t *testing.T) {
		expected := &Record{
			Status: 200,
			Body:   []byte("{\"message\":\"success\"}"),
			Header: map[string]string{"Content-Type": "application/json"},
		}

		err := mgr.Set(ctx, key, expected, time.Minute)
		require.NoError(t, err)

		got, err := mgr.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, expected.Status, got.Status)
		assert.Equal(t, expected.Body, got.Body)
		assert.Equal(t, expected.Header["Content-Type"], got.Header["Content-Type"])
	})

	t.Run("Expiration", func(t *testing.T) {
		// Use a short TTL.
		shortKey := "short-lived-key"
		record := &Record{Status: 201, Body: []byte("created")}

		ttl := 1 * time.Second
		err := mgr.Set(ctx, shortKey, record, ttl)
		require.NoError(t, err)

		// Get the actual key name in Redis.
		// Since we know the helper method is "idempotency:" + key.
		realKey := "idempotency:" + shortKey

		// Manually delete from Redis client to simulate expiration.
		err = c.Client.Del(ctx, realKey).Err()
		require.NoError(t, err)

		// Should fail and return ErrNotFound.
		_, err = mgr.Get(ctx, shortKey)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}
