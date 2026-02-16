package main

import (
	"os"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/infra/app"
	"github.com/shunwuse/go-hris/internal/pkg/config"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"go.uber.org/zap"
)

func initLogger(mgr *config.Manager[Config]) *logger.Logger {
	cfg := mgr.Config()

	log := logger.New(
		logger.WithConfig(logger.Config{
			Level:      cfg.Log.Level,
			FilePath:   cfg.Log.FilePath,
			MaxSize:    cfg.Log.MaxSize,
			MaxBackups: cfg.Log.MaxBackups,
			MaxAge:     cfg.Log.MaxAge,
			Compress:   cfg.Log.Compress,
		}),
		logger.WithFields(
			zap.String("instance_id", app.InstanceID),
			zap.String("hostname", app.Hostname),
			zap.Int("pid", os.Getpid()),
			zap.String("environment", cfg.Service.Environment),
		),
		logger.WithConsole(cfg.Service.Environment == constants.EnvDevelopment),
	)

	// Register reload hook for logger.
	mgr.OnChange(func(newCfg *Config) {
		log.UpdateLevel(newCfg.Log.Level)
	})

	return log
}
