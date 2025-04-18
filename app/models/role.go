package models

type Role struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`

	Users []User `gorm:"many2many:user_has_roles;" json:"users"`
}
