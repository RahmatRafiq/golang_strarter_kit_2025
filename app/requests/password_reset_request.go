package requests

// ForgotPasswordRequest represents the request to initiate password reset
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// ResetPasswordRequest represents the request to reset password with token
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required,min=64,max=64" example:"a1b2c3d4..."`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"newSecurePassword123"`
}
