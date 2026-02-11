package casts

import "time"

type Token struct {
	ExpiredAt    time.Time `json:"expired_at"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
	ExpiresIn    int64     `json:"expires_in,omitempty"`
}
