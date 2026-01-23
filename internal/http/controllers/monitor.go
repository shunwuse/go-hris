package controllers

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/infra/logger"
	"github.com/shunwuse/go-hris/internal/ports/service"
)

// MonitorController struct
type MonitorController struct {
	logger         *logger.Logger
	monitorService service.MonitorService
}

func NewMonitorController(
	log *logger.Logger,
	monitorService service.MonitorService,
) *MonitorController {
	return &MonitorController{
		logger:         log,
		monitorService: monitorService,
	}
}

func (c *MonitorController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	c.logger.WithContext(r.Context()).Info("health check controller invoked")

	health := c.monitorService.HealthCheck(r.Context())

	res := dtos.HealthResponse{
		Status: health.Status,
		Components: dtos.HealthComponentsResponse{
			Database: health.Components.Database,
			Redis:    health.Components.Redis,
		},
		Info: dtos.HealthInfoResponse{
			Version:     health.Info.Version,
			Environment: health.Info.Environment,
			Uptime:      health.Info.Uptime,
			InstanceID:  health.Info.InstanceID,
			Hostname:    health.Info.Hostname,
			GoVersion:   health.Info.GoVersion,
		},
	}

	if health.Status == constants.StatusDown {
		response.ServiceUnavailable(w, res)
		return
	}

	response.OK(w, res)
}
