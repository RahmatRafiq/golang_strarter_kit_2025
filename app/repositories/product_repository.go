package repositories

import (
	"errors"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/models/scopes"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"gorm.io/gorm"
)

// productRepository implements ProductRepositoryInterface
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) interfaces.ProductRepositoryInterface {
	return &productRepository{db: db}
}

// Create creates a new product
func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// Update updates existing product
func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete soft deletes a product by ID
func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// FindByID finds product by ID
func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// FindByIDWithCategory finds product by ID with category preloaded
func (r *productRepository) FindByIDWithCategory(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// List returns paginated list of products
func (r *productRepository) List(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// Count total
	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err := r.db.Scopes(scopes.Paginate(page, limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ListWithCategory returns products with preloaded category
func (r *productRepository) ListWithCategory(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// Count total
	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data with category
	err := r.db.Preload("Category").Scopes(scopes.Paginate(page, limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ListByCategory returns products by category ID
func (r *productRepository) ListByCategory(categoryID uint, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("category_id = ?", categoryID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err := query.Scopes(scopes.Paginate(page, limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Search searches products by name or description
func (r *productRepository) Search(keyword string, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	searchPattern := "%" + keyword + "%"
	query := r.db.Model(&models.Product{}).Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err := query.Preload("Category").Scopes(scopes.Paginate(page, limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Count returns total number of products
func (r *productRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Count(&count).Error
	return count, err
}

// CountByCategory returns number of products in a category
func (r *productRepository) CountByCategory(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

// UpdateStock updates product stock
func (r *productRepository) UpdateStock(productID uint, stock int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", productID).Update("stock", stock).Error
}

// GetLowStockProducts returns products with stock below threshold
func (r *productRepository) GetLowStockProducts(threshold int, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("stock < ?", threshold)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err := query.Preload("Category").Scopes(scopes.Paginate(page, limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
