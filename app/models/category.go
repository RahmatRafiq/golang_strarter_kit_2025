package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Products  *[]Product     `gorm:"foreignKey:CategoryID" json:"products,omitempty" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
	Category  string         `json:"category"`
	ID        uint           `gorm:"primaryKey" json:"id"`
}
