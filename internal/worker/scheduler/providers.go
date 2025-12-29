package scheduler

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/worker/scheduler/jobs"
)

var ProviderSet = wire.NewSet(
	jobs.ProviderSet,
	NewScheduler,
)
