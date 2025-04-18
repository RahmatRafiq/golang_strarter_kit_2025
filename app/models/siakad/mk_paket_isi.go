package models

type MKPaketIsi struct {
	MKPaketIsiID uint   `gorm:"primaryKey" json:"mk_paket_isi_id"`
	MKPaketID    uint   `gorm:"column:MKPaketID" json:"mk_paket_id"`
	MKID         uint   `gorm:"column:MKID" json:"mk_id"`
	ProdiID      string `gorm:"column:ProdiID" json:"prodi_id"`
	ProgramID    string `gorm:"column:ProgramID" json:"program_id"`

	// Relations
	MKPaket MKPaket `gorm:"foreignKey:MKPaketID" json:"mk_paket"`
}

func (MKPaketIsi) TableName() string {
	return "mkpaketisi"
}
