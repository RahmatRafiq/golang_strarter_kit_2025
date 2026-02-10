package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"strconv"

	"github.com/rs/zerolog/log"
)

type RoleService struct {
	roleRepo       interfaces.RoleRepositoryInterface
	permissionRepo interfaces.PermissionRepositoryInterface
}

func NewRoleService(roleRepo interfaces.RoleRepositoryInterface, permissionRepo interfaces.PermissionRepositoryInterface) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (s *RoleService) GetAll() ([]models.Role, error) {
	return s.roleRepo.GetAll()
}

func (s *RoleService) List(page, limit int) ([]models.Role, int64, error) {
	return s.roleRepo.List(page, limit)
}

func (s *RoleService) Put(updatedRole models.Role) (models.Role, error) {
	if updatedRole.ID == 0 {
		err := s.roleRepo.Create(&updatedRole)
		if err != nil {
			log.Error().
				Err(err).
				Str("name", updatedRole.Name).
				Msg("Failed to create role")
			return updatedRole, err
		}
		log.Info().
			Uint("role_id", updatedRole.ID).
			Str("name", updatedRole.Name).
			Msg("Role created successfully")
		return updatedRole, nil
	} else {
		err := s.roleRepo.Update(&updatedRole)
		if err != nil {
			log.Error().
				Err(err).
				Uint("role_id", updatedRole.ID).
				Msg("Failed to update role")
			return updatedRole, err
		}
		role, err := s.roleRepo.FindByID(updatedRole.ID)
		if err != nil {
			log.Error().
				Err(err).
				Uint("role_id", updatedRole.ID).
				Msg("Failed to fetch updated role")
			return updatedRole, err
		}
		log.Info().
			Uint("role_id", updatedRole.ID).
			Str("name", updatedRole.Name).
			Msg("Role updated successfully")
		return *role, nil
	}
}

func (s *RoleService) Create(role *models.Role) error {
	return s.roleRepo.Create(role)
}

func (s *RoleService) Update(role *models.Role) error {
	return s.roleRepo.Update(role)
}

func (s *RoleService) Delete(id string) error {
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", id).
			Msg("Invalid role ID format for deletion")
		return err
	}

	err = s.roleRepo.Delete(uint(roleID))
	if err != nil {
		log.Error().
			Err(err).
			Uint("role_id", uint(roleID)).
			Msg("Failed to delete role")
		return err
	}

	log.Info().
		Uint("role_id", uint(roleID)).
		Msg("Role deleted successfully")
	return nil
}

func (s *RoleService) DeleteByID(id uint) error {
	return s.roleRepo.Delete(id)
}

func (s *RoleService) AssignPermissionsToRole(roleId string, permissionIDs []uint) error {
	roleID, err := strconv.ParseUint(roleId, 10, 32)
	if err != nil {
		log.Error().
			Err(err).
			Str("role_id", roleId).
			Msg("Invalid role ID format for permission assignment")
		return err
	}

	err = s.roleRepo.AssignPermissions(uint(roleID), permissionIDs)
	if err != nil {
		log.Error().
			Err(err).
			Uint("role_id", uint(roleID)).
			Interface("permission_ids", permissionIDs).
			Msg("Failed to assign permissions to role")
		return err
	}

	log.Info().
		Uint("role_id", uint(roleID)).
		Interface("permission_ids", permissionIDs).
		Msg("Permissions assigned to role successfully")
	return nil
}

func (s *RoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}

func (s *RoleService) GetPermissionsByRoleId(roleId string) ([]models.Permission, error) {
	roleID, err := strconv.ParseUint(roleId, 10, 32)
	if err != nil {
		return nil, err
	}

	return s.roleRepo.GetPermissions(uint(roleID))
}

func (s *RoleService) GetPermissions(roleID uint) ([]models.Permission, error) {
	return s.roleRepo.GetPermissions(roleID)
}

func (s *RoleService) FindByID(id uint) (*models.Role, error) {
	return s.roleRepo.FindByID(id)
}

func (s *RoleService) FindByName(name string) (*models.Role, error) {
	return s.roleRepo.FindByName(name)
}

func (s *RoleService) ExistsByName(name string) (bool, error) {
	return s.roleRepo.ExistsByName(name)
}

func (s *RoleService) Count() (int64, error) {
	return s.roleRepo.Count()
}
