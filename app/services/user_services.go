package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"github.com/rs/zerolog/log"
)

type UserService struct {
	repo  interfaces.UserRepositoryInterface
	cache *CacheService
}

func NewUserService(repo interfaces.UserRepositoryInterface) *UserService {
	return &UserService{
		repo:  repo,
		cache: NewCacheService(),
	}
}

// GetAllUsers is DEPRECATED and REMOVED due to unsafe hardcoded limit.
// Use List(page, limit) with proper pagination instead.
//
// Reason for removal:
// - Hardcoded limit 1000 causes memory issues with large datasets
// - Misleading name ("All" but limited to 1000)
// - Not suitable for production use
//
// Migration guide:
//   Old: users, err := service.GetAllUsers()
//   New: users, total, err := service.List(1, 100) // with proper pagination
//
// REMOVED: func (s *UserService) GetAllUsers() ([]models.User, error)

func (s *UserService) List(page, limit int) ([]models.User, int64, error) {
	return s.repo.List(page, limit)
}

func (s *UserService) FindByID(id uint) (*models.User, error) {
	// Try cache first
	cacheKey := UserCacheKey(id)
	var user models.User

	if s.cache.IsEnabled() {
		if err := s.cache.Get(cacheKey, &user); err == nil {
			log.Debug().Uint("user_id", id).Msg("User retrieved from cache")
			return &user, nil
		}
	}

	// Cache miss, get from database
	userPtr, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Store in cache (log errors but don't fail)
	if s.cache.IsEnabled() {
		if err := s.cache.Set(cacheKey, userPtr, CacheTTLMedium); err != nil {
			log.Warn().Err(err).Uint("user_id", id).Msg("Failed to cache user")
		}
	}

	return userPtr, nil
}

func (s *UserService) FindByEmail(email string) (*models.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *UserService) Put(user models.User) (models.User, error) {
	if user.ID == 0 {
		err := s.repo.Create(&user)
		if err != nil {
			log.Error().
				Err(err).
				Str("email", user.Email).
				Msg("Failed to create user")
			return user, err
		}
		log.Info().
			Uint("user_id", user.ID).
			Str("email", user.Email).
			Msg("User created successfully")
		return user, nil
	}

	err := s.repo.Update(&user)
	if err != nil {
		log.Error().
			Err(err).
			Uint("user_id", user.ID).
			Msg("Failed to update user")
		return user, err
	}

	// Invalidate cache after update
	if s.cache.IsEnabled() {
		if err := s.cache.InvalidateUserCache(user.ID); err != nil {
			log.Warn().Err(err).Uint("user_id", user.ID).Msg("Failed to invalidate user cache")
		}
	}

	log.Info().
		Uint("user_id", user.ID).
		Str("email", user.Email).
		Msg("User updated successfully")
	return user, nil
}

func (s *UserService) Create(user *models.User) error {
	return s.repo.Create(user)
}

func (s *UserService) Update(user *models.User) error {
	err := s.repo.Update(user)
	if err != nil {
		return err
	}

	// Invalidate cache after update
	if s.cache.IsEnabled() {
		if err := s.cache.InvalidateUserCache(user.ID); err != nil {
			log.Warn().Err(err).Uint("user_id", user.ID).Msg("Failed to invalidate user cache")
		}
	}

	return nil
}

func (s *UserService) DeleteByID(id uint) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	// Invalidate cache after delete
	if s.cache.IsEnabled() {
		if err := s.cache.InvalidateUserCache(id); err != nil {
			log.Warn().Err(err).Uint("user_id", id).Msg("Failed to invalidate user cache")
		}
	}

	return nil
}

func (s *UserService) AssignRoles(userID uint, roleIDs []uint) error {
	err := s.repo.AssignRoles(userID, roleIDs)
	if err != nil {
		return err
	}

	// Invalidate user roles cache
	if s.cache.IsEnabled() {
		if err := s.cache.Delete(UserRolesCacheKey(userID)); err != nil {
			log.Warn().Err(err).Uint("user_id", userID).Msg("Failed to invalidate user roles cache")
		}
	}

	return nil
}

func (s *UserService) GetRoles(userID uint) ([]models.Role, error) {
	// Try cache first
	cacheKey := UserRolesCacheKey(userID)
	var roles []models.Role

	if s.cache.IsEnabled() {
		if err := s.cache.Get(cacheKey, &roles); err == nil {
			log.Debug().Uint("user_id", userID).Msg("User roles retrieved from cache")
			return roles, nil
		}
	}

	// Cache miss, get from database
	roles, err := s.repo.GetRoles(userID)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if s.cache.IsEnabled() {
		if err := s.cache.Set(cacheKey, roles, CacheTTLMedium); err != nil {
			log.Warn().Err(err).Uint("user_id", userID).Msg("Failed to cache user roles")
		}
	}

	return roles, nil
}

func (s *UserService) ExistsByEmail(email string) (bool, error) {
	return s.repo.ExistsByEmail(email)
}

func (s *UserService) Count() (int64, error) {
	return s.repo.Count()
}
