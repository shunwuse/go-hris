package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Reload(t *testing.T) {
	// 1. Initialize global state
	cfg, _ := Load()
	initialEnv := cfg.Service.Environment

	// 2. Prepare for reload with a new environment variable value using t.Setenv
	testEnvValue := "reload-test-env"
	t.Setenv("APP_ENV", testEnvValue)

	// 3. Register a hook to verify it's called during reload
	hookCalled := false
	var hookCfg *Config
	OnChange(func(cfg *Config) {
		hookCalled = true
		hookCfg = cfg
	})

	// 4. Trigger reload
	err := Reload()
	assert.NoError(t, err)

	// 5. Assertions
	assert.True(t, hookCalled, "Hook should be called on reload")
	assert.NotNil(t, hookCfg)
	assert.Equal(t, testEnvValue, hookCfg.Service.Environment, "Hook should receive the new value")

	// Verify current state returns the updated value
	current, _ := Load()
	assert.Equal(t, testEnvValue, current.Service.Environment, "Load should return the new value")
	assert.NotEqual(t, initialEnv, current.Service.Environment, "Value should have changed from initial")
}

func TestConfig_ConcurrentAccess(t *testing.T) {
	// Verify that Load is thread-safe with Reload
	for range 100 {
		go func() {
			_, _ = Load()
		}()
		go func() {
			_ = Reload()
		}()
	}
}
