package main

import (
	"github.com/shunwuse/go-hris/internal/infra/app"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
)

type Config struct {
	Log      logger.Config     `koanf:",squash"`
	Database database.Config   `koanf:",squash"`
	Cache    cache.Config      `koanf:",squash"`
	Service  app.ServiceConfig `koanf:",squash"`
	Auth     app.AuthConfig    `koanf:",squash"`
}
