package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"strconv"
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

func (s *PermissionService) Put(updatedPermission models.Permission) (models.Permission, error) {
	if updatedPermission.ID == 0 {
		err := s.repo.Create(&updatedPermission)
		return updatedPermission, err
	} else {
		err := s.repo.Update(&updatedPermission)
		if err != nil {
			return updatedPermission, err
		}
		permission, err := s.repo.FindByID(updatedPermission.ID)
		if err != nil {
			return updatedPermission, err
		}
		return *permission, nil
	}
}

func (s *PermissionService) Create(permission *models.Permission) error {
	return s.repo.Create(permission)
}

func (s *PermissionService) Update(permission *models.Permission) error {
	return s.repo.Update(permission)
}

func (s *PermissionService) Delete(id string) error {
	permissionID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}
	return s.repo.Delete(uint(permissionID))
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
