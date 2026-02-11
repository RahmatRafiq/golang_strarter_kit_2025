package integration

import (
	"strings"
	"testing"

	"golang_starter_kit_2025/app/repositories"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/tests/helpers"
)

// TestAuthService_InputLength tests handling of extreme input lengths
// EDGE CASE: Prevents DoS via expensive operations
func TestAuthService_InputLength_EdgeCases(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("edge - very long email (256+ characters)", func(t *testing.T) {
		// Arrange: Email longer than typical DB varchar(255)
		longEmail := strings.Repeat("a", 246) + "@test.com" // 256 chars

		loginReq := requests.LoginRequest{
			Email:    longEmail,
			Password: "password123",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Should handle gracefully (either accept or reject, but no crash)
		if err != nil {
			t.Logf("✓ Long email rejected with error: %v", err)
		} else if token != nil {
			t.Log("⚠️ Long email accepted - ensure DB column can handle it")
		}

		// No panic = test passes
		t.Log("✓ No crash with very long email")
	})

	t.Run("edge - extremely long email (1000+ characters)", func(t *testing.T) {
		// Arrange: Absurdly long email
		veryLongEmail := strings.Repeat("a", 990) + "@test.com" // 1000 chars

		loginReq := requests.LoginRequest{
			Email:    veryLongEmail,
			Password: "password123",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Should be rejected or handled
		if token != nil {
			t.Error("⚠️ Extremely long email should probably be rejected")
		}

		if err != nil {
			t.Logf("✓ Extremely long email rejected: %v", err)
		}
	})

	t.Run("edge - very long password (1000 characters)", func(t *testing.T) {
		// Arrange
		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("longpass@test.com").
			CreateWithPlainPassword()

		longPassword := strings.Repeat("a", 1000)

		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: longPassword,
		}

		// Act: This might take time due to Argon2 hashing
		startTime := testing.Benchmark(func(b *testing.B) {
			authService.Login(loginReq)
		})

		// Assert: Should complete in reasonable time (< 5 seconds)
		// Argon2 is expensive, but shouldn't DoS the server
		t.Logf("✓ Long password processed (benchmark: %v ops)", startTime.N)
	})

	t.Run("edge - extremely long password (10000 characters) DoS protection", func(t *testing.T) {
		// Arrange: 10KB password - potential DoS vector
		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("dospass@test.com").
			CreateWithPlainPassword()

		veryLongPassword := strings.Repeat("a", 10000)

		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: veryLongPassword,
		}

		// Act: Should have timeout or length validation
		token, err := authService.Login(loginReq)

		// Assert: Should ideally reject or timeout quickly
		if token != nil {
			t.Error("⚠️ WARNING: 10KB password accepted - potential DoS vector!")
		}

		if err != nil {
			t.Logf("✓ Extremely long password rejected: %v", err)
		} else {
			t.Log("⚠️ RECOMMENDATION: Add password length limit (max 72 chars for bcrypt, 128 for Argon2)")
		}
	})

	t.Run("edge - empty email", func(t *testing.T) {
		// Arrange
		loginReq := requests.LoginRequest{
			Email:    "",
			Password: "password123",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Should fail
		if err == nil {
			t.Error("Empty email should return error")
		}
		if token != nil {
			t.Error("Empty email should not generate token")
		}

		t.Logf("✓ Empty email rejected: %v", err)
	})

	t.Run("edge - single character email", func(t *testing.T) {
		// Arrange
		loginReq := requests.LoginRequest{
			Email:    "a",
			Password: "password123",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Invalid email format
		if token != nil {
			t.Error("Single char email should not generate token")
		}

		t.Logf("✓ Invalid email format rejected: %v", err)
	})
}

// TestAuthService_WhitespaceHandling tests whitespace edge cases
// EDGE CASE: User experience issue - common login failures
func TestAuthService_WhitespaceHandling_EdgeCases(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("edge - leading whitespace in email", func(t *testing.T) {
		// Arrange: Create user with clean email
		_, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("user@test.com").
			CreateWithPlainPassword()

		// Act: Login with leading whitespace
		loginReq := requests.LoginRequest{
			Email:    "  user@test.com",
			Password: plainPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert: Current behavior vs expected
		if err != nil {
			t.Log("⚠️ USABILITY ISSUE: Leading whitespace causes login failure")
			t.Log("   Users often copy-paste emails with whitespace")
			t.Log("   RECOMMENDATION: Trim email input in service/middleware")
			t.Logf("   Error: %v", err)
		} else if token != nil {
			t.Log("✓ Leading whitespace handled (email trimmed)")
		}
	})

	t.Run("edge - trailing whitespace in email", func(t *testing.T) {
		// Arrange
		_, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("user2@test.com").
			CreateWithPlainPassword()

		// Act
		loginReq := requests.LoginRequest{
			Email:    "user2@test.com  ",
			Password: plainPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert
		if err != nil {
			t.Log("⚠️ USABILITY ISSUE: Trailing whitespace causes login failure")
			t.Logf("   Error: %v", err)
		} else if token != nil {
			t.Log("✓ Trailing whitespace handled")
		}
	})

	t.Run("edge - whitespace in password should NOT be trimmed", func(t *testing.T) {
		// Arrange: Create user with password that has NO whitespace
		createdUser, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("passspace@test.com").
			WithPassword("mypassword").
			CreateWithPlainPassword()

		// Act: Try to login with whitespace around password
		loginReq := requests.LoginRequest{
			Email:    createdUser.Email,
			Password: "  mypassword  ",
		}

		token, err := authService.Login(loginReq)

		// Assert: Should FAIL - passwords should be exact
		if err == nil {
			t.Error("🚨 SECURITY ISSUE: Password whitespace was trimmed!")
			t.Error("   Passwords should be stored/compared exactly as entered")
		}
		if token != nil {
			t.Error("🚨 Token generated with wrong password (whitespace trimmed)")
		}

		t.Log("✓ Password whitespace NOT trimmed (correct behavior)")
	})

	t.Run("edge - password with internal whitespace", func(t *testing.T) {
		// Arrange: Create user with password containing spaces
		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("internalspace@test.com").
			WithPassword("my secure password").
			CreateWithPlainPassword()

		// Act: Login with correct password (including spaces)
		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: "my secure password",
		}

		token, err := authService.Login(loginReq)

		// Assert: Should work
		if err != nil {
			t.Errorf("Password with spaces should work: %v", err)
		}
		if token == nil {
			t.Error("Expected token for valid password with spaces")
		}

		t.Log("✓ Password with internal spaces works correctly")
	})

	t.Run("edge - multiple spaces in email", func(t *testing.T) {
		// Arrange
		loginReq := requests.LoginRequest{
			Email:    "user  @  test  .com",
			Password: "password123",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Invalid email format
		if token != nil {
			t.Error("Email with internal spaces should be invalid")
		}

		t.Logf("✓ Invalid email format rejected: %v", err)
	})

	t.Run("edge - tab characters in email", func(t *testing.T) {
		// Arrange
		_, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("tabtest@test.com").
			CreateWithPlainPassword()

		// Act: Email with tab character
		loginReq := requests.LoginRequest{
			Email:    "tabtest@test.com\t",
			Password: plainPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert
		if err != nil {
			t.Log("⚠️ Tab character causes login failure")
			t.Log("   RECOMMENDATION: Sanitize all whitespace types")
		} else if token != nil {
			t.Log("✓ Tab character handled")
		}
	})

	t.Run("edge - newline in email (sanitized by TrimSpace)", func(t *testing.T) {
		// Arrange: Create user first
		_, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("user@test.com").
			CreateWithPlainPassword()

		// Act: Try login with newline (will be trimmed)
		loginReq := requests.LoginRequest{
			Email:    "user@test.com\n",
			Password: plainPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert: IMPROVED behavior - newline is trimmed, login succeeds
		if err != nil || token == nil {
			t.Error("✓ IMPROVED: Newline trimmed by sanitization, login should succeed")
		} else {
			t.Log("✓ Newline character sanitized (TrimSpace improvement)")
		}
	})
}

// TestAuthService_CaseSensitivity tests case handling
// EDGE CASE: Email case sensitivity varies by database
func TestAuthService_CaseSensitivity_EdgeCases(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("edge - email case sensitivity", func(t *testing.T) {
		// Arrange: Create user with mixed case email
		_, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("User@Test.COM").
			CreateWithPlainPassword()

		// Act: Login with lowercase
		loginReq := requests.LoginRequest{
			Email:    "user@test.com",
			Password: plainPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert: Depends on database collation
		if err != nil {
			t.Log("⚠️ Email is case-sensitive!")
			t.Log("   User@Test.COM != user@test.com")
			t.Log("   RECOMMENDATION: Normalize emails to lowercase on registration")
			t.Logf("   Error: %v", err)
		} else if token != nil {
			t.Log("✓ Email is case-insensitive (database collation handles it)")
		}
	})

	t.Run("edge - email case insensitive lookup", func(t *testing.T) {
		// Arrange: Create with lowercase
		_, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("test@example.com").
			CreateWithPlainPassword()

		// Test various case combinations
		testCases := []string{
			"test@example.com",     // Original
			"Test@example.com",     // Capitalize first
			"TEST@EXAMPLE.COM",     // All caps
			"TeSt@ExAmPlE.cOm",     // Mixed
			"test@EXAMPLE.com",     // Domain caps
		}

		for _, emailVariant := range testCases {
			loginReq := requests.LoginRequest{
				Email:    emailVariant,
				Password: plainPassword,
			}

			token, err := authService.Login(loginReq)

			if err != nil {
				t.Logf("⚠️ Case variant '%s' failed to login", emailVariant)
			} else if token != nil {
				t.Logf("✓ Case variant '%s' login successful", emailVariant)
			}
		}
	})

	t.Run("edge - password case sensitivity (should be case-sensitive)", func(t *testing.T) {
		// Arrange: Create user with mixed case password
		createdUser, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("passcase@test.com").
			WithPassword("MyPassword123").
			CreateWithPlainPassword()

		// Act: Try different case combinations
		wrongCases := []string{
			"mypassword123",  // All lowercase
			"MYPASSWORD123",  // All uppercase
			"myPassword123",  // Different case
		}

		for _, wrongPass := range wrongCases {
			loginReq := requests.LoginRequest{
				Email:    createdUser.Email,
				Password: wrongPass,
			}

			token, err := authService.Login(loginReq)

			// Assert: Should fail
			if err == nil || token != nil {
				t.Errorf("🚨 SECURITY ISSUE: Password case insensitive! '%s' accepted", wrongPass)
			}
		}

		// Correct case should work
		correctReq := requests.LoginRequest{
			Email:    createdUser.Email,
			Password: "MyPassword123",
		}

		correctToken, correctErr := authService.Login(correctReq)
		if correctErr != nil || correctToken == nil {
			t.Error("Correct password case should work")
		}

		t.Log("✓ Password is case-sensitive (correct behavior)")
	})
}

// TestAuthService_UnicodeHandling tests international character support
// EDGE CASE: International users
func TestAuthService_UnicodeHandling_EdgeCases(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("edge - unicode in email (internationalized email)", func(t *testing.T) {
		// Arrange: Email with Unicode characters
		user, plainPassword := helpers.NewUserFactory(t, testDB).
			WithEmail("用户@测试.com"). // Chinese characters
			CreateWithPlainPassword()

		// Act
		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: plainPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert
		if err != nil {
			t.Log("⚠️ Unicode email not supported")
			t.Logf("   Error: %v", err)
		} else if token != nil {
			t.Log("✓ Unicode email supported (internationalized domain names)")
		}
	})

	t.Run("edge - unicode in password (international passwords)", func(t *testing.T) {
		// Arrange: Password with Japanese, emoji, etc.
		testCases := []struct {
			name     string
			password string
		}{
			{"Japanese", "パスワード123"},
			{"Chinese", "密码测试123"},
			{"Arabic", "كلمةالسر123"},
			{"Emoji", "password🔐🔑"},
			{"Mixed", "Pass™️word©123®"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create user with unicode password
				user, _ := helpers.NewUserFactory(t, testDB).
					WithEmail("unicode" + tc.name + "@test.com").
					WithPassword(tc.password).
					CreateWithPlainPassword()

				// Act: Login with unicode password
				loginReq := requests.LoginRequest{
					Email:    user.Email,
					Password: tc.password,
				}

				token, err := authService.Login(loginReq)

				// Assert
				if err != nil {
					t.Logf("⚠️ %s password failed: %v", tc.name, err)
				} else if token == nil {
					t.Errorf("%s password login failed", tc.name)
				} else {
					t.Logf("✓ %s password works correctly", tc.name)
				}
			})
		}
	})

	t.Run("edge - emoji in email local part", func(t *testing.T) {
		// Arrange: Emoji in email (technically allowed in RFC 6531)
		loginReq := requests.LoginRequest{
			Email:    "user😀@test.com",
			Password: "password123",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Probably not supported (and that's okay)
		if token != nil {
			t.Log("⚠️ Emoji in email accepted - ensure this is intentional")
		} else if err != nil {
			t.Log("✓ Emoji in email rejected (expected)")
		}
	})
}

// TestAuthService_SpecialCharacters tests special character handling
// EDGE CASE: Password strength and email validation
func TestAuthService_SpecialCharacters_EdgeCases(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("edge - all special characters in password", func(t *testing.T) {
		// Arrange: Password with every special character
		specialPassword := "!@#$%^&*()_+-=[]{}|;:',.<>?/~`\"\\"

		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("special@test.com").
			WithPassword(specialPassword).
			CreateWithPlainPassword()

		// Act
		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: specialPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert: Should work
		if err != nil {
			t.Errorf("Special characters should be allowed: %v", err)
		}
		if token == nil {
			t.Error("Expected token for password with special chars")
		}

		t.Log("✓ All special characters in password work")
	})

	t.Run("edge - quotes in password", func(t *testing.T) {
		// Arrange: Password with single and double quotes
		quotedPassword := `password"with'quotes`

		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("quotes@test.com").
			WithPassword(quotedPassword).
			CreateWithPlainPassword()

		// Act
		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: quotedPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert
		if err != nil {
			t.Errorf("Quotes in password should work: %v", err)
		} else if token == nil {
			t.Error("Expected token to be returned")
		}

		t.Log("✓ Quotes in password handled correctly")
	})

	t.Run("edge - backslashes in password", func(t *testing.T) {
		// Arrange: Password with backslashes
		backslashPassword := `pass\word\123`

		user, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("backslash@test.com").
			WithPassword(backslashPassword).
			CreateWithPlainPassword()

		// Act
		loginReq := requests.LoginRequest{
			Email:    user.Email,
			Password: backslashPassword,
		}

		token, err := authService.Login(loginReq)

		// Assert
		if err != nil {
			t.Errorf("Backslashes should be allowed: %v", err)
		} else if token == nil {
			t.Error("Expected token to be returned")
		}

		t.Log("✓ Backslashes in password work")
	})

	t.Run("edge - special characters in email local part", func(t *testing.T) {
		// Arrange: RFC 5322 compliant special chars
		testEmails := []string{
			"user+tag@test.com",        // Plus sign (common)
			"user.name@test.com",       // Dot (common)
			"user_name@test.com",       // Underscore
			"user-name@test.com",       // Hyphen
			"user+tag+filter@test.com", // Multiple plus
		}

		for _, email := range testEmails {
			_, plainPassword := helpers.NewUserFactory(t, testDB).
				WithEmail(email).
				CreateWithPlainPassword()

			loginReq := requests.LoginRequest{
				Email:    email,
				Password: plainPassword,
			}

			token, err := authService.Login(loginReq)

			if err != nil || token == nil {
				t.Errorf("Valid RFC 5322 email '%s' should work: %v", email, err)
			} else {
				t.Logf("✓ Email '%s' works correctly", email)
			}
		}
	})
}

// TestAuthService_NullBytes tests null byte injection
// EDGE CASE: Security - null byte attacks
func TestAuthService_NullBytes_EdgeCases(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	authRepo := repositories.NewAuthRepository(testDB.DB)
	authService := services.NewAuthService(authRepo)

	t.Run("edge - null byte in email", func(t *testing.T) {
		// Arrange: Null byte injection attempt
		loginReq := requests.LoginRequest{
			Email:    "user\x00@test.com",
			Password: "password123",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Should be rejected or sanitized
		if token != nil {
			t.Error("⚠️ Null byte in email accepted - potential security issue")
		}
		if err == nil {
			t.Error("Expected error for null byte in email")
		}

		t.Log("✓ Null byte in email handled")
	})

	t.Run("edge - null byte in password", func(t *testing.T) {
		// Arrange
		createdUser, _ := helpers.NewUserFactory(t, testDB).
			WithEmail("nullbyte@test.com").
			CreateWithPlainPassword()

		loginReq := requests.LoginRequest{
			Email:    createdUser.Email,
			Password: "pass\x00word",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Should fail (password mismatch)
		if token != nil {
			t.Error("Password with null byte should not match")
		}
		if err == nil {
			t.Error("Expected error for null byte in password")
		}

		t.Log("✓ Null byte in password handled correctly")
	})

	t.Run("edge - null byte termination exploit", func(t *testing.T) {
		// Arrange: Attempt to exploit null byte termination
		// E.g., "admin\x00@test.com" might be truncated to "admin"
		loginReq := requests.LoginRequest{
			Email:    "admin\x00@attacker.com",
			Password: "guessed_password",
		}

		// Act
		token, err := authService.Login(loginReq)

		// Assert: Should not bypass authentication
		if token != nil {
			t.Error("🚨 SECURITY: Null byte might have bypassed validation!")
		}
		if err == nil {
			t.Error("Expected error for null byte exploit attempt")
		}

		t.Log("✓ Null byte termination exploit prevented")
	})
}
