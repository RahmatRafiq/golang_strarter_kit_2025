package services

import (
	"golang_starter_kit_2025/app/helpers"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct{}

var jwtKey = []byte(helpers.GetEnv("APP_KEY", "your_secret_key"))

func (*JwtService) GenerateToken(claim jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtKey)
}

func (*JwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
}

func (*JwtService) ExtractClaims(token *jwt.Token) jwt.MapClaims {
	return token.Claims.(jwt.MapClaims)
}
