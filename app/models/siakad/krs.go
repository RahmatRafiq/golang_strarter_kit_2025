package models

import (
	"time"

	"gorm.io/gorm"
)

type KRS struct {
	KRSID       uint           `gorm:"primaryKey" json:"krs_id"`
	MahasiswaID string         `gorm:"column:MhswID" json:"mahasiswa_id"`
	TahunID     string         `gorm:"column:TahunID" json:"tahun_id"`
	MKID        uint           `gorm:"column:MKID" json:"mk_id"`
	KHSID       string         `gorm:"column:KHSID" json:"khs_id"`   // Menambahkan KHSID
	KodeID      string         `gorm:"column:KodeID" json:"kode_id"` // Menambahkan KodeID
	SKS         int            `gorm:"column:SKS" json:"sks"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
}

func (KRS) TableName() string {
	return "krs"
}
