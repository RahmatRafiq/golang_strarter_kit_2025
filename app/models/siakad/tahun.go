package models

type Tahun struct {
	TahunID       string `gorm:"column:TahunID;primaryKey" json:"tahun_id"`
	Nama          string `gorm:"column:Nama" json:"nama"`
	TglKRSMulai   string `gorm:"column:TglKRSMulai" json:"tgl_krs_mulai"`
	TglKRSSelesai string `gorm:"column:TglKRSSelesai" json:"tgl_krs_selesai"`
	KodeID        string `gorm:"column:KodeID" json:"kode_id"`
	ProdiID       string `gorm:"column:ProdiID" json:"prodi_id"`
	ProgramID     string `gorm:"column:ProgramID" json:"program_id"`
}

func (Tahun) TableName() string {
	return "tahun"
}
