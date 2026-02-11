package repositories

import (
	"time"

	"golang_starter_kit_2025/app/models"

	"gorm.io/gorm"
)

// PasswordResetRepositoryInterface defines the contract for password reset operations
type PasswordResetRepositoryInterface interface {
	Create(token *models.PasswordResetToken) error
	FindByToken(token string) (*models.PasswordResetToken, error)
	FindByEmail(email string) (*models.PasswordResetToken, error)
	Update(token *models.PasswordResetToken) error
	DeleteExpired() error
	DeleteByEmail(email string) error
}

// PasswordResetRepository handles password reset token database operations
type PasswordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepositoryInterface {
	return &PasswordResetRepository{db: db}
}

// Create creates a new password reset token
func (r *PasswordResetRepository) Create(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

// FindByToken finds a password reset token by token string
func (r *PasswordResetRepository) FindByToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("token = ?", token).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

// FindByEmail finds the latest password reset token for an email
func (r *PasswordResetRepository) FindByEmail(email string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("email = ?", email).
		Order("created_at DESC").
		First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

// Update updates a password reset token
func (r *PasswordResetRepository) Update(token *models.PasswordResetToken) error {
	return r.db.Save(token).Error
}

// DeleteExpired deletes all expired tokens
func (r *PasswordResetRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordResetToken{}).Error
}

// DeleteByEmail deletes all password reset tokens for an email
func (r *PasswordResetRepository) DeleteByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.PasswordResetToken{}).Error
}
