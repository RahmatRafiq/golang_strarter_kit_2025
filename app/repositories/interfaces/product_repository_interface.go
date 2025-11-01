package interfaces

import "golang_starter_kit_2025/app/models"

// ProductRepositoryInterface defines contract for product data operations
type ProductRepositoryInterface interface {
	// Create creates a new product
	Create(product *models.Product) error

	// Update updates existing product
	Update(product *models.Product) error

	// Delete soft deletes a product by ID
	Delete(id uint) error

	// FindByID finds product by ID
	FindByID(id uint) (*models.Product, error)

	// FindByIDWithCategory finds product by ID with category preloaded
	FindByIDWithCategory(id uint) (*models.Product, error)

	// List returns paginated list of products
	List(page, limit int) ([]models.Product, int64, error)

	// ListWithCategory returns products with preloaded category
	ListWithCategory(page, limit int) ([]models.Product, int64, error)

	// ListByCategory returns products by category ID
	ListByCategory(categoryID uint, page, limit int) ([]models.Product, int64, error)

	// Search searches products by name or description
	Search(keyword string, page, limit int) ([]models.Product, int64, error)

	// Count returns total number of products
	Count() (int64, error)

	// CountByCategory returns number of products in a category
	CountByCategory(categoryID uint) (int64, error)

	// UpdateStock updates product stock
	UpdateStock(productID uint, stock int) error

	// GetLowStockProducts returns products with stock below threshold
	GetLowStockProducts(threshold int, page, limit int) ([]models.Product, int64, error)
}
