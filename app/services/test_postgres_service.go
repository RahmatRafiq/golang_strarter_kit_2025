package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"
)

type TestService struct{}

func (s *TestService) GetAll() ([]models.Test, error) {
	conn, err := facades.PostgreSQL()
	if err != nil {
		return nil, err
	}
	var tests []models.Test
	if err := conn.DB.Find(&tests).Error; err != nil {
		return nil, err
	}
	return tests, nil
}

func (s *TestService) GetByID(id uint) (*models.Test, error) {
	conn, err := facades.PostgreSQL()
	if err != nil {
		return nil, err
	}
	var test models.Test
	if err := conn.DB.First(&test, id).Error; err != nil {
		return nil, err
	}
	return &test, nil
}

func (s *TestService) Create(test *models.Test) error {
	conn, err := facades.PostgreSQL()
	if err != nil {
		return err
	}
	return conn.DB.Create(test).Error
}

func (s *TestService) Update(test *models.Test) error {
	conn, err := facades.PostgreSQL()
	if err != nil {
		return err
	}
	return conn.DB.Save(test).Error
}

func (s *TestService) Delete(id uint) error {
	conn, err := facades.PostgreSQL()
	if err != nil {
		return err
	}
	return conn.DB.Delete(&models.Test{}, id).Error
}
