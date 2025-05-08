package models

import (
	"time"

	"golang_starter_kit_2025/app/helpers"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Reference   string    `gorm:"unique" json:"reference"`
	StoreID     uint      `json:"store_id"`
	CategoryID  uint      `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Margin      float64   `json:"margin"`
	Stock       int       `json:"stock"`
	Sold        int       `json:"sold"`
	Images      []string  `json:"images" gorm:"serializer:json"`
	ReceivedAt  time.Time `json:"received_at"`

	// Store    *Store    `json:"store"`
	Category *Category `json:"category"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
}

// BeforeCreate hook
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.Reference = helpers.GenerateReference("PRD")
	return
}

// AfterFind hook
func (m *Product) AfterFind(tx *gorm.DB) (err error) {
	if m.Images != nil && len(m.Images) > 0 {
		for i, Image := range m.Images {
			m.Images[i] = helpers.GetFileURL(Image, "member_lands")
		}
	}

	return
}

// AfterCreate hook
func (m *Product) AfterCreate(tx *gorm.DB) (err error) {
	if m.Images != nil && len(m.Images) > 0 {
		for i, image := range m.Images {
			m.Images[i] = helpers.GetFileURL(image, "member_lands")
		}
	}

	return
}

// AfterUpdate hook
func (m *Product) AfterUpdate(tx *gorm.DB) (err error) {
	if m.Images != nil && len(m.Images) > 0 {
		for i, image := range m.Images {
			m.Images[i] = helpers.GetFileURL(image, "member_lands")
		}
	}

	return
}
