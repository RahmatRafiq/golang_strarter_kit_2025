package casts

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	UserID    uint `json:"user_id"`
	ExpiredAt int64 `json:"expired_at"`
}

// set JWT claims
func NewJwtClaims(userID uint, expiredAt int64) jwt.MapClaims {
	return jwt.MapClaims{
		"user_id": userID,
		"expired_at": expiredAt,
	}
}

// get JWT claims and parse it
func ParseJwtClaims(claims jwt.Claims) (JwtClaims) {
	mapClaims := claims.(jwt.MapClaims)
	userID := uint(mapClaims["user_id"].(float64))
	expiredAt := int64(mapClaims["expired_at"].(float64))
	return JwtClaims{
		UserID: userID,
		ExpiredAt: expiredAt,
	}
}