package models

import (
	"time"

	"gorm.io/gorm"
)

type MKPaket struct {
	MKPaketID uint           `gorm:"primaryKey" json:"mk_paket_id"`
	Nama      string         `gorm:"column:nama" json:"nama"`
	ProdiID   string         `gorm:"column:ProdiID" json:"prodi_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
}

func (MKPaket) TableName() string {
	return "mkpaket"
}
