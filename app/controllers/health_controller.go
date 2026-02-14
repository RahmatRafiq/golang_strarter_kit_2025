package controllers

import (
	"net/http"
	"strconv"

	"golang_starter_kit_2025/app/casts"
	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/facades"

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

// GetHistory returns historical health data for services
// @Summary		Get service health history
// @Description	Returns aggregated daily health status for the last N days
// @Tags			Health
// @Produce		json
// @Param			service	query		string	false	"Service name (database, redis, api)"
// @Param			days	query		int		false	"Number of days (default: 90)"
// @Success		200		{array}		casts.ServiceHistoryResponse
// @Router			/health/history [get]
func (c *HealthController) GetHistory(ctx *gin.Context) {
	service := ctx.Query("service")
	days := 90 // default

	if daysParam := ctx.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	historyService := services.NewHealthHistoryService()

	if service != "" {
		// Single service
		history := historyService.GetDailyAggregated(service, days)
		ctx.JSON(http.StatusOK, history)
	} else {
		// All services
		serviceList := []string{"database", "redis", "api"}
		result := make([]casts.ServiceHistoryResponse, 0)
		for _, svc := range serviceList {
			result = append(result, historyService.GetDailyAggregated(svc, days))
		}
		ctx.JSON(http.StatusOK, result)
	}
}

// GetQueueHealth returns queue worker health and metrics
// @Summary		Queue health check
// @Description	Returns queue worker metrics and health status
// @Tags			Health
// @Produce		json
// @Success		200	{object}	casts.QueueHealthResponse
// @Success		503	{object}	casts.QueueHealthResponse	"Queue not enabled"
// @Router			/health/queue [get]
func (c *HealthController) GetQueueHealth(ctx *gin.Context) {
	// Check if queue is enabled
	if helpers.GetEnv("QUEUE_ENABLED", "false") != "true" {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message: "Queue system is not enabled",
		}, http.StatusServiceUnavailable)
		return
	}

	// Get worker manager from facades
	workerManager := facades.GetWorkerManager()
	if workerManager == nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message: "Worker manager not initialized",
		}, http.StatusServiceUnavailable)
		return
	}

	// Get metrics
	metrics := workerManager.GetMetrics()
	if metrics == nil {
		response := casts.QueueHealthResponse{
			Status:  "healthy",
			Message: "Queue is running (metrics disabled)",
		}
		helpers.ResponseSuccess(ctx, &helpers.ResponseParams[casts.QueueHealthResponse]{
			Item: &response,
		}, http.StatusOK)
		return
	}

	// Calculate health status
	status := "healthy"
	failureRate := float64(0)
	if metrics.TasksProcessed > 0 {
		failureRate = float64(metrics.TasksFailed) / float64(metrics.TasksProcessed) * 100
	}

	// Consider unhealthy if failure rate > 50%
	if failureRate > 50 {
		status = "degraded"
	}

	response := casts.QueueHealthResponse{
		Status:           status,
		TasksProcessed:   metrics.TasksProcessed,
		TasksFailed:      metrics.TasksFailed,
		TasksRetried:     metrics.TasksRetried,
		AverageLatencyMs: metrics.AverageLatencyMs,
		ActiveWorkers:    metrics.ActiveWorkers,
		FailureRate:      failureRate,
		LastUpdated:      metrics.LastUpdated,
		QueueDepth:       metrics.QueueDepth,
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[casts.QueueHealthResponse]{
		Item: &response,
	}, http.StatusOK)

	log.Debug().
		Str("status", status).
		Int64("tasks_processed", metrics.TasksProcessed).
		Float64("failure_rate", failureRate).
		Msg("Queue health check requested")
}
