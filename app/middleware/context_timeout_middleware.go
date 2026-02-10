package middleware

import (
	"context"
	"time"

	"golang_starter_kit_2025/app/helpers"

	"github.com/gin-gonic/gin"
)

// ContextTimeoutMiddleware adds request timeout to context
// Prevents long-running requests from holding resources indefinitely
func ContextTimeoutMiddleware() gin.HandlerFunc {
	// Get timeout from env, default 30 seconds
	timeoutSeconds := helpers.GetEnvInt("REQUEST_TIMEOUT_SECONDS", 30)
	timeout := time.Duration(timeoutSeconds) * time.Second

	return func(c *gin.Context) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to track if handler finished
		finished := make(chan struct{})

		go func() {
			c.Next()
			close(finished)
		}()

		select {
		case <-finished:
			// Handler completed successfully
			return
		case <-ctx.Done():
			// Timeout occurred
			if ctx.Err() == context.DeadlineExceeded {
				c.JSON(504, gin.H{
					"error":   "Request timeout",
					"message": "The request took too long to process",
				})
				c.Abort()
			}
		}
	}
}
