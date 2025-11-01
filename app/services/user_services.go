package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"strconv"
)

// UserService handles user business logic
type UserService struct {
	repo interfaces.UserRepositoryInterface
}

// NewUserService creates a new user service
func NewUserService(repo interfaces.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

// GetAllUsers returns all users (consider deprecating in favor of List with pagination)
func (s *UserService) GetAllUsers() ([]models.User, error) {
	// For backward compatibility, return first 1000 users
	users, _, err := s.repo.List(1, 1000)
	return users, err
}

// List returns paginated list of users
func (s *UserService) List(page, limit int) ([]models.User, int64, error) {
	return s.repo.List(page, limit)
}

// Find finds user by ID (string for backward compatibility)
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

// FindByID finds user by ID (uint version)
func (s *UserService) FindByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

// FindByEmail finds user by email
func (s *UserService) FindByEmail(email string) (*models.User, error) {
	return s.repo.FindByEmail(email)
}

// Put creates or updates a user (based on ID)
func (s *UserService) Put(user models.User) (models.User, error) {
	if user.ID == 0 {
		// Create new user
		err := s.repo.Create(&user)
		return user, err
	} else {
		// Update existing user
		err := s.repo.Update(&user)
		return user, err
	}
}

// Create creates a new user
func (s *UserService) Create(user *models.User) error {
	return s.repo.Create(user)
}

// Update updates existing user
func (s *UserService) Update(user *models.User) error {
	return s.repo.Update(user)
}

// Delete deletes user by ID (string for backward compatibility)
func (s *UserService) Delete(id string) error {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	return s.repo.Delete(uint(userID))
}

// DeleteByID deletes user by ID (uint version)
func (s *UserService) DeleteByID(id uint) error {
	return s.repo.Delete(id)
}

// AssignRolesToUser assigns roles to a user
func (s *UserService) AssignRolesToUser(userId string, roleIDs []uint) error {
	userID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return err
	}

	return s.repo.AssignRoles(uint(userID), roleIDs)
}

// AssignRoles assigns roles to a user (uint version)
func (s *UserService) AssignRoles(userID uint, roleIDs []uint) error {
	return s.repo.AssignRoles(userID, roleIDs)
}

// GetRolesByUserId gets roles for a user
func (s *UserService) GetRolesByUserId(userId string) ([]models.Role, error) {
	userID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return nil, err
	}

	return s.repo.GetRoles(uint(userID))
}

// GetRoles gets roles for a user (uint version)
func (s *UserService) GetRoles(userID uint) ([]models.Role, error) {
	return s.repo.GetRoles(userID)
}

// ExistsByEmail checks if user with email exists
func (s *UserService) ExistsByEmail(email string) (bool, error) {
	return s.repo.ExistsByEmail(email)
}

// Count returns total number of users
func (s *UserService) Count() (int64, error) {
	return s.repo.Count()
}
