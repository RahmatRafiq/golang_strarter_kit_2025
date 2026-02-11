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
	"golang_starter_kit_2025/app/workers/tasks"
	"golang_starter_kit_2025/facades"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type EmailVerificationService struct {
	verificationRepo repositories.EmailVerificationRepositoryInterface
	userRepo         interfaces.UserRepositoryInterface
	emailService     *EmailService
}

func NewEmailVerificationService(
	verificationRepo repositories.EmailVerificationRepositoryInterface,
	userRepo interfaces.UserRepositoryInterface,
) *EmailVerificationService {
	return &EmailVerificationService{
		verificationRepo: verificationRepo,
		userRepo:         userRepo,
		emailService:     NewEmailService(),
	}
}

func (s *EmailVerificationService) SendVerificationEmail(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.EmailVerifiedAt != nil {
		return fmt.Errorf("email already verified")
	}

	token, err := generateVerificationToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	if err := s.verificationRepo.DeleteByUserID(userID); err != nil {
		log.Warn().Err(err).Uint("user_id", userID).Msg("Failed to delete old verification tokens")
	}

	expiryHours := helpers.GetEnvInt("EMAIL_VERIFICATION_EXPIRY_HOURS", 24)
	verifyToken := &models.EmailVerificationToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(expiryHours) * time.Hour),
	}

	if err := s.verificationRepo.Create(verifyToken); err != nil {
		return fmt.Errorf("failed to create verification token: %w", err)
	}

	appURL := helpers.GetEnv("APP_URL", "http://localhost:3000")
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", appURL, token)

	queueClient := facades.GetQueueClient()
	if queueClient != nil {
		task, err := tasks.NewSendEmailTask(
			user.Email,
			"Verify Your Email Address",
			fmt.Sprintf("Click here to verify your email: %s", verifyURL),
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create email task")
		} else {
			if _, err := queueClient.EnqueueWithPriority(context.Background(), task, "high"); err != nil {
				log.Error().Err(err).Msg("Failed to enqueue verification email")
			}
		}
	}

	if s.emailService != nil && s.emailService.config.IsConfigured() {
		go func() {
			if err := s.emailService.SendVerificationEmail(user.Email, user.Username, verifyURL, expiryHours); err != nil {
				log.Error().Err(err).Str("email", user.Email).Msg("Failed to send verification email")
			}
		}()
	}

	log.Info().Uint("user_id", userID).Msg("Verification email sent")
	return nil
}

func (s *EmailVerificationService) VerifyEmail(token string) error {
	verifyToken, err := s.verificationRepo.FindByToken(token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid or expired verification token")
		}
		return fmt.Errorf("failed to find verification token: %w", err)
	}

	if verifyToken.IsExpired() {
		return fmt.Errorf("verification token has expired")
	}

	user, err := s.userRepo.FindByID(verifyToken.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.EmailVerifiedAt != nil {
		return fmt.Errorf("email already verified")
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	if err := s.verificationRepo.DeleteByUserID(user.ID); err != nil {
		log.Warn().Err(err).Uint("user_id", user.ID).Msg("Failed to delete verification token")
	}

	log.Info().Uint("user_id", user.ID).Msg("Email verified successfully")
	return nil
}

func generateVerificationToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
