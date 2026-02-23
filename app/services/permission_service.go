package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"github.com/rs/zerolog/log"
)

type PermissionService struct {
	repo interfaces.PermissionRepositoryInterface
}

func NewPermissionService(repo interfaces.PermissionRepositoryInterface) *PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) List(page, limit int) ([]models.Permission, int64, error) {
	return s.repo.List(page, limit)
}

func (s *PermissionService) Create(permission models.Permission) (models.Permission, error) {
	permission.ID = 0

	err := s.repo.Create(&permission)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", permission.Name).
			Msg("Failed to create permission")
		return permission, err
	}

	log.Info().
		Uint("permission_id", permission.ID).
		Str("name", permission.Name).
		Msg("Permission created successfully")
	return permission, nil
}

func (s *PermissionService) Update(id uint, permission models.Permission) (models.Permission, error) {
	permission.ID = id

	err := s.repo.Update(&permission)
	if err != nil {
		log.Error().
			Err(err).
			Uint("permission_id", permission.ID).
			Msg("Failed to update permission")
		return permission, err
	}

	updatedPermission, err := s.repo.FindByID(permission.ID)
	if err != nil {
		log.Error().
			Err(err).
			Uint("permission_id", permission.ID).
			Msg("Failed to fetch updated permission")
		return permission, err
	}

	log.Info().
		Uint("permission_id", permission.ID).
		Str("name", permission.Name).
		Msg("Permission updated successfully")
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
