package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// DatabaseType represents the type of database
type DatabaseType string

const (
	MySQL      DatabaseType = "mysql"
	PostgreSQL DatabaseType = "postgres"
	SQLite     DatabaseType = "sqlite"
	SQLServer  DatabaseType = "sqlserver"
)

// DatabaseConfig holds configuration for a database connection
type DatabaseConfig struct {
	Type            DatabaseType  `json:"type"`
	Host            string        `json:"host"`
	Port            string        `json:"port"`
	Database        string        `json:"database"`
	Username        string        `json:"username"`
	Password        string        `json:"password"`
	Charset         string        `json:"charset"`
	Timezone        string        `json:"timezone"`
	SSLMode         string        `json:"ssl_mode"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	MaxOpenConns    int           `json:"max_open_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

// DatabaseConfigs holds multiple database configurations
type DatabaseConfigs struct {
	Default     string                     `json:"default"`
	Connections map[string]*DatabaseConfig `json:"connections"`
}

// GetDatabaseConfigs loads database configurations from environment variables
func GetDatabaseConfigs() *DatabaseConfigs {
	configs := &DatabaseConfigs{
		Default:     getEnv("DB_CONNECTION", "mysql"),
		Connections: make(map[string]*DatabaseConfig),
	}

	// MySQL Configuration
	configs.Connections["mysql"] = &DatabaseConfig{
		Type:            MySQL,
		Host:            getEnv("MYSQL_HOST", "localhost"),
		Port:            getEnv("MYSQL_PORT", "3306"),
		Database:        getEnv("MYSQL_DB", ""),
		Username:        getEnv("MYSQL_USER", ""),
		Password:        getEnv("MYSQL_PASSWORD", ""),
		Charset:         getEnv("MYSQL_CHARSET", "utf8mb4"),
		Timezone:        getEnv("MYSQL_TIMEZONE", "Local"),
		MaxIdleConns:    getEnvAsInt("MYSQL_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getEnvAsInt("MYSQL_MAX_OPEN_CONNS", 200),
		ConnMaxLifetime: getEnvAsDuration("MYSQL_CONN_MAX_LIFETIME", 15*time.Minute),
		ConnMaxIdleTime: getEnvAsDuration("MYSQL_CONN_MAX_IDLE_TIME", 5*time.Minute),
	}

	// PostgreSQL Configuration
	configs.Connections["postgres"] = &DatabaseConfig{
		Type:            PostgreSQL,
		Host:            getEnv("POSTGRES_HOST", "localhost"),
		Port:            getEnv("POSTGRES_PORT", "5432"),
		Database:        getEnv("POSTGRES_DB", ""),
		Username:        getEnv("POSTGRES_USER", ""),
		Password:        getEnv("POSTGRES_PASSWORD", ""),
		Timezone:        getEnv("POSTGRES_TIMEZONE", "UTC"),
		SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
		MaxIdleConns:    getEnvAsInt("POSTGRES_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getEnvAsInt("POSTGRES_MAX_OPEN_CONNS", 200),
		ConnMaxLifetime: getEnvAsDuration("POSTGRES_CONN_MAX_LIFETIME", 15*time.Minute),
		ConnMaxIdleTime: getEnvAsDuration("POSTGRES_CONN_MAX_IDLE_TIME", 5*time.Minute),
	}

	// MySQL Secondary Configuration (for multiple MySQL instances)
	configs.Connections["mysql_secondary"] = &DatabaseConfig{
		Type:            MySQL,
		Host:            getEnv("MYSQL_SECONDARY_HOST", "localhost"),
		Port:            getEnv("MYSQL_SECONDARY_PORT", "3306"),
		Database:        getEnv("MYSQL_SECONDARY_DB", ""),
		Username:        getEnv("MYSQL_SECONDARY_USER", ""),
		Password:        getEnv("MYSQL_SECONDARY_PASSWORD", ""),
		Charset:         getEnv("MYSQL_SECONDARY_CHARSET", "utf8mb4"),
		Timezone:        getEnv("MYSQL_SECONDARY_TIMEZONE", "Local"),
		MaxIdleConns:    getEnvAsInt("MYSQL_SECONDARY_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getEnvAsInt("MYSQL_SECONDARY_MAX_OPEN_CONNS", 200),
		ConnMaxLifetime: getEnvAsDuration("MYSQL_SECONDARY_CONN_MAX_LIFETIME", 15*time.Minute),
		ConnMaxIdleTime: getEnvAsDuration("MYSQL_SECONDARY_CONN_MAX_IDLE_TIME", 5*time.Minute),
	}

	return configs
}

// BuildDSN builds the DSN string for the database connection
func (cfg *DatabaseConfig) BuildDSN() string {
	switch cfg.Type {
	case MySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset, cfg.Timezone)

	case PostgreSQL:
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode, cfg.Timezone)

	case SQLite:
		return cfg.Database

	case SQLServer:
		return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	default:
		return ""
	}
}

// Validate checks if the configuration is valid
func (cfg *DatabaseConfig) Validate() error {
	if cfg.Type == "" {
		return fmt.Errorf("database type is required")
	}

	switch cfg.Type {
	case MySQL, PostgreSQL, SQLServer:
		if cfg.Host == "" {
			return fmt.Errorf("host is required for %s", cfg.Type)
		}
		if cfg.Database == "" {
			return fmt.Errorf("database name is required for %s", cfg.Type)
		}
		if cfg.Username == "" {
			return fmt.Errorf("username is required for %s", cfg.Type)
		}
	case SQLite:
		if cfg.Database == "" {
			return fmt.Errorf("database file path is required for SQLite")
		}
	}

	return nil
}

// Helper functions
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}
