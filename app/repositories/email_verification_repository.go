package repositories

import (
	"golang_starter_kit_2025/app/models"
	"time"

	"gorm.io/gorm"
)

type EmailVerificationRepositoryInterface interface {
	Create(token *models.EmailVerificationToken) error
	FindByToken(token string) (*models.EmailVerificationToken, error)
	FindByUserID(userID uint) (*models.EmailVerificationToken, error)
	DeleteByUserID(userID uint) error
	DeleteExpired() error
}

type EmailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepositoryInterface {
	return &EmailVerificationRepository{db: db}
}

func (r *EmailVerificationRepository) Create(token *models.EmailVerificationToken) error {
	return r.db.Create(token).Error
}

func (r *EmailVerificationRepository) FindByToken(token string) (*models.EmailVerificationToken, error) {
	var verifyToken models.EmailVerificationToken
	err := r.db.Where("token = ?", token).First(&verifyToken).Error
	if err != nil {
		return nil, err
	}
	return &verifyToken, nil
}

func (r *EmailVerificationRepository) FindByUserID(userID uint) (*models.EmailVerificationToken, error) {
	var verifyToken models.EmailVerificationToken
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").First(&verifyToken).Error
	if err != nil {
		return nil, err
	}
	return &verifyToken, nil
}

func (r *EmailVerificationRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.EmailVerificationToken{}).Error
}

func (r *EmailVerificationRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.EmailVerificationToken{}).Error
}
