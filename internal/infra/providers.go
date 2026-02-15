package infra

import (
	"github.com/google/wire"
	"github.com/shunwuse/go-hris/internal/infra/alerter"
	"github.com/shunwuse/go-hris/internal/infra/cache"
	"github.com/shunwuse/go-hris/internal/infra/database"
	"github.com/shunwuse/go-hris/internal/infra/handler"
	"github.com/shunwuse/go-hris/internal/infra/idempotency"
	"github.com/shunwuse/go-hris/internal/infra/lifecycle"
	"github.com/shunwuse/go-hris/internal/infra/lock"
	"github.com/shunwuse/go-hris/internal/infra/metrics"
	"github.com/shunwuse/go-hris/internal/infra/routine"
)

var ProviderSet = wire.NewSet(
	database.New,
	database.NewTransactor,
	cache.New,
	lock.New,
	idempotency.New,
	metrics.New,
	handler.New,
	alerter.NewMultiAlerter,
	routine.New,
	lifecycle.New,
)
