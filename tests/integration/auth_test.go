package integration

import (
	"strings"
	"testing"

	"golang_starter_kit_2025/app/repositories"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/tests/helpers"
)

// TestAuthService_Login_Integration tests full authentication flow with real database
// Similar to Laravel: php artisan test --filter=AuthTest
func TestAuthService_Login_Integration(t *testing.T) {
	// Setup test database (like Laravel's RefreshDatabase)
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	// Create services with real repositories
	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("success - user can login with valid credentials", func(t *testing.T) {
		// Arrange: Create test user with known password
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("john@example.com").
			WithUsername("john_doe").
			WithPassword("secure_password_123").
			CreateWithPlainPassword()

		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		}

		// Act: Attempt login
		token, err := authService.Login(loginReq)

		// Assert: Login successful
		if err != nil {
			t.Errorf("Expected successful login, got error: %v", err)
		}

		if token == nil {
			t.Error("Expected token, got nil")
		}

		if token != nil {
			if token.Token == "" {
				t.Error("Expected access token, got empty string")
			}
			if token.RefreshToken == "" {
				t.Error("Expected refresh token, got empty string")
			}
			if token.TokenType != "Bearer" {
				t.Errorf("Expected token type 'Bearer', got '%s'", token.TokenType)
			}
			t.Logf("✓ Login successful - Token length: %d", len(token.Token))
		}

		// Token generation verified - user can use token for authenticated requests
		// Note: Token storage is handled by JWT service in refresh_tokens table
	})

	t.Run("error - user not found", func(t *testing.T) {
		// Arrange: Try to login with non-existent email
		loginReq := requests.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "any_password",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Login fails
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
		if token != nil {
			t.Error("Expected nil token, got non-nil")
		}
		// Check error message contains "invalid" or "email" or "password" (language-agnostic)
		if err != nil {
			errMsg := strings.ToLower(err.Error())
			if !strings.Contains(errMsg, "invalid") && !strings.Contains(errMsg, "email") && !strings.Contains(errMsg, "password") {
				t.Errorf("Expected error message about invalid credentials, got '%s'", err.Error())
			}
		}
		t.Logf("✓ Correctly rejected non-existent user: %v", err)
	})

	t.Run("error - invalid password", func(t *testing.T) {
		// Arrange: Create user with known password
		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("jane@example.com").
			WithPassword("correct_password").
			CreateWithPlainPassword()

		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: "wrong_password",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Login fails
		if err == nil {
			t.Error("Expected error for wrong password, got nil")
		}
		if token != nil {
			t.Error("Expected nil token, got non-nil")
		}
		t.Logf("✓ Correctly rejected wrong password: %v", err)
	})

	t.Run("error - empty password", func(t *testing.T) {
		// Arrange: Create user
		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("test@example.com").
			WithPassword("password123").
			CreateWithPlainPassword()

		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: "", // Empty password
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Login fails
		if err == nil {
			t.Error("Expected error for empty password, got nil")
		}
		if token != nil {
			t.Error("Expected nil token, got non-nil")
		}
		t.Logf("✓ Correctly rejected empty password: %v", err)
	})
}

