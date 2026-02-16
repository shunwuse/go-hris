package config

import (
	"log"
	"sync"

	"github.com/knadh/koanf/v2"
)

type Manager[T any] struct {
	target *T
	config *config
	mu     sync.RWMutex
	hooks  []func(*T)
}

// NewManager creates a new Manager for a generic configuration type T.
func NewManager[T any](opts ...Option) *Manager[T] {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &Manager[T]{
		target: new(T),
		config: o,
	}
}

// Load performs the actual loading of configuration from all sources.
func (m *Manager[T]) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	k := koanf.New(".")

	for _, src := range m.config.sources {
		if err := k.Load(src.provider, src.parser); err != nil {
			log.Printf("warning: failed to load source: %v", err)
		}
	}

	if err := k.Unmarshal("", m.target); err != nil {
		return err
	}

	// Trigger hooks.
	for _, hook := range m.hooks {
		hook(m.target)
	}

	return nil
}

// Config returns the thread-safe configuration target.
func (m *Manager[T]) Config() *T {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.target
}

// OnChange registers a hook to be called when the configuration is reloaded.
func (m *Manager[T]) OnChange(hook func(*T)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.hooks = append(m.hooks, hook)
}
