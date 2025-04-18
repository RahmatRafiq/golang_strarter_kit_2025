package models

import (
	"time"

	"gorm.io/gorm"
)

type Mahasiswa struct {
	MhswID            string `gorm:"column:MhswID;primaryKey" json:"mhsw_id"`
	Foto              string `gorm:"column:Foto" json:"foto"`
	TahunID           string `gorm:"column:TahunID" json:"tahun_id"`
	Nama              string `gorm:"column:Nama" json:"nama"`
	StatusAwalID      string `gorm:"column:StatusAwalID" json:"status_awal_id"`
	StatusMhswID      string `gorm:"column:StatusMhswID" json:"status_mhsw_id"`
	ProgramID         string `gorm:"column:ProgramID" json:"program_id"`
	ProdiID           string `gorm:"column:ProdiID" json:"prodi_id"`
	PenasehatAkademik string `gorm:"column:PA" json:"penasehat_akademik"`
	Alamat            string `gorm:"column:Alamat" json:"alamat"`
	Negara            string `gorm:"column:KTP" json:"negara"`
	Handphone         string `gorm:"column:Handphone" json:"handphone"`
	Email             string `gorm:"column:Email" json:"email"`
	NamaIbu           string `gorm:"column:NamaIbu" json:"nama_ibu"`
	AgamaID           string `gorm:"column:Agama" json:"agama_id"`
	HandphoneOrtu     string `gorm:"column:HandphoneOrtu" json:"handphone_ortu"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`

	Agama *[]Agama `gorm:"foreignKey:Agama" json:"agama"`
}

func (Mahasiswa) TableName() string {
	return "mhsw"
}
