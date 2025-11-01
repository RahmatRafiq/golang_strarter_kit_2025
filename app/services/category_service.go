package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"strconv"
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
		return category, err
	} else {
		err := s.categoryRepo.Update(&category)
		return category, err
	}
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
		return err
	}
	return s.categoryRepo.Delete(uint(categoryID))
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
