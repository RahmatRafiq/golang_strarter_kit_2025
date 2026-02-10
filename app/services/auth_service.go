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
	user, err := auth.repo.FindUserByEmail(request.Email)
	if err != nil {
		return nil, errors.New("Email atau password salah")
	}

	check, err := helpers.ComparePasswordArgon2(request.Password, user.Password)
	if err != nil {
		return nil, errors.New("Email atau password salah")
	}
	if !check {
		return nil, errors.New("Email atau password salah")
	}

	// Generate token pair (access + refresh token)
	tokenPair, err := auth.jwt.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, err
	}

	// Store access token in user table for backward compatibility
	user.JwtToken = tokenPair.AccessToken
	if err := auth.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &casts.Token{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiredAt:    time.Now().Add(15 * time.Minute),
		ExpiresIn:    tokenPair.ExpiresIn,
		TokenType:    tokenPair.TokenType,
	}, nil
}

func (auth *AuthService) Logout(tokenString string) error {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	var userID uint
	switch v := userId.(type) {
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

	// Revoke all refresh tokens for this user
	if err := auth.jwt.RevokeAllUserTokens(userID); err != nil {
		return err
	}

	return nil
}

func (auth *AuthService) RefreshToken(refreshTokenString string) (*casts.Token, error) {
	// Use the new refresh token mechanism
	tokenPair, err := auth.jwt.RefreshAccessToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Extract user ID from new access token to update user table
	token, err := jwt.Parse(tokenPair.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["user_id"]

	var userID uint
	switch v := userId.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	default:
		return nil, errors.New("invalid user id")
	}

	user, err := auth.repo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update user's access token for backward compatibility
	user.JwtToken = tokenPair.AccessToken
	if err := auth.repo.UpdateUser(user); err != nil {
		return nil, err
	}

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
