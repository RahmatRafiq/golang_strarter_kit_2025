package middleware

import (
	"net/http"

	"golang_starter_kit_2025/app/helpers"

	"github.com/gin-gonic/gin"
)

// RequestSizeLimit middleware limits the size of incoming request bodies
// to prevent DoS attacks via large payloads
// Default: 10MB (configurable via RATE_LIMIT_MAX_REQUEST_SIZE_MB env var)
func RequestSizeLimit() gin.HandlerFunc {
	// Get max size from env (in MB), default to 10MB
	maxSizeMB := helpers.GetEnvInt64("RATE_LIMIT_MAX_REQUEST_SIZE_MB", 10)
	maxBytes := maxSizeMB * 1024 * 1024 // Convert MB to bytes

	return func(c *gin.Context) {
		// Set max bytes reader to prevent reading large bodies
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

		// Continue to next handler
		c.Next()

		// Check if body was too large
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":       "request too large",
				"message":     "Request body exceeds maximum allowed size",
				"max_size_mb": maxSizeMB,
			})
			c.Abort()
		}
	}
}
