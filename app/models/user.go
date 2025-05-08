package models

import (
	"time"

	"golang_starter_kit_2025/app/helpers"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Reference string         `gorm:"type:varchar(100);uniqueIndex" json:"reference"`
	Username  string         `gorm:"type:varchar(100);uniqueIndex" json:"username"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Password  string         `gorm:"type:varchar(255)" json:"password"`
	JwtToken  string         `gorm:"type:varchar(255)" json:"jwt_token" swaggerignore:"true"`
	FcmToken  string         `gorm:"type:varchar(255)" json:"fcm_token" swaggerignore:"true"`
	Pin       string         `gorm:"type:varchar(255)" json:"pin"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" swaggerignore:"true"`

	Roles []Role `gorm:"many2many:user_has_roles;" json:"roles" swaggerignore:"true"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	reference := helpers.GenerateReference("USR")
	password, err := helpers.HashPasswordArgon2(u.Password, helpers.DefaultParams)
	if err != nil {
		println(err.Error())
		return
	}
	// pin, err := helpers.HashPasswordBcrypt(u.Pin)
	pin, err := helpers.HashPasswordArgon2(u.Pin, helpers.DefaultParams)
	if err != nil {
		println(err.Error())
	}
	tx.Statement.SetColumn("reference", reference)
	tx.Statement.SetColumn("password", password)
	tx.Statement.SetColumn("pin", pin)

	return
}
