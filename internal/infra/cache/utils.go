package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Fetch is a helper to implement the Cache-Aside pattern.
// It tries to get data from Redis, and if it fails, it calls the fetcher function,
// saves the result to Redis, and returns it.
func Fetch[T any](ctx context.Context, c *Cache, key string, ttl time.Duration, fetcher func() (T, error)) (T, error) {
	var zero T

	if c == nil || c.Client == nil {
		return fetcher()
	}

	// Try to get from cache first.
	if result, ok := fetchCachedValue[T](ctx, c, key); ok {
		return result, nil
	}

	v, err, _ := c.sf.Do(key, func() (any, error) {
		// Double check cache.
		if result, ok := fetchCachedValue[T](ctx, c, key); ok {
			return result, nil
		}

		result, err := fetcher()
		if err != nil {
			return zero, err
		}

		// Store the result in cache.
		storeCachedValue(ctx, c, key, ttl, result)

		return result, nil
	})

	if err != nil {
		return zero, err
	}

	return v.(T), nil
}

func fetchCachedValue[T any](ctx context.Context, c *Cache, key string) (T, bool) {
	var zero T

	val, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			c.warn(ctx, "failed to read cache", zap.String("key", key), zap.Error(err))
		}

		return zero, false
	}

	var result T
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		c.warn(ctx, "failed to decode cached value", zap.String("key", key), zap.Error(err))

		if delErr := c.Client.Del(ctx, key).Err(); delErr != nil {
			c.warn(ctx, "failed to delete corrupted cache entry", zap.String("key", key), zap.Error(delErr))
		}

		return zero, false
	}

	return result, true
}

func storeCachedValue[T any](ctx context.Context, c *Cache, key string, ttl time.Duration, value T) {
	data, err := json.Marshal(value)
	if err != nil {
		c.warn(ctx, "failed to encode cached value", zap.String("key", key), zap.Error(err))
		return
	}

	if err := c.Client.Set(ctx, key, data, ttl).Err(); err != nil {
		c.warn(ctx, "failed to write cache", zap.String("key", key), zap.Error(err))
	}
}
