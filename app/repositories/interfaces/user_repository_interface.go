package interfaces

import "golang_starter_kit_2025/app/models"

// UserRepositoryInterface defines contract for user data operations
type UserRepositoryInterface interface {
	// Create creates a new user
	Create(user *models.User) error

	// Update updates existing user
	Update(user *models.User) error

	// Delete soft deletes a user by ID
	Delete(id uint) error

	// FindByID finds user by ID
	FindByID(id uint) (*models.User, error)

	// FindByEmail finds user by email
	FindByEmail(email string) (*models.User, error)

	// List returns paginated list of users
	List(page, limit int) ([]models.User, int64, error)

	// ListWithRoles returns users with preloaded roles
	ListWithRoles(page, limit int) ([]models.User, int64, error)

	// AssignRoles assigns roles to user
	AssignRoles(userID uint, roleIDs []uint) error

	// GetRoles gets user's roles
	GetRoles(userID uint) ([]models.Role, error)

	// ExistsByEmail checks if user with email exists
	ExistsByEmail(email string) (bool, error)

	// Count returns total number of users
	Count() (int64, error)
}
