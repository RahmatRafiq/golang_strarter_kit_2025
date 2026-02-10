package casts

// HealthHistoryPoint represents a single health check data point
type HealthHistoryPoint struct {
	Timestamp int64   `json:"timestamp"`
	Status    string  `json:"status"` // "up" or "down"
	LatencyMs float64 `json:"latency_ms"`
	Message   string  `json:"message,omitempty"`
}

// DailyStatus represents aggregated health status for a single day
type DailyStatus struct {
	Date   string  `json:"date"`   // "2024-01-15"
	Status string  `json:"status"` // "up", "degraded", "down"
	Uptime float64 `json:"uptime"` // Percentage for the day (0-100)
}

// ServiceHistoryResponse contains historical health data for a service
type ServiceHistoryResponse struct {
	Service       string        `json:"service"`
	DailyHistory  []DailyStatus `json:"daily_history"`
	OverallUptime float64       `json:"overall_uptime"` // Overall percentage (0-100)
}
