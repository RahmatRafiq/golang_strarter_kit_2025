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

func (s *RoleService) Put(updatedRole models.Role) (models.Role, error) {
	if updatedRole.ID == 0 {
		err := s.roleRepo.Create(&updatedRole)
		if err != nil {
			log.Error().
				Err(err).
				Str("name", updatedRole.Name).
				Msg("Failed to create role")
			return updatedRole, err
		}
		log.Info().
			Uint("role_id", updatedRole.ID).
			Str("name", updatedRole.Name).
			Msg("Role created successfully")
		return updatedRole, nil
	}

	err := s.roleRepo.Update(&updatedRole)
	if err != nil {
		log.Error().
			Err(err).
			Uint("role_id", updatedRole.ID).
			Msg("Failed to update role")
		return updatedRole, err
	}
	role, err := s.roleRepo.FindByID(updatedRole.ID)
	if err != nil {
		log.Error().
			Err(err).
			Uint("role_id", updatedRole.ID).
			Msg("Failed to fetch updated role")
		return updatedRole, err
	}
	log.Info().
		Uint("role_id", updatedRole.ID).
		Str("name", updatedRole.Name).
		Msg("Role updated successfully")
	return *role, nil
}

func (s *RoleService) Create(role *models.Role) error {
	return s.roleRepo.Create(role)
}

func (s *RoleService) Update(role *models.Role) error {
	err := s.roleRepo.Update(role)
	if err != nil {
		return err
	}

	// Invalidate cache after update
	if s.cache.IsEnabled() {
		if err := s.cache.InvalidateRoleCache(role.ID); err != nil {
			log.Warn().Err(err).Uint("role_id", role.ID).Msg("Failed to invalidate role cache")
		}
	}

	return nil
}

func (s *RoleService) DeleteByID(id uint) error {
	err := s.roleRepo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate cache after delete
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

	// Invalidate role permissions cache
	if s.cache.IsEnabled() {
		if err := s.cache.Delete(RolePermissionsCacheKey(roleID)); err != nil {
			log.Warn().Err(err).Uint("role_id", roleID).Msg("Failed to invalidate role permissions cache")
		}
		// Also invalidate all user caches since permissions changed
		if err := s.cache.DeletePattern("user:*:roles"); err != nil {
			log.Warn().Err(err).Msg("Failed to invalidate user roles cache pattern")
		}
	}

	return nil
}

func (s *RoleService) GetPermissions(roleID uint) ([]models.Permission, error) {
	// Try cache first
	cacheKey := RolePermissionsCacheKey(roleID)
	var permissions []models.Permission

	if s.cache.IsEnabled() {
		if err := s.cache.Get(cacheKey, &permissions); err == nil {
			log.Debug().Uint("role_id", roleID).Msg("Role permissions retrieved from cache")
			return permissions, nil
		}
	}

	// Cache miss, get from database
	permissions, err := s.roleRepo.GetPermissions(roleID)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if s.cache.IsEnabled() {
		if err := s.cache.Set(cacheKey, permissions, CacheTTLLong); err != nil {
			log.Warn().Err(err).Uint("role_id", roleID).Msg("Failed to cache role permissions")
		}
	}

	return permissions, nil
}

func (s *RoleService) FindByID(id uint) (*models.Role, error) {
	// Try cache first
	cacheKey := RoleCacheKey(id)
	var role models.Role

	if s.cache.IsEnabled() {
		if err := s.cache.Get(cacheKey, &role); err == nil {
			log.Debug().Uint("role_id", id).Msg("Role retrieved from cache")
			return &role, nil
		}
	}

	// Cache miss, get from database
	rolePtr, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Store in cache
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
