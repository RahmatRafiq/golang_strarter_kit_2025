package seeds

import (
	"golang_starter_kit_2025/app/database/factories"
	"golang_starter_kit_2025/app/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func SeedUserSeeder(db *gorm.DB) error {
	log.Debug().Msg("Seeding UserSeeder")

	// Initialize user factory
	userFactory := factories.NewUserFactory(db)

	// Create admin user
	_, err := userFactory.Create(map[string]interface{}{
		"username": "admin",
		"email":    "admin@example.com",
	})
	if err != nil {
		return err
	}
	log.Debug().Msg("Created admin user")

	// Create 10 random test users using factory
	users, err := userFactory.CreateMany(10)
	if err != nil {
		return err
	}
	log.Debug().Int("count", len(users)).Msg("Created test users with random data")

	return nil
}

func RollbackUserSeeder(db *gorm.DB) error {
	log.Debug().Msg("Rolling back UserSeeder")

	// Delete all users (admin + test users)
	return db.Unscoped().
		Where("username = ? OR username LIKE ?", "admin", "%").
		Delete(&models.User{}).
		Error
}
