package config

import (
	"sync"
)

var (
	configReloadMu    sync.Mutex
	configReloadHooks []func(*Config)
)

// Reload reloads the configuration from all sources and triggers registered hooks.
func Reload() error {
	configReloadMu.Lock()
	defer configReloadMu.Unlock()

	newConfig, err := load()
	if err != nil {
		return err
	}
	instance.Store(newConfig)

	// Trigger hooks.
	for _, hook := range configReloadHooks {
		hook(newConfig)
	}

	return nil
}

// OnChange registers a hook to be called when the configuration is reloaded.
func OnChange(hook func(*Config)) {
	configReloadMu.Lock()
	defer configReloadMu.Unlock()
	configReloadHooks = append(configReloadHooks, hook)
}
