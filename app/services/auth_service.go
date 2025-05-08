//go:generate mockgen -source=auth_service.go -destination=mocks/auth_service.go -package=mocks
package services

import (
	"errors"
	"strings"
	"time"

	"golang_starter_kit_2025/app/casts"
	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/facades"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwt *JwtService
}

func (auth *AuthService) Login(request requests.LoginRequest) (*casts.Token, error) {
	var user models.User
	if err := facades.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		return nil, errors.New("Email atau password salah")
	}

	// if !CheckPasswordHash(request.Password, user.Password) {
	// 	return "", errors.New("Email atau password salah")
	// }
	check, err := helpers.ComparePasswordArgon2(request.Password, user.Password)
	if err != nil {
		return nil, errors.New("Email atau password salah")
	}
	if !check {
		return nil, errors.New("Email atau password salah")
	}

	// NOW: user can login multiple times
	// if user.JwtToken != "" {
	// 	return "", errors.New("Logout terlebih dahulu")
	// }

	expires := helpers.GetEnvInt("JWT_EXPIRE_MINUTES", 60)
	expireAt := time.Now().Add(time.Minute * time.Duration(expires)).Unix()
	// Generate JWT token
	tokenString, err := auth.jwt.GenerateToken(casts.NewJwtClaims(user.ID, expireAt))

	if err != nil {
		return nil, err
	}

	// Update user with the new token
	user.JwtToken = tokenString
	if err := facades.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &casts.Token{Token: tokenString, ExpiredAt: time.Unix(expireAt, 0)}, nil
}

func (auth *AuthService) Logout(tokenString string) error {
	// Hapus "Bearer " dari token jika ada
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Validasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	// Ambil user ID dari token
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	// Hapus JWT token dari database
	var user models.User
	if err := facades.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	user.JwtToken = ""
	if err := facades.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (auth *AuthService) RefreshToken(tokenString string) (*casts.Token, error) {
	// Refresh token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	// Ambil user dari database
	var user models.User
	if err := facades.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Generate JWT token
	expires := helpers.GetEnvInt("JWT_EXPIRE_MINUTES", 60)
	expireAt := time.Now().Add(time.Minute * time.Duration(expires)).Unix()
	tokenString, err = auth.jwt.GenerateToken(casts.NewJwtClaims(user.ID, expireAt))

	// Update user dengan token baru
	user.JwtToken = tokenString
	if err := facades.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &casts.Token{Token: tokenString, ExpiredAt: time.Unix(expireAt, 0)}, nil
}

func CheckPasswordHash(passwordOrPin, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordOrPin))
	return err == nil
}
