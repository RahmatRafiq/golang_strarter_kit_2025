package models

import (
	"time"
)

type Test struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Age         int       `json:"age"`
	Price       float64   `json:"price"`
	IsActive    bool      `json:"is_active"`
	BirthDate   time.Time `json:"birth_date"`
	LoginTime   string    `json:"login_time"`
	IpAddress   string    `json:"ip_address"`
	DataJson    string    `json:"data_json"`
	FileBytea   []byte    `json:"file_bytea"`
}

func (Test) TableName() string {
	return "test"
}
