package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"github.com/rs/zerolog/log"
)

type RoleService struct {
	roleRepo       interfaces.RoleRepositoryInterface
	permissionRepo interfaces.PermissionRepositoryInterface
	cache          *CacheService
}

func NewRoleService(roleRepo interfaces.RoleRepositoryInterface, permissionRepo interfaces.PermissionRepositoryInterface) *RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		cache:          NewCacheService(),
	}
}

func (s *RoleService) GetAll() ([]models.Role, error) {
	return s.roleRepo.GetAll()
}

func (s *RoleService) List(page, limit int) ([]models.Role, int64, error) {
	return s.roleRepo.List(page, limit)
}

func (s *RoleService) Create(role models.Role) (models.Role, error) {
	role.ID = 0

	err := s.roleRepo.Create(&role)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", role.Name).
			Msg("Failed to create role")
		return role, err
	}

	log.Info().
		Uint("role_id", role.ID).
		Str("name", role.Name).
		Msg("Role created successfully")
	return role, nil
}

func (s *RoleService) Update(id uint, role models.Role) (models.Role, error) {
	role.ID = id

	err := s.roleRepo.Update(&role)
	if err != nil {
		log.Error().
			Err(err).
			Uint("role_id", role.ID).
			Msg("Failed to update role")
		return role, err
	}

	updatedRole, err := s.roleRepo.FindByID(role.ID)
	if err != nil {
		log.Error().
			Err(err).
			Uint("role_id", role.ID).
			Msg("Failed to fetch updated role")
		return role, err
	}

	if s.cache.IsEnabled() {
		if err := s.cache.InvalidateRoleCache(role.ID); err != nil {
			log.Warn().Err(err).Uint("role_id", role.ID).Msg("Failed to invalidate role cache")
		}
	}

	log.Info().
		Uint("role_id", role.ID).
		Str("name", role.Name).
		Msg("Role updated successfully")
	return *updatedRole, nil
}

func (s *RoleService) DeleteByID(id uint) error {
	err := s.roleRepo.Delete(id)
	if err != nil {
		return err
	}

	if s.cache.IsEnabled() {
		if err := s.cache.InvalidateRoleCache(id); err != nil {
			log.Warn().Err(err).Uint("role_id", id).Msg("Failed to invalidate role cache")
		}
	}

	return nil
}

func (s *RoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	err := s.roleRepo.AssignPermissions(roleID, permissionIDs)
	if err != nil {
		return err
	}

	if s.cache.IsEnabled() {
		if err := s.cache.Delete(RolePermissionsCacheKey(roleID)); err != nil {
			log.Warn().Err(err).Uint("role_id", roleID).Msg("Failed to invalidate role permissions cache")
		}
		if err := s.cache.DeletePattern("user:*:roles"); err != nil {
			log.Warn().Err(err).Msg("Failed to invalidate user roles cache pattern")
		}
	}

	return nil
}

func (s *RoleService) GetPermissions(roleID uint) ([]models.Permission, error) {
	cacheKey := RolePermissionsCacheKey(roleID)
	var permissions []models.Permission

	if s.cache.IsEnabled() {
		if err := s.cache.Get(cacheKey, &permissions); err == nil {
			log.Debug().Uint("role_id", roleID).Msg("Role permissions retrieved from cache")
			return permissions, nil
		}
	}

	permissions, err := s.roleRepo.GetPermissions(roleID)
	if err != nil {
		return nil, err
	}

	if s.cache.IsEnabled() {
		if err := s.cache.Set(cacheKey, permissions, CacheTTLLong); err != nil {
			log.Warn().Err(err).Uint("role_id", roleID).Msg("Failed to cache role permissions")
		}
	}

	return permissions, nil
}

func (s *RoleService) FindByID(id uint) (*models.Role, error) {
	cacheKey := RoleCacheKey(id)
	var role models.Role

	if s.cache.IsEnabled() {
		if err := s.cache.Get(cacheKey, &role); err == nil {
			log.Debug().Uint("role_id", id).Msg("Role retrieved from cache")
			return &role, nil
		}
	}

	rolePtr, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if s.cache.IsEnabled() {
		if err := s.cache.Set(cacheKey, rolePtr, CacheTTLLong); err != nil {
			log.Warn().Err(err).Uint("role_id", id).Msg("Failed to cache role")
		}
	}

	return rolePtr, nil
}

func (s *RoleService) FindByName(name string) (*models.Role, error) {
	return s.roleRepo.FindByName(name)
}

func (s *RoleService) ExistsByName(name string) (bool, error) {
	return s.roleRepo.ExistsByName(name)
}

func (s *RoleService) Count() (int64, error) {
	return s.roleRepo.Count()
}
