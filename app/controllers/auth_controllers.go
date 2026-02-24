package controllers

import (
	"net/http"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// @Summary		Login
// @Description	API for user login with email and password
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			body	body		requests.LoginRequest	true	"Login data"
// @Success		200		{object}	helpers.ResponseParams[any]
// @Router			/auth/login [put]
func (c *AuthController) Login(ctx *gin.Context) {
	// Check if using legacy endpoint
	if ctx.Request.Method == "PUT" && ctx.FullPath() == "/auth/login" {
		log.Warn().
			Str("endpoint", "PUT /auth/login").
			Str("replacement", "POST /api/v1/auth/login").
			Str("sunset_date", "v2.0.0 (Q3 2026)").
			Msg("DEPRECATED: Using legacy auth endpoint")
	}

	var loginData requests.LoginRequest
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		log.Warn().
			Err(err).
			Str("email", loginData.Email).
			Msg("Login request binding failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Login(loginData)
	if err != nil {
		log.Warn().
			Err(err).
			Str("email", loginData.Email).
			Msg("Login attempt failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Str("email", loginData.Email).
		Msg("User logged in via controller")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Token: token}, 200)
}

// @Summary		Logout
// @Description	API for user logout, requires valid token
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	helpers.ResponseParams[any]
// @Router			/auth/logout [get]
func (c *AuthController) Logout(ctx *gin.Context) {
	// get token from context
	tokenString, exists := ctx.Get("token")
	if !exists {
		log.Warn().Msg("Logout attempt without token")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
		return
	}

	token, ok := tokenString.(string)
	if !ok {
		log.Warn().Msg("Logout attempt with invalid token format")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	err := c.service.Logout(token)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Logout failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Info().Msg("User logged out successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Message: "Logout successful"}, 200)
}

// @Summary		Refresh Token
// @Description	API to refresh access token using refresh token
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			body	body		requests.RefreshTokenRequest	true	"Refresh token"
// @Success		200		{object}	helpers.ResponseParams[any]
// @Router			/auth/refresh [post]
func (c *AuthController) Refresh(ctx *gin.Context) {
	// Check if using legacy endpoint
	if ctx.FullPath() == "/auth/refresh" {
		log.Warn().
			Str("endpoint", "POST /auth/refresh").
			Str("replacement", "POST /api/v1/auth/refresh").
			Str("sunset_date", "v2.0.0 (Q3 2026)").
			Msg("DEPRECATED: Using legacy auth endpoint")
	}

	var refreshRequest requests.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&refreshRequest); err != nil {
		log.Warn().
			Err(err).
			Msg("Refresh token request binding failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	token, err := c.service.RefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Token refresh failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Info().Msg("Token refreshed successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Token: token}, 200)
}