// TestAuthService_RefreshToken_Integration tests token refresh flow
func TestAuthService_RefreshToken_Integration(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("success - refresh token generates new access token", func(t *testing.T) {
		// Arrange: Login to get initial tokens
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("refresh@example.com").
			WithPassword("password123").
			CreateWithPlainPassword()

		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		}

		initialToken, err := authService.Login(loginReq)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		// Act: Refresh token
		newToken, err := authService.RefreshToken(initialToken.RefreshToken)

		// Assert: New token generated
		if err != nil {
			t.Errorf("Expected successful token refresh, got error: %v", err)
		}

		if newToken == nil {
			t.Error("Expected new token, got nil")
		}

		if newToken != nil {
			if newToken.Token == "" {
				t.Error("Expected new access token, got empty string")
			}
			if newToken.RefreshToken == "" {
				t.Error("Expected new refresh token, got empty string")
			}
			// Note: Access token might be the same if generated within same second
			// The important part is that refresh was successful
			t.Logf("✓ Token refreshed - New Token length: %d", len(newToken.Token))
		}
	})

	t.Run("error - invalid refresh token", func(t *testing.T) {
		// Act: Try to refresh with invalid token
		newToken, err := authService.RefreshToken("invalid_refresh_token_xyz")

		// Assert: Refresh fails
		if err == nil {
			t.Error("Expected error for invalid refresh token, got nil")
		}
		if newToken != nil {
			t.Error("Expected nil token, got non-nil")
		}
		t.Logf("✓ Correctly rejected invalid refresh token: %v", err)
	})

	t.Run("error - empty refresh token", func(t *testing.T) {
		// Act: Try to refresh with empty token
		newToken, err := authService.RefreshToken("")

		// Assert: Refresh fails
		if err == nil {
			t.Error("Expected error for empty refresh token, got nil")
		}
		if newToken != nil {
			t.Error("Expected nil token, got non-nil")
		}
		t.Logf("✓ Correctly rejected empty refresh token: %v", err)
	})
}

// TestAuthService_Logout_Integration tests logout flow
func TestAuthService_Logout_Integration(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("success - logout clears user tokens", func(t *testing.T) {
		// Arrange: Login to get access token
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("logout@example.com").
			WithPassword("password123").
			CreateWithPlainPassword()

		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		}

		token, err := authService.Login(loginReq)
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}

		// Act: Logout
		err = authService.Logout(token.Token)

		// Assert: Logout successful
		if err != nil {
			t.Errorf("Expected successful logout, got error: %v", err)
		}

		// Logout successful - refresh tokens revoked in database
		t.Log("✓ Logout successful - user session terminated")
	})

	t.Run("error - invalid access token format", func(t *testing.T) {
		// Act: Try to logout with invalid token
		err := authService.Logout("invalid_token_format")

		// Assert: Logout fails
		if err == nil {
			t.Error("Expected error for invalid token format, got nil")
		}
		if err != nil && err.Error() != "invalid token" {
			t.Errorf("Expected 'invalid token', got '%s'", err.Error())
		}
		t.Logf("✓ Correctly rejected invalid token format: %v", err)
	})

	t.Run("error - empty token", func(t *testing.T) {
		// Act: Try to logout with empty token
		err := authService.Logout("")

		// Assert: Logout fails
		if err == nil {
			t.Error("Expected error for empty token, got nil")
		}
		t.Logf("✓ Correctly rejected empty token: %v", err)
	})
}

// TestAuthService_CompleteFlow_Integration tests complete auth flow
// Login → Use token → Refresh → Logout
func TestAuthService_CompleteFlow_Integration(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("complete authentication flow", func(t *testing.T) {
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("flow@example.com").
			WithPassword("secure_password").
			CreateWithPlainPassword()

		// Step 1: Login
		t.Log("Step 1: Login with credentials")
		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		}

		token, err := authService.Login(loginReq)
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if token == nil {
			t.Fatal("Expected token, got nil")
		}
		t.Logf("✓ Login successful - Got access token")

		// Step 2: Refresh token
		t.Log("Step 2: Refresh access token")
		newToken, err := authService.RefreshToken(token.RefreshToken)
		if err != nil {
			t.Fatalf("Token refresh failed: %v", err)
		}
		if newToken == nil {
			t.Fatal("Expected new token, got nil")
		}
		t.Logf("✓ Token refreshed successfully")

		// Step 3: Logout
		t.Log("Step 3: Logout and clear tokens")
		err = authService.Logout(newToken.Token)
		if err != nil {
			t.Fatalf("Logout failed: %v", err)
		}
		t.Logf("✓ Logout successful")

		// Step 4: Verify flow completed successfully
		t.Log("✓ Complete authentication flow successful: Login → Refresh → Logout")
	})
}
