package helpers

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var RedisClient *redis.Client

// InitRedis initializes Redis client with configuration from environment variables
func InitRedis() error {
	addr := GetEnv("REDIS_ADDR", "localhost:6379")
	password := GetEnv("REDIS_PASSWORD", "")
	db := GetEnvInt("REDIS_DB", 0)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Error().Err(err).Msg("Failed to connect to Redis")
		return err
	}

	log.Info().
		Str("addr", addr).
		Int("db", db).
		Msg("Redis connected successfully")

	return nil
}

// CloseRedis closes Redis connection
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
