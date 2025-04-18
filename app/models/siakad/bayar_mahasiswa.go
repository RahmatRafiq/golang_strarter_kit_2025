package models

import (
	"time"

	"gorm.io/gorm"
)

type BayarMahasiswa struct {
	RekeningID     string         `gorm:"column:RekeningID" json:"rekening_id"`
	MhswID         string         `gorm:"column:MhswID" json:"mhsw_id"`
	Tanggal        time.Time      `gorm:"column:Tanggal" json:"tanggal"`
	TahunID        string         `gorm:"column:TahunID" json:"tahun_id"`
	NamaPembayaran string         `gorm:"column:NamaPembayaran" json:"nama_pembayaran"`
	Jumlah         float64        `gorm:"column:Jumlah" json:"jumlah"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`

	BipotMhsw BipotMahasiswa `gorm:"foreignKey:RekeningID;references:BipotMhswid" json:"bipotmhsw"`
}

func (BayarMahasiswa) TableName() string {
	return "bayarmhsw"
}
