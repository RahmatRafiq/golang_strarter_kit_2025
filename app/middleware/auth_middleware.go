package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"golang_starter_kit_2025/app/casts"
	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("your_secret_key")

var jwtService services.JwtService

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		skipAuth := os.Getenv("SKIP_AUTH")
		if skipAuth == "true" {
			c.Set("user_id", uint(1))
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
