package controllers

import (
	"net/http"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type StorageController struct {
	storageService *services.StorageService
}

func NewStorageController(storageService *services.StorageService) *StorageController {
	return &StorageController{storageService: storageService}
}

// UploadFile godoc
// @Summary      Upload file to storage
// @Description  Upload a file to S3/MinIO storage
// @Tags         Storage
// @Accept       multipart/form-data
// @Produce      json
// @Param        file    formData  file    true  "File to upload"
// @Param        folder  formData  string  false "Upload folder" default(general)
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /api/v1/storage/upload [post]
func (c *StorageController) UploadFile(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message: "No file uploaded",
			Errors:  map[string]string{"file": "File is required"},
		}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	folder := ctx.DefaultPostForm("folder", "general")

	filename, err := c.storageService.UploadFile(file, header, folder)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload file")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message: "Upload failed",
			Errors:  map[string]string{"upload": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	url, _ := c.storageService.GetFileURL(filename, 24*time.Hour)

	item := any(map[string]interface{}{"filename": filename, "url": url})
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{
		Message: "File uploaded successfully",
		Item:    &item,
	}, http.StatusOK)
}

// GetFileURL godoc
// @Summary      Get presigned file URL
// @Description  Get a temporary presigned URL for accessing a file
// @Tags         Storage
// @Produce      json
// @Param        filename  path      string  true  "Filename"
// @Success      200       {object}  map[string]interface{}
// @Failure      404       {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /api/v1/storage/url/{filename} [get]
func (c *StorageController) GetFileURL(ctx *gin.Context) {
	filename := ctx.Param("filename")
	expiry := 24 * time.Hour

	url, err := c.storageService.GetFileURL(filename, expiry)
	if err != nil {
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message: "File not found",
			Errors:  map[string]string{"filename": "File does not exist"},
		}, http.StatusNotFound)
		return
	}

	item := any(map[string]interface{}{"url": url, "filename": filename})
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{
		Message: "File URL generated successfully",
		Item:    &item,
	}, http.StatusOK)
}

// DeleteFile godoc
// @Summary      Delete file from storage
// @Description  Delete a file from S3/MinIO storage
// @Tags         Storage
// @Produce      json
// @Param        filename  path      string  true  "Filename"
// @Success      200       {object}  map[string]interface{}
// @Failure      500       {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /api/v1/storage/{filename} [delete]
func (c *StorageController) DeleteFile(ctx *gin.Context) {
	filename := ctx.Param("filename")

	if err := c.storageService.DeleteFile(filename); err != nil {
		log.Error().Err(err).Msg("Failed to delete file")
		helpers.ResponseError(ctx, &helpers.ResponseParams[any]{
			Message: "Delete failed",
			Errors:  map[string]string{"delete": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	item := any(map[string]interface{}{"filename": filename})
	helpers.ResponseSuccess(ctx, &helpers.ResponseParams[any]{
		Message: "File deleted successfully",
		Item:    &item,
	}, http.StatusOK)
}
