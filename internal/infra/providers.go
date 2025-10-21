package infra

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	GetConfig,
	GetLogger,
	GetDatabase,
	NewRequestHandler,
)
