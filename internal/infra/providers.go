package infra

import (
	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	NewEnv,
	GetLogger,
	GetDatabase,
	NewRequestHandler,
)
