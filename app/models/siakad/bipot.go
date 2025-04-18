package models

type Bipot struct {
	BipotID string `gorm:"column:BIPOTID;primaryKey" json:"bipot_id"`
	Nama    string `gorm:"column:Nama" json:"nama"`
}

func (Bipot) TableName() string {
	return "bipot"
}
