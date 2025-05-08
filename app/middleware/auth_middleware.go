package middleware

import (
	"net/http"
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
		tokenString, shouldReturn := CheckTokenExist(c)
		if shouldReturn {
			return
		}

		// Token format: "Bearer <token>"
		shouldReturn1 := CheckBearerTokenPrefix(tokenString, c)
		if shouldReturn1 {
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// token, shouldReturn2 := CheckTokenValidity(tokenString, c)
		token, shouldReturn2 := CheckTokenValidity(tokenString, c)
		if shouldReturn2 {
			return
		}

		// NOW: user can login multiple times
		claims := casts.ParseJwtClaims(jwtService.ExtractClaims(token))

		if claims.ExpiredAt < time.Now().Unix() {
			helpers.ResponseError(c, &helpers.ResponseParams[any]{
				Reference: "ERROR-4",
				Message:   "Token sudah kadaluarsa",
			}, http.StatusUnauthorized)
			c.Abort()
			return
		}

		// set token and user id to context
		c.Set("token", tokenString)
		c.Set("user_id", claims.UserID)
		// c.JSON(http.StatusOK, gin.H{"user_id": c.GetString("user_id")})
		// c.Request.WithContext(context.WithValue(c.Request.Context(), "user_id", claims.UserID))
		// var user models.User
		// if err := facades.DB.Where("id = ? AND jwt_token = ?", userId, tokenString).First(&user).Error; err != nil {
		// 	helpers.ResponseError(c, http.StatusUnauthorized, "Token tidak valid", "error_4")
		// 	c.Abort()
		// 	return
		// }

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
