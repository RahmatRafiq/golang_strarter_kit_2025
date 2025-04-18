package models

type Agama struct {
	Agama string `gorm:"column:Agama;primaryKey" json:"agama"`
	Nama  string `gorm:"column:Nama" json:"nama"`
}

func (Agama) TableName() string {
	return "agama"
}
