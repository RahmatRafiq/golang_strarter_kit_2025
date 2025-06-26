package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Client *redis.Client
	Ctx    context.Context
}

var Redis *RedisConfig

func InitRedis() {
	addr := os.Getenv("REDIS_HOST")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")

	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		if dbNum, err := strconv.Atoi(dbStr); err == nil {
			db = dbNum
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}

	Redis = &RedisConfig{
		Client: rdb,
		Ctx:    ctx,
	}

	fmt.Println("✅ Redis connected successfully")
}

func CloseRedis() {
	if Redis != nil && Redis.Client != nil {
		Redis.Client.Close()
	}
}

// Helper methods
func (r *RedisConfig) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(r.Ctx, key, value, expiration).Err()
}

func (r *RedisConfig) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *RedisConfig) Del(key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}

func (r *RedisConfig) LPush(key string, values ...interface{}) error {
	return r.Client.LPush(r.Ctx, key, values...).Err()
}

func (r *RedisConfig) RPop(key string) (string, error) {
	return r.Client.RPop(r.Ctx, key).Result()
}

func (r *RedisConfig) BRPop(timeout time.Duration, keys ...string) ([]string, error) {
	return r.Client.BRPop(r.Ctx, timeout, keys...).Result()
}

func (r *RedisConfig) LLen(key string) (int64, error) {
	return r.Client.LLen(r.Ctx, key).Result()
}
