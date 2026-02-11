package config

import (
	"os"
	"strconv"

	"github.com/hibiken/asynq"
)

// QueueConfig holds the configuration for job queue system
type QueueConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Concurrency   int
	Queues        map[string]int
}

// GetQueueConfig returns queue configuration from environment variables
func GetQueueConfig() *QueueConfig {
	redisDB, _ := strconv.Atoi(os.Getenv("QUEUE_REDIS_DB"))
	if redisDB == 0 {
		redisDB = 1 // Default to DB 1 for queue (separate from cache)
	}

	concurrency, _ := strconv.Atoi(os.Getenv("QUEUE_CONCURRENCY"))
	if concurrency == 0 {
		concurrency = 10 // Default 10 concurrent workers
	}

	return &QueueConfig{
		RedisAddr:     os.Getenv("QUEUE_REDIS_ADDR"),
		RedisPassword: os.Getenv("QUEUE_REDIS_PASSWORD"),
		RedisDB:       redisDB,
		Concurrency:   concurrency,
		Queues: map[string]int{
			"critical": 6, // 60% of workers
			"default":  3, // 30% of workers
			"low":      1, // 10% of workers
		},
	}
}

// GetRedisClientOpt returns asynq.RedisClientOpt from config
func (c *QueueConfig) GetRedisClientOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	}
}
