package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Category string `json:"category"`

	Products *[]Product `gorm:"foreignKey:CategoryID" json:"products"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
}
