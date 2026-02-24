package models

import (
	"time"

	"golang_starter_kit_2025/app/helpers"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type User struct {
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" swaggerignore:"true"`
	Reference       string         `gorm:"type:varchar(100);uniqueIndex" json:"reference"`
	Username        string         `gorm:"type:varchar(100);uniqueIndex" json:"username"`
	Email           string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	Password        string         `gorm:"type:varchar(255)" json:"password"`
	JwtToken        string         `gorm:"type:varchar(255)" json:"jwt_token" swaggerignore:"true"`
	FcmToken        string         `gorm:"type:varchar(255)" json:"fcm_token" swaggerignore:"true"`
	Pin             string         `gorm:"type:varchar(255)" json:"pin"`
	Roles           []Role         `gorm:"many2many:user_has_roles;" json:"roles" swaggerignore:"true"`
	ID              uint           `gorm:"primaryKey" json:"id"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	reference, err := helpers.GenerateReference("USR")
	if err != nil {
		log.Error().Err(err).Str("username", u.Username).Msg("Failed to generate user reference")
		return err
	}

	password, err := helpers.HashPasswordArgon2(u.Password, helpers.DefaultParams)
	if err != nil {
		log.Error().Err(err).Str("username", u.Username).Msg("Failed to hash password in BeforeCreate")
		return err
	}

	// pin, err := helpers.HashPasswordBcrypt(u.Pin)
	pin, err := helpers.HashPasswordArgon2(u.Pin, helpers.DefaultParams)
	if err != nil {
		log.Error().Err(err).Str("username", u.Username).Msg("Failed to hash PIN in BeforeCreate")
		return err
	}

	tx.Statement.SetColumn("reference", reference)
	tx.Statement.SetColumn("password", password)
	tx.Statement.SetColumn("pin", pin)

	return nil
}
