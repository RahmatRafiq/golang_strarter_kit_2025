package integration

import (
	"testing"

	"golang_starter_kit_2025/app/repositories"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/tests/helpers"
)

// TestOAuthService_Integration tests OAuth service with real database
// Note: OAuth callback tests require real OAuth credentials and token exchange
// These tests focus on URL generation and user management
func TestOAuthService_GoogleAuthURL_Integration(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	userRepo := repositories.NewUserRepository(testDB.DB)
	oauthService := services.NewOAuthService(userRepo)

	t.Run("success - generates valid Google OAuth URL with state", func(t *testing.T) {
		state := "integration-test-state-google-123"

		url := oauthService.GetGoogleAuthURL(state)

		if url == "" {
			t.Fatal("Expected non-empty URL")
		}

		// Verify URL structure
		requiredComponents := []string{
			"https://accounts.google.com",
			"oauth2",
			"state=" + state,
			"response_type=code",
			"client_id=",
			"redirect_uri=",
			"scope=",
		}

		for _, component := range requiredComponents {
			if !contains(url, component) {
				t.Errorf("URL missing required component: %s", component)
			}
		}

		// Verify Google-specific scopes
		if !contains(url, "userinfo.email") || !contains(url, "userinfo.profile") {
			t.Error("URL missing Google userinfo scopes")
		}

		t.Logf("✓ Generated Google OAuth URL: %s...", url[:100])
	})

	t.Run("success - different states generate different URLs", func(t *testing.T) {
		state1 := "state-session-1"
		state2 := "state-session-2"

		url1 := oauthService.GetGoogleAuthURL(state1)
		url2 := oauthService.GetGoogleAuthURL(state2)

		if url1 == url2 {
			t.Error("Expected different URLs for different states")
		}

		if !contains(url1, state1) || !contains(url2, state2) {
			t.Error("URLs must contain their respective states")
		}
	})
}

func TestOAuthService_GitHubAuthURL_Integration(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	userRepo := repositories.NewUserRepository(testDB.DB)
	oauthService := services.NewOAuthService(userRepo)

	t.Run("success - generates valid GitHub OAuth URL with state", func(t *testing.T) {
		state := "integration-test-state-github-456"

		url := oauthService.GetGitHubAuthURL(state)

		if url == "" {
			t.Fatal("Expected non-empty URL")
		}

		// Verify URL structure
		requiredComponents := []string{
			"https://github.com",
			"login/oauth/authorize",
			"state=" + state,
			"client_id=",
			"redirect_uri=",
			"scope=",
		}

		for _, component := range requiredComponents {
			if !contains(url, component) {
				t.Errorf("URL missing required component: %s", component)
			}
		}

		// Verify GitHub-specific scope (URL encoded)
		if !contains(url, "user") && !contains(url, "email") {
			t.Error("URL missing GitHub user/email scope")
		}

		t.Logf("✓ Generated GitHub OAuth URL: %s...", url[:100])
	})

	t.Run("success - consistent URL format across multiple generations", func(t *testing.T) {
		state := "consistent-state-test"

		url1 := oauthService.GetGitHubAuthURL(state)
		url2 := oauthService.GetGitHubAuthURL(state)

		// URLs should be identical for same state
		if url1 != url2 {
			t.Error("Expected consistent URLs for same state")
		}
	})
}

// NOTE: Testing HandleGoogleCallback and HandleGitHubCallback requires:
// 1. Valid OAuth client credentials (GOOGLE_CLIENT_ID, GITHUB_CLIENT_ID, etc.)
// 2. Real authorization code from OAuth provider
// 3. Network access to OAuth provider APIs
//
// These are end-to-end tests and should be run separately with:
// - Real OAuth app credentials configured
// - Manual authorization code retrieval
// - Network connectivity
//
// Example manual test:
// 1. Get auth URL: url := oauthService.GetGoogleAuthURL("test-state")
// 2. Open URL in browser, authorize
// 3. Copy authorization code from callback URL
// 4. Call: user, token, err := oauthService.HandleGoogleCallback(code)

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
