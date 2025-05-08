package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"

	"gorm.io/gorm/clause"
)

type UserService struct{}

func (*UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := facades.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (*UserService) Find(id string) (models.User, error) {
	var user models.User
	if err := facades.DB.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (*UserService) Put(user models.User) (models.User, error) {

	if err := facades.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"username", "email", "password", "fcm_token", "updated_at"}),
	}).Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (*UserService) Delete(id string) error {
	var user models.User
	if err := facades.DB.First(&user, id).Error; err != nil {
		return err
	}
	return facades.DB.Delete(&user).Error
}

func (*UserService) AssignRolesToUser(userId string, roles []uint) error {
	var user models.User
	if err := facades.DB.First(&user, userId).Error; err != nil {
		return err
	}

	// Clear existing roles for the user
	facades.DB.Where("user_id = ?", user.ID).Delete(&models.UserHasRole{})

	// Assign new roles
	for _, roleId := range roles {
		userRole := models.UserHasRole{
			UserID: user.ID,
			RoleID: roleId,
		}
		if err := facades.DB.Create(&userRole).Error; err != nil {
			return err
		}
	}

	return nil
}
func (*UserService) GetRolesByUserId(userId string) ([]models.Role, error) {
	var roles []models.Role
	if err := facades.DB.Table("roles").
		Select("roles.*").
		Joins("join user_has_roles on roles.id = user_has_roles.role_id").
		Where("user_has_roles.user_id = ?", userId).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
