package services_test

import (
	"errors"
	"testing"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/mocks"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"
	testhelpers "golang_starter_kit_2025/tests/helpers"

	"go.uber.org/mock/gomock"
)

// NOTE: These unit tests require database connection for JWT token generation
// The AuthService uses JwtService which stores refresh tokens in database via facades.DB
// For pure unit tests without database, see the _WithoutDB tests below
// For complete integration tests, see tests/integration/auth_test.go

// createTestPasswordHash creates a valid Argon2 password hash for testing
func createTestPasswordHash(t *testing.T, plainPassword string) string {
	t.Helper()

	hash, err := helpers.HashPasswordArgon2(plainPassword, helpers.DefaultParams)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}
	return hash
}

func TestAuthService_Login_Unit(t *testing.T) {
	// Setup test database (required for JWT token generation)
	testDB := testhelpers.SetupTestDB(t)
	defer testhelpers.CleanupTestDB(t, testDB)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepositoryInterface(ctrl)
	authService := services.NewAuthService(mockRepo)

	t.Run("success - valid credentials", func(t *testing.T) {
		plainPassword := "correctpassword"

		// Create user with plain password - BeforeCreate hook will hash it
		userToCreate := &models.User{
			Email:    "test@example.com",
			Password: plainPassword,
			Username: "testuser",
		}

		// Create actual user in database for foreign key constraint
		if err := testDB.DB.Create(userToCreate).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// After creation, password is hashed and ID is assigned
		expectedUser := userToCreate

		loginReq := requests.LoginRequest{
			Email:    "test@example.com",
			Password: plainPassword,
		}

		// Expect FindUserByEmail to be called
		mockRepo.EXPECT().
			FindUserByEmail(loginReq.Email).
			Return(expectedUser, nil)

		// Expect UpdateUser to be called to store JWT token
		mockRepo.EXPECT().
			UpdateUser(gomock.Any()).
			DoAndReturn(func(user *models.User) error {
				// Verify token was set
				if user.JwtToken == "" {
					t.Error("Expected JWT token to be set")
				}
				return nil
			})

		token, err := authService.Login(loginReq)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == nil {
			t.Fatal("expected token, got nil")
		}
		if token.Token == "" {
			t.Error("expected access token, got empty string")
		}
		if token.RefreshToken == "" {
			t.Error("expected refresh token, got empty string")
		}
		if token.TokenType != "Bearer" {
			t.Errorf("expected token type 'Bearer', got '%s'", token.TokenType)
		}
	})

	t.Run("error - user not found", func(t *testing.T) {
		loginReq := requests.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "anypassword",
		}

		mockRepo.EXPECT().
			FindUserByEmail(loginReq.Email).
			Return(nil, errors.New("user not found"))

		token, err := authService.Login(loginReq)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if token != nil {
			t.Error("expected nil token, got non-nil")
		}
		if err.Error() != "invalid email or password" {
			t.Errorf("expected 'invalid email or password', got %v", err)
		}
	})

	t.Run("error - invalid password", func(t *testing.T) {
		plainPassword := "correctpassword"
		hashedPassword := createTestPasswordHash(t, plainPassword)

		expectedUser := &models.User{
			ID:       1,
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		loginReq := requests.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword", // Wrong password
		}

		mockRepo.EXPECT().
			FindUserByEmail(loginReq.Email).
			Return(expectedUser, nil)

		token, err := authService.Login(loginReq)

		if err == nil {
			t.Error("expected error for wrong password, got nil")
		}
		if token != nil {
			t.Error("expected nil token, got non-nil")
		}
		if err.Error() != "invalid email or password" {
			t.Errorf("expected 'invalid email or password', got %v", err)
		}
	})

	t.Run("error - empty password", func(t *testing.T) {
		plainPassword := "correctpassword"
		hashedPassword := createTestPasswordHash(t, plainPassword)

		expectedUser := &models.User{
			ID:       1,
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		loginReq := requests.LoginRequest{
			Email:    "test@example.com",
			Password: "", // Empty password
		}

		mockRepo.EXPECT().
			FindUserByEmail(loginReq.Email).
			Return(expectedUser, nil)

		token, err := authService.Login(loginReq)

		if err == nil {
			t.Error("expected error for empty password, got nil")
		}
		if token != nil {
			t.Error("expected nil token, got non-nil")
		}
	})

	t.Run("error - failed to update user token", func(t *testing.T) {
		plainPassword := "correctpassword"

		// Create user with plain password - BeforeCreate hook will hash it
		userToCreate := &models.User{
			Email:    "test2@example.com",
			Password: plainPassword,
			Username: "testuser2",
		}

		// Create actual user in database for foreign key constraint
		if err := testDB.DB.Create(userToCreate).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		// After creation, password is hashed and ID is assigned
		expectedUser := userToCreate

		loginReq := requests.LoginRequest{
			Email:    "test2@example.com",
			Password: plainPassword,
		}

		mockRepo.EXPECT().
			FindUserByEmail(loginReq.Email).
			Return(expectedUser, nil)

		mockRepo.EXPECT().
			UpdateUser(gomock.Any()).
			Return(errors.New("database error"))

		token, err := authService.Login(loginReq)

		if err == nil {
			t.Error("expected error when update fails, got nil")
		}
		if token != nil {
			t.Error("expected nil token, got non-nil")
		}
		if err.Error() != "database error" {
			t.Errorf("expected 'database error', got %v", err)
		}
	})
}

