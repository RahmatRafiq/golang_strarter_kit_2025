package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"golang_starter_kit_2025/app/workers"
	"golang_starter_kit_2025/app/workers/tasks"
	"golang_starter_kit_2025/facades"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// PasswordResetService handles password reset operations
type PasswordResetService struct {
	passwordResetRepo repositories.PasswordResetRepositoryInterface
	userRepo          interfaces.UserRepositoryInterface
	emailService      *EmailService
	queueClient       *workers.Client
}

// NewPasswordResetService creates a new password reset service
func NewPasswordResetService(
	passwordResetRepo repositories.PasswordResetRepositoryInterface,
	userRepo interfaces.UserRepositoryInterface,
) *PasswordResetService {
	return &PasswordResetService{
		passwordResetRepo: passwordResetRepo,
		userRepo:          userRepo,
		emailService:      NewEmailService(),
		queueClient:       facades.GetQueueClient(),
	}
}

// RequestPasswordReset creates a password reset token and sends email
func (s *PasswordResetService) RequestPasswordReset(email string) error {
	// Check if user exists
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Don't reveal if email exists or not (security best practice)
			log.Warn().Str("email", email).Msg("Password reset requested for non-existent email")
			return nil
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Generate secure random token (64 chars hex = 32 bytes random)
	token, err := generatePasswordResetToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Delete any existing tokens for this email
	if err := s.passwordResetRepo.DeleteByEmail(email); err != nil {
		log.Warn().Err(err).Str("email", email).Msg("Failed to delete old reset tokens")
	}

	// Create reset token record
	// Note: Tokens are stored as plain hex strings for lookup efficiency
	// Security is maintained through:
	// 1. Cryptographically secure random generation (32 bytes = 256 bits)
	// 2. Short expiry time (1 hour default)
	// 3. One-time use enforcement
	// 4. HTTPS transmission
	expiryHours := helpers.GetEnvInt("PASSWORD_RESET_EXPIRY_HOURS", 1)
	resetToken := &models.PasswordResetToken{
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(expiryHours) * time.Hour),
	}

	if err := s.passwordResetRepo.Create(resetToken); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// Generate reset URL
	appURL := helpers.GetEnv("APP_URL", "http://localhost:3000")
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", appURL, token)

	// Send reset email via queue
	if s.queueClient != nil {
		task, err := tasks.NewSendEmailTask(
			email,
			"Reset Your Password",
			fmt.Sprintf("Click here to reset your password: %s\n\nThis link will expire in %d hour(s).", resetURL, expiryHours),
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create email task")
		} else {
			if _, err := s.queueClient.EnqueueWithPriority(context.Background(), task, "critical"); err != nil {
				log.Error().Err(err).Msg("Failed to enqueue password reset email")
			}
		}
	}

	// Also send via email service directly as fallback
	if s.emailService != nil && s.emailService.config.IsConfigured() {
		go func() {
			_ = s.emailService.SendPasswordResetEmail(email, user.Username, resetURL, expiryHours)
		}()
	}

	log.Info().Str("email", email).Msg("Password reset requested")
	return nil
}

// ResetPassword validates token and resets the password
func (s *PasswordResetService) ResetPassword(token, newPassword string) error {
	// Find reset token
	resetToken, err := s.passwordResetRepo.FindByToken(token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or expired reset token")
		}
		return fmt.Errorf("failed to find reset token: %w", err)
	}

	// Validate token
	if !resetToken.IsValid() {
		if resetToken.IsExpired() {
			return fmt.Errorf("reset token has expired")
		}
		if resetToken.IsUsed() {
			return fmt.Errorf("reset token has already been used")
		}
	}

	// Get user
	user, err := s.userRepo.FindByEmail(resetToken.Email)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Update password
	user.Password = newPassword // BeforeUpdate hook will hash it
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	resetToken.MarkAsUsed()
	if err := s.passwordResetRepo.Update(resetToken); err != nil {
		log.Warn().Err(err).Msg("Failed to mark reset token as used")
	}

	log.Info().Str("email", resetToken.Email).Msg("Password reset successful")
	return nil
}

// CleanupExpiredTokens removes expired reset tokens
func (s *PasswordResetService) CleanupExpiredTokens() error {
	return s.passwordResetRepo.DeleteExpired()
}

// generatePasswordResetToken generates a cryptographically secure random token
func generatePasswordResetToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
