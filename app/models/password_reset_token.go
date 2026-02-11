package models

import "time"

// PasswordResetToken represents a password reset token record
type PasswordResetToken struct {
	UsedAt    *time.Time `gorm:"index" json:"used_at,omitempty"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	Email     string     `gorm:"size:255;not null;index" json:"email"`
	Token     string     `gorm:"size:255;not null;uniqueIndex" json:"token"`
	ID        uint       `gorm:"primaryKey" json:"id"`
}

// TableName specifies the table name for GORM
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// IsExpired checks if the token has expired
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if the token has been used
func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

// IsValid checks if the token is valid (not expired and not used)
func (t *PasswordResetToken) IsValid() bool {
	return !t.IsExpired() && !t.IsUsed()
}

// MarkAsUsed marks the token as used
func (t *PasswordResetToken) MarkAsUsed() {
	now := time.Now()
	t.UsedAt = &now
}
