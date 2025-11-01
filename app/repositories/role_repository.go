package repositories

import (
	"errors"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/models/scopes"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"gorm.io/gorm"
)

// roleRepository implements RoleRepositoryInterface
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) interfaces.RoleRepositoryInterface {
	return &roleRepository{db: db}
}

// Create creates a new role
func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// Update updates existing role
func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete soft deletes a role by ID
func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// FindByID finds role by ID
func (r *roleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

// FindByName finds role by name
func (r *roleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

// List returns paginated list of roles
func (r *roleRepository) List(page, limit int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	// Count total
	if err := r.db.Model(&models.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err := r.db.Scopes(scopes.Paginate(page, limit)).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// ListWithPermissions returns roles with preloaded permissions
func (r *roleRepository) ListWithPermissions(page, limit int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	// Count total
	if err := r.db.Model(&models.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data with permissions
	err := r.db.Preload("Permissions").Scopes(scopes.Paginate(page, limit)).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// AssignPermissions assigns permissions to role
func (r *roleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Find role
		var role models.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		// Find permissions
		var permissions []models.Permission
		if err := tx.Find(&permissions, permissionIDs).Error; err != nil {
			return err
		}

		// Check if all permissions found
		if len(permissions) != len(permissionIDs) {
			return errors.New("some permissions not found")
		}

		// Clear existing permissions and assign new ones
		if err := tx.Model(&role).Association("Permissions").Replace(permissions); err != nil {
			return err
		}

		return nil
	})
}

// GetPermissions gets role's permissions
func (r *roleRepository) GetPermissions(roleID uint) ([]models.Permission, error) {
	var role models.Role
	if err := r.db.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

// ExistsByName checks if role with name exists
func (r *roleRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Role{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count returns total number of roles
func (r *roleRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Role{}).Count(&count).Error
	return count, err
}

// GetAll returns all roles without pagination
func (r *roleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).Error
	return roles, err
}
