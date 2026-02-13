package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"golang_starter_kit_2025/app/helpers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// responseWriter is a custom writer to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// generateCacheKey creates a unique cache key based on URL and query params
func generateCacheKey(url string) string {
	hash := sha256.Sum256([]byte(url))
	return "cache:" + hex.EncodeToString(hash[:])
}

// CacheMiddleware provides HTTP response caching with Redis
// ttl: time-to-live for cached responses
func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Skip if Redis is not connected
		if helpers.RedisClient == nil {
			c.Next()
			return
		}

		// Generate cache key from full URL (includes query params)
		cacheKey := generateCacheKey(c.Request.URL.String())

		// Try to get from cache
		cached, err := helpers.RedisClient.Get(c.Request.Context(), cacheKey).Result()
		if err == nil {
			// Cache hit - return cached response
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(cached))
			c.Abort()
			return
		}

		// Cache miss or error
		if err != redis.Nil {
			// Log Redis error but continue
			log.Warn().Err(err).Str("cache_key", cacheKey).Msg("Redis cache error")
		}

		// Create custom writer to capture response
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = writer

		// Continue to handler
		c.Next()

		// IMPROVED: Validate response before caching
		if c.Writer.Status() == http.StatusOK {
			c.Header("X-Cache", "MISS")

			// Validate response is cacheable
			if !isCacheable(c, writer.body.Bytes()) {
				log.Debug().
					Str("path", c.Request.URL.Path).
					Msg("Response not cacheable, skipping cache")
				return
			}

			// Store in Redis with TTL
			err := helpers.RedisClient.Set(
				c.Request.Context(),
				cacheKey,
				writer.body.String(),
				ttl,
			).Err()

			if err != nil {
				log.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to cache response")
			}
		}
	}
}

// isCacheable validates if response should be cached
// Prevents caching error responses, non-JSON responses, etc.
func isCacheable(c *gin.Context, body []byte) bool {
	contentType := c.Writer.Header().Get("Content-Type")

	// Only cache JSON responses
	if !strings.Contains(contentType, "application/json") {
		return false
	}

	// Don't cache error responses (even if HTTP 200)
	// Check for common error patterns in JSON
	bodyStr := string(body)
	if strings.Contains(bodyStr, `"error":true`) ||
		strings.Contains(bodyStr, `"error":"`) ||
		strings.Contains(bodyStr, `"status":"error"`) {
		return false
	}

	// Don't cache empty responses
	if len(body) == 0 || bodyStr == "{}" || bodyStr == "[]" {
		return false
	}

	return true
}

// InvalidateCache removes a specific cache entry
// FIXED: Use context.Background() instead of TODO()
func InvalidateCache(url string) error {
	if helpers.RedisClient == nil {
		return nil
	}

	cacheKey := generateCacheKey(url)
	return helpers.RedisClient.Del(context.Background(), cacheKey).Err()
}

// InvalidateCachePattern removes cache entries matching a pattern
func InvalidateCachePattern(pattern string) error {
	if helpers.RedisClient == nil {
		return nil
	}

	ctx := context.Background()
	iter := helpers.RedisClient.Scan(ctx, 0, "cache:"+pattern+"*", 0).Iterator()

	for iter.Next(ctx) {
		if err := helpers.RedisClient.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}
