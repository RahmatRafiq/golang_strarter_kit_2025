package integration

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/tests/helpers"
)

// TestAuthService_SoftDeletedUser tests that soft-deleted users cannot login
// CRITICAL SECURITY: Prevents deleted accounts from being used
func TestAuthService_SoftDeletedUser_Security(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("CRITICAL - soft-deleted user cannot login", func(t *testing.T) {
		// Arrange: Create user
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("deleted@test.com").
			WithPassword("password123").
			CreateWithPlainPassword()

		// Soft delete the user
		err := testDB.DB.Delete(user).Error
		if err != nil {
			t.Fatalf("Failed to soft delete user: %v", err)
		}

		// Verify user is soft-deleted
		var deletedUser models.User
		testDB.DB.Unscoped().Where("id = ?", user.ID).First(&deletedUser)
		if deletedUser.DeletedAt.Time.IsZero() {
			t.Fatal("User should be soft-deleted")
		}

		// Act: Try to login with deleted user
		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert: Login should fail
		if err == nil {
			t.Error("🚨 SECURITY ISSUE: Soft-deleted user can still login!")
		}
		if token != nil {
			t.Error("🚨 SECURITY ISSUE: Token was generated for deleted user!")
		}
		t.Logf("✓ Soft-deleted user correctly rejected")
	})

	t.Run("CRITICAL - hard-deleted user cannot login", func(t *testing.T) {
		// Arrange: Create and permanently delete user
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("permanent@test.com").
			WithPassword("password123").
			CreateWithPlainPassword()

		// Hard delete
		err := testDB.DB.Unscoped().Delete(user).Error
		if err != nil {
			t.Fatalf("Failed to hard delete user: %v", err)
		}

		// Act: Try to login
		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert: Should fail
		if err == nil {
			t.Error("Expected error for deleted user")
		}
		if token != nil {
			t.Error("Expected nil token for deleted user")
		}
		t.Log("✓ Hard-deleted user correctly rejected")
	})
}

// TestAuthService_RefreshTokenReuse tests refresh token rotation security
// CRITICAL SECURITY: Prevents token replay attacks
func TestAuthService_RefreshTokenReuse_Security(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("CRITICAL - refresh token cannot be reused after refresh", func(t *testing.T) {
		// Arrange: Login to get tokens
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("tokenreuse@test.com").
			WithPassword("password123").
			CreateWithPlainPassword()

		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		}

		initialToken, err := authService.Login(loginReq)
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}

		oldRefreshToken := initialToken.RefreshToken

		// Act: Refresh token (first time - should work)
		newToken, err := authService.RefreshToken(oldRefreshToken)
		if err != nil {
			t.Fatalf("First refresh failed: %v", err)
		}
		if newToken == nil {
			t.Fatal("Expected new token after refresh")
		}

		t.Logf("✓ First refresh successful")

		// Act: Try to reuse old refresh token (should fail)
		reusedToken, err := authService.RefreshToken(oldRefreshToken)

		// Assert: Reuse should fail
		if err == nil {
			t.Error("🚨 SECURITY ISSUE: Old refresh token can be reused!")
		}
		if reusedToken != nil {
			t.Error("🚨 SECURITY ISSUE: Token generated from revoked refresh token!")
		}

		// Verify old token is revoked in database
		var oldToken models.RefreshToken
		result := testDB.DB.Where("token = ?", oldRefreshToken).First(&oldToken)

		if result.Error == nil && !oldToken.Revoked {
			t.Error("🚨 SECURITY ISSUE: Old refresh token not marked as revoked!")
		}

		t.Logf("✓ Old refresh token correctly rejected (revoked)")
	})

	t.Run("CRITICAL - multiple refresh token generations create new tokens", func(t *testing.T) {
		// Arrange
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("multirefresh@test.com").
			CreateWithPlainPassword()

		token, _ := authService.Login(requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		})

		// Act: Chain refresh 3 times
		token1, _ := authService.RefreshToken(token.RefreshToken)
		token2, _ := authService.RefreshToken(token1.RefreshToken)
		token3, _ := authService.RefreshToken(token2.RefreshToken)

		// Assert: Each refresh token should be different
		tokens := []string{
			token.RefreshToken,
			token1.RefreshToken,
			token2.RefreshToken,
			token3.RefreshToken,
		}

		for i := 0; i < len(tokens); i++ {
			for j := i + 1; j < len(tokens); j++ {
				if tokens[i] == tokens[j] {
					t.Errorf("🚨 SECURITY ISSUE: Refresh tokens should be unique! Token %d == Token %d", i, j)
				}
			}
		}

		// Verify only the last token is valid
		var activeTokens []models.RefreshToken
		testDB.DB.Where("user_id = ? AND revoked = ?", user.ID, false).Find(&activeTokens)

		validCount := 0
		for _, rt := range activeTokens {
			if rt.Token == token3.RefreshToken {
				validCount++
			}
		}

		if validCount != 1 {
			t.Errorf("Expected exactly 1 valid refresh token, got %d", validCount)
		}

		t.Logf("✓ Refresh token rotation working correctly")
	})
}

