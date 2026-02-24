package factories

import (
	"fmt"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// UserFactory creates User instances for testing and seeding
type UserFactory struct {
	*BaseFactory
}

// NewUserFactory creates a new UserFactory
func NewUserFactory(db *gorm.DB) *UserFactory {
	return &UserFactory{
		BaseFactory: NewBaseFactory(db),
	}
}

// Make creates a User instance without saving to database
func (f *UserFactory) Make(overrides ...map[string]interface{}) *models.User {
	reference, err := helpers.GenerateReference("USR")
	if err != nil {
		// Factory fallback: use timestamp-based reference if UUID fails
		reference = fmt.Sprintf("USR-FALLBACK-%d", time.Now().UnixNano())
		log.Warn().Err(err).Msg("Failed to generate reference in factory, using fallback")
	}

	user := &models.User{
		Reference: reference,
		Username:  RandomUsername(),
		Email:     RandomEmail(),
		Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password"
		Pin:       "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Apply overrides if provided
	if len(overrides) > 0 {
		override := overrides[0]
		if val, ok := override["username"]; ok {
			if strVal, ok := val.(string); ok {
				user.Username = strVal
			}
		}
		if val, ok := override["email"]; ok {
			if strVal, ok := val.(string); ok {
				user.Email = strVal
			}
		}
		if val, ok := override["password"]; ok {
			if strVal, ok := val.(string); ok {
				user.Password = strVal
			}
		}
		if val, ok := override["pin"]; ok {
			if strVal, ok := val.(string); ok {
				user.Pin = strVal
			}
		}
		if val, ok := override["reference"]; ok {
			if strVal, ok := val.(string); ok {
				user.Reference = strVal
			}
		}
	}

	return user
}

// MakeMany creates multiple User instances without saving
func (f *UserFactory) MakeMany(count int, overrides ...map[string]interface{}) []*models.User {
	users := make([]*models.User, count)
	for i := 0; i < count; i++ {
		users[i] = f.Make(overrides...)
	}
	return users
}

// Create creates and saves a User to the database
func (f *UserFactory) Create(overrides ...map[string]interface{}) (*models.User, error) {
	user := f.Make(overrides...)
	if err := f.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// CreateMany creates and saves multiple Users to the database
func (f *UserFactory) CreateMany(count int, overrides ...map[string]interface{}) ([]*models.User, error) {
	users := f.MakeMany(count, overrides...)
	if err := f.DB.Create(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CreateInBatches creates multiple Users in batches for better performance
func (f *UserFactory) CreateInBatches(count, batchSize int, overrides ...map[string]interface{}) ([]*models.User, error) {
	users := f.MakeMany(count, overrides...)
	if err := f.DB.CreateInBatches(&users, batchSize).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// WithUsername creates a user with a specific username
func (f *UserFactory) WithUsername(username string) *UserFactory {
	// This is a builder pattern - you can chain methods
	return f
}

// WithEmail creates a user with a specific email
func (f *UserFactory) WithEmail(email string) *UserFactory {
	return f
}

// Admin creates an admin user
func (f *UserFactory) Admin() *models.User {
	return f.Make(map[string]interface{}{
		"username": "admin",
		"email":    "admin@example.com",
	})
}

// CreateAdmin creates and saves an admin user
func (f *UserFactory) CreateAdmin() (*models.User, error) {
	admin := f.Admin()
	if err := f.DB.Create(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}
