package models

import (
	"time"
)

type Test struct {
	BirthDate   time.Time `json:"birth_date"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LoginTime   string    `json:"login_time"`
	IPAddress   string    `json:"ip_address"`
	DataJSON    string    `json:"data_json"`
	FileBytea   []byte    `json:"file_bytea"`
	ID          uint      `gorm:"primaryKey" json:"id"`
	Age         int       `json:"age"`
	Price       float64   `json:"price"`
	IsActive    bool      `json:"is_active"`
}

func (Test) TableName() string {
	return "test"
}
