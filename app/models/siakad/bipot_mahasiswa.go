package models

import (
	"time"

	"gorm.io/gorm"
)

type BipotMahasiswa struct {
	BipotID     string         `gorm:"column:BipotID;primaryKey" json:"bipot_id"`
	MahasiswaID string         `gorm:"column:MhswID" json:"mahasiswa_id"`
	TahunID     string         `gorm:"column:TahunID" json:"tahun_id"`
	NamaTagihan string         `gorm:"column:Nama" json:"nama_tagihan"`
	Jumlah      int            `gorm:"column:Jumlah" json:"jumlah"`
	Besar       float64        `gorm:"column:Besar" json:"besar"`
	Dibayar     float64        `gorm:"column:Dibayar" json:"dibayar"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
}

func (BipotMahasiswa) TableName() string {
	return "bipotmhsw"
}
