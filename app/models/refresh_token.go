package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index:idx_user_token" json:"user_id"`
	Token     string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time      `gorm:"not null;index:idx_expires" json:"expires_at"`
	Revoked   bool           `gorm:"default:false;not null" json:"revoked"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
