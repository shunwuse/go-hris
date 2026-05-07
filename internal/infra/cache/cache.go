package cache

import (
	"context"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type Cache struct {
	logger *logger.Logger
	Client *redis.Client

	miniRedis *miniredis.Miniredis
	sf        singleflight.Group
}

func New(cfg *Config, log *logger.Logger) *Cache {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal("failed to connect to redis", zap.Error(err))
	}

	return &Cache{
		logger:    log,
		Client:    rdb,
		miniRedis: miniRedis,
		sf:        singleflight.Group{},
	}
}

func (c *Cache) Close() error {
	if c == nil {
		return nil
	}

	if c.Client == nil {
		if c.miniRedis != nil {
			c.miniRedis.Close()
		}
		return nil
	}

	if err := c.Client.Close(); err != nil {
		return err
	}
	if c.miniRedis != nil {
		c.miniRedis.Close()
	}
	return nil
}

func (c *Cache) warn(ctx context.Context, msg string, fields ...zap.Field) {
	if c == nil || c.logger == nil {
		return
	}

	c.logger.WithContext(ctx).Warn(msg, fields...)
}
