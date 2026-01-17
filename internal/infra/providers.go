package infra

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	GetConfig,
	GetLogger,
	GetDatabase,
	NewTransactor,
	NewCache,
	NewMetrics,
	NewRequestHandler,
)
