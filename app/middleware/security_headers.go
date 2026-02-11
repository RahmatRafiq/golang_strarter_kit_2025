package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security-related HTTP headers to protect against common web vulnerabilities
// These headers provide defense-in-depth security for production environments
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-Content-Type-Options: Prevents MIME type sniffing
		// Browsers won't try to guess content type, reducing XSS risk
		c.Header("X-Content-Type-Options", "nosniff")

		// X-Frame-Options: Prevents clickjacking attacks
		// DENY = page cannot be displayed in iframe/frame/embed/object
		c.Header("X-Frame-Options", "DENY")

		// X-XSS-Protection: Enables browser's XSS filter (legacy but still useful)
		// 1; mode=block = detect and block XSS attacks
		c.Header("X-XSS-Protection", "1; mode=block")

		// Strict-Transport-Security (HSTS): Enforces HTTPS for 1 year
		// Only set if request is already HTTPS (prevents warnings)
		if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
			// max-age=31536000 = 1 year in seconds
			// includeSubDomains = apply to all subdomains
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// Content-Security-Policy: Controls what resources can be loaded
		// This is a restrictive policy - adjust based on your frontend needs
		appEnv := os.Getenv("APP_ENV")
		if appEnv == "production" {
			// Production: Strict CSP
			c.Header("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
					"style-src 'self' 'unsafe-inline'; "+
					"img-src 'self' data: https:; "+
					"font-src 'self' data:; "+
					"connect-src 'self'; "+
					"frame-ancestors 'none'")
		} else {
			// Development: More permissive for debugging
			c.Header("Content-Security-Policy", "default-src 'self' 'unsafe-inline' 'unsafe-eval'")
		}

		// Referrer-Policy: Controls how much referrer information is sent
		// strict-origin-when-cross-origin = send full URL for same-origin, only origin for cross-origin
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy: Controls browser features and APIs
		// Deny access to sensitive features like geolocation, camera, microphone
		c.Header("Permissions-Policy",
			"geolocation=(), "+
				"microphone=(), "+
				"camera=(), "+
				"payment=(), "+
				"usb=(), "+
				"magnetometer=(), "+
				"gyroscope=(), "+
				"accelerometer=()")

		// X-Permitted-Cross-Domain-Policies: Restrict Adobe Flash/PDF cross-domain requests
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		c.Next()
	}
}
