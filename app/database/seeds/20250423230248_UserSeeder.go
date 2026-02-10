package seeds

import (
	"golang_starter_kit_2025/app/database/factories"
	"golang_starter_kit_2025/app/models"
	"log"

	"gorm.io/gorm"
)

func SeedUserSeeder(db *gorm.DB) error {
	log.Println("🌱 Seeding UserSeeder...")

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
	log.Println("  ✓ Created admin user")

	// Create 10 random test users using factory
	users, err := userFactory.CreateMany(10)
	if err != nil {
		return err
	}
	log.Printf("  ✓ Created %d test users with random data\n", len(users))

	return nil
}

func RollbackUserSeeder(db *gorm.DB) error {
	log.Println("🗑️ Rolling back UserSeeder...")

	// Delete all users (admin + test users)
	return db.Unscoped().
		Where("username = ? OR username LIKE ?", "admin", "%").
		Delete(&models.User{}).
		Error
}
