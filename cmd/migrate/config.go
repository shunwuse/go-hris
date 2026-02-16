package main

import (
	"github.com/shunwuse/go-hris/internal/infra/database"
)

type Config struct {
	Database database.Config `koanf:",squash"`
}
