package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ZerologMiddleware replaces Gin's default logger with structured zerolog
func ZerologMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Get or generate request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Calculate request duration
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Build base logger with common fields
		logger := log.With().
			Str("request_id", requestID).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("ip", clientIP).
			Str("user_agent", userAgent).
			Logger()

		// Add error details if present
		if len(c.Errors) > 0 {
			logger = logger.With().Str("error", c.Errors.String()).Logger()
		}

		// Log with appropriate level based on status code
		switch {
		case statusCode >= 500:
			logger.Error().Msg("Server error")
		case statusCode >= 400:
			logger.Warn().Msg("Client error")
		case statusCode >= 300:
			logger.Info().Msg("Redirect")
		default:
			logger.Info().Msg("Request completed")
		}
	}
}

// generateRequestID creates a unique request identifier
func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
