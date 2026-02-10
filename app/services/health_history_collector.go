package services

import (
	"time"

	"golang_starter_kit_2025/app/helpers"

	"github.com/rs/zerolog/log"
)

// StartHealthHistoryCollector starts a background goroutine to collect and store health checks
func StartHealthHistoryCollector() {
	interval := time.Duration(helpers.GetEnvInt("HEALTH_CHECK_INTERVAL", 5)) * time.Minute
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		// Run immediately on startup
		collectAndStoreHealth()

		// Then run periodically
		for range ticker.C {
			collectAndStoreHealth()
		}
	}()

	log.Info().
		Dur("interval", interval).
		Msg("Health history collector started")
}

// collectAndStoreHealth collects current health status and stores in Redis
func collectAndStoreHealth() {
	// Only run if Redis is available
	if helpers.RedisClient == nil {
		log.Debug().Msg("Skipping health history collection (Redis unavailable)")
		return
	}

	// Run health checks
	healthService := NewHealthService()
	health := healthService.GetDetailedHealth()

	// Store each service status
	historyService := NewHealthHistoryService()
	for serviceName, serviceData := range health.Services {
		err := historyService.StoreHealthCheck(
			serviceName,
			string(serviceData.Status),
			serviceData.LatencyMs,
			serviceData.Message,
		)

		if err != nil {
			log.Error().
				Err(err).
				Str("service", serviceName).
				Msg("Failed to store health check")
		} else {
			log.Debug().
				Str("service", serviceName).
				Str("status", string(serviceData.Status)).
				Float64("latency_ms", serviceData.LatencyMs).
				Msg("Health check stored")
		}
	}
}
