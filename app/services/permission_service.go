package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
)

type PermissionService struct {
	repo interfaces.PermissionRepositoryInterface
}

func NewPermissionService(repo interfaces.PermissionRepositoryInterface) *PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) GetAll() ([]models.Permission, error) {
	return s.repo.GetAll()
}

func (s *PermissionService) List(page, limit int) ([]models.Permission, int64, error) {
	return s.repo.List(page, limit)
}

func (s *PermissionService) Create(permission models.Permission) (models.Permission, error) {
	permission.ID = 0

	err := s.repo.Create(&permission)
	if err != nil {
		return permission, err
	}

	return permission, nil
}

func (s *PermissionService) Update(id uint, permission models.Permission) (models.Permission, error) {
	permission.ID = id

	err := s.repo.Update(&permission)
	if err != nil {
		return permission, err
	}

	updatedPermission, err := s.repo.FindByID(permission.ID)
	if err != nil {
		return permission, err
	}

	return *updatedPermission, nil
}

func (s *PermissionService) DeleteByID(id uint) error {
	return s.repo.Delete(id)
}

func (s *PermissionService) FindByID(id uint) (*models.Permission, error) {
	return s.repo.FindByID(id)
}

func (s *PermissionService) FindByName(name string) (*models.Permission, error) {
	return s.repo.FindByName(name)
}

func (s *PermissionService) FindByIDs(ids []uint) ([]models.Permission, error) {
	return s.repo.FindByIDs(ids)
}

func (s *PermissionService) ExistsByName(name string) (bool, error) {
	return s.repo.ExistsByName(name)
}

func (s *PermissionService) Count() (int64, error) {
	return s.repo.Count()
}