func TestAuthService_Logout_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepositoryInterface(ctrl)
	authService := services.NewAuthService(mockRepo)

	t.Run("error - invalid token format", func(t *testing.T) {
		invalidToken := "invalid-token-string"

		err := authService.Logout(invalidToken)

		if err == nil {
			t.Error("expected error for invalid token, got nil")
		}
		if err.Error() != "invalid token" {
			t.Errorf("expected 'invalid token', got %v", err)
		}
	})

	t.Run("error - empty token", func(t *testing.T) {
		err := authService.Logout("")

		if err == nil {
			t.Error("expected error for empty token, got nil")
		}
		if err.Error() != "invalid token" {
			t.Errorf("expected 'invalid token', got %v", err)
		}
	})

	t.Run("error - malformed token", func(t *testing.T) {
		malformedToken := "Bearer malformed.jwt.token"

		err := authService.Logout(malformedToken)

		if err == nil {
			t.Error("expected error for malformed token, got nil")
		}
		if err.Error() != "invalid token" {
			t.Errorf("expected 'invalid token', got %v", err)
		}
	})

	// Note: Testing successful logout requires generating a valid JWT token
	// which needs JwtService setup. This is better covered in integration tests.
}

func TestAuthService_RefreshToken_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepositoryInterface(ctrl)
	authService := services.NewAuthService(mockRepo)

	t.Run("error - invalid refresh token", func(t *testing.T) {
		invalidRefreshToken := "invalid-refresh-token"

		token, err := authService.RefreshToken(invalidRefreshToken)

		if err == nil {
			t.Error("expected error for invalid refresh token, got nil")
		}
		if token != nil {
			t.Error("expected nil token, got non-nil")
		}
	})

	t.Run("error - empty refresh token", func(t *testing.T) {
		token, err := authService.RefreshToken("")

		if err == nil {
			t.Error("expected error for empty refresh token, got nil")
		}
		if token != nil {
			t.Error("expected nil token, got non-nil")
		}
	})

	t.Run("error - malformed refresh token", func(t *testing.T) {
		malformedToken := "malformed.token.string"

		token, err := authService.RefreshToken(malformedToken)

		if err == nil {
			t.Error("expected error for malformed token, got nil")
		}
		if token != nil {
			t.Error("expected nil token, got non-nil")
		}
	})

	// Note: Testing successful token refresh requires:
	// 1. Valid refresh token from database
	// 2. JwtService fully mocked
	// 3. Database state setup
	// This is better covered in integration tests.
}

func TestAuthService_CheckPasswordHash_Unit(t *testing.T) {
	t.Run("success - correct password matches hash", func(t *testing.T) {
		plainPassword := "testpassword123"
		hashedPassword := createTestPasswordHash(t, plainPassword)

		match, err := helpers.ComparePasswordArgon2(plainPassword, hashedPassword)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !match {
			t.Error("expected password to match hash")
		}
	})

	t.Run("error - incorrect password does not match hash", func(t *testing.T) {
		plainPassword := "correctpassword"
		wrongPassword := "wrongpassword"
		hashedPassword := createTestPasswordHash(t, plainPassword)

		match, err := helpers.ComparePasswordArgon2(wrongPassword, hashedPassword)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if match {
			t.Error("expected password to not match hash")
		}
	})

	t.Run("error - invalid hash format", func(t *testing.T) {
		invalidHash := "not-a-valid-hash"

		match, err := helpers.ComparePasswordArgon2("anypassword", invalidHash)

		if err == nil {
			t.Error("expected error for invalid hash format, got nil")
		}
		if match {
			t.Error("expected no match for invalid hash")
		}
	})

	t.Run("error - empty password", func(t *testing.T) {
		hashedPassword := createTestPasswordHash(t, "somepassword")

		match, err := helpers.ComparePasswordArgon2("", hashedPassword)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if match {
			t.Error("expected empty password to not match")
		}
	})
}
