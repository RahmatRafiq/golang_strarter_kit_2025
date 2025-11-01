package interfaces

import "golang_starter_kit_2025/app/models"

// RoleRepositoryInterface defines contract for role data operations
type RoleRepositoryInterface interface {
	// Create creates a new role
	Create(role *models.Role) error

	// Update updates existing role
	Update(role *models.Role) error

	// Delete soft deletes a role by ID
	Delete(id uint) error

	// FindByID finds role by ID
	FindByID(id uint) (*models.Role, error)

	// FindByName finds role by name
	FindByName(name string) (*models.Role, error)

	// List returns paginated list of roles
	List(page, limit int) ([]models.Role, int64, error)

	// ListWithPermissions returns roles with preloaded permissions
	ListWithPermissions(page, limit int) ([]models.Role, int64, error)

	// AssignPermissions assigns permissions to role
	AssignPermissions(roleID uint, permissionIDs []uint) error

	// GetPermissions gets role's permissions
	GetPermissions(roleID uint) ([]models.Permission, error)

	// ExistsByName checks if role with name exists
	ExistsByName(name string) (bool, error)

	// Count returns total number of roles
	Count() (int64, error)

	// GetAll returns all roles without pagination
	GetAll() ([]models.Role, error)
}
