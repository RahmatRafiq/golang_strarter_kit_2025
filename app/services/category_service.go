package services

import (
	"strconv"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"github.com/rs/zerolog/log"
)

type CategoryService struct {
	categoryRepo interfaces.CategoryRepositoryInterface
	productRepo  interfaces.ProductRepositoryInterface
}

func NewCategoryService(categoryRepo interfaces.CategoryRepositoryInterface, productRepo interfaces.ProductRepositoryInterface) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *CategoryService) List(page, limit int) ([]models.Category, int64, error) {
	return s.categoryRepo.List(page, limit)
}

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

func (s *CategoryService) FindByID(id uint) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *CategoryService) FindByIDWithProducts(id uint) (*models.Category, error) {
	return s.categoryRepo.FindByIDWithProducts(id)
}

func (s *CategoryService) PutCategory(category models.Category) (models.Category, error) {
	if category.ID == 0 {
		err := s.categoryRepo.Create(&category)
		if err != nil {
			log.Error().
				Err(err).
				Str("category", category.Category).
				Msg("Failed to create category")
			return category, err
		}
		log.Info().
			Uint("category_id", category.ID).
			Str("category", category.Category).
			Msg("Category created successfully")
		return category, nil
	}

	err := s.categoryRepo.Update(&category)
	if err != nil {
		log.Error().
			Err(err).
			Uint("category_id", category.ID).
			Msg("Failed to update category")
		return category, err
	}
	log.Info().
		Uint("category_id", category.ID).
		Str("category", category.Category).
		Msg("Category updated successfully")
	return category, nil
}

func (s *CategoryService) Create(category *models.Category) error {
	return s.categoryRepo.Create(category)
}

func (s *CategoryService) Update(category *models.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *CategoryService) DeleteCategory(id string) error {
	categoryID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", id).
			Msg("Invalid category ID format for deletion")
		return err
	}

	err = s.categoryRepo.Delete(uint(categoryID))
	if err != nil {
		log.Error().
			Err(err).
			Uint("category_id", uint(categoryID)).
			Msg("Failed to delete category")
		return err
	}

	log.Info().
		Uint("category_id", uint(categoryID)).
		Msg("Category deleted successfully")
	return nil
}

func (s *CategoryService) DeleteByID(id uint) error {
	return s.categoryRepo.Delete(id)
}

func (s *CategoryService) FindByName(name string) (*models.Category, error) {
	return s.categoryRepo.FindByName(name)
}

func (s *CategoryService) ExistsByName(name string) (bool, error) {
	return s.categoryRepo.ExistsByName(name)
}

func (s *CategoryService) Count() (int64, error) {
	return s.categoryRepo.Count()
}

func (s *CategoryService) ListWithProducts(page, limit int) ([]models.Category, int64, error) {
	return s.categoryRepo.ListWithProducts(page, limit)
}
