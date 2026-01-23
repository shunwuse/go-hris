package config

import (
	"sync"
)

var (
	configReloadMu    sync.Mutex
	configReloadHooks []func(*Config)
)

// ReloadConfig reloads the configuration from all sources and triggers registered hooks.
func ReloadConfig() error {
	configReloadMu.Lock()
	defer configReloadMu.Unlock()

	newConfig := loadConfig()
	globalConfig.Store(newConfig)

	// Trigger hooks
	for _, hook := range configReloadHooks {
		hook(newConfig)
	}

	return nil
}

// OnConfigChange registers a hook to be called when the configuration is reloaded.
func OnConfigChange(hook func(*Config)) {
	configReloadMu.Lock()
	defer configReloadMu.Unlock()
	configReloadHooks = append(configReloadHooks, hook)
}
