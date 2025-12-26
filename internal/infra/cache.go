package infra

import (
	"context"
	"encoding/json"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Cache struct {
	Client *redis.Client
	Prefix string

	miniRedis *miniredis.Miniredis
}

func NewCache(config Config, logger *Logger) *Cache {
	addr := config.RedisAddr
	password := config.RedisPassword
	db := config.RedisDB
	prefix := config.RedisPrefix

	var miniRedis *miniredis.Miniredis

	if config.UseMiniredis {
		mini, err := miniredis.Run()
		if err != nil {
			logger.Fatal("failed to start miniredis", zap.Error(err))
		}

		miniRedis = mini

		addr = miniRedis.Addr()
		password = ""
		db = 0

		logger.Info("using embedded miniredis", zap.String("addr", addr))
	} else {
		logger.Info("connecting to external redis", zap.String("addr", addr))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if prefix != "" {
		rdb.AddHook(&prefixHook{prefix: prefix})
	}

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}

	return &Cache{
		Client:    rdb,
		Prefix:    prefix,
		miniRedis: miniRedis,
	}
}

// prefixHook is a redis.Hook that automatically adds a prefix to keys.
type prefixHook struct {
	prefix string
}

func (h *prefixHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (h *prefixHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		h.addPrefix(cmd)
		return next(ctx, cmd)
	}
}

func (h *prefixHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		for _, cmd := range cmds {
			h.addPrefix(cmd)
		}
		return next(ctx, cmds)
	}
}

func (h *prefixHook) addPrefix(cmd redis.Cmder) {
	args := cmd.Args()
	if len(args) < 2 {
		return
	}

	// Most Redis commands have the key at index 1.
	// e.g., GET key, SET key val, DEL key1 key2...
	// For commands with multiple keys like DEL, MGET, we handle them specifically if needed.
	// Here we handle the most common case: the first argument after the command name is the key.
	if key, ok := args[1].(string); ok {
		args[1] = h.prefix + ":" + key
	}
}

// Key prepends the global prefix to the given key.
func (c *Cache) Key(key string) string {
	if c.Prefix == "" {
		return key
	}
	return c.Prefix + ":" + key
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
// Note: The prefix is automatically added by the Redis Hook.
func CacheGetOrSet[T any](ctx context.Context, c *Cache, key string, ttl time.Duration, fetcher func() (T, error)) (T, error) {
	// 1. Try to get from cache.
	val, err := c.Client.Get(ctx, key).Result()
	if err == nil {
		var result T
		if err := json.Unmarshal([]byte(val), &result); err == nil {
			return result, nil
		}
	}

	// 2. Cache miss, call fetcher.
	result, err := fetcher()
	if err != nil {
		return result, err
	}

	// 3. Save to cache.
	if data, err := json.Marshal(result); err == nil {
		c.Client.Set(ctx, key, data, ttl)
	}

	return result, nil
}
