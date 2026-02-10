package casts

import "time"

type Token struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiredAt    time.Time `json:"expired_at"`
	ExpiresIn    int64     `json:"expires_in,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
}
