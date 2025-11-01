package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"strconv"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo interfaces.CategoryRepositoryInterface
	productRepo  interfaces.ProductRepositoryInterface
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo interfaces.CategoryRepositoryInterface, productRepo interfaces.ProductRepositoryInterface) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
}

// GetAllCategories returns all categories without pagination
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

// List returns paginated list of categories
func (s *CategoryService) List(page, limit int) ([]models.Category, int64, error) {
	return s.categoryRepo.List(page, limit)
}

// GetCategoryByID finds category by ID (string for backward compatibility)
func (s *CategoryService) GetCategoryByID(id string) (models.Category, error) {
	categoryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return models.Category{}, err
	}

	category, err := s.categoryRepo.FindByID(uint(categoryID))
	if err != nil {
		return models.Category{}, err
	}
	return *category, nil
}

// FindByID finds category by ID (uint version)
func (s *CategoryService) FindByID(id uint) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

// FindByIDWithProducts finds category by ID with products preloaded
func (s *CategoryService) FindByIDWithProducts(id uint) (*models.Category, error) {
	return s.categoryRepo.FindByIDWithProducts(id)
}

// PutCategory creates or updates a category (based on ID)
func (s *CategoryService) PutCategory(category models.Category) (models.Category, error) {
	if category.ID == 0 {
		// Create new category
		err := s.categoryRepo.Create(&category)
		return category, err
	} else {
		// Update existing category
		err := s.categoryRepo.Update(&category)
		return category, err
	}
}

// Create creates a new category
func (s *CategoryService) Create(category *models.Category) error {
	return s.categoryRepo.Create(category)
}

// Update updates existing category
func (s *CategoryService) Update(category *models.Category) error {
	return s.categoryRepo.Update(category)
}

// DeleteCategory deletes category by ID (string for backward compatibility)
func (s *CategoryService) DeleteCategory(id string) error {
	categoryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}
	return s.categoryRepo.Delete(uint(categoryID))
}

// DeleteByID deletes category by ID (uint version)
func (s *CategoryService) DeleteByID(id uint) error {
	return s.categoryRepo.Delete(id)
}

// FindByName finds category by name
func (s *CategoryService) FindByName(name string) (*models.Category, error) {
	return s.categoryRepo.FindByName(name)
}

// ExistsByName checks if category with name exists
func (s *CategoryService) ExistsByName(name string) (bool, error) {
	return s.categoryRepo.ExistsByName(name)
}

// Count returns total number of categories
func (s *CategoryService) Count() (int64, error) {
	return s.categoryRepo.Count()
}

// ListWithProducts returns categories with preloaded products
func (s *CategoryService) ListWithProducts(page, limit int) ([]models.Category, int64, error) {
	return s.categoryRepo.ListWithProducts(page, limit)
}
