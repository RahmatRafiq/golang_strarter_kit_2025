package interfaces

import "golang_starter_kit_2025/app/models"

// CategoryRepositoryInterface defines contract for category data operations
type CategoryRepositoryInterface interface {
	// Create creates a new category
	Create(category *models.Category) error

	// Update updates existing category
	Update(category *models.Category) error

	// Delete soft deletes a category by ID
	Delete(id uint) error

	// FindByID finds category by ID
	FindByID(id uint) (*models.Category, error)

	// FindByIDWithProducts finds category by ID with products preloaded
	FindByIDWithProducts(id uint) (*models.Category, error)

	// FindByName finds category by name
	FindByName(name string) (*models.Category, error)

	// List returns paginated list of categories
	List(page, limit int) ([]models.Category, int64, error)

	// ListWithProducts returns categories with preloaded products
	ListWithProducts(page, limit int) ([]models.Category, int64, error)

	// ExistsByName checks if category with name exists
	ExistsByName(name string) (bool, error)

	// Count returns total number of categories
	Count() (int64, error)

	// GetAll returns all categories without pagination
	GetAll() ([]models.Category, error)
}
