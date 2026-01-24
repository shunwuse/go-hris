package infra

import (
	"github.com/google/wire"
	termalerter "github.com/shunwuse/go-hris/internal/infra/alerter/provider/terminal"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/infra/config"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/infra/handler"
	"github.com/shunwuse/go-hris/internal/infra/idempotency"
	"github.com/shunwuse/go-hris/internal/infra/lifecycle"
	"github.com/shunwuse/go-hris/internal/infra/lock"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/shunwuse/go-hris/internal/infra/metrics"
	"github.com/shunwuse/go-hris/internal/infra/routine"
)

var ProviderSet = wire.NewSet(
	config.Get,
	logger.L,
	database.DB,
	database.NewTransactor,
	cache.New,
	lock.New,
	idempotency.New,
	metrics.New,
	handler.New,
	termalerter.New,
	routine.New,
	lifecycle.New,
)
