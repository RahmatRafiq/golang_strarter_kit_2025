package integration

import (
	"testing"

	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/config"
)

// TestStorageService_Integration tests storage service initialization
// Full file upload/download tests require MinIO running:
// docker run -p 9000:9000 -e "MINIO_ROOT_USER=minioadmin" -e "MINIO_ROOT_PASSWORD=minioadmin" minio/minio server /data

func TestStorageService_Initialization_Integration(t *testing.T) {
	t.Run("success - service initializes with config", func(t *testing.T) {
		storageService, err := services.NewStorageService()

		// Service may fail if MinIO not available - this is expected in most test environments
		if err != nil {
			t.Logf("Storage service unavailable (expected without MinIO): %v", err)
			t.Skip("MinIO not running - skipping integration test")
			return
		}

		if storageService == nil {
			t.Error("Expected non-nil service when no error")
		}

		t.Log("✓ Storage service initialized successfully")
	})
}

func TestStorageConfig_Integration(t *testing.T) {
	t.Run("success - config loads with all required fields", func(t *testing.T) {
		cfg := config.GetStorageConfig()

		if cfg == nil {
			t.Fatal("Expected non-nil config")
		}

		// Verify all required fields are present
		requiredFields := map[string]string{
			"Endpoint":        cfg.Endpoint,
			"AccessKeyID":     cfg.AccessKeyID,
			"SecretAccessKey": cfg.SecretAccessKey,
			"BucketName":      cfg.BucketName,
			"Region":          cfg.Region,
		}

		for fieldName, fieldValue := range requiredFields {
			if fieldValue == "" {
				t.Errorf("Required field %s is empty", fieldName)
			}
		}

		t.Logf("✓ Storage config loaded: Endpoint=%s, Bucket=%s, Region=%s, UseSSL=%v",
			cfg.Endpoint, cfg.BucketName, cfg.Region, cfg.UseSSL)
	})

	t.Run("success - default values are sensible", func(t *testing.T) {
		cfg := config.GetStorageConfig()

		// Verify defaults (may be overridden by env)
		if cfg.Endpoint == "" {
			t.Error("Endpoint should have default value")
		}
		if cfg.BucketName == "" {
			t.Error("BucketName should have default value")
		}
		if cfg.Region == "" {
			t.Error("Region should have default value")
		}

		// Log actual values for verification
		t.Logf("Configuration values:")
		t.Logf("  Endpoint: %s", cfg.Endpoint)
		t.Logf("  Bucket: %s", cfg.BucketName)
		t.Logf("  Region: %s", cfg.Region)
		t.Logf("  UseSSL: %v", cfg.UseSSL)
	})
}

func TestStorageService_BucketConfiguration_Integration(t *testing.T) {
	storageService, err := services.NewStorageService()
	if err != nil {
		t.Skipf("Skipping bucket test - MinIO not available: %v", err)
	}

	t.Run("success - service is ready for operations", func(t *testing.T) {
		if storageService == nil {
			t.Error("Service should not be nil when initialization succeeds")
		}

		t.Log("✓ Storage service ready for file operations")
		t.Log("  Note: Actual file operations (upload/download/delete) require MinIO running")
	})
}

// NOTE: Full integration tests for file operations (upload/download/delete) require:
// 1. MinIO server running on localhost:9000 
// 2. Valid credentials (minioadmin/minioadmin by default)
// 3. Network connectivity
//
// To run full integration tests:
// 1. Start MinIO: docker run -p 9000:9000 -e "MINIO_ROOT_USER=minioadmin" -e "MINIO_ROOT_PASSWORD=minioadmin" minio/minio server /data
// 2. Run tests: go test ./tests/integration/storage_test.go -v
//
// These tests focus on configuration and initialization which don't require external services
