package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler will set `X-Response-Time` header in response.
func Response(c *gin.Context) {
	startTime := time.Now()
	

	c.Next()
	duration := int(time.Since(startTime).Milliseconds())
	c.Writer.Header().Set("X-Response-Time", strconv.Itoa(duration))
}
