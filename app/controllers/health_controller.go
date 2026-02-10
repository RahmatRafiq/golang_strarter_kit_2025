package controllers

import (
	"net/http"

	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HealthController struct {
	service *services.HealthService
}

func NewHealthController() *HealthController {
	return &HealthController{
		service: services.NewHealthService(),
	}
}

// GetHealth returns basic health status
// @Summary		Basic health check
// @Description	Returns basic API health status for load balancer health checks
// @Tags			Health
// @Produce		json
// @Success		200	{object}	casts.HealthResponse
// @Router			/health [get]
func (c *HealthController) GetHealth(ctx *gin.Context) {
	health := c.service.GetBasicHealth()

	log.Debug().
		Str("status", string(health.Status)).
		Msg("Basic health check requested")

	ctx.JSON(http.StatusOK, health)
}

// GetDetailedHealth returns comprehensive health status
// @Summary		Detailed health check
// @Description	Returns comprehensive health status including all service dependencies
// @Tags			Health
// @Produce		json
// @Success		200	{object}	casts.DetailedHealthResponse
// @Success		503	{object}	casts.DetailedHealthResponse	"Service unavailable"
// @Router			/health/detailed [get]
func (c *HealthController) GetDetailedHealth(ctx *gin.Context) {
	health := c.service.GetDetailedHealth()

	log.Info().
		Str("status", string(health.Status)).
		Str("uptime", health.Uptime).
		Interface("services", health.Services).
		Msg("Detailed health check requested")

	// Return 503 if unhealthy
	statusCode := http.StatusOK
	if health.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	ctx.JSON(statusCode, health)
}
