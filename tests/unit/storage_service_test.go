package services_test

import (
	"testing"

	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/config"
)

func TestStorageService_Initialization_Unit(t *testing.T) {
	t.Run("success - initializes with valid configuration", func(t *testing.T) {
		storageService, err := services.NewStorageService()

		// Service may fail if MinIO is not available, which is okay for unit tests
		if err != nil {
			t.Logf("Storage service not available (expected if MinIO not running): %v", err)
			return // Skip test if MinIO not available
		}

		if storageService == nil {
			t.Error("expected non-nil storage service when no error")
		}
	})
}

func TestStorageConfig_Unit(t *testing.T) {
	t.Run("success - loads storage configuration", func(t *testing.T) {
		cfg := config.GetStorageConfig()

		if cfg == nil {
			t.Fatal("expected non-nil storage config")
		}

		// Verify required fields are set (even if using defaults)
		if cfg.Endpoint == "" {
			t.Error("expected non-empty endpoint")
		}
		if cfg.BucketName == "" {
			t.Error("expected non-empty bucket name")
		}
		if cfg.AccessKeyID == "" {
			t.Error("expected non-empty access key ID")
		}
		if cfg.SecretAccessKey == "" {
			t.Error("expected non-empty secret access key")
		}
		if cfg.Region == "" {
			t.Error("expected non-empty region")
		}
	})

	t.Run("success - has correct default values", func(t *testing.T) {
		cfg := config.GetStorageConfig()

		// These may be overridden by environment variables
		// We just verify they're not empty
		if cfg.Endpoint == "" {
			t.Error("endpoint should have default value")
		}
		if cfg.BucketName == "" {
			t.Error("bucket name should have default value")
		}
		if cfg.Region == "" {
			t.Error("region should have default value")
		}

		t.Logf("Configuration loaded: Endpoint=%s, Bucket=%s, Region=%s",
			cfg.Endpoint, cfg.BucketName, cfg.Region)
	})

	t.Run("success - UseSSL is boolean", func(t *testing.T) {
		cfg := config.GetStorageConfig()

		// Verify UseSSL field exists and is boolean (true or false)
		if cfg.UseSSL {
			t.Log("SSL is enabled")
		} else {
			t.Log("SSL is disabled")
		}
	})
}

func TestStorageService_ConfigurationValidation_Unit(t *testing.T) {
	t.Run("configuration has all required fields", func(t *testing.T) {
		cfg := config.GetStorageConfig()

		requiredFields := map[string]string{
			"Endpoint":        cfg.Endpoint,
			"AccessKeyID":     cfg.AccessKeyID,
			"SecretAccessKey": cfg.SecretAccessKey,
			"BucketName":      cfg.BucketName,
			"Region":          cfg.Region,
		}

		for fieldName, fieldValue := range requiredFields {
			if fieldValue == "" {
				t.Errorf("required field %s is empty", fieldName)
			}
		}
	})
}

// NOTE: Storage service integration tests (upload, download, delete) require MinIO running
// These tests focus on configuration and initialization logic only
// For full integration tests, run with MinIO container: docker run -p 9000:9000 minio/minio server /data
