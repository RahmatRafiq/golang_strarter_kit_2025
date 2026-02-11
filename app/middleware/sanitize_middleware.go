package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
)

var (
	// SQL injection patterns
	sqlInjectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(union\s+select|union\s+all\s+select)`),
		regexp.MustCompile(`(?i)(insert\s+into|update\s+.+\s+set|delete\s+from)`),
		regexp.MustCompile(`(?i)(drop\s+table|drop\s+database|create\s+table)`),
		regexp.MustCompile(`(?i)(--|#|\/\*|\*\/)`),
		regexp.MustCompile(`(?i)(xp_|sp_|exec\s+|execute\s+)`),
		regexp.MustCompile(`(?i)(or\s+['"]?1['"]?\s*=\s*['"]?1|and\s+['"]?1['"]?\s*=\s*['"]?1)`),
	}

	// XSS patterns
	xssPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(<script|</script>|javascript:|onerror=|onload=)`),
		regexp.MustCompile(`(?i)(<iframe|<object|<embed|<applet)`),
	}
)

// SanitizeInput sanitizes all JSON inputs to prevent XSS
func SanitizeInput() gin.HandlerFunc {
	policy := bluemonday.UGCPolicy() // User Generated Content policy (more permissive than Strict)

	return func(c *gin.Context) {
		// Only sanitize if content type is JSON
		if !strings.Contains(c.GetHeader("Content-Type"), "application/json") {
			c.Next()
			return
		}

		// Read body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}

		// Restore body for next handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Try to parse as JSON
		var body interface{}
		if jsonErr := json.Unmarshal(bodyBytes, &body); jsonErr != nil {
			c.Next()
			return
		}

		// Sanitize the data
		sanitized := sanitizeValue(body, policy)

		// Store sanitized body in context
		c.Set("sanitized_body", sanitized)

		// Update request body with sanitized data
		sanitizedBytes, err := json.Marshal(sanitized)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal sanitized data")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitizedBytes))
		c.Request.ContentLength = int64(len(sanitizedBytes))

		c.Next()
	}
}

// SQLInjectionCheck checks for SQL injection patterns in query params and path
func SQLInjectionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check query parameters
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if containsSQLPattern(value) {
					c.JSON(400, gin.H{
						"error":   "Invalid input detected",
						"message": "Suspicious pattern found in request",
					})
					c.Abort()
					return
				}
			}
		}

		// Check URL path parameters
		for _, param := range c.Params {
			if containsSQLPattern(param.Value) {
				c.JSON(400, gin.H{
					"error":   "Invalid input detected",
					"message": "Suspicious pattern found in request",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// XSSProtection checks for XSS patterns
func XSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check query parameters
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if containsXSSPattern(value) {
					c.JSON(400, gin.H{
						"error":   "Invalid input detected",
						"message": "Potentially malicious content detected",
					})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

// sanitizeValue recursively sanitizes any value (string, map, slice)
func sanitizeValue(value interface{}, policy *bluemonday.Policy) interface{} {
	switch v := value.(type) {
	case string:
		return policy.Sanitize(v)
	case map[string]interface{}:
		return sanitizeMap(v, policy)
	case []interface{}:
		return sanitizeSlice(v, policy)
	default:
		return value
	}
}

// sanitizeMap sanitizes all string values in a map
func sanitizeMap(data map[string]interface{}, policy *bluemonday.Policy) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		result[key] = sanitizeValue(value, policy)
	}

	return result
}

// sanitizeSlice sanitizes all string values in a slice
func sanitizeSlice(data []interface{}, policy *bluemonday.Policy) []interface{} {
	result := make([]interface{}, len(data))

	for i, value := range data {
		result[i] = sanitizeValue(value, policy)
	}

	return result
}

// containsSQLPattern checks if string contains SQL injection patterns
func containsSQLPattern(s string) bool {
	for _, pattern := range sqlInjectionPatterns {
		if pattern.MatchString(s) {
			return true
		}
	}
	return false
}

// containsXSSPattern checks if string contains XSS patterns
func containsXSSPattern(s string) bool {
	for _, pattern := range xssPatterns {
		if pattern.MatchString(s) {
			return true
		}
	}
	return false
}
