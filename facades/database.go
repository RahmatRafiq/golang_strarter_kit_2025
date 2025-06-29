package facades

import (
	"log"

	"golang_starter_kit_2025/database"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	manager *database.Manager
	DB      *gorm.DB // Backward compatibility for existing code
)

// ConnectDB initializes the database manager and connects to the default database
func ConnectDB(envFiles ...string) *gorm.DB {
	// Load environment variables from .env files if provided
	if len(envFiles) > 0 {
		err := godotenv.Load(envFiles...)
		if err != nil {
			log.Printf("Warning: No .env file found. Using environment variables instead. Error: %v", err)
		} else {
			log.Println(".env file loaded successfully")
		}
	}

	// Initialize database manager
	manager = database.GetManager()

	// Connect to default database
	conn, err := manager.GetDefaultConnection()
	if err != nil {
		log.Fatalf("Error: failed to connect to default database: %v", err)
	}

	// Set global DB for backward compatibility
	DB = conn.DB

	log.Println("Database manager initialized successfully")
	return DB
}

// GetDB returns the default database connection (backward compatibility)
func GetDB() *gorm.DB {
	if DB == nil {
		ConnectDB()
	}
	return DB
}

// GetConnection returns a specific database connection
func GetConnection(connectionName string) (*database.Connection, error) {
	if manager == nil {
		ConnectDB()
	}
	return manager.GetConnection(connectionName)
}

// GetDefaultConnection returns the default database connection
func GetDefaultConnection() (*database.Connection, error) {
	if manager == nil {
		ConnectDB()
	}
	return manager.GetDefaultConnection()
}

// GetManager returns the database manager instance
func GetManager() *database.Manager {
	if manager == nil {
		ConnectDB()
	}
	return manager
}

// CloseDB closes all database connections
func CloseDB() {
	if manager != nil {
		manager.CloseAllConnections()
		log.Println("All database connections closed successfully")
	}
}

// CloseConnection closes a specific database connection
func CloseConnection(connectionName string) error {
	if manager == nil {
		return nil
	}
	return manager.CloseConnection(connectionName)
}

// ListConnections returns all available connection names
func ListConnections() []string {
	if manager == nil {
		ConnectDB()
	}
	return manager.ListConnections()
}

// IsConnected checks if a connection is healthy
func IsConnected(connectionName string) bool {
	if manager == nil {
		return false
	}
	return manager.IsConnected(connectionName)
}

// MySQL returns the MySQL connection
func MySQL() (*database.Connection, error) {
	return GetConnection("mysql")
}

// PostgreSQL returns the PostgreSQL connection
func PostgreSQL() (*database.Connection, error) {
	return GetConnection("postgres")
}

// MySQLSecondary returns the secondary MySQL connection
func MySQLSecondary() (*database.Connection, error) {
	return GetConnection("mysql_secondary")
}
