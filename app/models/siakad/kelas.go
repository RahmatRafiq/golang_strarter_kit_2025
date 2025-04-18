package models

type Kelas struct {
	KelasID string `gorm:"column:KelasID;primaryKey" json:"kelas_id"`
	Nama    string `gorm:"column:Nama" json:"nama"`
}

func (Kelas) TableName() string {
	return "kelas"
}
