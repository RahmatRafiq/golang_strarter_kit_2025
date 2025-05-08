package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"
)

type PermissionService struct{}

func (*PermissionService) GetAll() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := facades.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (*PermissionService) Put(updatedPermission models.Permission) (models.Permission, error) {
	var permission models.Permission

	if count := facades.DB.Model(&models.Permission{}).Where("id = ?", updatedPermission.ID).Find(&map[string]interface{}{}).RowsAffected; count == 0 {
		if err := facades.DB.Create(&updatedPermission).Error; err != nil {
			return permission, err
		}
	} else {
		if err := facades.DB.Where("id = ?", updatedPermission.ID).Updates(&updatedPermission).Error; err != nil {
			return permission, err
		}

		if err := facades.DB.First(&permission, updatedPermission.ID).Error; err != nil {
			return permission, err
		}
	}

	return permission, nil
}

func (*PermissionService) Delete(id string) error {
	var permission models.Permission
	if err := facades.DB.First(&permission, id).Error; err != nil {
		return err
	}
	return facades.DB.Delete(&permission).Error
}
