package models

type Role struct {
	Name        string       `json:"name"`
	Group       string       `json:"group"`
	Users       []User       `gorm:"many2many:user_has_roles;" json:"users"`
	Permissions []Permission `gorm:"many2many:role_has_permissions;" json:"permissions"`
	ID          uint         `gorm:"primaryKey" json:"id"`
}
