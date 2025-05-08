package services

import (
	"errors"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"
)

type RoleService struct{}

func (*RoleService) GetAll() ([]models.Role, error) {
	var roles []models.Role
	if err := facades.DB.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (*RoleService) Put(updatedRole models.Role) (models.Role, error) {
	var role models.Role

	if count := facades.DB.Model(&models.Role{}).Where("id = ?", updatedRole.ID).Find(&map[string]interface{}{}).RowsAffected; count == 0 {
		if err := facades.DB.Create(&updatedRole).Error; err != nil {
			return role, err
		}
	} else {
		if err := facades.DB.Where("id = ?", updatedRole.ID).Updates(&updatedRole).Error; err != nil {
			return role, err
		}

		if err := facades.DB.First(&role, updatedRole.ID).Error; err != nil {
			return role, err
		}
	}

	return role, nil
}

func (*RoleService) Delete(id string) error {
	var role models.Role
	if err := facades.DB.First(&role, id).Error; err != nil {
		return err
	}
	return facades.DB.Delete(&role).Error
}

func (*RoleService) AssignPermissionsToRole(roleId string, permissions []uint) error {
	var role models.Role
	if err := facades.DB.First(&role, roleId).Error; err != nil {
		return err
	}

	// Validasi permissions sebelum diassign
	var validPermissions []uint
	facades.DB.Table("permissions").Where("id IN ?", permissions).Pluck("id", &validPermissions)

	if len(validPermissions) != len(permissions) {
		return errors.New("one or more permission IDs are invalid")
	}

	// Clear existing permissions for the role
	facades.DB.Where("role_id = ?", role.ID).Delete(&models.RoleHasPermissions{})

	// Assign new permissions
	for _, permId := range validPermissions {
		rolePerm := models.RoleHasPermissions{
			RoleID:       role.ID,
			PermissionID: permId,
		}
		if err := facades.DB.Create(&rolePerm).Error; err != nil {
			return err
		}
	}

	return nil
}

func (*RoleService) GetPermissionsByRoleId(roleId string) ([]models.Permission, error) {
	var permissions []models.Permission
	if err := facades.DB.Table("permissions").
		Select("permissions.*").
		Joins("join role_has_permissions on permissions.id = role_has_permissions.permission_id").
		Where("role_has_permissions.role_id = ?", roleId).
		Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
