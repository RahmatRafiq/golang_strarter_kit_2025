package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"golang_starter_kit_2025/app/controllers"
	"golang_starter_kit_2025/app/repositories"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/facades"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OAuthController", Ordered, func() {
	var (
		router     *gin.Engine
		controller *controllers.OAuthController
	)

	BeforeAll(func() {
		// Setup test database
		wd, err := os.Getwd()
		if err != nil {
			Skip(fmt.Sprintf("Error getting working directory: %v", err))
		}

		envPath := fmt.Sprintf("%s/../../.env.test", wd)
		if err := godotenv.Load(envPath); err != nil {
			Skip("Error loading .env.test file")
		}

		facades.ConnectDB()
	})

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()

		// Setup dependencies
		userRepo := repositories.NewUserRepository(facades.DB)
		oauthService := services.NewOAuthService(userRepo)
		controller = controllers.NewOAuthController(oauthService)
	})

	AfterAll(func() {
		facades.CloseDB()
	})

	Describe("GoogleLogin", func() {
		It("should redirect to Google OAuth URL", func() {
			router.GET("/auth/google", controller.GoogleLogin)

			req, _ := http.NewRequest("GET", "/auth/google", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusTemporaryRedirect))
			Expect(resp.Header().Get("Location")).To(ContainSubstring("accounts.google.com"))
		})

		It("should set oauth_state cookie", func() {
			router.GET("/auth/google", controller.GoogleLogin)

			req, _ := http.NewRequest("GET", "/auth/google", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			cookies := resp.Result().Cookies()
			var stateCookie *http.Cookie
			for _, cookie := range cookies {
				if cookie.Name == "oauth_state" {
					stateCookie = cookie
					break
				}
			}

			Expect(stateCookie).NotTo(BeNil())
			Expect(stateCookie.Value).NotTo(BeEmpty())
			Expect(stateCookie.MaxAge).To(Equal(600))
			Expect(stateCookie.HttpOnly).To(BeTrue())
		})
	})

	Describe("GitHubLogin", func() {
		It("should redirect to GitHub OAuth URL", func() {
			router.GET("/auth/github", controller.GitHubLogin)

			req, _ := http.NewRequest("GET", "/auth/github", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			Expect(resp.Code).To(Equal(http.StatusTemporaryRedirect))
			Expect(resp.Header().Get("Location")).To(ContainSubstring("github.com"))
		})

		It("should set oauth_state cookie", func() {
			router.GET("/auth/github", controller.GitHubLogin)

			req, _ := http.NewRequest("GET", "/auth/github", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			cookies := resp.Result().Cookies()
			var stateCookie *http.Cookie
			for _, cookie := range cookies {
				if cookie.Name == "oauth_state" {
					stateCookie = cookie
					break
				}
			}

			Expect(stateCookie).NotTo(BeNil())
			Expect(stateCookie.Value).NotTo(BeEmpty())
		})
	})

	Describe("GoogleCallback", func() {
		Context("with invalid state", func() {
			It("should return error", func() {
				router.GET("/auth/google/callback", controller.GoogleCallback)

				req, _ := http.NewRequest("GET", "/auth/google/callback?state=invalid&code=test", nil)
				req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "different"})
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

				var result map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &result)

				Expect(result["status"]).To(Equal("error"))
				Expect(result["message"]).To(Equal("Invalid state parameter"))
			})
		})

		Context("with missing state", func() {
			It("should return error", func() {
				router.GET("/auth/google/callback", controller.GoogleCallback)

				req, _ := http.NewRequest("GET", "/auth/google/callback", nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("GitHubCallback", func() {
		Context("with invalid state", func() {
			It("should return error", func() {
				router.GET("/auth/github/callback", controller.GitHubCallback)

				req, _ := http.NewRequest("GET", "/auth/github/callback?state=invalid&code=test", nil)
				req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "different"})
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

				var result map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &result)

				Expect(result["status"]).To(Equal("error"))
				Expect(result["message"]).To(Equal("Invalid state parameter"))
			})
		})
	})
})

var _ = Describe("OAuthService Integration", Ordered, func() {
	var oauthService *services.OAuthService

	BeforeAll(func() {
		wd, err := os.Getwd()
		if err != nil {
			Skip(fmt.Sprintf("Error getting working directory: %v", err))
		}

		envPath := fmt.Sprintf("%s/../../.env.test", wd)
		if err := godotenv.Load(envPath); err != nil {
			Skip("Error loading .env.test file")
		}

		facades.ConnectDB()

		userRepo := repositories.NewUserRepository(facades.DB)
		oauthService = services.NewOAuthService(userRepo)
	})

	AfterAll(func() {
		facades.CloseDB()
	})

	Describe("GetGoogleAuthURL", func() {
		It("should generate valid Google OAuth URL", func() {
			state := "test-state-123"
			url := oauthService.GetGoogleAuthURL(state)

			Expect(url).To(ContainSubstring("accounts.google.com/o/oauth2/v2/auth"))
			Expect(url).To(ContainSubstring("state=" + state))
			Expect(url).To(ContainSubstring("client_id="))
		})
	})

	Describe("GetGitHubAuthURL", func() {
		It("should generate valid GitHub OAuth URL", func() {
			state := "test-state-456"
			url := oauthService.GetGitHubAuthURL(state)

			Expect(url).To(ContainSubstring("github.com/login/oauth/authorize"))
			Expect(url).To(ContainSubstring("state=" + state))
			Expect(url).To(ContainSubstring("client_id="))
		})
	})
})
