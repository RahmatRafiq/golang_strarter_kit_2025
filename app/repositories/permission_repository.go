package repositories

import (
	"errors"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/models/scopes"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"gorm.io/gorm"
)

// permissionRepository implements PermissionRepositoryInterface
type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(db *gorm.DB) interfaces.PermissionRepositoryInterface {
	return &permissionRepository{db: db}
}

// Create creates a new permission
func (r *permissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

// Update updates existing permission
func (r *permissionRepository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

// Delete soft deletes a permission by ID
func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

// FindByID finds permission by ID
func (r *permissionRepository) FindByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}
	return &permission, nil
}

// FindByName finds permission by name
func (r *permissionRepository) FindByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}
	return &permission, nil
}

// List returns paginated list of permissions
func (r *permissionRepository) List(page, limit int) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	// Count total
	if err := r.db.Model(&models.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err := r.db.Scopes(scopes.Paginate(page, limit)).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// ExistsByName checks if permission with name exists
func (r *permissionRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Permission{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count returns total number of permissions
func (r *permissionRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Permission{}).Count(&count).Error
	return count, err
}

// GetAll returns all permissions without pagination
func (r *permissionRepository) GetAll() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

// FindByIDs finds permissions by multiple IDs
func (r *permissionRepository) FindByIDs(ids []uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Find(&permissions, ids).Error
	return permissions, err
}
