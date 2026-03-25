package lock

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/shunwuse/go-hris/internal/infra/app"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
)

const MinLockTTL = 100 * time.Millisecond

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
	if ttl < MinLockTTL {
		return nil, ErrInvalidLockTTL
	}

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
		lock:   lock,
		ttl:    ttl,
		key:    key,
		logger: l.logger,
	}, nil
}

// ObtainAutoRefresh obtains a lock and immediately starts auto-refresh.
// This minimizes the gap between lock acquisition and keepalive startup.
func (l *Locker) ObtainAutoRefresh(ctx context.Context, key string, ttl time.Duration, options *redislock.Options) (*Lock, error) {
	lock, err := l.Obtain(ctx, key, ttl, options)
	if err != nil {
		return nil, err
	}

	lock.autoRefresh(ctx, true)

	return lock, nil
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

// Do executes under lock protection.
func (l *Locker) Do(ctx context.Context, key string, ttl time.Duration, fn func(context.Context) error) error {
	return l.DoWithOptions(ctx, key, ttl, nil, fn)
}

// DoWithOptions executes under lock protection with custom lock options.
func (l *Locker) DoWithOptions(ctx context.Context, key string, ttl time.Duration, options *redislock.Options, fn func(context.Context) error) error {
	if fn == nil {
		return ErrNilLockFunction
	}
	if l == nil {
		return ErrNilLocker
	}

	lock, err := l.ObtainAutoRefresh(ctx, key, ttl, options)
	if err != nil {
		return err
	}

	defer func() {
		releaseErr := lock.Release(ctx)
		if releaseErr != nil && !errors.Is(releaseErr, redislock.ErrLockNotHeld) {
			l.logger.WithContext(ctx).Warn("Failed to release lock in lock helper.")
		}
	}()

	return fn(ctx)
}

var (
	// ErrLockNotObtained is a proxy for redislock.ErrNotObtained.
	ErrLockNotObtained = redislock.ErrNotObtained

	// ErrInvalidLockTTL indicates lock TTL is too short for safe execution.
	ErrInvalidLockTTL = errors.New("lock: invalid ttl")

	// ErrLockLost indicates the lock was lost during execution.
	ErrLockLost = errors.New("lock: lock lost during execution")

	// ErrNilLockFunction indicates nil callback passed to lock helpers.
	ErrNilLockFunction = errors.New("lock: nil function")

	// ErrNilLocker indicates nil locker passed to lock helpers.
	ErrNilLocker = errors.New("lock: nil locker")
)
