package consumer

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/worker/consumer/handlers"
)

var ProviderSet = wire.NewSet(
	handlers.ProviderSet,
	NewConsumer,
)
