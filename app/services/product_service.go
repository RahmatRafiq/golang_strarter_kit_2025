package services

import (
	"log"
	"strconv"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"golang_starter_kit_2025/app/requests"

	"github.com/gin-gonic/gin"
)

// ProductService handles product business logic
type ProductService struct {
	productRepo  interfaces.ProductRepositoryInterface
	categoryRepo interfaces.CategoryRepositoryInterface
	fileService  FileService
}

// NewProductService creates a new product service
func NewProductService(productRepo interfaces.ProductRepositoryInterface, categoryRepo interfaces.CategoryRepositoryInterface) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		fileService:  FileService{},
	}
}

// GetAll returns filtered list of products
func (s *ProductService) GetAll(filters requests.FilterRequest) ([]models.Product, error) {
	// For backward compatibility, use repository List method
	// Note: Advanced filtering should be added to repository
	products, _, err := s.productRepo.ListWithCategory(1, 1000)
	return products, err
}

// List returns paginated list of products
func (s *ProductService) List(page, limit int) ([]models.Product, int64, error) {
	return s.productRepo.ListWithCategory(page, limit)
}

// GetByID finds product by ID (string for backward compatibility)
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

// FindByID finds product by ID (uint version)
func (s *ProductService) FindByID(id uint) (*models.Product, error) {
	return s.productRepo.FindByIDWithCategory(id)
}

// Put creates or updates a product with file handling
func (s *ProductService) Put(ctx *gin.Context, request requests.ProductRequest) (*models.Product, error) {
	var product models.Product

	// Upload product images if provided
	var filenames []string
	if request.Images != nil {
		for _, image := range request.Images {
			filename, err := s.fileService.StoreBase64File(image, "images", "products")
			if err != nil {
				return nil, err
			}
			filenames = append(filenames, *filename)
			log.Println(&filename)
		}
	}

	// Build product model from request
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

	// Create or update using repository
	if request.ID == 0 {
		// Create new product
		err := s.productRepo.Create(&product)
		return &product, err
	} else {
		// Update existing product
		err := s.productRepo.Update(&product)
		if err != nil {
			return &product, err
		}
		// Reload product
		updated, err := s.productRepo.FindByID(request.ID)
		if err != nil {
			return &product, err
		}
		return updated, nil
	}
}

// Create creates a new product
func (s *ProductService) Create(product *models.Product) error {
	return s.productRepo.Create(product)
}

// Update updates existing product
func (s *ProductService) Update(product *models.Product) error {
	return s.productRepo.Update(product)
}

// Delete deletes product by ID (string for backward compatibility)
func (s *ProductService) Delete(id string) error {
	productID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}
	return s.productRepo.Delete(uint(productID))
}

// DeleteByID deletes product by ID (uint version)
func (s *ProductService) DeleteByID(id uint) error {
	return s.productRepo.Delete(id)
}

// ListByCategory returns products by category ID
func (s *ProductService) ListByCategory(categoryID uint, page, limit int) ([]models.Product, int64, error) {
	return s.productRepo.ListByCategory(categoryID, page, limit)
}

// Search searches products by keyword
func (s *ProductService) Search(keyword string, page, limit int) ([]models.Product, int64, error) {
	return s.productRepo.Search(keyword, page, limit)
}

// UpdateStock updates product stock
func (s *ProductService) UpdateStock(productID uint, stock int) error {
	return s.productRepo.UpdateStock(productID, stock)
}

// GetLowStockProducts returns products with low stock
func (s *ProductService) GetLowStockProducts(threshold int, page, limit int) ([]models.Product, int64, error) {
	return s.productRepo.GetLowStockProducts(threshold, page, limit)
}

// Count returns total number of products
func (s *ProductService) Count() (int64, error) {
	return s.productRepo.Count()
}

// CountByCategory returns number of products in a category
func (s *ProductService) CountByCategory(categoryID uint) (int64, error) {
	return s.productRepo.CountByCategory(categoryID)
}
