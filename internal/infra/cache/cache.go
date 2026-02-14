package cache

import (
	"context"

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

func New(cfg *config.Config, log *logger.Logger) *Cache {
	addr := cfg.Cache.RedisAddr
	password := cfg.Cache.RedisPassword
	db := cfg.Cache.RedisDB

	var miniRedis *miniredis.Miniredis

	if cfg.Cache.UseMiniredis {
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
