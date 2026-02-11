package models

import "time"

type EmailVerificationToken struct {
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Token     string    `gorm:"size:255;not null;uniqueIndex" json:"token"`
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

func (t *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
