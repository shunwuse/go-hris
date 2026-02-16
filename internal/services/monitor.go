package services

import (
	"context"
	"time"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/domains"
	"github.com/shunwuse/go-hris/internal/infra/app"
	"github.com/shunwuse/go-hris/internal/pkg/logger"
	"github.com/shunwuse/go-hris/internal/ports/repository"
	"github.com/shunwuse/go-hris/internal/ports/service"
)

type monitorService struct {
	config            *app.ServiceConfig
	logger            *logger.Logger
	monitorRepository repository.MonitorRepository
}

func NewMonitorService(
	cfg *app.ServiceConfig,
	log *logger.Logger,
	monitorRepository repository.MonitorRepository,
) service.MonitorService {
	return &monitorService{
		config:            cfg,
		logger:            log,
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

	uptime := time.Since(app.AppStartTime).Round(time.Second).String()

	return &domains.Health{
		Status: status,
		Components: domains.HealthComponents{
			Database: dbStatus,
			Redis:    redisStatus,
		},
		Info: domains.HealthInfo{
			Version:     app.Version,
			Environment: s.config.Environment,
			Uptime:      uptime,
			InstanceID:  app.InstanceID,
			Hostname:    app.Hostname,
			GoVersion:   app.GoVersion,
		},
	}
}
