package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/shunwuse/go-hris/internal/infra/config"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type Cache struct {
	Client *redis.Client

	miniRedis *miniredis.Miniredis
	sf        singleflight.Group
}

func NewCache(cfg config.Config, log *logger.Logger) *Cache {
	addr := cfg.RedisAddr
	password := cfg.RedisPassword
	db := cfg.RedisDB

	var miniRedis *miniredis.Miniredis

	if cfg.UseMiniredis {
		mini, err := miniredis.Run()
		if err != nil {
			log.Fatal("failed to start miniredis", zap.Error(err))
		}

		miniRedis = mini

		addr = miniRedis.Addr()
		password = ""
		db = 0

		log.Info("using embedded miniredis", zap.String("addr", addr))
	} else {
		log.Info("connecting to external redis", zap.String("addr", addr))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatal("failed to connect to redis", zap.Error(err))
	}

	return &Cache{
		Client:    rdb,
		miniRedis: miniRedis,
		sf:        singleflight.Group{},
	}
}

func (c *Cache) Close() error {
	if err := c.Client.Close(); err != nil {
		return err
	}
	if c.miniRedis != nil {
		c.miniRedis.Close()
	}
	return nil
}

// CacheGetOrSet is a helper to implement the Cache-Aside pattern.
// It tries to get data from Redis, and if it fails, it calls the fetcher function,
// saves the result to Redis, and returns it.
func CacheGetOrSet[T any](ctx context.Context, c *Cache, key string, ttl time.Duration, fetcher func() (T, error)) (T, error) {
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
