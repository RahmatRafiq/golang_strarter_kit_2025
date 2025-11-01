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

	expires := helpers.GetEnvInt("JWT_EXPIRE_MINUTES", 60)
	expireAt := time.Now().Add(time.Minute * time.Duration(expires)).Unix()
	tokenString, err := auth.jwt.GenerateToken(casts.NewJwtClaims(user.ID, expireAt))

	if err != nil {
		return nil, err
	}

	user.JwtToken = tokenString
	if err := auth.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &casts.Token{Token: tokenString, ExpiredAt: time.Unix(expireAt, 0)}, nil
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

	user.JwtToken = ""
	if err := auth.repo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (auth *AuthService) RefreshToken(tokenString string) (*casts.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
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

	expires := helpers.GetEnvInt("JWT_EXPIRE_MINUTES", 60)
	expireAt := time.Now().Add(time.Minute * time.Duration(expires)).Unix()
	tokenString, err = auth.jwt.GenerateToken(casts.NewJwtClaims(user.ID, expireAt))

	user.JwtToken = tokenString
	if err := auth.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &casts.Token{Token: tokenString, ExpiredAt: time.Unix(expireAt, 0)}, nil
}

func CheckPasswordHash(passwordOrPin, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordOrPin))
	return err == nil
}
