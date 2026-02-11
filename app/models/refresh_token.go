package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ExpiresAt time.Time      `gorm:"not null;index:idx_expires" json:"expires_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Token     string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"token"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index:idx_user_token" json:"user_id"`
	Revoked   bool           `gorm:"default:false;not null" json:"revoked"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
