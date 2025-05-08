package services

import (
	"path/filepath"

	"golang_starter_kit_2025/app/helpers"

	"github.com/gin-gonic/gin"
)

type FileService struct{}

func (service FileService) UploadFile(ctx *gin.Context, key string, path string) (*string, error) {
	file, err := ctx.FormFile(key)
	if err != nil {
		return nil, err
	}

	// rename file with uuid as filename with extension
	fileName := helpers.GenerateReference(key) + filepath.Ext(file.Filename)

	// save file
	if err := ctx.SaveUploadedFile(file, helpers.StoragePath()+filepath.Join(path, fileName)); err != nil {
		return nil, err
	}

	return &fileName, nil
}

func (service FileService) StoreBase64File(base64 string, key string, path string) (*string, error) {
	fileName := helpers.GenerateReference(key) + ".png"

	if err := helpers.StoreBase64File(base64, path, fileName); err != nil {
		return nil, err
	}

	return &fileName, nil
}
