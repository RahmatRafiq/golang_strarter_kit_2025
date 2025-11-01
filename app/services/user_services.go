package services

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"strconv"
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
		return user, err
	} else {
		err := s.repo.Update(&user)
		return user, err
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
		return err
	}

	return s.repo.Delete(uint(userID))
}

func (s *UserService) DeleteByID(id uint) error {
	return s.repo.Delete(id)
}

func (s *UserService) AssignRolesToUser(userId string, roleIDs []uint) error {
	userID, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return err
	}

	return s.repo.AssignRoles(uint(userID), roleIDs)
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
