package config

type StorageConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Region          string
	UseSSL          bool
}

func GetStorageConfig() *StorageConfig {
	return &StorageConfig{
		Endpoint:        getEnv("S3_ENDPOINT", "localhost:9000"),
		AccessKeyID:     getEnv("S3_ACCESS_KEY", "minioadmin"),
		SecretAccessKey: getEnv("S3_SECRET_KEY", "minioadmin"),
		UseSSL:          getEnv("S3_USE_SSL", "false") == "true",
		BucketName:      getEnv("S3_BUCKET", "uploads"),
		Region:          getEnv("S3_REGION", "us-east-1"),
	}
}
