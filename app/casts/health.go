package casts

import "time"

// HealthStatus represents the overall health status
type HealthStatus string

const (
	StatusHealthy   HealthStatus = "healthy"
	StatusDegraded  HealthStatus = "degraded"
	StatusUnhealthy HealthStatus = "unhealthy"
)

// ServiceStatus represents individual service health
type ServiceStatus string

const (
	ServiceUp       ServiceStatus = "up"
	ServiceDown     ServiceStatus = "down"
	ServiceDegraded ServiceStatus = "degraded"
)

// HealthResponse is the basic health check response
type HealthResponse struct {
	Status    HealthStatus `json:"status" example:"healthy"`
	Timestamp time.Time    `json:"timestamp" example:"2025-02-11T10:30:00Z"`
	Version   string       `json:"version" example:"1.0.0"`
}

// DetailedHealthResponse provides comprehensive health information
type DetailedHealthResponse struct {
	Services  map[string]Service `json:"services"`
	Timestamp time.Time          `json:"timestamp" example:"2025-02-11T10:30:00Z"`
	Status    HealthStatus       `json:"status" example:"healthy"`
	Uptime    string             `json:"uptime" example:"5d 3h 12m"`
	Version   string             `json:"version" example:"1.0.0"`
}

// Service represents health information for a specific service
type Service struct {
	Details   any           `json:"details,omitempty"`
	Message   string        `json:"message,omitempty" example:""`
	Status    ServiceStatus `json:"status" example:"up"`
	LatencyMs float64       `json:"latency_ms,omitempty" example:"2.5"`
}

// DatabaseDetails provides database-specific health information
type DatabaseDetails struct {
	Active int `json:"active" example:"5"`
	Idle   int `json:"idle" example:"15"`
	Max    int `json:"max" example:"20"`
}

// RedisDetails provides Redis-specific health information
type RedisDetails struct {
	Memory  string  `json:"memory,omitempty" example:"125MB"`
	HitRate float64 `json:"hit_rate,omitempty" example:"0.85"`
}

// QueueHealthResponse provides queue worker health and metrics
type QueueHealthResponse struct {
	QueueDepth       map[string]int64 `json:"queue_depth,omitempty"`
	LastUpdated      time.Time        `json:"last_updated" example:"2026-02-14T15:30:00Z"`
	Status           string           `json:"status" example:"healthy"`
	Message          string           `json:"message,omitempty" example:"Queue is running (metrics disabled)"`
	TasksProcessed   int64            `json:"tasks_processed" example:"1234"`
	TasksFailed      int64            `json:"tasks_failed" example:"12"`
	TasksRetried     int64            `json:"tasks_retried" example:"45"`
	AverageLatencyMs float64          `json:"average_latency_ms" example:"125.5"`
	FailureRate      float64          `json:"failure_rate" example:"0.97"`
	ActiveWorkers    int              `json:"active_workers" example:"10"`
}
