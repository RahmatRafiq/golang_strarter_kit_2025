package models

import "time"

type RoleHasPermissions struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `json:"role_id"`
	PermissionID uint      `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}