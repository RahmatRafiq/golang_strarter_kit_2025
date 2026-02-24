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

var jwtService services.JwtService

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		skipAuth := os.Getenv("SKIP_AUTH")

		// Development bypass (with security warning)
		// Note: Production check moved to startup validation in bootstrap/main.go
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

		mapClaims, err := jwtService.ExtractClaims(token)
		if err != nil {
			helpers.ResponseError(c, &helpers.ResponseParams[any]{
				Reference: "ERROR-4",
				Message:   "Invalid token",
			}, http.StatusUnauthorized)
			c.Abort()
			return
		}

		claims, err := casts.ParseJwtClaims(mapClaims)
		if err != nil {
			helpers.ResponseError(c, &helpers.ResponseParams[any]{
				Reference: "ERROR-4",
				Message:   "Invalid token",
			}, http.StatusUnauthorized)
			c.Abort()
			return
		}

		if claims.ExpiredAt < time.Now().Unix() {
			helpers.ResponseError(c, &helpers.ResponseParams[any]{
				Reference: "ERROR-4",
				Message:   "Token has expired",
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
			Message:   "Invalid token",
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
			Message:   "Invalid token",
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
			Message:   "Token required",
		}, http.StatusUnauthorized)
		c.Abort()
		return "", true
	}
	return tokenString, false
}
