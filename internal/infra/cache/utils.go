package cache

import (
	"context"
	"encoding/json"
	"time"
)

// Fetch is a helper to implement the Cache-Aside pattern.
// It tries to get data from Redis, and if it fails, it calls the fetcher function,
// saves the result to Redis, and returns it.
func Fetch[T any](ctx context.Context, c *Cache, key string, ttl time.Duration, fetcher func() (T, error)) (T, error) {
	// 1. Try to get from cache.
	val, err := c.Client.Get(ctx, key).Result()
	if err == nil {
		var result T
		if err := json.Unmarshal([]byte(val), &result); err == nil {
			return result, nil
		}
	}

	// 2. Cache miss, call fetcher with singleflight.
	v, err, _ := c.sf.Do(key, func() (any, error) {
		val, err := c.Client.Get(ctx, key).Result() // Double check cache
		if err == nil {
			var result T
			if err := json.Unmarshal([]byte(val), &result); err == nil {
				return result, nil
			}
		}

		result, err := fetcher()
		if err != nil {
			return nil, err
		}

		// 3. Save to cache.
		if data, err := json.Marshal(result); err == nil {
			c.Client.Set(ctx, key, data, ttl)
		}

		return result, nil
	})

	if err != nil {
		var zero T
		return zero, err
	}

	return v.(T), nil
}
