package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type OAuthController struct {
	oauthService *services.OAuthService
}

func NewOAuthController(oauthService *services.OAuthService) *OAuthController {
	return &OAuthController{oauthService: oauthService}
}

// GoogleLogin godoc
// @Summary      Login with Google OAuth2
// @Description  Redirects user to Google OAuth2 consent page
// @Tags         OAuth
// @Produce      json
// @Success      307  {string}  string  "Redirect to Google OAuth"
// @Router       /api/v1/auth/google [get]
func (c *OAuthController) GoogleLogin(ctx *gin.Context) {
	state := generateState()
	ctx.SetCookie("oauth_state", state, 600, "/", "", false, true)
	ctx.Redirect(http.StatusTemporaryRedirect, c.oauthService.GetGoogleAuthURL(state))
}

func (c *OAuthController) GoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	cookie, _ := ctx.Cookie("oauth_state")

	if state != cookie {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message: "Invalid state parameter",
			Errors:  map[string]string{"state": "CSRF validation failed"},
		}, http.StatusBadRequest)
		return
	}

	code := ctx.Query("code")
	user, token, err := c.oauthService.HandleGoogleCallback(code)
	if err != nil {
		log.Error().Err(err).Msg("Google OAuth failed")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message: "OAuth authentication failed",
			Errors:  map[string]string{"oauth": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	item := any(map[string]interface{}{"user": user, "token": token})
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{
		Message: "Login successful",
		Item:    &item,
	}, http.StatusOK)
}

// GitHubLogin godoc
// @Summary      Login with GitHub OAuth2
// @Description  Redirects user to GitHub OAuth2 authorization page
// @Tags         OAuth
// @Produce      json
// @Success      307  {string}  string  "Redirect to GitHub OAuth"
// @Router       /api/v1/auth/github [get]
func (c *OAuthController) GitHubLogin(ctx *gin.Context) {
	state := generateState()
	ctx.SetCookie("oauth_state", state, 600, "/", "", false, true)
	ctx.Redirect(http.StatusTemporaryRedirect, c.oauthService.GetGitHubAuthURL(state))
}

func (c *OAuthController) GitHubCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	cookie, _ := ctx.Cookie("oauth_state")

	if state != cookie {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message: "Invalid state parameter",
			Errors:  map[string]string{"state": "CSRF validation failed"},
		}, http.StatusBadRequest)
		return
	}

	code := ctx.Query("code")
	user, token, err := c.oauthService.HandleGitHubCallback(code)
	if err != nil {
		log.Error().Err(err).Msg("GitHub OAuth failed")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message: "OAuth authentication failed",
			Errors:  map[string]string{"oauth": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	item := any(map[string]interface{}{"user": user, "token": token})
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{
		Message: "Login successful",
		Item:    &item,
	}, http.StatusOK)
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
