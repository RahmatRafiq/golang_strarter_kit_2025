package repositories

import (
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/services"

	"gorm.io/gorm"
)

// UserRepository handles user data operations across multiple databases
type UserRepository struct {
	dbService *services.DatabaseService
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		dbService: services.NewDatabaseService(),
	}
}

// CreateOnMySQL creates a user in MySQL database
func (r *UserRepository) CreateOnMySQL(user *models.User) error {
	return r.dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
		return db.Create(user).Error
	})
}

// CreateOnPostgreSQL creates a user in PostgreSQL database
func (r *UserRepository) CreateOnPostgreSQL(user *models.User) error {
	return r.dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
		return db.Create(user).Error
	})
}

// CreateOnBothDatabases creates a user in both MySQL and PostgreSQL
func (r *UserRepository) CreateOnBothDatabases(user *models.User) error {
	return r.dbService.SyncData(func(mysql, postgres *gorm.DB) error {
		// Start transaction on MySQL
		if err := mysql.Create(user).Error; err != nil {
			return err
		}

		// Start transaction on PostgreSQL
		if err := postgres.Create(user).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetFromMySQL retrieves a user from MySQL database
func (r *UserRepository) GetFromMySQL(id uint) (*models.User, error) {
	var user models.User
	err := r.dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
		return db.First(&user, id).Error
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetFromPostgreSQL retrieves a user from PostgreSQL database
func (r *UserRepository) GetFromPostgreSQL(id uint) (*models.User, error) {
	var user models.User
	err := r.dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
		return db.First(&user, id).Error
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllFromMySQL retrieves all users from MySQL database
func (r *UserRepository) GetAllFromMySQL() ([]models.User, error) {
	var users []models.User
	err := r.dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
		return db.Find(&users).Error
	})
	return users, err
}

// GetAllFromPostgreSQL retrieves all users from PostgreSQL database
func (r *UserRepository) GetAllFromPostgreSQL() ([]models.User, error) {
	var users []models.User
	err := r.dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
		return db.Find(&users).Error
	})
	return users, err
}

// UpdateOnMySQL updates a user in MySQL database
func (r *UserRepository) UpdateOnMySQL(user *models.User) error {
	return r.dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
		return db.Save(user).Error
	})
}

// UpdateOnPostgreSQL updates a user in PostgreSQL database
func (r *UserRepository) UpdateOnPostgreSQL(user *models.User) error {
	return r.dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
		return db.Save(user).Error
	})
}

// DeleteFromMySQL deletes a user from MySQL database
func (r *UserRepository) DeleteFromMySQL(id uint) error {
	return r.dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
		return db.Delete(&models.User{}, id).Error
	})
}

// DeleteFromPostgreSQL deletes a user from PostgreSQL database
func (r *UserRepository) DeleteFromPostgreSQL(id uint) error {
	return r.dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
		return db.Delete(&models.User{}, id).Error
	})
}

// SyncUserFromMySQLToPostgreSQL syncs a specific user from MySQL to PostgreSQL
func (r *UserRepository) SyncUserFromMySQLToPostgreSQL(userID uint) error {
	user, err := r.GetFromMySQL(userID)
	if err != nil {
		return err
	}

	return r.dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
		// Check if user exists in PostgreSQL
		var existingUser models.User
		err := db.First(&existingUser, userID).Error

		if err == gorm.ErrRecordNotFound {
			// Create new user
			return db.Create(user).Error
		} else if err != nil {
			return err
		} else {
			// Update existing user
			return db.Save(user).Error
		}
	})
}

// SyncAllUsersFromMySQLToPostgreSQL syncs all users from MySQL to PostgreSQL
func (r *UserRepository) SyncAllUsersFromMySQLToPostgreSQL() error {
	users, err := r.GetAllFromMySQL()
	if err != nil {
		return err
	}

	return r.dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
		for _, user := range users {
			var existingUser models.User
			err := db.First(&existingUser, user.ID).Error

			if err == gorm.ErrRecordNotFound {
				// Create new user
				if err := db.Create(&user).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				// Update existing user
				if err := db.Save(&user).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
