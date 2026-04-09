package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	AppID    string `koanf:"app_id"`
	Port     int    `koanf:"port"`
	LogLevel string `koanf:"log_level"`
	DB_Host  string `koanf:"db_host"`
	DB_User  string `koanf:"db_user"`
}

func TestConfig_Precedence(t *testing.T) {
	// Create temporary env files.
	tmpDir := t.TempDir()
	defaultEnv := filepath.Join(tmpDir, "default.env")
	dotEnv := filepath.Join(tmpDir, ".env")

	err := os.WriteFile(defaultEnv, []byte(`
APP_ID=default-id
PORT=8080
LOG_LEVEL=info
DB_HOST=localhost
`), 0644)
	require.NoError(t, err)

	err = os.WriteFile(dotEnv, []byte(`
PORT=9090
LOG_LEVEL=debug
`), 0644)
	require.NoError(t, err)

	// Set OS environment variable (highest precedence).
	t.Setenv("APP_LOG_LEVEL", "error")
	t.Setenv("APP_DB_USER", "postgres")

	mgr := NewManager[TestConfig](
		WithDotEnv(defaultEnv),
		WithDotEnv(dotEnv),
		WithEnv("APP_", "."),
	)
	require.NoError(t, mgr.Load())
	cfg := mgr.Config()

	// 1. AppID should come from default.env (not overridden).
	assert.Equal(t, "default-id", cfg.AppID)
	// 2. Port should come from .env (overrides default.env).
	assert.Equal(t, 9090, cfg.Port)
	// 3. LogLevel should come from OS env (overrides all).
	assert.Equal(t, "error", cfg.LogLevel)
	// 4. DB_Host from default.env.
	assert.Equal(t, "localhost", cfg.DB_Host)
	// 5. DB_User from OS env.
	assert.Equal(t, "postgres", cfg.DB_User)
}

func TestConfig_HotReload_Concurrency(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	err := os.WriteFile(envFile, []byte(`PORT=1000`), 0644)
	require.NoError(t, err)

	mgr := NewManager[TestConfig](
		WithDotEnv(envFile),
		WithWatch(true),
	)

	err = mgr.Load()
	require.NoError(t, err)

	initialCfg := mgr.Config()
	assert.Equal(t, 1000, initialCfg.Port)

	// Test concurrency reading while updating.
	var wg sync.WaitGroup
	start := make(chan struct{})

	// Readers.
	for range 50 {
		wg.Go(func() {
			<-start
			for range 100 {
				_ = mgr.Config().Port
			}
		})
	}

	close(start)

	// Simulation of external update.
	err = os.WriteFile(envFile, []byte(`PORT=2000`), 0644)
	require.NoError(t, err)

	// Manual reload to simulate watcher trigger.
	err = mgr.Load()
	require.NoError(t, err)

	wg.Wait()

	assert.Equal(t, 2000, mgr.Config().Port)
}

func TestConfig_OnChange(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	err := os.WriteFile(envFile, []byte(`PORT=1000`), 0644)
	require.NoError(t, err)

	mgr := NewManager[TestConfig](WithDotEnv(envFile))
	err = mgr.Load()
	require.NoError(t, err)

	changed := false
	mgr.OnChange(func(c *TestConfig) {
		if c.Port == 2000 {
			changed = true
		}
	})

	err = os.WriteFile(envFile, []byte(`PORT=2000`), 0644)
	require.NoError(t, err)

	// Manual reload to simulate watcher trigger.
	err = mgr.Load()
	require.NoError(t, err)

	assert.True(t, changed)
	assert.Equal(t, 2000, mgr.Config().Port)
}

func TestConfig_MissingFile(t *testing.T) {
	// Should not fail if a file is missing unless we want strict mode.
	mgr := NewManager[TestConfig](
		WithDotEnv("non-existent.env"),
	)
	assert.NoError(t, mgr.Load())
}

func TestConfig_Reload_UsesSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")

	err := os.WriteFile(envFile, []byte(`PORT=1000`), 0644)
	require.NoError(t, err)

	mgr := NewManager[TestConfig](WithDotEnv(envFile))
	require.NoError(t, mgr.Load())

	oldCfg := mgr.Config()
	require.Equal(t, 1000, oldCfg.Port)

	err = os.WriteFile(envFile, []byte(`PORT=2000`), 0644)
	require.NoError(t, err)

	require.NoError(t, mgr.Load())

	newCfg := mgr.Config()
	assert.NotSame(t, oldCfg, newCfg)
	assert.Equal(t, 1000, oldCfg.Port)
	assert.Equal(t, 2000, newCfg.Port)
}
