package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct{}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

var jwtKey = []byte(helpers.GetEnv("JWT_SECRET_KEY", ""))

func (*JwtService) GenerateToken(claim jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtKey)
}

func (*JwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
}

func (*JwtService) ExtractClaims(token *jwt.Token) jwt.MapClaims {
	return token.Claims.(jwt.MapClaims)
}

// GenerateTokenPair generates access token and refresh token
func (j *JwtService) GenerateTokenPair(userID uint) (*TokenPair, error) {
	// Access token expires in 15 minutes
	accessTokenExpiry := time.Now().Add(15 * time.Minute)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     accessTokenExpiry.Unix(),
		"iat":     time.Now().Unix(),
		"type":    "access",
	}

	accessToken, err := j.GenerateToken(accessClaims)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenString, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	// Store refresh token in database (expires in 7 days)
	refreshToken := models.RefreshToken{
		UserID:    userID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Revoked:   false,
	}

	if err := facades.DB.Create(&refreshToken).Error; err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    900, // 15 minutes in seconds
		TokenType:    "Bearer",
	}, nil
}

// RefreshAccessToken generates new access token using refresh token
func (j *JwtService) RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	var refreshToken models.RefreshToken

	// Find and validate refresh token
	if err := facades.DB.Where("token = ? AND revoked = ? AND expires_at > ?",
		refreshTokenString, false, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Revoke old refresh token (token rotation)
	refreshToken.Revoked = true
	if err := facades.DB.Save(&refreshToken).Error; err != nil {
		return nil, err
	}

	// Generate new token pair
	return j.GenerateTokenPair(refreshToken.UserID)
}

// RevokeRefreshToken revokes a refresh token
func (*JwtService) RevokeRefreshToken(refreshTokenString string) error {
	result := facades.DB.Model(&models.RefreshToken{}).
		Where("token = ?", refreshTokenString).
		Update("revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("refresh token not found")
	}

	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (*JwtService) RevokeAllUserTokens(userID uint) error {
	return facades.DB.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true).Error
}

// CleanupExpiredTokens removes expired tokens from database (should be run periodically)
func (*JwtService) CleanupExpiredTokens() error {
	return facades.DB.Where("expires_at < ? OR revoked = ?", time.Now(), true).
		Delete(&models.RefreshToken{}).Error
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
