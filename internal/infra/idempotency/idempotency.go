package idempotency

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shunwuse/go-hris/internal/infra/cache"
)

var (
	ErrNotFound = errors.New("idempotency record not found")
)

type Manager struct {
	cache *cache.Cache
}

func New(cache *cache.Cache) *Manager {
	return &Manager{cache: cache}
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

		return nil, err
	}

	var record Record
	if err := json.Unmarshal(val, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (m *Manager) Set(ctx context.Context, key string, record *Record, ttl time.Duration) error {
	val, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return m.cache.Client.Set(ctx, m.key(key), val, ttl).Err()
}

func (m *Manager) key(k string) string {
	return "idempotency:" + k
}
