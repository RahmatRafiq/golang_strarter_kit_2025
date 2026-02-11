package models

import (
	"time"
)

type Permission struct {
	Name  string `json:"name"`
	Group string `json:"group"`
	ID    uint   `gorm:"primaryKey" json:"id"`
}

type UserHasPermissions struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id"`
	PermissionID uint      `json:"permission_id"`
}
