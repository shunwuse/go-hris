package infra

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	GetEnv,
	GetLogger,
	GetDatabase,
	NewRequestHandler,
)
