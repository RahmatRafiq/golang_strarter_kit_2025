//go:generate mockgen -source=auth_service.go -destination=mocks/auth_service.go -package=mocks
package services

import (
	"errors"
	"strings"
	"time"

	"golang_starter_kit_2025/app/casts"
	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"golang_starter_kit_2025/app/requests"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwt  *JwtService
	repo interfaces.AuthRepositoryInterface
}

func NewAuthService(repo interfaces.AuthRepositoryInterface) *AuthService {
	return &AuthService{
		jwt:  &JwtService{},
		repo: repo,
	}
}

func (auth *AuthService) Login(request requests.LoginRequest) (*casts.Token, error) {
	// IMPROVED: Sanitize email input (trim whitespace for better UX)
	sanitizedEmail := strings.TrimSpace(request.Email)

	user, err := auth.repo.FindUserByEmail(sanitizedEmail)
	if err != nil {
		log.Warn().
			Str("email", sanitizedEmail).
			Msg("Login attempt with non-existent email")
		return nil, errors.New("email atau password salah")
	}

	// IMPROVED: Validate password length (max 128 chars for Argon2)
	if len(request.Password) > 128 {
		log.Warn().
			Uint("user_id", user.ID).
			Str("email", sanitizedEmail).
			Int("password_length", len(request.Password)).
			Msg("Login attempt with excessively long password (DoS protection)")
		return nil, errors.New("email atau password salah")
	}

	check, err := helpers.ComparePasswordArgon2(request.Password, user.Password)
	if err != nil || !check {
		log.Warn().
			Uint("user_id", user.ID).
			Str("email", sanitizedEmail).
			Msg("Login attempt with invalid password")
		return nil, errors.New("email atau password salah")
	}

	// Generate token pair (access + refresh token)
	tokenPair, err := auth.jwt.GenerateTokenPair(user.ID)
	if err != nil {
		log.Error().
			Err(err).
			Uint("user_id", user.ID).
			Msg("Failed to generate token pair")
		return nil, err
	}

	// Store access token in user table for backward compatibility
	user.JwtToken = tokenPair.AccessToken
	if err := auth.repo.UpdateUser(user); err != nil {
		log.Error().
			Err(err).
			Uint("user_id", user.ID).
			Msg("Failed to update user token")
		return nil, err
	}

	log.Info().
		Uint("user_id", user.ID).
		Str("email", sanitizedEmail).
		Msg("User logged in successfully")

	return &casts.Token{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiredAt:    time.Now().Add(15 * time.Minute),
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    tokenPair.TokenType,
	}, nil
}

// Logout logs out the user from current device only (preserves other sessions)
// IMPROVED: Now supports per-device logout instead of terminating all sessions
func (auth *AuthService) Logout(tokenString string) error {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		log.Warn().
			Err(err).
			Msg("Logout attempt with invalid token")
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Error().Msg("Invalid claims type in token")
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
		log.Error().
			Interface("user_id", userIDVal).
			Msg("Invalid user ID type in token")
		return errors.New("invalid user id")
	}

	user, err := auth.repo.FindUserByID(userID)
	if err != nil {
		log.Warn().
			Uint("user_id", userID).
			Msg("Logout attempt for non-existent user")
		return errors.New("user not found")
	}

	// Clear JWT token in user table
	user.JwtToken = ""
	if err := auth.repo.UpdateUser(user); err != nil {
		log.Error().
			Err(err).
			Uint("user_id", userID).
			Msg("Failed to clear user token")
		return err
	}

	// IMPROVED: Revoke only current session token (per-device logout)
	// This allows users to stay logged in on other devices
	if err := auth.jwt.RevokeCurrentSessionToken(tokenString); err != nil {
		log.Error().
			Err(err).
			Uint("user_id", userID).
			Msg("Failed to revoke current session token")
		return err
	}

	log.Info().
		Uint("user_id", userID).
		Msg("User logged out from current device successfully")

	return nil
}

// LogoutAllDevices logs out the user from ALL devices (security incident response)
// Use this for security purposes when user wants to terminate all sessions
func (auth *AuthService) LogoutAllDevices(tokenString string) error {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
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

	user, err := auth.repo.FindUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Clear JWT token in user table
	user.JwtToken = ""
	if err := auth.repo.UpdateUser(user); err != nil {
		return err
	}

	// Revoke ALL refresh tokens for security
	if err := auth.jwt.RevokeAllUserTokens(userID); err != nil {
		return err
	}

	log.Info().
		Uint("user_id", userID).
		Msg("User logged out from ALL devices (security action)")

	return nil
}

func (auth *AuthService) RefreshToken(refreshTokenString string) (*casts.Token, error) {
	// Use the new refresh token mechanism
	tokenPair, err := auth.jwt.RefreshAccessToken(refreshTokenString)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Token refresh failed with invalid refresh token")
		return nil, err
	}

	// Extract user ID from new access token to update user table
	token, err := jwt.Parse(tokenPair.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to parse newly generated access token")
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Error().Msg("Invalid claims type in refreshed token")
		return nil, errors.New("invalid claims type")
	}

	userIDVal := claims["user_id"]

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		log.Error().
			Interface("user_id", userIDVal).
			Msg("Invalid user ID type in refreshed token")
		return nil, errors.New("invalid user id")
	}

	user, err := auth.repo.FindUserByID(userID)
	if err != nil {
		log.Warn().
			Uint("user_id", userID).
			Msg("Token refresh for non-existent user")
		return nil, errors.New("user not found")
	}

	// Update user's access token for backward compatibility
	user.JwtToken = tokenPair.AccessToken
	if err := auth.repo.UpdateUser(user); err != nil {
		log.Error().
			Err(err).
			Uint("user_id", userID).
			Msg("Failed to update user token after refresh")
		return nil, err
	}

	log.Info().
		Uint("user_id", userID).
		Msg("Token refreshed successfully")

	return &casts.Token{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiredAt:    time.Now().Add(15 * time.Minute),
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    tokenPair.TokenType,
	}, nil
}

func CheckPasswordHash(passwordOrPin, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordOrPin))
	return err == nil
}
