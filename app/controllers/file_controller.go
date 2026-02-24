package controllers

import (
	"net/http"
	"time"

	"golang_starter_kit_2025/app/helpers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

type FileController struct {
}

func NewFileController() *FileController {
	return &FileController{}
}

// @Summary		Serve file
// @Description	Serve file
// @Tags			File
// @Accept			json
// @Produce		jpeg
// @Param			signature	query		string	false	"Signature"
// @Param			key			path		string	true	"File key"
// @Param			path		path		string	true	"File path"
// @Success		200			{string}	string	"File"
// @Router			/file/{key}/{filename} [get]
func (controller FileController) ServeFile(ctx *gin.Context) {
	key := ctx.Param("key")
	filename := ctx.Param("filename")
	signature := ctx.Query("signature")

	if key == "" || filename == "" || signature == "" {
		log.Warn().
			Str("key", key).
			Str("filename", filename).
			Bool("has_signature", signature != "").
			Msg("File serve request missing parameters")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Missing required parameters"},
			Message:   "File not found",
			Reference: "ERROR-2",
		}, http.StatusNotFound)
		return
	}

	var jwtKey = []byte(helpers.GetEnv("APP_KEY", "your_secret_key"))
	token, err := jwt.Parse(signature, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("filename", filename).
			Msg("Invalid file signature")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Invalid signature"},
			Message:   "File not found",
			Reference: "ERROR-4",
		}, http.StatusBadRequest)
		return
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Warn().
			Str("filename", filename).
			Bool("claims_ok", ok).
			Bool("token_valid", token.Valid).
			Msg("Invalid token claims or token not valid")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Invalid token"},
			Message:   "File not found",
			Reference: "ERROR-4",
		}, http.StatusBadRequest)
		return
	}

	expiredAtFloat, ok := tokenClaims["expired_at"].(float64)
	if !ok {
		log.Warn().
			Str("filename", filename).
			Msg("Invalid token format - expired_at field missing or wrong type")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Invalid token format"},
			Message:   "Invalid token format",
			Reference: "ERROR-4",
		}, http.StatusBadRequest)
		return
	}

	expiredAt := int64(expiredAtFloat)
	if expiredAt < time.Now().Unix() {
		log.Warn().
			Str("filename", filename).
			Int64("expired_at", expiredAt).
			Int64("current_time", time.Now().Unix()).
			Msg("File access token expired")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Token expired"},
			Message:   "File not found",
			Reference: "ERROR-4",
		}, http.StatusBadRequest)
		return
	}

	log.Info().
		Str("filename", filename).
		Str("key", key).
		Msg("File served successfully")
	ctx.File("storage/" + key + "/" + filename)
}

// @Summary		Serve file without authentication
// @Description	Serve file directly without JWT authentication
// @Tags			File
// @Accept			json
// @Produce		jpeg
// @Param			key			path		string				true	"File key"
// @Param			filename	path		string				true	"File name"
// @Success		200			{file}		string				"File"
// @Failure		404			{object}	map[string]string	"File not found"
// @Router			/file/public/{key}/{filename} [get]
func (controller FileController) ServePublicFile(ctx *gin.Context) {
	key := ctx.Param("key")
	filename := ctx.Param("filename")

	// Validasi parameter
	if key == "" || filename == "" {
		log.Warn().
			Str("key", key).
			Str("filename", filename).
			Msg("Public file serve request missing parameters")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Errors:    map[string]string{"error": "Missing required parameters"},
			Message:   "File not found",
			Reference: "ERROR-2",
		}, http.StatusNotFound)
		return
	}

	// Menyajikan file tanpa autentikasi
	filePath := "storage/" + key + "/" + filename
	log.Info().
		Str("filename", filename).
		Str("key", key).
		Str("file_path", filePath).
		Msg("Public file served successfully")
	ctx.File(filePath)
}
