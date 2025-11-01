package repositories

import (
	"errors"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/models/scopes"
	"golang_starter_kit_2025/app/repositories/interfaces"

	"gorm.io/gorm"
)

// userRepository implements UserRepositoryInterface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) interfaces.UserRepositoryInterface {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update updates existing user
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user by ID
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// FindByID finds user by ID
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds user by email
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// List returns paginated list of users
func (r *userRepository) List(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	err := r.db.Scopes(scopes.Paginate(page, limit)).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListWithRoles returns users with preloaded roles
func (r *userRepository) ListWithRoles(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data with roles
	err := r.db.Preload("Roles").Scopes(scopes.Paginate(page, limit)).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// AssignRoles assigns roles to user
func (r *userRepository) AssignRoles(userID uint, roleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Find user
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		// Find roles
		var roles []models.Role
		if err := tx.Find(&roles, roleIDs).Error; err != nil {
			return err
		}

		// Check if all roles found
		if len(roles) != len(roleIDs) {
			return errors.New("some roles not found")
		}

		// Clear existing roles and assign new ones
		if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
			return err
		}

		return nil
	})
}

// GetRoles gets user's roles
func (r *userRepository) GetRoles(userID uint) ([]models.Role, error) {
	var user models.User
	if err := r.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Roles, nil
}

// ExistsByEmail checks if user with email exists
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count returns total number of users
func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}
