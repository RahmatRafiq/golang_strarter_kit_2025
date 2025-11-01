package repositories

import (
	"errors"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/models/scopes"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"gorm.io/gorm"
)

// categoryRepository implements CategoryRepositoryInterface
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) interfaces.CategoryRepositoryInterface {
	return &categoryRepository{db: db}
}

// Create creates a new category
func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// Update updates existing category
func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete soft deletes a category by ID
func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// FindByID finds category by ID
func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// FindByIDWithProducts finds category by ID with products preloaded
func (r *categoryRepository) FindByIDWithProducts(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Products").First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// FindByName finds category by name
func (r *categoryRepository) FindByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// List returns paginated list of categories
func (r *categoryRepository) List(page, limit int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// Count total
	if err := r.db.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err := r.db.Scopes(scopes.PaginateSimple(page, limit)).Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// ListWithProducts returns categories with preloaded products
func (r *categoryRepository) ListWithProducts(page, limit int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// Count total
	if err := r.db.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data with products
	err := r.db.Preload("Products").Scopes(scopes.PaginateSimple(page, limit)).Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// ExistsByName checks if category with name exists
func (r *categoryRepository) ExistsByName(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count returns total number of categories
func (r *categoryRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Count(&count).Error
	return count, err
}

// GetAll returns all categories without pagination
func (r *categoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}
