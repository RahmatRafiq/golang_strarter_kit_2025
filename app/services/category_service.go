package services

import (
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

func (s *CategoryService) GetCategoryByIDUint(id uint) (models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
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

func (s *CategoryService) DeleteCategoryByID(id uint) error {
	err := s.categoryRepo.Delete(id)
	if err != nil {
		log.Error().
			Err(err).
			Uint("category_id", id).
			Msg("Failed to delete category")
		return err
	}

	log.Info().
		Uint("category_id", id).
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
