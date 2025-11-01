package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"strconv"
)

// PermissionService handles permission business logic
type PermissionService struct {
	repo interfaces.PermissionRepositoryInterface
}

// NewPermissionService creates a new permission service
func NewPermissionService(repo interfaces.PermissionRepositoryInterface) *PermissionService {
	return &PermissionService{repo: repo}
}

// GetAll returns all permissions without pagination
func (s *PermissionService) GetAll() ([]models.Permission, error) {
	return s.repo.GetAll()
}

// List returns paginated list of permissions
func (s *PermissionService) List(page, limit int) ([]models.Permission, int64, error) {
	return s.repo.List(page, limit)
}

// Put creates or updates a permission (based on ID)
func (s *PermissionService) Put(updatedPermission models.Permission) (models.Permission, error) {
	if updatedPermission.ID == 0 {
		// Create new permission
		err := s.repo.Create(&updatedPermission)
		return updatedPermission, err
	} else {
		// Update existing permission
		err := s.repo.Update(&updatedPermission)
		if err != nil {
			return updatedPermission, err
		}
		// Return updated permission
		permission, err := s.repo.FindByID(updatedPermission.ID)
		if err != nil {
			return updatedPermission, err
		}
		return *permission, nil
	}
}

// Create creates a new permission
func (s *PermissionService) Create(permission *models.Permission) error {
	return s.repo.Create(permission)
}

// Update updates existing permission
func (s *PermissionService) Update(permission *models.Permission) error {
	return s.repo.Update(permission)
}

// Delete deletes permission by ID (string for backward compatibility)
func (s *PermissionService) Delete(id string) error {
	permissionID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}
	return s.repo.Delete(uint(permissionID))
}

// DeleteByID deletes permission by ID (uint version)
func (s *PermissionService) DeleteByID(id uint) error {
	return s.repo.Delete(id)
}

// FindByID finds permission by ID
func (s *PermissionService) FindByID(id uint) (*models.Permission, error) {
	return s.repo.FindByID(id)
}

// FindByName finds permission by name
func (s *PermissionService) FindByName(name string) (*models.Permission, error) {
	return s.repo.FindByName(name)
}

// FindByIDs finds permissions by multiple IDs
func (s *PermissionService) FindByIDs(ids []uint) ([]models.Permission, error) {
	return s.repo.FindByIDs(ids)
}

// ExistsByName checks if permission with name exists
func (s *PermissionService) ExistsByName(name string) (bool, error) {
	return s.repo.ExistsByName(name)
}

// Count returns total number of permissions
func (s *PermissionService) Count() (int64, error) {
	return s.repo.Count()
}
