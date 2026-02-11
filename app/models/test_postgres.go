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
	IPAddress   string    `json:"ip_address"`
	DataJSON    string    `json:"data_json"`
	FileBytea   []byte    `json:"file_bytea"`
}

func (Test) TableName() string {
	return "test"
}
