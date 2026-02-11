package services

import (
	"strconv"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"golang_starter_kit_2025/app/requests"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ProductService struct {
	productRepo  interfaces.ProductRepositoryInterface
	categoryRepo interfaces.CategoryRepositoryInterface
	fileService  FileService
}

func NewProductService(productRepo interfaces.ProductRepositoryInterface, categoryRepo interfaces.CategoryRepositoryInterface) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		fileService:  FileService{},
	}
}

func (s *ProductService) GetAll(filters requests.FilterRequest) ([]models.Product, error) {
	products, _, err := s.productRepo.ListWithCategory(1, 1000)
	return products, err
}

func (s *ProductService) List(page, limit int) ([]models.Product, int64, error) {
	return s.productRepo.ListWithCategory(page, limit)
}

func (s *ProductService) GetByID(id string) (models.Product, error) {
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return models.Product{}, err
	}

	product, err := s.productRepo.FindByIDWithCategory(uint(productID))
	if err != nil {
		return models.Product{}, err
	}
	return *product, nil
}

func (s *ProductService) FindByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByIDWithCategory(id)
}

func (s *ProductService) Put(ctx *gin.Context, request requests.ProductRequest) (*models.Product, error) {
	var product models.Product

	var filenames []string
	if request.Images != nil {
		for _, image := range request.Images {
			filename, err := s.fileService.StoreBase64File(image, "images", "products")
			if err != nil {
				return nil, err
			}
			filenames = append(filenames, *filename)
			log.Debug().Str("filename", *filename).Msg("Stored product image")
		}
	}

	if request.ID != 0 {
		product.ID = request.ID
	}
	if request.CategoryID != 0 {
		product.CategoryID = request.CategoryID
	}
	if request.Name != "" {
		product.Name = request.Name
	}
	if request.Description != "" {
		product.Description = request.Description
	}
	if request.Price != 0 {
		product.Price = request.Price
	}
	if request.Margin != 0 {
		product.Margin = request.Margin
	}
	if request.Stock != 0 {
		product.Stock = request.Stock
	}
	if request.Sold != 0 {
		product.Sold = request.Sold
	}
	if !request.ReceivedAt.IsZero() {
		product.ReceivedAt = request.ReceivedAt
	}
	if request.Images != nil {
		product.Images = filenames
	}

	if request.ID == 0 {
		err := s.productRepo.Create(&product)
		if err != nil {
			log.Error().
				Err(err).
				Str("name", product.Name).
				Msg("Failed to create product")
			return &product, err
		}
		log.Info().
			Uint("product_id", product.ID).
			Str("name", product.Name).
			Msg("Product created successfully")
		return &product, nil
	} else {
		err := s.productRepo.Update(&product)
		if err != nil {
			log.Error().
				Err(err).
				Uint("product_id", product.ID).
				Msg("Failed to update product")
			return &product, err
		}
		updated, err := s.productRepo.FindByID(request.ID)
		if err != nil {
			log.Error().
				Err(err).
				Uint("product_id", request.ID).
				Msg("Failed to fetch updated product")
			return &product, err
		}
		log.Info().
			Uint("product_id", product.ID).
			Str("name", product.Name).
			Msg("Product updated successfully")
		return updated, nil
	}
}

func (s *ProductService) Create(product *models.Product) error {
	return s.productRepo.Create(product)
}

func (s *ProductService) Update(product *models.Product) error {
	return s.productRepo.Update(product)
}

func (s *ProductService) Delete(id string) error {
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", id).
			Msg("Invalid product ID format for deletion")
		return err
	}

	err = s.productRepo.Delete(uint(productID))
	if err != nil {
		log.Error().
			Err(err).
			Uint("product_id", uint(productID)).
			Msg("Failed to delete product")
		return err
	}

	log.Info().
		Uint("product_id", uint(productID)).
		Msg("Product deleted successfully")
	return nil
}

func (s *ProductService) DeleteByID(id uint) error {
	return s.productRepo.Delete(id)
}

func (s *ProductService) ListByCategory(categoryID uint, page, limit int) ([]models.Product, int64, error) {
	return s.productRepo.ListByCategory(categoryID, page, limit)
}

func (s *ProductService) Search(keyword string, page, limit int) ([]models.Product, int64, error) {
	return s.productRepo.Search(keyword, page, limit)
}

func (s *ProductService) UpdateStock(productID uint, stock int) error {
	return s.productRepo.UpdateStock(productID, stock)
}

func (s *ProductService) GetLowStockProducts(threshold int, page, limit int) ([]models.Product, int64, error) {
	return s.productRepo.GetLowStockProducts(threshold, page, limit)
}

func (s *ProductService) Count() (int64, error) {
	return s.productRepo.Count()
}

func (s *ProductService) CountByCategory(categoryID uint) (int64, error) {
	return s.productRepo.CountByCategory(categoryID)
}