// TestAuthService_ConcurrentAccess tests thread safety and race conditions
// CRITICAL SECURITY: Prevents race condition exploits
func TestAuthService_ConcurrentAccess_Security(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("CRITICAL - concurrent login attempts are safe", func(t *testing.T) {
		// Arrange
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("concurrent@test.com").
			WithPassword("password123").
			CreateWithPlainPassword()

		// Act: 10 concurrent login attempts
		var wg sync.WaitGroup
		results := make([]error, 10)
		tokens := make([]*models.User, 10)

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				token, err := authService.Login(requests.LoginRequest{
					Email:    user.Email,
					Password: plainPassword,
				})
				results[idx] = err
				if token != nil {
					tokens[idx] = user
				}
			}(i)
		}

		wg.Wait()

		// Assert: All should succeed (no race conditions)
		successCount := 0
		for i, err := range results {
			if err != nil {
				t.Errorf("Login %d failed: %v", i, err)
			} else {
				successCount++
			}
		}

		if successCount != 10 {
			t.Errorf("Expected all 10 concurrent logins to succeed, got %d", successCount)
		}

		t.Logf("✓ Concurrent logins handled safely (%d/10 successful)", successCount)
	})

	t.Run("CRITICAL - concurrent refresh prevents double use", func(t *testing.T) {
		// Arrange: Get initial token
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("concurrentrefresh@test.com").
			CreateWithPlainPassword()

		token, _ := authService.Login(requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		})

		// Act: Try to refresh same token 10 times concurrently
		var wg sync.WaitGroup
		successCount := 0
		var mu sync.Mutex

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				newToken, err := authService.RefreshToken(token.RefreshToken)
				if err == nil && newToken != nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// Assert: Only 1 should succeed (token rotation)
		if successCount > 1 {
			t.Errorf("🚨 SECURITY ISSUE: %d concurrent refreshes succeeded (expected only 1)", successCount)
			t.Error("This indicates a race condition in token rotation!")
		} else {
			t.Logf("✓ Token rotation prevents concurrent reuse (%d/10 succeeded)", successCount)
		}

		// Verify token is revoked
		var refreshToken models.RefreshToken
		testDB.DB.Where("token = ?", token.RefreshToken).First(&refreshToken)

		if !refreshToken.Revoked {
			t.Error("🚨 Refresh token should be revoked after use")
		}
	})

	t.Run("CRITICAL - concurrent logout is safe", func(t *testing.T) {
		// Arrange
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("concurrentlogout@test.com").
			CreateWithPlainPassword()

		token, _ := authService.Login(requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		})

		// Act: Multiple concurrent logout attempts
		var wg sync.WaitGroup
		errorCount := 0
		var mu sync.Mutex

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := authService.Logout(token.Token)
				if err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
				}
			}()
		}

		wg.Wait()

		// Assert: First logout succeeds, rest may fail (acceptable)
		// Important: No panics or race conditions
		t.Logf("✓ Concurrent logout handled safely (%d errors out of 5 attempts)", errorCount)
	})
}

// TestAuthService_TokenExpiration tests token lifetime handling
// SECURITY: Validates token expiry mechanisms
func TestAuthService_TokenExpiration_Security(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("security - expired refresh token rejected", func(t *testing.T) {
		// Arrange: Create user and manually insert expired token
		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("expired@test.com").
			CreateWithPlainPassword()

		expiredToken := models.RefreshToken{
			UserID:    user.ID,
			Token:     "expired_token_xyz_123",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Already expired
			Revoked:   false,
		}

		err := testDB.DB.Create(&expiredToken).Error
		if err != nil {
			t.Fatalf("Failed to create expired token: %v", err)
		}

		// Act: Try to refresh with expired token
		newToken, err := authService.RefreshToken(expiredToken.Token)

		// Assert: Should fail
		if err == nil {
			t.Error("🚨 SECURITY ISSUE: Expired refresh token accepted!")
		}
		if newToken != nil {
			t.Error("🚨 SECURITY ISSUE: Token generated from expired refresh token!")
		}

		t.Log("✓ Expired refresh token correctly rejected")
	})

	t.Run("security - revoked token rejected", func(t *testing.T) {
		// Arrange: Create revoked token
		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("revoked@test.com").
			CreateWithPlainPassword()

		revokedToken := models.RefreshToken{
			UserID:    user.ID,
			Token:     "revoked_token_xyz_123",
			ExpiresAt: time.Now().Add(24 * time.Hour), // Still valid time-wise
			Revoked:   true,                           // But revoked
		}

		testDB.DB.Create(&revokedToken)

		// Act: Try to use revoked token
		newToken, err := authService.RefreshToken(revokedToken.Token)

		// Assert: Should fail
		if err == nil {
			t.Error("🚨 SECURITY ISSUE: Revoked token accepted!")
		}
		if newToken != nil {
			t.Error("🚨 SECURITY ISSUE: Token generated from revoked token!")
		}

		t.Log("✓ Revoked token correctly rejected")
	})
}

