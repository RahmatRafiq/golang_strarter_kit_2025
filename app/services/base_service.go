package services

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

// BaseService provides common service layer functionality
// MED-8: Extract base service patterns for code reuse
type BaseService struct {
	cache     *CacheService
	tableName string
}

// NewBaseService creates a new base service
func NewBaseService(tableName string) *BaseService {
	return &BaseService{
		cache:     NewCacheService(),
		tableName: tableName,
	}
}

// GetCache returns the cache service
func (s *BaseService) GetCache() *CacheService {
	return s.cache
}

// IsCacheEnabled checks if caching is available
func (s *BaseService) IsCacheEnabled() bool {
	return s.cache != nil && s.cache.IsEnabled()
}

// CacheGet retrieves value from cache with logging
func (s *BaseService) CacheGet(key string, dest interface{}) error {
	if !s.IsCacheEnabled() {
		return fmt.Errorf("cache not enabled")
	}

	err := s.cache.Get(key, dest)
	if err == nil {
		log.Debug().
			Str("table", s.tableName).
			Str("key", key).
			Msg("Cache hit")
	}
	return err
}

// CacheSet stores value in cache with logging
func (s *BaseService) CacheSet(key string, value interface{}, ttl time.Duration) error {
	if !s.IsCacheEnabled() {
		return nil // Don't fail if cache unavailable
	}

	if err := s.cache.Set(key, value, ttl); err != nil {
		log.Warn().
			Err(err).
			Str("table", s.tableName).
			Str("key", key).
			Msg("Failed to cache value")
		return err
	}

	log.Debug().
		Str("table", s.tableName).
		Str("key", key).
		Dur("ttl", ttl).
		Msg("Value cached")
	return nil
}

// CacheInvalidate removes cache entry with logging
func (s *BaseService) CacheInvalidate(key string) error {
	if !s.IsCacheEnabled() {
		return nil
	}

	if err := s.cache.Delete(key); err != nil {
		log.Warn().
			Err(err).
			Str("table", s.tableName).
			Str("key", key).
			Msg("Failed to invalidate cache")
		return err
	}

	log.Debug().
		Str("table", s.tableName).
		Str("key", key).
		Msg("Cache invalidated")
	return nil
}

// CacheInvalidatePattern removes cache entries matching pattern
func (s *BaseService) CacheInvalidatePattern(pattern string) error {
	if !s.IsCacheEnabled() {
		return nil
	}

	if err := s.cache.DeletePattern(pattern); err != nil {
		log.Warn().
			Err(err).
			Str("table", s.tableName).
			Str("pattern", pattern).
			Msg("Failed to invalidate cache pattern")
		return err
	}

	log.Debug().
		Str("table", s.tableName).
		Str("pattern", pattern).
		Msg("Cache pattern invalidated")
	return nil
}

// GetTableName returns the table name
func (s *BaseService) GetTableName() string {
	return s.tableName
}

// CacheKeyForID generates cache key for single record
func (s *BaseService) CacheKeyForID(id uint) string {
	return fmt.Sprintf("%s:%d", s.tableName, id)
}

// CacheKeyForList generates cache key for list query
func (s *BaseService) CacheKeyForList(page, limit int) string {
	return ListCacheKey(s.tableName, page, limit)
}

// CacheKeyForCount generates cache key for count query
func (s *BaseService) CacheKeyForCount(conditions ...string) string {
	return CountCacheKey(s.tableName, conditions...)
}
