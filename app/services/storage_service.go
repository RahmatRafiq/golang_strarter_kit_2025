package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"golang_starter_kit_2025/config"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

type StorageService struct {
	client *minio.Client
	config *config.StorageConfig
}

func NewStorageService() (*StorageService, error) {
	cfg := config.GetStorageConfig()

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// FIXED: Add timeout to prevent goroutine leaks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to check bucket existence")
	} else if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{Region: cfg.Region})
		if err != nil {
			log.Warn().Err(err).Msg("Failed to create bucket")
		} else {
			log.Info().Str("bucket", cfg.BucketName).Msg("Bucket created")
		}
	}

	return &StorageService{
		client: client,
		config: cfg,
	}, nil
}

func (s *StorageService) UploadFile(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// FIXED: Add timeout to prevent goroutine leaks
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := s.client.PutObject(
		ctx,
		s.config.BucketName,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: header.Header.Get("Content-Type")},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return filename, nil
}

func (s *StorageService) GetFileURL(filename string, expiry time.Duration) (string, error) {
	// FIXED: Add timeout to prevent goroutine leaks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url, err := s.client.PresignedGetObject(
		ctx,
		s.config.BucketName,
		filename,
		expiry,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

func (s *StorageService) DeleteFile(filename string) error {
	// FIXED: Add timeout to prevent goroutine leaks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.client.RemoveObject(
		ctx,
		s.config.BucketName,
		filename,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *StorageService) DownloadFile(filename string, writer io.Writer) error {
	// FIXED: Add timeout to prevent goroutine leaks
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	object, err := s.client.GetObject(
		ctx,
		s.config.BucketName,
		filename,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()

	_, err = io.Copy(writer, object)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	return nil
}
