package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"strconv"
)

// RoleService handles role business logic
type RoleService struct {
	roleRepo       interfaces.RoleRepositoryInterface
	permissionRepo interfaces.PermissionRepositoryInterface
}

// NewRoleService creates a new role service
func NewRoleService(roleRepo interfaces.RoleRepositoryInterface, permissionRepo interfaces.PermissionRepositoryInterface) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// GetAll returns all roles without pagination
func (s *RoleService) GetAll() ([]models.Role, error) {
	return s.roleRepo.GetAll()
}

// List returns paginated list of roles
func (s *RoleService) List(page, limit int) ([]models.Role, int64, error) {
	return s.roleRepo.List(page, limit)
}

// Put creates or updates a role (based on ID)
func (s *RoleService) Put(updatedRole models.Role) (models.Role, error) {
	if updatedRole.ID == 0 {
		// Create new role
		err := s.roleRepo.Create(&updatedRole)
		return updatedRole, err
	} else {
		// Update existing role
		err := s.roleRepo.Update(&updatedRole)
		if err != nil {
			return updatedRole, err
		}
		// Return updated role
		role, err := s.roleRepo.FindByID(updatedRole.ID)
		if err != nil {
			return updatedRole, err
		}
		return *role, nil
	}
}

// Create creates a new role
func (s *RoleService) Create(role *models.Role) error {
	return s.roleRepo.Create(role)
}

// Update updates existing role
func (s *RoleService) Update(role *models.Role) error {
	return s.roleRepo.Update(role)
}

// Delete deletes role by ID (string for backward compatibility)
func (s *RoleService) Delete(id string) error {
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}
	return s.roleRepo.Delete(uint(roleID))
}

// DeleteByID deletes role by ID (uint version)
func (s *RoleService) DeleteByID(id uint) error {
	return s.roleRepo.Delete(id)
}

// AssignPermissionsToRole assigns permissions to a role
func (s *RoleService) AssignPermissionsToRole(roleId string, permissionIDs []uint) error {
	roleID, err := strconv.ParseUint(roleId, 10, 32)
	if err != nil {
		return err
	}

	return s.roleRepo.AssignPermissions(uint(roleID), permissionIDs)
}

// AssignPermissions assigns permissions to a role (uint version)
func (s *RoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}

// GetPermissionsByRoleId gets permissions for a role
func (s *RoleService) GetPermissionsByRoleId(roleId string) ([]models.Permission, error) {
	roleID, err := strconv.ParseUint(roleId, 10, 32)
	if err != nil {
		return nil, err
	}

	return s.roleRepo.GetPermissions(uint(roleID))
}

// GetPermissions gets permissions for a role (uint version)
func (s *RoleService) GetPermissions(roleID uint) ([]models.Permission, error) {
	return s.roleRepo.GetPermissions(roleID)
}

// FindByID finds role by ID
func (s *RoleService) FindByID(id uint) (*models.Role, error) {
	return s.roleRepo.FindByID(id)
}

// FindByName finds role by name
func (s *RoleService) FindByName(name string) (*models.Role, error) {
	return s.roleRepo.FindByName(name)
}

// ExistsByName checks if role with name exists
func (s *RoleService) ExistsByName(name string) (bool, error) {
	return s.roleRepo.ExistsByName(name)
}

// Count returns total number of roles
func (s *RoleService) Count() (int64, error) {
	return s.roleRepo.Count()
}
