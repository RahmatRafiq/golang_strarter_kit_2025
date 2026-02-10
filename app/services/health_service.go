package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang_starter_kit_2025/app/casts"
	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/facades"

	"github.com/rs/zerolog/log"
)

var startTime = time.Now()

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

// GetBasicHealth returns basic health status
func (s *HealthService) GetBasicHealth() *casts.HealthResponse {
	return &casts.HealthResponse{
		Status:    casts.StatusHealthy,
		Timestamp: time.Now(),
		Version:   getVersion(),
	}
}

// GetDetailedHealth returns comprehensive health status
func (s *HealthService) GetDetailedHealth() *casts.DetailedHealthResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	services := make(map[string]casts.Service)
	overallStatus := casts.StatusHealthy

	// Check database
	dbService := s.checkDatabase(ctx)
	services["database"] = dbService
	if dbService.Status != casts.ServiceUp {
		overallStatus = casts.StatusDegraded
		log.Warn().Str("service", "database").Msg("Database health check failed")
	}

	// Check Redis
	redisService := s.checkRedis(ctx)
	services["redis"] = redisService
	if redisService.Status != casts.ServiceUp {
		// Redis is optional, so only degrade if database is also down
		if overallStatus == casts.StatusDegraded {
			overallStatus = casts.StatusUnhealthy
		}
		log.Warn().Str("service", "redis").Msg("Redis health check failed")
	}

	// API is healthy if we can respond
	services["api"] = casts.Service{
		Status:    casts.ServiceUp,
		LatencyMs: 0,
		Message:   "API server operational",
	}

	return &casts.DetailedHealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Uptime:    formatUptime(time.Since(startTime)),
		Version:   getVersion(),
		Services:  services,
	}
}

// checkDatabase verifies database connectivity
func (s *HealthService) checkDatabase(ctx context.Context) casts.Service {
	start := time.Now()

	// Get default database connection
	db, err := facades.GetConnection("mysql")
	if err != nil {
		return casts.Service{
			Status:    casts.ServiceDown,
			LatencyMs: time.Since(start).Seconds() * 1000,
			Message:   fmt.Sprintf("Failed to get connection: %v", err),
		}
	}

	// Test connection with simple query
	sqlDB, err := db.DB.DB()
	if err != nil {
		return casts.Service{
			Status:    casts.ServiceDown,
			LatencyMs: time.Since(start).Seconds() * 1000,
			Message:   fmt.Sprintf("Failed to get SQL DB: %v", err),
		}
	}

	// Ping database with context timeout
	if err := sqlDB.PingContext(ctx); err != nil {
		return casts.Service{
			Status:    casts.ServiceDown,
			LatencyMs: time.Since(start).Seconds() * 1000,
			Message:   fmt.Sprintf("Database ping failed: %v", err),
		}
	}

	// Get connection pool stats
	stats := sqlDB.Stats()
	latency := time.Since(start).Seconds() * 1000

	return casts.Service{
		Status:    casts.ServiceUp,
		LatencyMs: latency,
		Message:   "Database connection healthy",
		Details: casts.DatabaseDetails{
			Active: stats.InUse,
			Idle:   stats.Idle,
			Max:    stats.MaxOpenConnections,
		},
	}
}

// checkRedis verifies Redis connectivity
func (s *HealthService) checkRedis(ctx context.Context) casts.Service {
	start := time.Now()

	// Check if Redis is initialized
	if helpers.RedisClient == nil {
		return casts.Service{
			Status:    casts.ServiceDown,
			LatencyMs: 0,
			Message:   "Redis client not initialized",
		}
	}

	// Ping Redis with context timeout
	_, err := helpers.RedisClient.Ping(ctx).Result()
	latency := time.Since(start).Seconds() * 1000

	if err != nil {
		return casts.Service{
			Status:    casts.ServiceDown,
			LatencyMs: latency,
			Message:   fmt.Sprintf("Redis ping failed: %v", err),
		}
	}

	return casts.Service{
		Status:    casts.ServiceUp,
		LatencyMs: latency,
		Message:   "Redis connection healthy",
	}
}

// formatUptime formats duration into human-readable uptime
func formatUptime(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// getVersion returns application version from environment or default
func getVersion() string {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		return "1.0.0"
	}
	return version
}
