package models

type Dosen struct {
	DosenID string `gorm:"column:DosenID;primaryKey" json:"dosen_id"`
	Nama    string `gorm:"column:Nama" json:"nama"`
	Gelar   string `gorm:"column:Gelar" json:"gelar"`
}

func (Dosen) TableName() string {
	return "dosen"
}
