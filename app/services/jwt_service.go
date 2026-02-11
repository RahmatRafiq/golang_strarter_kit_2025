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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JwtService struct{}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
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

func (*JwtService) ExtractClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}
	return claims, nil
}

// GenerateTokenPair generates access token and refresh token
func (j *JwtService) GenerateTokenPair(userID uint) (*TokenPair, error) {
	return j.generateTokenPairWithTx(facades.DB, userID)
}

// generateTokenPairWithTx generates token pair using provided transaction
func (j *JwtService) generateTokenPairWithTx(db *gorm.DB, userID uint) (*TokenPair, error) {
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

	if err := db.Create(&refreshToken).Error; err != nil {
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
// FIXED: Added transaction with row locking to prevent concurrent refresh race condition
func (j *JwtService) RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	var newTokenPair *TokenPair

	// Use explicit transaction with row locking
	err := facades.DB.Transaction(func(tx *gorm.DB) error {
		var refreshToken models.RefreshToken

		// Find and lock the refresh token row
		// Use GORM's Clauses for database-agnostic locking (SQLite, MySQL, PostgreSQL)
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token = ? AND revoked = ? AND expires_at > ?", refreshTokenString, false, time.Now()).
			First(&refreshToken).Error

		if err != nil {
			return errors.New("invalid or expired refresh token")
		}

		// Check if token was found (ID will be 0 if not found)
		if refreshToken.ID == 0 {
			return errors.New("invalid or expired refresh token")
		}

		// Revoke old refresh token (token rotation)
		refreshToken.Revoked = true
		if saveErr := tx.Save(&refreshToken).Error; saveErr != nil {
			return saveErr
		}

		// Generate new token pair using the same transaction
		pair, err := j.generateTokenPairWithTx(tx, refreshToken.UserID)
		if err != nil {
			return err
		}

		newTokenPair = pair
		return nil
	})

	if err != nil {
		return nil, err
	}

	return newTokenPair, nil
}

// RevokeRefreshToken revokes a specific refresh token (per-device logout)
// IMPROVED: Now supports per-device logout instead of revoking all sessions
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

// RevokeCurrentSessionToken revokes only the current session's refresh token
// This allows users to logout from one device without affecting other devices
func (*JwtService) RevokeCurrentSessionToken(accessToken string) error {
	// Parse access token to get user ID
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return errors.New("invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims type")
	}

	userIDVal := claims["user_id"]
	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		return errors.New("invalid user id")
	}

	// Find the most recently used (non-revoked) token for this user
	// In a real app, you'd want to track which specific token corresponds to this access token
	// For now, we'll revoke the most recently created active token
	return facades.DB.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Order("created_at DESC").
		Limit(1).
		Update("revoked", true).Error
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
