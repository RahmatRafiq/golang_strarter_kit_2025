package repositories

import (
	"errors"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) interfaces.AuthRepositoryInterface {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}
