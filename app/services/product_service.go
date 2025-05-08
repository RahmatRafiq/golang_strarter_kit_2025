package services

import (
	"log"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/facades"

	"github.com/gin-gonic/gin"
)

type ProductService struct {
	fileService FileService
}

func NewProductService() *ProductService {
	return &ProductService{
		fileService: FileService{},
	}
}

func (service *ProductService) GetAll(filters requests.FilterRequest) ([]models.Product, error) {
	var products []models.Product
	query := facades.DB

	if filters.Search != nil {
		query = query.Where("name LIKE ?", "%"+*filters.Search+"%").
			Or("description LIKE ?", "%"+*filters.Search+"%").
			Or("reference LIKE ?", "%"+*filters.Search+"%")
	}

	if filters.OrderBy != nil {
		query = query.Order(*filters.OrderBy + " " + *filters.OrderDirection)
	} else {
		query = query.Order("updated_at desc")
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (service *ProductService) GetByID(id string) (models.Product, error) {
	var product models.Product
	if err := facades.DB.First(&product, id).Error; err != nil {
		return product, err
	}
	return product, nil
}

func (service *ProductService) Put(ctx *gin.Context, request requests.ProductRequest) (*models.Product, error) {
	var product models.Product

	// Mengunggah file gambar produk jika ada
	var filenames []string
	if request.Images != nil {
		for _, image := range request.Images {
			filename, err := service.fileService.StoreBase64File(image, "images", "products")
			if err != nil {
				return nil, err
			}
			filenames = append(filenames, *filename)
			log.Println(&filename)
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

	if count := facades.DB.Model(&models.Product{}).Where("id = ?", request.ID).Find(&map[string]interface{}{}).RowsAffected; count == 0 {
		if err := facades.DB.Create(&product).Error; err != nil {
			return &product, err
		}
	} else {
		if err := facades.DB.Model(&models.Product{}).Where("id = ?", request.ID).Updates(&product).Error; err != nil {
			return &product, err
		}
		if err := facades.DB.First(&product, request.ID).Error; err != nil {
			return &product, err
		}
	}

	return &product, nil
}

func (service *ProductService) Delete(id string) error {
	result := facades.DB.Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}