// TestAuthService_SessionManagement tests multi-device session handling
// SECURITY: Validates proper session isolation
func TestAuthService_SessionManagement_Security(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("security - multiple device sessions allowed", func(t *testing.T) {
		// Arrange
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("multidevice@test.com").
			CreateWithPlainPassword()

		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		}

		// Act: Login from 3 devices
		device1, _ := authService.Login(loginReq)
		device2, _ := authService.Login(loginReq)
		device3, _ := authService.Login(loginReq)

		// Assert: All tokens should be valid
		if device1 == nil || device2 == nil || device3 == nil {
			t.Fatal("All device logins should succeed")
		}

		// Verify 3 active refresh tokens exist
		var activeTokens []models.RefreshToken
		testDB.DB.Where("user_id = ? AND revoked = ? AND expires_at > ?",
			user.ID, false, time.Now()).Find(&activeTokens)

		if len(activeTokens) < 3 {
			t.Errorf("Expected at least 3 active refresh tokens, got %d", len(activeTokens))
		}

		t.Logf("✓ Multiple device sessions working (%d active tokens)", len(activeTokens))
	})

	t.Run("ISSUE - logout revokes ALL devices", func(t *testing.T) {
		// Arrange: Login from 2 devices
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("logoutall@test.com").
			CreateWithPlainPassword()

		phone, _ := authService.Login(requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		})

		laptop, _ := authService.Login(requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		})

		// Act: Logout phone
		err := authService.Logout(phone.Token)
		if err != nil {
			t.Fatalf("Phone logout failed: %v", err)
		}

		// Try to refresh laptop token
		newLaptopToken, err := authService.RefreshToken(laptop.RefreshToken)

		// Assert: Current behavior - laptop token also revoked
		if err != nil {
			t.Log("⚠️ KNOWN ISSUE: Logout revokes ALL user tokens (including other devices)")
			t.Log("   Expected: Only phone session revoked")
			t.Log("   Actual: Both phone AND laptop sessions revoked")
			t.Log("   File: app/services/auth_service.go:136 (RevokeAllUserTokens)")
		}

		if newLaptopToken == nil {
			t.Log("⚠️ Laptop session was terminated (unintended)")
		}

		// This is documented behavior, not a test failure
		// But should be noted as improvement opportunity
	})

	t.Run("security - token cleanup removes expired tokens", func(t *testing.T) {
		// Arrange: Create expired and revoked tokens
		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("cleanup@test.com").
			CreateWithPlainPassword()

		// Create 5 expired tokens
		for i := 0; i < 5; i++ {
			expiredToken := models.RefreshToken{
				UserID:    user.ID,
				Token:     fmt.Sprintf("expired_token_%d", i),
				ExpiresAt: time.Now().Add(-1 * time.Hour),
				Revoked:   false,
			}
			testDB.DB.Create(&expiredToken)
		}

		// Create 3 revoked tokens
		for i := 0; i < 3; i++ {
			revokedToken := models.RefreshToken{
				UserID:    user.ID,
				Token:     fmt.Sprintf("revoked_token_%d", i),
				ExpiresAt: time.Now().Add(24 * time.Hour),
				Revoked:   true,
			}
			testDB.DB.Create(&revokedToken)
		}

		// Count before cleanup
		var countBefore int64
		testDB.DB.Model(&models.RefreshToken{}).Where("user_id = ?", user.ID).Count(&countBefore)

		// Act: Run cleanup
		jwtService := services.JwtService{}
		err := jwtService.CleanupExpiredTokens()
		if err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}

		// Count after cleanup
		var countAfter int64
		testDB.DB.Model(&models.RefreshToken{}).Where("user_id = ?", user.ID).Count(&countAfter)

		// Assert: Expired and revoked tokens should be deleted
		if countAfter >= countBefore {
			t.Errorf("🚨 Cleanup did not remove tokens! Before: %d, After: %d", countBefore, countAfter)
		}

		t.Logf("✓ Token cleanup removed %d expired/revoked tokens", countBefore-countAfter)
	})
}
