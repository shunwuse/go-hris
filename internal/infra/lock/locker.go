package lock

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/shunwuse/go-hris/internal/infra/app"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/infra/logger"
)

// Locker provides distributed locking capabilities using Redis.
type Locker struct {
	logger *logger.Logger

	client *redislock.Client
	cache  *cache.Cache

	defaultMetadata string
}

// New creates a new Locker instance.
func New(cache *cache.Cache, log *logger.Logger) *Locker {
	return &Locker{
		client:          redislock.New(cache.Client),
		cache:           cache,
		logger:          log,
		defaultMetadata: app.Hostname + ":" + app.InstanceID,
	}
}

// Obtain tries to obtain a lock with the given key and TTL.
// It returns a Lock instance if successful, or an error if the lock cannot be obtained.
func (l *Locker) Obtain(ctx context.Context, key string, ttl time.Duration, options *redislock.Options) (*Lock, error) {
	if options == nil {
		options = &redislock.Options{}
	}
	if options.Metadata == "" {
		options.Metadata = l.defaultMetadata
	}

	lock, err := l.client.Obtain(ctx, key, ttl, options)
	if err != nil {
		return nil, err
	}

	return &Lock{
		lock: lock,
		ttl:  ttl,
	}, nil
}

// PeekMetadata looks up the metadata of a lock without acquiring it.
// This allows you to see "who" currently holds the lock if they set the Metadata field.
func (l *Locker) PeekMetadata(ctx context.Context, key string) (string, error) {
	val, err := l.cache.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	// bsm/redislock stores its token (22 chars) followed by metadata.
	if len(val) <= 22 {
		return "", nil
	}

	return val[22:], nil
}

var (
	// ErrLockNotObtained is a proxy for redislock.ErrNotObtained.
	ErrLockNotObtained = redislock.ErrNotObtained
)
