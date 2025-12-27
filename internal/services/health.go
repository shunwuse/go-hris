package services

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
	"github.com/shunwuse/go-hris/internal/repositories"
)

type healthService struct {
	logger           *infra.Logger
	healthRepository *repositories.HealthRepository
}

func NewHealthService(
	logger *infra.Logger,
	healthRepository *repositories.HealthRepository,
) service.HealthService {
	return &healthService{
		logger:           logger,
		healthRepository: healthRepository,
	}
}

func (s *healthService) Check(ctx context.Context) *domains.Health {
	dbStatus := constants.StatusUp
	if ok := s.healthRepository.CheckDatabase(ctx); !ok {
		dbStatus = constants.StatusDown
	}

	redisStatus := constants.StatusUp
	if ok := s.healthRepository.CheckRedis(ctx); !ok {
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
