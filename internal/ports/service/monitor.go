package service

import (
	"context"

	"github.com/shunwuse/go-hris/internal/domains"
)

type MonitorService interface {
	HealthCheck(ctx context.Context) *domains.Health
}
