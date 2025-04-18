package models

type StatusMahasiswa struct {
	StatusMhswID string `gorm:"column:StatusMhswID;primaryKey" json:"status_mhsw_id"`
	Nama         string `gorm:"column:Nama" json:"nama"`
}

func (StatusMahasiswa) TableName() string {
	return "statusmhsw"
}
