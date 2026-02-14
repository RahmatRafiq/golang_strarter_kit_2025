package config

import (
	"os"
	"strconv"
)

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// GetRedisConfig returns Redis configuration from environment variables
// This is the centralized Redis config, similar to Laravel's config/database.php
func GetRedisConfig() *RedisConfig {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		db = 0 // Default Redis DB
	}

	return &RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
}

// GetTestRedisConfig returns Redis config for testing (DB 15)
func GetTestRedisConfig() *RedisConfig {
	cfg := GetRedisConfig()
	cfg.DB = 15 // Dedicated test database
	return cfg
}
