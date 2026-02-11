package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// HTTPSRedirect redirects HTTP requests to HTTPS in production environment
// This middleware enforces secure connections to protect sensitive data
func HTTPSRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip HTTPS enforcement in development and testing environments
		appEnv := os.Getenv("APP_ENV")
		if appEnv == "development" || appEnv == "testing" {
			c.Next()
			return
		}

		// Check if request came via HTTPS
		// X-Forwarded-Proto is set by load balancers (nginx, HAProxy, AWS ALB, etc)
		proto := c.Request.Header.Get("X-Forwarded-Proto")
		if proto == "" {
			// Fallback to request URL scheme if no load balancer
			proto = c.Request.URL.Scheme
		}

		// Also check X-Forwarded-Ssl header (some load balancers use this)
		if proto == "" && c.Request.Header.Get("X-Forwarded-Ssl") == "on" {
			proto = "https"
		}

		// Redirect to HTTPS if request is not secure
		if proto != "https" && proto != "" {
			httpsURL := "https://" + c.Request.Host + c.Request.RequestURI

			log.Warn().
				Str("original_url", c.Request.URL.String()).
				Str("redirect_url", httpsURL).
				Str("client_ip", c.ClientIP()).
				Msg("Redirecting HTTP to HTTPS")

			c.Redirect(301, httpsURL) // 301 Permanent Redirect
			c.Abort()
			return
		}

		c.Next()
	}
}
