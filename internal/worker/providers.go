package worker

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/worker/consumer"
	"github.com/shunwuse/go-hris/internal/worker/scheduler"
)

var ProviderSet = wire.NewSet(
	scheduler.ProviderSet,
	consumer.ProviderSet,
)
