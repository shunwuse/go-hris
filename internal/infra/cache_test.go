package infra

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheGetOrSet(t *testing.T) {
	logger := GetLogger()
	config := Config{UseMiniredis: true}
	cache := NewCache(config, logger)
	defer func() { _ = cache.Close() }()

	ctx := context.Background()

	type TestUser struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	t.Run("CacheMissAndHit", func(t *testing.T) {
		key := "user:1"
		expected := TestUser{ID: 1, Name: "Alice"}
		fetcherCalled := 0
		fetcher := func() (TestUser, error) {
			fetcherCalled++
			return expected, nil
		}

		// 1. First call: Cache Miss
		got, err := CacheGetOrSet(ctx, cache, key, time.Minute, fetcher)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcherCalled)
		assert.Equal(t, expected, got)

		// 2. Second call: Cache Hit
		got, err = CacheGetOrSet(ctx, cache, key, time.Minute, fetcher)
		require.NoError(t, err)
		assert.Equal(t, 1, fetcherCalled)
		assert.Equal(t, expected, got)
	})

	t.Run("Expiration", func(t *testing.T) {
		key := "user:2"
		expected := TestUser{ID: 2, Name: "Bob"}
		fetcherCalled := 0
		fetcher := func() (TestUser, error) {
			fetcherCalled++
			return expected, nil
		}

		// 1. Set cache with short TTL
		ttl := time.Second
		_, _ = CacheGetOrSet(ctx, cache, key, ttl, fetcher)

		// 2. Fast forward time in the cache's internal miniredis
		cache.miniRedis.FastForward(ttl + time.Second)

		// 3. Call again: should be a Cache Miss
		_, err := CacheGetOrSet(ctx, cache, key, ttl, fetcher)
		require.NoError(t, err)
		assert.Equal(t, 2, fetcherCalled)
	})

	t.Run("FetcherError", func(t *testing.T) {
		key := "user:error"
		fetcher := func() (TestUser, error) {
			return TestUser{}, context.DeadlineExceeded // Using a standard error for comparison
		}

		_, err := CacheGetOrSet(ctx, cache, key, time.Minute, fetcher)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
