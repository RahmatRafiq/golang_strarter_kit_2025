package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"strconv"
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
		return updatedRole, err
	} else {
		err := s.roleRepo.Update(&updatedRole)
		if err != nil {
			return updatedRole, err
		}
		role, err := s.roleRepo.FindByID(updatedRole.ID)
		if err != nil {
			return updatedRole, err
		}
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
		return err
	}
	return s.roleRepo.Delete(uint(roleID))
}

func (s *RoleService) DeleteByID(id uint) error {
	return s.roleRepo.Delete(id)
}

func (s *RoleService) AssignPermissionsToRole(roleId string, permissionIDs []uint) error {
	roleID, err := strconv.ParseUint(roleId, 10, 32)
	if err != nil {
		return err
	}

	return s.roleRepo.AssignPermissions(uint(roleID), permissionIDs)
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
