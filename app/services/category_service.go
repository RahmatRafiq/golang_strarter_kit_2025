package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"

	"gorm.io/gorm/clause"
)

type CategoryService struct{}

func (*CategoryService) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := facades.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (*CategoryService) GetCategoryByID(id string) (models.Category, error) {
	var category models.Category
	if err := facades.DB.First(&category, id).Error; err != nil {
		return category, err
	}
	return category, nil
}

// Menggabungkan Create dan Update dalam satu fungsi PutCategory
func (*CategoryService) PutCategory(category models.Category) (models.Category, error) {
	if err := facades.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // Kolom yang digunakan untuk menentukan konflik
		DoUpdates: clause.AssignmentColumns([]string{"category", "updated_at"}),
	}).Create(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (*CategoryService) DeleteCategory(id string) error {
	var category models.Category
	if err := facades.DB.First(&category, id).Error; err != nil {
		return err
	}
	return facades.DB.Delete(&category).Error
}
