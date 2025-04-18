package models

import (
	"time"

	"gorm.io/gorm"
)

type Jadwal struct {
	JadwalID   uint           `gorm:"primaryKey" json:"jadwal_id"`
	MKID       uint           `gorm:"column:MKID" json:"mk_id"`
	HariID     uint           `gorm:"column:HariID" json:"hari_id"`
	JamMulai   string         `gorm:"column:JamMulai" json:"jam_mulai"`
	JamSelesai string         `gorm:"column:JamSelesai" json:"jam_selesai"`
	Kapasitas  int            `gorm:"column:Kapasitas" json:"kapasitas"`
	DosenID    uint           `gorm:"column:DosenID" json:"dosen_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`

	// Relations
	Dosen Dosen `gorm:"foreignKey:DosenID" json:"dosen"`
}

func (Jadwal) TableName() string {
	return "jadwal"
}
