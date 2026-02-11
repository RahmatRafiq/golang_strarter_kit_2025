package casts

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserID    uint  `json:"user_id"`
	ExpiredAt int64 `json:"expired_at"`
}

func NewJwtClaims(userID uint, expiredAt int64) jwt.MapClaims {
	return jwt.MapClaims{
		"user_id":    userID,
		"expired_at": expiredAt,
	}
}

func ParseJwtClaims(claims jwt.Claims) (JwtClaims, error) {
	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return JwtClaims{}, fmt.Errorf("invalid claims type")
	}

	userIDFloat, ok := mapClaims["user_id"].(float64)
	if !ok {
		return JwtClaims{}, fmt.Errorf("invalid user_id type")
	}

	expiredAtFloat, ok := mapClaims["expired_at"].(float64)
	if !ok {
		return JwtClaims{}, fmt.Errorf("invalid expired_at type")
	}

	return JwtClaims{
		UserID:    uint(userIDFloat),
		ExpiredAt: int64(expiredAtFloat),
	}, nil
}
