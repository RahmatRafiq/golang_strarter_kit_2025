package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang_starter_kit_2025/app/helpers"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// CacheService provides caching functionality with Redis backend
type CacheService struct {
	client *redis.Client
	ctx    context.Context
}

// Default TTL values for different cache types
const (
	CacheTTLShort  = 5 * time.Minute  // For frequently changing data
	CacheTTLMedium = 30 * time.Minute // For moderate data
	CacheTTLLong   = 2 * time.Hour    // For stable data
	CacheTTLDay    = 24 * time.Hour   // For very stable data
)

// NewCacheService creates a new cache service instance
func NewCacheService() *CacheService {
	client := helpers.GetRedisClient()
	if client == nil {
		log.Warn().Msg("Redis client not available, cache service disabled")
		return &CacheService{
			client: nil,
			ctx:    context.Background(),
		}
	}

	return &CacheService{
		client: client,
		ctx:    context.Background(),
	}
}

// IsEnabled checks if cache service is available
func (s *CacheService) IsEnabled() bool {
	return s.client != nil
}

// Get retrieves a value from cache
func (s *CacheService) Get(key string, dest interface{}) error {
	if !s.IsEnabled() {
		return fmt.Errorf("cache not enabled")
	}

	val, err := s.client.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss")
	}
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to get from cache")
		return err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to unmarshal cache value")
		return err
	}

	log.Debug().Str("key", key).Msg("Cache hit")
	return nil
}

// Set stores a value in cache
func (s *CacheService) Set(key string, value interface{}, ttl time.Duration) error {
	if !s.IsEnabled() {
		return fmt.Errorf("cache not enabled")
	}

	data, err := json.Marshal(value)
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to marshal cache value")
		return err
	}

	if err := s.client.Set(s.ctx, key, data, ttl).Err(); err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to set cache")
		return err
	}

	log.Debug().Str("key", key).Dur("ttl", ttl).Msg("Cache set")
	return nil
}

// Delete removes a value from cache
func (s *CacheService) Delete(key string) error {
	if !s.IsEnabled() {
		return nil
	}

	if err := s.client.Del(s.ctx, key).Err(); err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to delete from cache")
		return err
	}

	log.Debug().Str("key", key).Msg("Cache deleted")
	return nil
}

// DeletePattern deletes all keys matching a pattern
func (s *CacheService) DeletePattern(pattern string) error {
	if !s.IsEnabled() {
		return nil
	}

	iter := s.client.Scan(s.ctx, 0, pattern, 0).Iterator()
	var keys []string

	for iter.Next(s.ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		log.Error().Err(err).Str("pattern", pattern).Msg("Failed to scan cache keys")
		return err
	}

	if len(keys) > 0 {
		if err := s.client.Del(s.ctx, keys...).Err(); err != nil {
			log.Error().Err(err).Str("pattern", pattern).Msg("Failed to delete cache keys")
			return err
		}
		log.Debug().Str("pattern", pattern).Int("count", len(keys)).Msg("Cache keys deleted")
	}

	return nil
}

// Remember caches the result of a callback function
func (s *CacheService) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error) {
	// Try to get from cache first
	var result interface{}
	if err := s.Get(key, &result); err == nil {
		return result, nil
	}

	// Cache miss or error, execute callback
	result, err := callback()
	if err != nil {
		return nil, err
	}

	// Store in cache (log cache errors but don't fail the operation)
	if setErr := s.Set(key, result, ttl); setErr != nil {
		log.Warn().Err(setErr).Str("key", key).Msg("Failed to cache result in Remember")
	}

	return result, nil
}

// RememberQuery is a typed version of Remember for query result caching
// MED-3: Query result caching layer for database queries
// Usage: cache.RememberQuery("users:list:page:1", CacheTTLMedium, func() ([]User, error) { ... })
func RememberQuery[T any](s *CacheService, key string, ttl time.Duration, query func() (T, error)) (T, error) {
	var result T

	// Try to get from cache first
	if err := s.Get(key, &result); err == nil {
		log.Debug().Str("key", key).Msg("Query result from cache")
		return result, nil
	}

	// Cache miss, execute query
	result, err := query()
	if err != nil {
		return result, err
	}

	// Store in cache (log errors but don't fail query)
	if setErr := s.Set(key, result, ttl); setErr != nil {
		log.Warn().Err(setErr).Str("key", key).Msg("Failed to cache query result")
	}

	log.Debug().Str("key", key).Dur("ttl", ttl).Msg("Query result cached")
	return result, nil
}

// Exists checks if a key exists in cache
func (s *CacheService) Exists(key string) bool {
	if !s.IsEnabled() {
		return false
	}

	exists, err := s.client.Exists(s.ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Failed to check cache existence")
		return false
	}

	return exists > 0
}

// Flush clears all cache
func (s *CacheService) Flush() error {
	if !s.IsEnabled() {
		return nil
	}

	if err := s.client.FlushAll(s.ctx).Err(); err != nil {
		log.Error().Err(err).Msg("Failed to flush cache")
		return err
	}

	log.Info().Msg("Cache flushed")
	return nil
}

// Helper functions for common cache key patterns

// CountCacheKey generates cache key for COUNT queries
// MED-4: Optimize COUNT queries with caching
func CountCacheKey(tableName string, conditions ...string) string {
	if len(conditions) == 0 {
		return fmt.Sprintf("count:%s:all", tableName)
	}
	return fmt.Sprintf("count:%s:%v", tableName, conditions)
}

// ListCacheKey generates cache key for list queries
func ListCacheKey(tableName string, page, limit int, filters ...string) string {
	if len(filters) == 0 {
		return fmt.Sprintf("list:%s:page:%d:limit:%d", tableName, page, limit)
	}
	return fmt.Sprintf("list:%s:page:%d:limit:%d:filter:%v", tableName, page, limit, filters)
}

// UserCacheKey generates cache key for user data
func UserCacheKey(userID uint) string {
	return fmt.Sprintf("user:%d", userID)
}

// UserRolesCacheKey generates cache key for user roles
func UserRolesCacheKey(userID uint) string {
	return fmt.Sprintf("user:%d:roles", userID)
}

// RoleCacheKey generates cache key for role data
func RoleCacheKey(roleID uint) string {
	return fmt.Sprintf("role:%d", roleID)
}

// RolePermissionsCacheKey generates cache key for role permissions
func RolePermissionsCacheKey(roleID uint) string {
	return fmt.Sprintf("role:%d:permissions", roleID)
}

// InvalidateUserCache invalidates all user-related cache
func (s *CacheService) InvalidateUserCache(userID uint) error {
	patterns := []string{
		UserCacheKey(userID),
		UserRolesCacheKey(userID),
	}

	for _, pattern := range patterns {
		if err := s.Delete(pattern); err != nil {
			log.Error().Err(err).Str("pattern", pattern).Msg("Failed to invalidate user cache")
		}
	}

	return nil
}

// InvalidateRoleCache invalidates all role-related cache
func (s *CacheService) InvalidateRoleCache(roleID uint) error {
	patterns := []string{
		RoleCacheKey(roleID),
		RolePermissionsCacheKey(roleID),
		"user:*:roles", // Invalidate all user roles cache
	}

	for _, pattern := range patterns {
		if err := s.DeletePattern(pattern); err != nil {
			log.Error().Err(err).Str("pattern", pattern).Msg("Failed to invalidate role cache")
		}
	}

	return nil
}
