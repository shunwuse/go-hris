package services

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/repository"
	"github.com/shunwuse/go-hris/internal/ports/service"
)

type monitorService struct {
	logger            *infra.Logger
	monitorRepository repository.MonitorRepository
}

func NewMonitorService(
	logger *infra.Logger,
	monitorRepository repository.MonitorRepository,
) service.MonitorService {
	return &monitorService{
		logger:            logger,
		monitorRepository: monitorRepository,
	}
}

func (s *monitorService) HealthCheck(ctx context.Context) *domains.Health {
	dbStatus := constants.StatusUp
	if ok := s.monitorRepository.CheckDatabase(ctx); !ok {
		dbStatus = constants.StatusDown
	}

	redisStatus := constants.StatusUp
	if ok := s.monitorRepository.CheckRedis(ctx); !ok {
		redisStatus = constants.StatusDown
	}

	status := constants.StatusUp
	if dbStatus == constants.StatusDown ||
		redisStatus == constants.StatusDown {
		status = constants.StatusDown
	}

	uptime := time.Since(infra.AppStartTime).Round(time.Second).String()

	return &domains.Health{
		Status: status,
		Components: domains.HealthComponents{
			Database: dbStatus,
			Redis:    redisStatus,
		},
		Info: domains.HealthInfo{
			Version:     infra.Version,
			Environment: infra.GetConfig().Environment,
			Uptime:      uptime,
			InstanceID:  infra.InstanceID,
			Hostname:    infra.Hostname,
			GoVersion:   infra.GoVersion,
		},
	}
}
