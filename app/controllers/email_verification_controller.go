package controllers

import (
	"net/http"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"

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
		log.Warn().Msg("Unauthorized verification email request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Unauthorized",
			Reference: "ERROR-1",
		}, http.StatusUnauthorized)
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		log.Error().Msg("Invalid user_id type")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "Invalid user ID",
			Reference: "ERROR-2",
		}, http.StatusInternalServerError)
		return
	}

	if err := c.verificationService.SendVerificationEmail(userIDUint); err != nil {
		log.Error().Err(err).Uint("user_id", userIDUint).Msg("Failed to send verification email")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Failed to send verification email",
			Reference: "ERROR-3",
		}, http.StatusInternalServerError)
		return
	}

	log.Info().Uint("user_id", userIDUint).Msg("Verification email sent successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{
		Message: "Verification email sent successfully",
	}, http.StatusOK)
}

func (c *EmailVerificationController) VerifyEmail(ctx *gin.Context) {
	var req requests.VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid verification request")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Invalid request",
			Reference: "ERROR-4",
		}, http.StatusBadRequest)
		return
	}

	if err := c.verificationService.VerifyEmail(req.Token); err != nil {
		log.Error().Err(err).Str("token", req.Token).Msg("Failed to verify email")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": err.Error()},
			Message:   "Email verification failed",
			Reference: "ERROR-3",
		}, http.StatusBadRequest)
		return
	}

	log.Info().Str("token", req.Token).Msg("Email verified successfully")
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{
		Message: "Email verified successfully",
	}, http.StatusOK)
}
