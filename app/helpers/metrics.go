package helpers

import (
	"time"

	"golang_starter_kit_2025/facades"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

var (
	// Database connection pool metrics
	dbConnectionPoolStats = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_connection_pool_stats",
			Help: "Database connection pool statistics",
		},
		[]string{"connection", "metric"},
	)
)

// CollectDBMetrics starts a background goroutine to collect database metrics
func CollectDBMetrics() {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			collectCurrentDBStats()
		}
	}()

	log.Info().Msg("Database metrics collector started")
}

// collectCurrentDBStats collects current database connection pool statistics
func collectCurrentDBStats() {
	// Get current active connection from facades
	sqlDB, err := facades.DB.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get database connection for metrics")
		return
	}

	stats := sqlDB.Stats()
	connectionName := GetEnv("DB_CONNECTION", "default")

	// Record metrics
	dbConnectionPoolStats.WithLabelValues(connectionName, "max_open_connections").Set(float64(stats.MaxOpenConnections))
	dbConnectionPoolStats.WithLabelValues(connectionName, "open_connections").Set(float64(stats.OpenConnections))
	dbConnectionPoolStats.WithLabelValues(connectionName, "in_use").Set(float64(stats.InUse))
	dbConnectionPoolStats.WithLabelValues(connectionName, "idle").Set(float64(stats.Idle))
	dbConnectionPoolStats.WithLabelValues(connectionName, "wait_count").Set(float64(stats.WaitCount))
	dbConnectionPoolStats.WithLabelValues(connectionName, "wait_duration_ms").Set(float64(stats.WaitDuration.Milliseconds()))
	dbConnectionPoolStats.WithLabelValues(connectionName, "max_idle_closed").Set(float64(stats.MaxIdleClosed))
	dbConnectionPoolStats.WithLabelValues(connectionName, "max_lifetime_closed").Set(float64(stats.MaxLifetimeClosed))
}
