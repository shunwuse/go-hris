package controllers

import (
	"net/http"

	"github.com/shunwuse/go-hris/internal/constants"
	"github.com/shunwuse/go-hris/internal/dtos"
	"github.com/shunwuse/go-hris/internal/http/response"
	"github.com/shunwuse/go-hris/internal/infra"
	"github.com/shunwuse/go-hris/internal/ports/service"
)

// HealthController struct
type HealthController struct {
	logger        *infra.Logger
	healthService service.HealthService
}

func NewHealthController(
	logger *infra.Logger,
	healthService service.HealthService,
) *HealthController {
	return &HealthController{
		logger:        logger,
		healthService: healthService,
	}
}

func (c *HealthController) Check(w http.ResponseWriter, r *http.Request) {
	c.logger.WithContext(r.Context()).Info("health check controller invoked")

	health := c.healthService.Check(r.Context())

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
