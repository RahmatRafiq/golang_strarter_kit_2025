package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConfig holds configuration for MongoDB connection
type MongoDBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// GetMongoDBConfig loads MongoDB configuration from environment variables
func GetMongoDBConfig() *MongoDBConfig {
	return &MongoDBConfig{
		Host:     getEnv("MONGO_HOST", "localhost"),
		Port:     getEnv("MONGO_PORT", "27017"),
		Username: getEnv("MONGO_USERNAME", "root"),
		Password: getEnv("MONGO_PASSWORD", "password"),
		Database: getEnv("MONGO_DB", "admin"),
	}
}

// BuildMongoURI builds the MongoDB URI string
func (cfg *MongoDBConfig) BuildMongoURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

// Validate checks if the MongoDB configuration is valid
func (cfg *MongoDBConfig) Validate() error {
	if cfg.Host == "" {
		return fmt.Errorf("MongoDB host is required")
	}
	if cfg.Port == "" {
		return fmt.Errorf("MongoDB port is required")
	}
	if cfg.Username == "" {
		return fmt.Errorf("MongoDB username is required")
	}
	if cfg.Database == "" {
		return fmt.Errorf("MongoDB database is required")
	}
	return nil
}

var MongoClient *mongo.Client

// ConnectMongo initializes the MongoDB connection using MongoDBConfig
func ConnectMongo() {
	cfg := GetMongoDBConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatal("MongoDB config error:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := cfg.BuildMongoURI()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping error:", err)
	}

	fmt.Println("✅ Connected to MongoDB!")
	MongoClient = client
}

// Hapus fungsi getEnv karena sudah ada di database.go
