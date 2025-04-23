package seeds

import (
	"golang_strarter_kit_2025/app/helpers"
	"golang_strarter_kit_2025/app/models"
	"log"
	"time"

	"gorm.io/gorm"
)

func SeedUserSeeder(db *gorm.DB) error {
	log.Println("🌱 Seeding UserSeeder...")

	password, err := helpers.HashPasswordArgon2("password123", helpers.DefaultParams)
	if err != nil {
		return err
	}
	pin, err := helpers.HashPasswordArgon2("123456", helpers.DefaultParams)
	if err != nil {
		return err
	}

	data := models.User{
		Reference: helpers.GenerateReference("USR"),
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  password,
		Pin:       pin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}
func RollbackUserSeeder(db *gorm.DB) error {
	log.Println("🗑️ Rolling back UserSeeder…")
	return db.Unscoped().
		Where("username = ?", "admin").
		Delete(&models.User{}).
		Error
}
