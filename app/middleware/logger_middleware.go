package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Proses request
		c.Next()

		// Hitung waktu respon dan log hasil
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		log.Printf("Status: %d | Latency: %v | Path: %s | Method: %s",
			statusCode, latency, c.Request.URL.Path, c.Request.Method)
	}
}
