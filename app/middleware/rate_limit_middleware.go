package middleware

import (
	"fmt"
	"strconv"
	"time"

	"golang_starter_kit_2025/app/helpers"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// GlobalRateLimiter limits requests per IP globally to prevent DDoS
// Default: 100 requests per minute per IP
// IMPROVED: Uses Redis for distributed rate limiting (horizontal scaling safe)
func GlobalRateLimiter() gin.HandlerFunc {
	limit := helpers.GetEnvInt64("RATE_LIMIT_GLOBAL", 100)
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  limit,
	}

	// Try Redis first (distributed), fallback to memory (single instance)
	store := createRateLimitStore()
	instance := limiter.New(store, rate)
	middleware := mgin.NewMiddleware(instance)

	return middleware
}

// AuthRateLimiter limits login/register attempts to prevent brute force
// Default: 5 attempts per 15 minutes per IP
// IMPROVED: Uses Redis for distributed rate limiting (horizontal scaling safe)
func AuthRateLimiter() gin.HandlerFunc {
	limit := helpers.GetEnvInt64("RATE_LIMIT_AUTH", 5)
	periodMinutes := helpers.GetEnvInt("RATE_LIMIT_AUTH_PERIOD_MINUTES", 15)

	rate := limiter.Rate{
		Period: time.Duration(periodMinutes) * time.Minute,
		Limit:  limit,
	}

	// Try Redis first (distributed), fallback to memory (single instance)
	store := createRateLimitStore()
	instance := limiter.New(store, rate)
	middleware := mgin.NewMiddleware(instance)

	return middleware
}

// UserRateLimiter limits requests per authenticated user to prevent API abuse
// Default: 1000 requests per hour per user
// IMPROVED: Uses Redis for distributed rate limiting (horizontal scaling safe)
func UserRateLimiter() gin.HandlerFunc {
	limit := helpers.GetEnvInt64("RATE_LIMIT_USER", 1000)
	periodMinutes := helpers.GetEnvInt("RATE_LIMIT_USER_PERIOD_MINUTES", 60)

	rate := limiter.Rate{
		Period: time.Duration(periodMinutes) * time.Minute,
		Limit:  limit,
	}

	// Try Redis first (distributed), fallback to memory (single instance)
	store := createRateLimitStore()

	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			// If no user ID, skip user-based rate limiting
			c.Next()
			return
		}

		// Create unique key for this user
		key := fmt.Sprintf("user:%v", userID)
		limiterInstance := limiter.New(store, rate)

		// Check rate limit
		res, err := limiterInstance.Get(c.Request.Context(), key)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   "rate limit error",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.FormatInt(res.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(res.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(res.Reset, 10))

		// Check if limit reached
		if res.Reached {
			retryAfter := time.Until(time.Unix(res.Reset, 0))
			c.Header("Retry-After", strconv.FormatInt(int64(retryAfter.Seconds()), 10))

			c.JSON(429, gin.H{
				"error":       "rate limit exceeded",
				"message":     fmt.Sprintf("Too many requests. Please try again in %v", retryAfter.Round(time.Second)),
				"retry_after": res.Reset,
				"limit":       res.Limit,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIRateLimiter is a general purpose rate limiter for API endpoints
// Can be configured per route group
// IMPROVED: Uses Redis for distributed rate limiting (horizontal scaling safe)
func APIRateLimiter(limit int64, period time.Duration) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}

	// Try Redis first (distributed), fallback to memory (single instance)
	store := createRateLimitStore()
	instance := limiter.New(store, rate)
	middleware := mgin.NewMiddleware(instance)

	return middleware
}

// createRateLimitStore creates Redis store if available, falls back to memory store
// Redis store is required for horizontal scaling (multi-instance deployments)
func createRateLimitStore() limiter.Store {
	redisClient := helpers.GetRedisClient()

	if redisClient == nil {
		log.Warn().Msg("Redis not available, using memory store for rate limiting (not distributed)")
		return memory.NewStore()
	}

	// Create Redis store for distributed rate limiting
	store, err := sredis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix:   "rate_limit",
		MaxRetry: 3,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to create Redis rate limit store, falling back to memory")
		return memory.NewStore()
	}

	log.Info().Msg("Rate limiting using Redis store (distributed)")
	return store
}
