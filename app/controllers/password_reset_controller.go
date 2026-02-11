package controllers

import (
	"net/http"

	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// PasswordResetController handles password reset operations
type PasswordResetController struct {
	passwordResetService *services.PasswordResetService
}

// NewPasswordResetController creates a new password reset controller
func NewPasswordResetController(passwordResetService *services.PasswordResetService) *PasswordResetController {
	return &PasswordResetController{
		passwordResetService: passwordResetService,
	}
}

// ForgotPassword handles password reset requests
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body requests.ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} map[string]interface{} "Reset email sent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/forgot-password [post]
func (c *PasswordResetController) ForgotPassword(ctx *gin.Context) {
	var req requests.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid forgot password request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Request password reset
	if err := c.passwordResetService.RequestPasswordReset(req.Email); err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to request password reset")
		// Don't reveal internal errors to prevent user enumeration
		ctx.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword handles password reset with token
// @Summary Reset password with token
// @Description Reset user password using valid reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body requests.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} map[string]interface{} "Password reset successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/reset-password [post]
func (c *PasswordResetController) ResetPassword(ctx *gin.Context) {
	var req requests.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid reset password request")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Reset password
	if err := c.passwordResetService.ResetPassword(req.Token, req.NewPassword); err != nil {
		log.Error().Err(err).Msg("Failed to reset password")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Password reset failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password has been reset successfully. You can now login with your new password.",
	})
}
