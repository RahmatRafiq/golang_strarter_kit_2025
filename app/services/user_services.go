package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"strconv"

	"github.com/rs/zerolog/log"
)

type UserService struct {
	repo interfaces.UserRepositoryInterface
}

func NewUserService(repo interfaces.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, _, err := s.repo.List(1, 1000)
	return users, err
}

func (s *UserService) List(page, limit int) ([]models.User, int64, error) {
	return s.repo.List(page, limit)
}

func (s *UserService) Find(id string) (models.User, error) {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return models.User{}, err
	}

	user, err := s.repo.FindByID(uint(userID))
	if err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (s *UserService) FindByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
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
	} else {
		err := s.repo.Update(&user)
		if err != nil {
			log.Error().
				Err(err).
				Uint("user_id", user.ID).
				Msg("Failed to update user")
			return user, err
		}
		log.Info().
			Uint("user_id", user.ID).
			Str("email", user.Email).
			Msg("User updated successfully")
		return user, nil
	}
}

func (s *UserService) Create(user *models.User) error {
	return s.repo.Create(user)
}

func (s *UserService) Update(user *models.User) error {
	return s.repo.Update(user)
}

func (s *UserService) Delete(id string) error {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Error().
			Err(err).
			Str("id", id).
			Msg("Invalid user ID format for deletion")
		return err
	}

	err = s.repo.Delete(uint(userID))
	if err != nil {
		log.Error().
			Err(err).
			Uint("user_id", uint(userID)).
			Msg("Failed to delete user")
		return err
	}

	log.Info().
		Uint("user_id", uint(userID)).
		Msg("User deleted successfully")
	return nil
}

func (s *UserService) DeleteByID(id uint) error {
	return s.repo.Delete(id)
}

func (s *UserService) AssignRolesToUser(userId string, roleIDs []uint) error {
	userID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userId).
			Msg("Invalid user ID format for role assignment")
		return err
	}

	err = s.repo.AssignRoles(uint(userID), roleIDs)
	if err != nil {
		log.Error().
			Err(err).
			Uint("user_id", uint(userID)).
			Interface("role_ids", roleIDs).
			Msg("Failed to assign roles to user")
		return err
	}

	log.Info().
		Uint("user_id", uint(userID)).
		Interface("role_ids", roleIDs).
		Msg("Roles assigned to user successfully")
	return nil
}

func (s *UserService) AssignRoles(userID uint, roleIDs []uint) error {
	return s.repo.AssignRoles(userID, roleIDs)
}

func (s *UserService) GetRolesByUserId(userId string) ([]models.Role, error) {
	userID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return nil, err
	}

	return s.repo.GetRoles(uint(userID))
}

func (s *UserService) GetRoles(userID uint) ([]models.Role, error) {
	return s.repo.GetRoles(userID)
}

func (s *UserService) ExistsByEmail(email string) (bool, error) {
	return s.repo.ExistsByEmail(email)
}

func (s *UserService) Count() (int64, error) {
	return s.repo.Count()
}
