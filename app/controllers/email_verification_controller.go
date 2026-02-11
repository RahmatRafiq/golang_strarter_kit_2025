package controllers

import (
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type EmailVerificationController struct {
	verificationService *services.EmailVerificationService
}

func NewEmailVerificationController(verificationService *services.EmailVerificationService) *EmailVerificationController {
	return &EmailVerificationController{
		verificationService: verificationService,
	}
}

func (c *EmailVerificationController) SendVerification(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.verificationService.SendVerificationEmail(userID.(uint)); err != nil {
		log.Error().Err(err).Msg("Failed to send verification email")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send verification email",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent successfully",
	})
}

func (c *EmailVerificationController) VerifyEmail(ctx *gin.Context) {
	var req requests.VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid verification request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	if err := c.verificationService.VerifyEmail(req.Token); err != nil {
		log.Error().Err(err).Msg("Failed to verify email")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Email verification failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}
