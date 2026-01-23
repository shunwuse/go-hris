package infra

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/infra/config"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/infra/handler"
	"github.com/shunwuse/go-hris/internal/infra/lifecycle"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/shunwuse/go-hris/internal/infra/metrics"
	"github.com/shunwuse/go-hris/internal/infra/routine"
)

var ProviderSet = wire.NewSet(
	config.GetConfig,
	logger.GetLogger,
	database.GetDatabase,
	database.NewTransactor,
	cache.NewCache,
	metrics.NewMetrics,
	handler.NewRequestHandler,
	routine.NewRoutineGroup,
	lifecycle.NewLifecycle,
)
