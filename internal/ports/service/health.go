package service

import (
	"context"

	"github.com/shunwuse/go-hris/internal/domains"
)

type HealthService interface {
	Check(ctx context.Context) *domains.Health
}
