package repositories

import (
	"errors"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/models/scopes"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) interfaces.ProductRepositoryInterface {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

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

func (r *productRepository) List(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Scopes(scopes.PaginateSimple(page, limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) ListWithCategory(page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Category").Scopes(scopes.PaginateSimple(page, limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) ListByCategory(categoryID uint, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("category_id = ?", categoryID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(scopes.PaginateSimple(page, limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Search(keyword string, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	searchPattern := "%" + keyword + "%"
	query := r.db.Model(&models.Product{}).Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Category").Scopes(scopes.PaginateSimple(page, limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *productRepository) CountByCategory(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

func (r *productRepository) UpdateStock(productID uint, stock int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", productID).Update("stock", stock).Error
}

func (r *productRepository) GetLowStockProducts(threshold int, page, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("stock < ?", threshold)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Category").Scopes(scopes.PaginateSimple(page, limit)).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
