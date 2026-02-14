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
	Details   interface{}   `json:"details,omitempty"`
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
