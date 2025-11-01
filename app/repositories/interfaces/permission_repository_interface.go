package interfaces

import "golang_starter_kit_2025/app/models"

// PermissionRepositoryInterface defines contract for permission data operations
type PermissionRepositoryInterface interface {
	// Create creates a new permission
	Create(permission *models.Permission) error

	// Update updates existing permission
	Update(permission *models.Permission) error

	// Delete soft deletes a permission by ID
	Delete(id uint) error

	// FindByID finds permission by ID
	FindByID(id uint) (*models.Permission, error)

	// FindByName finds permission by name
	FindByName(name string) (*models.Permission, error)

	// List returns paginated list of permissions
	List(page, limit int) ([]models.Permission, int64, error)

	// ExistsByName checks if permission with name exists
	ExistsByName(name string) (bool, error)

	// Count returns total number of permissions
	Count() (int64, error)

	// GetAll returns all permissions without pagination
	GetAll() ([]models.Permission, error)

	// FindByIDs finds permissions by multiple IDs
	FindByIDs(ids []uint) ([]models.Permission, error)
}
