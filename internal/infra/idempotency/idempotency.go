package idempotency

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("idempotency record not found")
)

type Manager struct {
	log   *logger.Logger
	cache *cache.Cache
}

func New(log *logger.Logger, cache *cache.Cache) *Manager {
	return &Manager{
		log:   log,
		cache: cache,
	}
}

type Record struct {
	Status int               `json:"status"`
	Body   []byte            `json:"body"`
	Header map[string]string `json:"header"`
}

func (m *Manager) Get(ctx context.Context, key string) (*Record, error) {
	val, err := m.cache.Client.Get(ctx, m.key(key)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}

		// Log error to help trace potential connection or timeout issues.
		m.log.WithContext(ctx).Error("failed to get idempotency record",
			zap.String("key", key),
			zap.Error(err),
		)

		return nil, err
	}

	var record Record
	if err := json.Unmarshal(val, &record); err != nil {
		// Data corruption or format changes should be logged as error.
		m.log.WithContext(ctx).Error("failed to unmarshal idempotency record",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, err
	}

	return &record, nil
}

func (m *Manager) Set(ctx context.Context, key string, record *Record, ttl time.Duration) error {
	val, err := json.Marshal(record)
	if err != nil {
		// Log serialization error for debugging.
		m.log.WithContext(ctx).Error("failed to marshal idempotency record",
			zap.String("key", key),
			zap.Error(err),
		)
		return err
	}

	if err := m.cache.Client.Set(ctx, m.key(key), val, ttl).Err(); err != nil {
		// Redis write failures should be tracked.
		m.log.WithContext(ctx).Error("failed to set idempotency record",
			zap.String("key", key),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (m *Manager) key(k string) string {
	return "idempotency:" + k
}
