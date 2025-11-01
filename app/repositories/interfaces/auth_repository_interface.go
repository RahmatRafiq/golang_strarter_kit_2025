package interfaces

import "golang_starter_kit_2025/app/models"

type AuthRepositoryInterface interface {
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
}
