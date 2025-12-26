package infra

import (
	"context"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Cache struct {
	Client *redis.Client

	miniRedis *miniredis.Miniredis
}

func NewCache(config Config, logger *Logger) *Cache {
	addr := config.RedisAddr
	password := config.RedisPassword
	db := config.RedisDB

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

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}

	return &Cache{
		Client:    rdb,
		miniRedis: miniRedis,
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
