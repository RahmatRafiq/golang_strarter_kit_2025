package controllers

import (
	"net/http"
	"time"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	folder := ctx.DefaultPostForm("folder", "general")

	filename, err := c.storageService.UploadFile(file, header, folder)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upload file")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	url, _ := c.storageService.GetFileURL(filename, 24*time.Hour)

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": filename,
		"url":      url,
	})
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": url})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
