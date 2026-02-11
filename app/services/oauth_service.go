package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories/interfaces"
	"golang_starter_kit_2025/config"

	"github.com/rs/zerolog/log"
)

type OAuthService struct {
	userRepo interfaces.UserRepositoryInterface
	config   *config.OAuthConfig
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type GitHubUserInfo struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	ID        int64  `json:"id"`
}

func NewOAuthService(userRepo interfaces.UserRepositoryInterface) *OAuthService {
	return &OAuthService{
		userRepo: userRepo,
		config:   config.GetOAuthConfig(),
	}
}

func (s *OAuthService) GetGoogleAuthURL(state string) string {
	return s.config.GoogleConfig.AuthCodeURL(state)
}

func (s *OAuthService) GetGitHubAuthURL(state string) string {
	return s.config.GitHubConfig.AuthCodeURL(state)
}

func (s *OAuthService) HandleGoogleCallback(code string) (*models.User, string, error) {
	token, err := s.config.GoogleConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange token: %w", err)
	}

	client := s.config.GoogleConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, "", fmt.Errorf("failed to parse user info: %w", err)
	}

	return s.findOrCreateUser(userInfo.Email, userInfo.Name, "google")
}

func (s *OAuthService) HandleGitHubCallback(code string) (*models.User, string, error) {
	token, err := s.config.GitHubConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange token: %w", err)
	}

	client := s.config.GitHubConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, "", fmt.Errorf("failed to parse user info: %w", err)
	}

	email := userInfo.Email
	if email == "" {
		email = fmt.Sprintf("%s@github.oauth", userInfo.Login)
	}

	return s.findOrCreateUser(email, userInfo.Name, "github")
}

func (s *OAuthService) findOrCreateUser(email, name, provider string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err == nil {
		now := time.Now()
		if user.EmailVerifiedAt == nil {
			user.EmailVerifiedAt = &now
			if updateErr := s.userRepo.Update(user); updateErr != nil {
				log.Warn().Err(updateErr).Msg("Failed to update email verification")
			}
		}

		jwtService := &JwtService{}
		tokenPair, tokenErr := jwtService.GenerateTokenPair(user.ID)
		if tokenErr != nil {
			return nil, "", fmt.Errorf("failed to generate token: %w", tokenErr)
		}

		log.Info().Str("email", email).Str("provider", provider).Msg("OAuth login successful")
		return user, tokenPair.AccessToken, nil
	}

	username := email
	if name != "" {
		username = name
	}

	now := time.Now()
	user = &models.User{
		Username:        username,
		Email:           email,
		Password:        "", // OAuth users don't have password
		EmailVerifiedAt: &now,
	}

	if createErr := s.userRepo.Create(user); createErr != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", createErr)
	}

	jwtService := &JwtService{}
	tokenPair, tokenErr := jwtService.GenerateTokenPair(user.ID)
	if tokenErr != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", tokenErr)
	}

	log.Info().Str("email", email).Str("provider", provider).Msg("OAuth user created")
	return user, tokenPair.AccessToken, nil
}
