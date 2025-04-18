package models

import (
	"time"

	"gorm.io/gorm"
)

type KHS struct {
	KHSID       string         `gorm:"column:KHSID;primaryKey" json:"khs_id"`
	MahasiswaID string         `gorm:"column:MhswID" json:"mahasiswa_id"`
	TahunID     string         `gorm:"column:TahunID" json:"tahun_id"`
	Sesi        int            `gorm:"column:Sesi" json:"sesi"`
	SKS         int            `gorm:"column:SKS" json:"sks"`
	IPS         float64        `gorm:"column:IPS" json:"ips"`
	Biaya       float64        `gorm:"column:Biaya" json:"biaya"`
	Potongan    float64        `gorm:"column:Potongan" json:"potongan"`
	Bayar       float64        `gorm:"column:Bayar" json:"bayar"`
	Tarik       float64        `gorm:"column:Tarik" json:"tarik"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`

	// Relations
	Mahasiswa Mahasiswa `gorm:"foreignKey:MahasiswaID" json:"mahasiswa"`
}

func (KHS) TableName() string {
	return "khs"
}
