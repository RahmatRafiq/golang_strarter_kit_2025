package controllers

import (
	"net/http"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type FileController struct {
	jwtService services.JwtService
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
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "File not found",
			Reference: "ERROR-7",
		}, 404)
		return
	}

	var jwtKey = []byte(helpers.GetEnv("APP_KEY", "your_secret_key"))
	token, err := jwt.Parse(signature, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "File not found",
			Reference: "ERROR-8",
		}, 400)
		return
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "File not found",
			Reference: "ERROR-8",
		}, 400)
		return
	}

	expiredAt := int64(tokenClaims["expired_at"].(float64))
	if expiredAt < time.Now().Unix() {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "File not found",
			Reference: "ERROR-9",
		}, 400)
		return
	}

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
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message:   "File not found",
			Reference: "ERROR-7",
		}, http.StatusNotFound)
		return
	}

	// Menyajikan file tanpa autentikasi
	filePath := "storage/" + key + "/" + filename
	ctx.File(filePath)
}
