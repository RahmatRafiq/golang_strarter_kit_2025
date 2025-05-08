package controllers

import (
	"net/http"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{service: service}
}

// @Summary		Login
// @Description	API untuk login dengan email dan password
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			body	body		requests.LoginRequest	true	"Login data"
// @Success		200		{object}	helpers.ResponseParams[any]
// @Router			/auth/login [put]
func (c *AuthController) Login(ctx *gin.Context) {
	var loginData requests.LoginRequest
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.service.Login(loginData)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Token: token}, 200)
}

// @Summary		Logout
// @Description	API untuk logout, membutuhkan token yang valid
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	helpers.ResponseParams[any]
// @Router			/auth/logout [get]
func (c *AuthController) Logout(ctx *gin.Context) {
	// get token from context
	tokenString, _ := ctx.Get("token")

	err := c.service.Logout(tokenString.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Message: "Berhasil logout"}, 200)
}

// @Summary		Refresh Token
// @Description	API untuk refresh token, membutuhkan token yang valid
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200	{object}	helpers.ResponseParams[any]
// @Router			/auth/refresh [get]
func (c *AuthController) Refresh(ctx *gin.Context) {
	tokenString, _ := ctx.Get("token")

	token, err := c.service.RefreshToken(tokenString.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{Token: token}, 200)
}
