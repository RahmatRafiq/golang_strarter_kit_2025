package services_test

import (
	"testing"

	"golang_starter_kit_2025/app/mocks"
	"golang_starter_kit_2025/app/services"

	"go.uber.org/mock/gomock"
)

func TestOAuthService_GetGoogleAuthURL_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	oauthService := services.NewOAuthService(mockRepo)

	t.Run("success - generates valid Google OAuth URL", func(t *testing.T) {
		state := "test-state-123"
		url := oauthService.GetGoogleAuthURL(state)

		if url == "" {
			t.Error("expected non-empty URL, got empty string")
		}

		// Verify URL contains required OAuth components
		requiredComponents := []string{
			"accounts.google.com",
			"state=" + state,
			"response_type=code",
			"scope=",
		}

		for _, component := range requiredComponents {
			if !contains(url, component) {
				t.Errorf("URL missing required component: %s", component)
			}
		}
	})

	t.Run("success - generates different URLs for different states", func(t *testing.T) {
		state1 := "state-abc-123"
		state2 := "state-xyz-789"

		url1 := oauthService.GetGoogleAuthURL(state1)
		url2 := oauthService.GetGoogleAuthURL(state2)

		if !contains(url1, state1) {
			t.Error("URL 1 should contain state1")
		}
		if !contains(url2, state2) {
			t.Error("URL 2 should contain state2")
		}
		if url1 == url2 {
			t.Error("URLs should be different with different states")
		}
	})
}

func TestOAuthService_GetGitHubAuthURL_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	oauthService := services.NewOAuthService(mockRepo)

	t.Run("success - generates valid GitHub OAuth URL", func(t *testing.T) {
		state := "test-state-456"
		url := oauthService.GetGitHubAuthURL(state)

		if url == "" {
			t.Error("expected non-empty URL, got empty string")
		}

		// Verify URL contains required OAuth components
		requiredComponents := []string{
			"github.com",
			"state=" + state,
			"scope=",
		}

		for _, component := range requiredComponents {
			if !contains(url, component) {
				t.Errorf("URL missing required component: %s", component)
			}
		}
	})

	t.Run("success - generates different URLs for different states", func(t *testing.T) {
		state1 := "state-def-456"
		state2 := "state-uvw-012"

		url1 := oauthService.GetGitHubAuthURL(state1)
		url2 := oauthService.GetGitHubAuthURL(state2)

		if !contains(url1, state1) {
			t.Error("URL 1 should contain state1")
		}
		if !contains(url2, state2) {
			t.Error("URL 2 should contain state2")
		}
		if url1 == url2 {
			t.Error("URLs should be different with different states")
		}
	})
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
