package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang_starter_kit_2025/app/casts"
	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

var JwtKey = []byte("your_secret_key")

var jwtService services.JwtService

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		skipAuth := os.Getenv("SKIP_AUTH")
		appEnv := os.Getenv("APP_ENV")

		// SECURITY: Prevent SKIP_AUTH in production
		if skipAuth == "true" && appEnv == "production" {
			log.Fatal().Msg("SECURITY ERROR: SKIP_AUTH cannot be 'true' in production environment")
		}

		// Development bypass (with security warning)
		if skipAuth == "true" {
			log.Warn().
				Str("ip", c.ClientIP()).
				Str("path", c.Request.URL.Path).
				Msg("AUTH BYPASS: Using SKIP_AUTH=true (development only)")

			// Use user_id from header if provided (for testing different users)
			userIDStr := c.GetHeader("X-Test-User-ID")
			if userIDStr == "" {
				userIDStr = "1"
			}

			// Convert string to uint
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				userID = 1 // Default to 1 if parse fails
			}

			c.Set("user_id", uint(userID))
			c.Set("auth_bypassed", true)
			c.Next()
			return
		}

		tokenString, shouldReturn := CheckTokenExist(c)
		if shouldReturn {
			return
		}

		shouldReturn1 := CheckBearerTokenPrefix(tokenString, c)
		if shouldReturn1 {
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, shouldReturn2 := CheckTokenValidity(tokenString, c)
		if shouldReturn2 {
			return
		}

		claims := casts.ParseJwtClaims(jwtService.ExtractClaims(token))

		if claims.ExpiredAt < time.Now().Unix() {
			helpers.ResponseError(c, &helpers.ResponseParams[any]{
				Reference: "ERROR-4",
				Message:   "Token sudah kadaluarsa",
			}, http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("token", tokenString)
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func CheckTokenValidity(tokenString string, c *gin.Context) (*jwt.Token, bool) {
	token, err := jwtService.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Reference: "ERROR-3",
			Message:   "Token tidak valid",
		}, http.StatusUnauthorized)
		c.Abort()
		return nil, true
	}
	return token, false
}

func CheckBearerTokenPrefix(tokenString string, c *gin.Context) bool {
	if !strings.HasPrefix(tokenString, "Bearer ") {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Reference: "ERROR-2",
			Message:   "Token tidak valid",
		}, http.StatusUnauthorized)
		c.Abort()
		return true
	}
	return false
}

func CheckTokenExist(c *gin.Context) (string, bool) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Reference: "ERROR-1",
			Message:   "Membutuhkan token",
		}, http.StatusUnauthorized)
		c.Abort()
		return "", true
	}
	return tokenString, false
}
