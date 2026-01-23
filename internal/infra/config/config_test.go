package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Reload(t *testing.T) {
	// 1. Initialize global state
	_ = GetConfig()
	initialEnv := GetConfig().Environment

	// 2. Prepare for reload with a new environment variable value using t.Setenv
	testEnvValue := "reload-test-env"
	t.Setenv("APP_ENV", testEnvValue)

	// 3. Register a hook to verify it's called during reload
	hookCalled := false
	var hookCfg *Config
	OnConfigChange(func(cfg *Config) {
		hookCalled = true
		hookCfg = cfg
	})

	// 4. Trigger reload
	err := ReloadConfig()
	assert.NoError(t, err)

	// 5. Assertions
	assert.True(t, hookCalled, "Hook should be called on reload")
	assert.NotNil(t, hookCfg)
	assert.Equal(t, testEnvValue, hookCfg.Environment, "Hook should receive the new value")

	// Verify GetConfig returns the updated value
	current := GetConfig()
	assert.Equal(t, testEnvValue, current.Environment, "GetConfig should return the new value")
	assert.NotEqual(t, initialEnv, current.Environment, "Value should have changed from initial")
}

func TestConfig_ConcurrentAccess(t *testing.T) {
	// Verify that GetConfig is thread-safe with Reload
	for range 100 {
		go func() {
			_ = GetConfig()
		}()
		go func() {
			_ = ReloadConfig()
		}()
	}
}
