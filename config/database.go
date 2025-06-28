package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// DatabaseConfig represents configuration for a single database connection
type DatabaseConfig struct {
	Driver          string            `json:"driver"`
	Host            string            `json:"host"`
	Port            string            `json:"port"`
	Database        string            `json:"database"`
	Username        string            `json:"username"`
	Password        string            `json:"password"`
	Charset         string            `json:"charset"`
	Collation       string            `json:"collation"`
	Prefix          string            `json:"prefix"`
	Timezone        string            `json:"timezone"`
	SSLMode         string            `json:"ssl_mode"`
	MaxIdleConns    int               `json:"max_idle_conns"`
	MaxOpenConns    int               `json:"max_open_conns"`
	ConnMaxLifetime time.Duration     `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration     `json:"conn_max_idle_time"`
	Options         map[string]string `json:"options"`
	MongoOptions    *MongoConfig      `json:"mongo_options,omitempty"`
}

// MongoConfig represents MongoDB specific configuration
type MongoConfig struct {
	AuthSource             string `json:"auth_source"`
	AuthMechanism          string `json:"auth_mechanism"`
	ReplicaSet             string `json:"replica_set"`
	ReadPreference         string `json:"read_preference"`
	WriteConcern           string `json:"write_concern"`
	ReadConcern            string `json:"read_concern"`
	MaxPoolSize            int    `json:"max_pool_size"`
	MinPoolSize            int    `json:"min_pool_size"`
	MaxConnIdleTime        int    `json:"max_conn_idle_time"`
	ServerSelectionTimeout int    `json:"server_selection_timeout"`
}

// DatabasesConfig holds all database connections configuration
type DatabasesConfig struct {
	Default     string                    `json:"default"`
	Connections map[string]DatabaseConfig `json:"connections"`
}

// GetDatabaseConfig returns the database configuration
func GetDatabaseConfig() *DatabasesConfig {
	defaultConnection := getEnv("DB_CONNECTION", "mysql")

	config := &DatabasesConfig{
		Default: defaultConnection,
		Connections: map[string]DatabaseConfig{
			"mysql": {
				Driver:          "mysql",
				Host:            getEnv("DB_HOST", "localhost"),
				Port:            getEnv("DB_PORT", "3306"),
				Database:        getEnv("DB_DATABASE", "laravel"),
				Username:        getEnv("DB_USERNAME", "root"),
				Password:        getEnv("DB_PASSWORD", ""),
				Charset:         getEnv("DB_CHARSET", "utf8mb4"),
				Collation:       getEnv("DB_COLLATION", "utf8mb4_unicode_ci"),
				Prefix:          getEnv("DB_PREFIX", ""),
				Timezone:        getEnv("DB_TIMEZONE", "Local"),
				MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
				MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 200),
				ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME", 15)) * time.Minute,
				ConnMaxIdleTime: time.Duration(getEnvAsInt("DB_CONN_MAX_IDLE_TIME", 5)) * time.Minute,
				Options: map[string]string{
					"parseTime": "True",
					"loc":       "Local",
				},
			},
			"mysql_secondary": {
				Driver:          "mysql",
				Host:            getEnv("MYSQL_SECONDARY_HOST", "localhost"),
				Port:            getEnv("MYSQL_SECONDARY_PORT", "3306"),
				Database:        getEnv("MYSQL_SECONDARY_DB", "secondary"),
				Username:        getEnv("MYSQL_SECONDARY_USER", "root"),
				Password:        getEnv("MYSQL_SECONDARY_PASSWORD", ""),
				Charset:         getEnv("MYSQL_SECONDARY_CHARSET", "utf8mb4"),
				Collation:       getEnv("MYSQL_SECONDARY_COLLATION", "utf8mb4_unicode_ci"),
				Prefix:          getEnv("MYSQL_SECONDARY_PREFIX", ""),
				Timezone:        getEnv("MYSQL_SECONDARY_TIMEZONE", "Local"),
				MaxIdleConns:    getEnvAsInt("MYSQL_SECONDARY_MAX_IDLE_CONNS", 10),
				MaxOpenConns:    getEnvAsInt("MYSQL_SECONDARY_MAX_OPEN_CONNS", 200),
				ConnMaxLifetime: time.Duration(getEnvAsInt("MYSQL_SECONDARY_CONN_MAX_LIFETIME", 15)) * time.Minute,
				ConnMaxIdleTime: time.Duration(getEnvAsInt("MYSQL_SECONDARY_CONN_MAX_IDLE_TIME", 5)) * time.Minute,
				Options: map[string]string{
					"parseTime": "True",
					"loc":       "Local",
				},
			},
			"postgres": {
				Driver:          "postgres",
				Host:            getEnv("POSTGRES_HOST", "localhost"),
				Port:            getEnv("POSTGRES_PORT", "5432"),
				Database:        getEnv("POSTGRES_DB", "postgres"),
				Username:        getEnv("POSTGRES_USER", "postgres"),
				Password:        getEnv("POSTGRES_PASSWORD", ""),
				Charset:         getEnv("POSTGRES_CHARSET", "utf8"),
				SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
				Timezone:        getEnv("POSTGRES_TIMEZONE", "UTC"),
				MaxIdleConns:    getEnvAsInt("POSTGRES_MAX_IDLE_CONNS", 10),
				MaxOpenConns:    getEnvAsInt("POSTGRES_MAX_OPEN_CONNS", 200),
				ConnMaxLifetime: time.Duration(getEnvAsInt("POSTGRES_CONN_MAX_LIFETIME", 15)) * time.Minute,
				ConnMaxIdleTime: time.Duration(getEnvAsInt("POSTGRES_CONN_MAX_IDLE_TIME", 5)) * time.Minute,
			},
			"sqlite": {
				Driver:          "sqlite",
				Database:        getEnv("SQLITE_DATABASE", "database/database.sqlite"),
				MaxIdleConns:    getEnvAsInt("SQLITE_MAX_IDLE_CONNS", 1),
				MaxOpenConns:    getEnvAsInt("SQLITE_MAX_OPEN_CONNS", 1),
				ConnMaxLifetime: time.Duration(getEnvAsInt("SQLITE_CONN_MAX_LIFETIME", 15)) * time.Minute,
				ConnMaxIdleTime: time.Duration(getEnvAsInt("SQLITE_CONN_MAX_IDLE_TIME", 5)) * time.Minute,
			},
			"sqlserver": {
				Driver:          "sqlserver",
				Host:            getEnv("SQLSERVER_HOST", "localhost"),
				Port:            getEnv("SQLSERVER_PORT", "1433"),
				Database:        getEnv("SQLSERVER_DB", "master"),
				Username:        getEnv("SQLSERVER_USER", "sa"),
				Password:        getEnv("SQLSERVER_PASSWORD", ""),
				MaxIdleConns:    getEnvAsInt("SQLSERVER_MAX_IDLE_CONNS", 10),
				MaxOpenConns:    getEnvAsInt("SQLSERVER_MAX_OPEN_CONNS", 200),
				ConnMaxLifetime: time.Duration(getEnvAsInt("SQLSERVER_CONN_MAX_LIFETIME", 15)) * time.Minute,
				ConnMaxIdleTime: time.Duration(getEnvAsInt("SQLSERVER_CONN_MAX_IDLE_TIME", 5)) * time.Minute,
			},
			"mongodb": {
				Driver:   "mongodb",
				Host:     getEnv("MONGO_HOST", "localhost"),
				Port:     getEnv("MONGO_PORT", "27017"),
				Database: getEnv("MONGO_DATABASE", "laravel"),
				Username: getEnv("MONGO_USERNAME", ""),
				Password: getEnv("MONGO_PASSWORD", ""),
				MongoOptions: &MongoConfig{
					AuthSource:             getEnv("MONGO_AUTH_SOURCE", "admin"),
					AuthMechanism:          getEnv("MONGO_AUTH_MECHANISM", ""),
					ReplicaSet:             getEnv("MONGO_REPLICA_SET", ""),
					ReadPreference:         getEnv("MONGO_READ_PREFERENCE", "primary"),
					WriteConcern:           getEnv("MONGO_WRITE_CONCERN", "majority"),
					ReadConcern:            getEnv("MONGO_READ_CONCERN", "local"),
					MaxPoolSize:            getEnvAsInt("MONGO_MAX_POOL_SIZE", 100),
					MinPoolSize:            getEnvAsInt("MONGO_MIN_POOL_SIZE", 0),
					MaxConnIdleTime:        getEnvAsInt("MONGO_MAX_CONN_IDLE_TIME", 0),
					ServerSelectionTimeout: getEnvAsInt("MONGO_SERVER_SELECTION_TIMEOUT", 30),
				},
			},
		},
	}

	return config
}

// GetConnectionConfig returns configuration for a specific connection
func (dc *DatabasesConfig) GetConnectionConfig(name string) (DatabaseConfig, error) {
	if name == "" {
		name = dc.Default
	}

	config, exists := dc.Connections[name]
	if !exists {
		return DatabaseConfig{}, fmt.Errorf("database connection '%s' not found", name)
	}

	return config, nil
}

// GetDefaultConfig returns the default database configuration
func (dc *DatabasesConfig) GetDefaultConfig() (DatabaseConfig, error) {
	return dc.GetConnectionConfig(dc.Default)
}

// AddConnection adds a new database connection configuration
func (dc *DatabasesConfig) AddConnection(name string, config DatabaseConfig) {
	if dc.Connections == nil {
		dc.Connections = make(map[string]DatabaseConfig)
	}
	dc.Connections[name] = config
}

// RemoveConnection removes a database connection configuration
func (dc *DatabasesConfig) RemoveConnection(name string) {
	delete(dc.Connections, name)
}

// ListConnections returns all available connection names
func (dc *DatabasesConfig) ListConnections() []string {
	var connections []string
	for name := range dc.Connections {
		connections = append(connections, name)
	}
	return connections
}

// BuildDSN builds Data Source Name for SQL databases
func (config DatabaseConfig) BuildDSN() string {
	switch config.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
			config.Charset,
		)

		// Add additional options
		for key, value := range config.Options {
			dsn += fmt.Sprintf("&%s=%s", key, value)
		}

		return dsn

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			config.Host,
			config.Username,
			config.Password,
			config.Database,
			config.Port,
			config.SSLMode,
			config.Timezone,
		)
		return dsn

	case "sqlite":
		return config.Database

	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
		)
		return dsn

	default:
		log.Printf("Warning: Unknown driver '%s', returning empty DSN", config.Driver)
		return ""
	}
}

// BuildMongoURI builds MongoDB connection URI
func (config DatabaseConfig) BuildMongoURI() string {
	if config.Driver != "mongodb" {
		return ""
	}

	var uri string
	if config.Username != "" && config.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
		)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)
	}

	// Add additional MongoDB options
	if config.MongoOptions != nil {
		params := []string{}

		if config.MongoOptions.AuthSource != "" {
			params = append(params, fmt.Sprintf("authSource=%s", config.MongoOptions.AuthSource))
		}
		if config.MongoOptions.AuthMechanism != "" {
			params = append(params, fmt.Sprintf("authMechanism=%s", config.MongoOptions.AuthMechanism))
		}
		if config.MongoOptions.ReplicaSet != "" {
			params = append(params, fmt.Sprintf("replicaSet=%s", config.MongoOptions.ReplicaSet))
		}
		if config.MongoOptions.ReadPreference != "" {
			params = append(params, fmt.Sprintf("readPreference=%s", config.MongoOptions.ReadPreference))
		}
		if config.MongoOptions.WriteConcern != "" {
			params = append(params, fmt.Sprintf("w=%s", config.MongoOptions.WriteConcern))
		}
		if config.MongoOptions.ReadConcern != "" {
			params = append(params, fmt.Sprintf("readConcernLevel=%s", config.MongoOptions.ReadConcern))
		}
		if config.MongoOptions.MaxPoolSize > 0 {
			params = append(params, fmt.Sprintf("maxPoolSize=%d", config.MongoOptions.MaxPoolSize))
		}
		if config.MongoOptions.MinPoolSize > 0 {
			params = append(params, fmt.Sprintf("minPoolSize=%d", config.MongoOptions.MinPoolSize))
		}
		if config.MongoOptions.MaxConnIdleTime > 0 {
			params = append(params, fmt.Sprintf("maxIdleTimeMS=%d", config.MongoOptions.MaxConnIdleTime*1000))
		}
		if config.MongoOptions.ServerSelectionTimeout > 0 {
			params = append(params, fmt.Sprintf("serverSelectionTimeoutMS=%d", config.MongoOptions.ServerSelectionTimeout*1000))
		}

		if len(params) > 0 {
			uri += "/?" + fmt.Sprintf("%s", params[0])
			for _, param := range params[1:] {
				uri += "&" + param
			}
		}
	}

	return uri
}

// getEnvAsInt gets environment variable as integer with default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s: %s, using default: %d", key, valueStr, defaultValue)
		return defaultValue
	}

	return value
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
