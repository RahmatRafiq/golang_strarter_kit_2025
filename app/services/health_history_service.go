package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang_starter_kit_2025/app/casts"
	"golang_starter_kit_2025/app/helpers"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type HealthHistoryService struct{}

func NewHealthHistoryService() *HealthHistoryService {
	return &HealthHistoryService{}
}

// StoreHealthCheck stores a single health check point in Redis
func (s *HealthHistoryService) StoreHealthCheck(service, status string, latencyMs float64, message string) error {
	if helpers.RedisClient == nil {
		return nil // Gracefully skip if Redis unavailable
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	point := casts.HealthHistoryPoint{
		Timestamp: time.Now().Unix(),
		Status:    status,
		LatencyMs: latencyMs,
		Message:   message,
	}

	jsonData, err := json.Marshal(point)
	if err != nil {
		log.Error().Err(err).Str("service", service).Msg("Failed to marshal health check point")
		return err
	}

	key := "health:history:" + service
	score := float64(point.Timestamp)

	// Store in sorted set
	if err := helpers.RedisClient.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: string(jsonData),
	}).Err(); err != nil {
		log.Error().Err(err).Str("service", service).Msg("Failed to store health check in Redis")
		return err
	}

	// Cleanup old data (keep last 90 days)
	retentionDays := helpers.GetEnvInt("HEALTH_HISTORY_RETENTION_DAYS", 90)
	if err := s.CleanupOldData(service, retentionDays); err != nil {
		log.Warn().Err(err).Str("service", service).Msg("Failed to cleanup old data")
		// Non-critical error, continue
	}

	return nil
}

// GetServiceHistory retrieves raw health check points for a service
func (s *HealthHistoryService) GetServiceHistory(service string, days int) ([]casts.HealthHistoryPoint, error) {
	if helpers.RedisClient == nil {
		return []casts.HealthHistoryPoint{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := "health:history:" + service
	endTime := time.Now().Unix()
	startTime := endTime - int64(days*24*60*60)

	// Query sorted set by timestamp range
	results, err := helpers.RedisClient.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", startTime),
		Max: fmt.Sprintf("%d", endTime),
	}).Result()

	if err != nil {
		log.Error().Err(err).Str("service", service).Msg("Failed to retrieve health history from Redis")
		return []casts.HealthHistoryPoint{}, err
	}

	points := make([]casts.HealthHistoryPoint, 0, len(results))
	for _, jsonData := range results {
		var point casts.HealthHistoryPoint
		if err := json.Unmarshal([]byte(jsonData), &point); err != nil {
			log.Warn().Err(err).Msg("Failed to unmarshal health check point")
			continue
		}
		points = append(points, point)
	}

	return points, nil
}

// GetDailyAggregated aggregates health checks into daily status
func (s *HealthHistoryService) GetDailyAggregated(service string, days int) casts.ServiceHistoryResponse {
	points, err := s.GetServiceHistory(service, days)
	if err != nil || len(points) == 0 {
		return casts.ServiceHistoryResponse{
			Service:       service,
			DailyHistory:  []casts.DailyStatus{},
			OverallUptime: 0,
		}
	}

	// Group by day
	dayMap := make(map[string][]casts.HealthHistoryPoint)
	for _, point := range points {
		date := time.Unix(point.Timestamp, 0).Format("2006-01-02")
		dayMap[date] = append(dayMap[date], point)
	}

	// Calculate daily status
	dailyHistory := make([]casts.DailyStatus, 0)
	totalUpChecks := 0
	totalChecks := 0

	// Get last N days in order
	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dayPoints := dayMap[date]

		if len(dayPoints) == 0 {
			// No data for this day
			dailyHistory = append(dailyHistory, casts.DailyStatus{
				Date:   date,
				Status: "unknown",
				Uptime: 0,
			})
			continue
		}

		// Count up vs down
		upCount := 0
		for _, p := range dayPoints {
			totalChecks++
			if p.Status == "up" {
				upCount++
				totalUpChecks++
			}
		}

		uptimePercent := float64(upCount) / float64(len(dayPoints)) * 100

		// Determine day status
		dayStatus := "down"
		if uptimePercent == 100 {
			dayStatus = "up"
		} else if uptimePercent >= 50 {
			dayStatus = "degraded"
		}

		dailyHistory = append(dailyHistory, casts.DailyStatus{
			Date:   date,
			Status: dayStatus,
			Uptime: uptimePercent,
		})
	}

	// Calculate overall uptime
	overallUptime := 0.0
	if totalChecks > 0 {
		overallUptime = float64(totalUpChecks) / float64(totalChecks) * 100
	}

	return casts.ServiceHistoryResponse{
		Service:       service,
		DailyHistory:  dailyHistory,
		OverallUptime: overallUptime,
	}
}

// CleanupOldData removes health checks older than retention days
func (s *HealthHistoryService) CleanupOldData(service string, retentionDays int) error {
	if helpers.RedisClient == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	key := "health:history:" + service
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays).Unix()

	// Remove old entries
	if err := helpers.RedisClient.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", cutoffTime)).Err(); err != nil {
		log.Error().Err(err).Str("service", service).Msg("Failed to cleanup old health history")
		return err
	}

	return nil
}
