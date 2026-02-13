package database

import (
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// QueryMonitorPlugin monitors query performance and logs slow queries
// MED-6: Connection pool health monitoring
// MED-7: Slow query logging
type QueryMonitorPlugin struct {
	SlowQueryThreshold time.Duration
	EnableStackTrace   bool
}

// Name returns the plugin name
func (p *QueryMonitorPlugin) Name() string {
	return "query_monitor_plugin"
}

// Initialize initializes the plugin
func (p *QueryMonitorPlugin) Initialize(db *gorm.DB) error {
	// Register callback before query to capture start time
	db.Callback().Query().Before("gorm:query").Register("query_monitor:before", p.beforeQuery)
	db.Callback().Create().Before("gorm:create").Register("query_monitor:before", p.beforeQuery)
	db.Callback().Update().Before("gorm:update").Register("query_monitor:before", p.beforeQuery)
	db.Callback().Delete().Before("gorm:delete").Register("query_monitor:before", p.beforeQuery)
	db.Callback().Row().Before("gorm:row").Register("query_monitor:before", p.beforeQuery)
	db.Callback().Raw().Before("gorm:raw").Register("query_monitor:before", p.beforeQuery)

	// Register callback after query execution
	db.Callback().Query().After("gorm:query").Register("query_monitor:after", p.afterQuery)
	db.Callback().Create().After("gorm:create").Register("query_monitor:after", p.afterQuery)
	db.Callback().Update().After("gorm:update").Register("query_monitor:after", p.afterQuery)
	db.Callback().Delete().After("gorm:delete").Register("query_monitor:after", p.afterQuery)
	db.Callback().Row().After("gorm:row").Register("query_monitor:after", p.afterQuery)
	db.Callback().Raw().After("gorm:raw").Register("query_monitor:after", p.afterQuery)

	log.Info().
		Dur("slow_threshold", p.SlowQueryThreshold).
		Bool("stack_trace", p.EnableStackTrace).
		Msg("Query monitor plugin initialized")

	return nil
}

// beforeQuery captures query start time
func (p *QueryMonitorPlugin) beforeQuery(db *gorm.DB) {
	db.InstanceSet("query_start_time", time.Now())
}

// afterQuery monitors query execution and logs slow queries
func (p *QueryMonitorPlugin) afterQuery(db *gorm.DB) {
	// Get query execution time
	startTime, exists := db.InstanceGet("query_start_time")
	if !exists {
		return
	}
	elapsed := time.Since(startTime.(time.Time))

	// Log slow queries
	if elapsed >= p.SlowQueryThreshold {
		logger := log.Warn().
			Dur("duration", elapsed).
			Str("sql", db.Statement.SQL.String()).
			Interface("vars", db.Statement.Vars)

		if db.Error != nil {
			logger = logger.Err(db.Error)
		}

		if db.Statement.RowsAffected >= 0 {
			logger = logger.Int64("rows_affected", db.Statement.RowsAffected)
		}

		if p.EnableStackTrace {
			// Add stack trace for debugging
			logger = logger.Stack()
		}

		logger.Msg("Slow query detected")
	}

	// Log all queries in debug mode
	log.Debug().
		Dur("duration", elapsed).
		Str("sql", db.Statement.SQL.String()).
		Int64("rows_affected", db.Statement.RowsAffected).
		Msg("Query executed")
}

// NewQueryMonitorPlugin creates a new query monitor plugin
// Default slow threshold: 1 second
func NewQueryMonitorPlugin(slowThreshold time.Duration, enableStackTrace bool) *QueryMonitorPlugin {
	if slowThreshold == 0 {
		slowThreshold = 1 * time.Second
	}

	return &QueryMonitorPlugin{
		SlowQueryThreshold: slowThreshold,
		EnableStackTrace:   enableStackTrace,
	}
}

// PoolStatsMonitor monitors connection pool statistics
// MED-6: Connection pool health monitoring
type PoolStatsMonitor struct {
	CheckInterval time.Duration
	WarnThreshold float64 // Warn if pool usage > threshold (e.g., 0.8 = 80%)
}

// Start begins monitoring connection pool statistics
func (m *PoolStatsMonitor) Start(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get SQL DB for pool monitoring")
		return
	}

	ticker := time.NewTicker(m.CheckInterval)
	go func() {
		for range ticker.C {
			stats := sqlDB.Stats()

			// Calculate pool usage percentage
			usage := float64(stats.InUse) / float64(stats.MaxOpenConnections)

			// Warn if pool usage is high
			if usage >= m.WarnThreshold {
				log.Warn().
					Int("open_connections", stats.OpenConnections).
					Int("in_use", stats.InUse).
					Int("idle", stats.Idle).
					Int("max_open", stats.MaxOpenConnections).
					Int64("max_idle_closed", stats.MaxIdleClosed).
					Int64("wait_count", stats.WaitCount).
					Dur("wait_duration", stats.WaitDuration).
					Float64("usage_percent", usage*100).
					Msg("High database connection pool usage")
			} else {
				log.Debug().
					Int("open_connections", stats.OpenConnections).
					Int("in_use", stats.InUse).
					Int("idle", stats.Idle).
					Int("max_open", stats.MaxOpenConnections).
					Int64("max_idle_closed", stats.MaxIdleClosed).
					Int64("wait_count", stats.WaitCount).
					Dur("wait_duration", stats.WaitDuration).
					Float64("usage_percent", usage*100).
					Msg("Connection pool stats")
			}
		}
	}()

	log.Info().
		Dur("interval", m.CheckInterval).
		Float64("warn_threshold", m.WarnThreshold*100).
		Msg("Connection pool monitoring started")
}

// NewPoolStatsMonitor creates a new pool stats monitor
// Default: check every 30 seconds, warn at 80% usage
func NewPoolStatsMonitor(interval time.Duration, warnThreshold float64) *PoolStatsMonitor {
	if interval == 0 {
		interval = 30 * time.Second
	}
	if warnThreshold == 0 {
		warnThreshold = 0.8
	}

	return &PoolStatsMonitor{
		CheckInterval: interval,
		WarnThreshold: warnThreshold,
	}
}
