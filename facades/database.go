package facades

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB  // Global GORM database instance
	SqlDB *sql.DB   // Global raw SQL DB instance
	once  sync.Once // Singleton pattern to ensure only one connection
)

// ConnectDB initializes and returns the GORM DB connection.
// It loads environment variables, configures connection pooling,
// and ensures the connection is only opened once.
func ConnectDB(envFiles ...string) *gorm.DB {
	once.Do(func() {
		// Load environment variables from .env files if provided
		if len(envFiles) > 0 {
			err := godotenv.Load(envFiles...)
			if err != nil {
				log.Printf("Warning: No .env file found. Using environment variables instead. Error: %v", err)
			} else {
				log.Println(".env file loaded successfully")
			}
		}

		// Validate required environment variables
		requiredEnvVars := []string{"MYSQL_HOST", "MYSQL_PORT", "MYSQL_DB", "MYSQL_USER"}
		for _, v := range requiredEnvVars {
			if os.Getenv(v) == "" {
				log.Fatalf("Error: Required environment variable %s is not set", v)
			}
		}

		// Build DSN (Data Source Name) for MySQL
		dsn := fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DB"),
		)
		// Build DSN (Data Source Name) for MySQL without password

		// Open the database connection using GORM with custom logger
		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: 500 * time.Millisecond, // Log slow queries
					LogLevel:      logger.Warn,            // Warning level logging
					Colorful:      true,                   // Colorful log output
				},
			),
			PrepareStmt: true, // Prepared statement caching
		})
		if err != nil {
			log.Fatalf("Error: failed to connect to the database: %v", err)
		}

		// Get the underlying SQL DB object for connection pooling configuration
		SqlDB, err = DB.DB()
		if err != nil {
			log.Fatalf("Error: failed to get SQL DB object from GORM: %v", err)
		}

		// Configure connection pooling
		SqlDB.SetMaxIdleConns(10)
		SqlDB.SetMaxOpenConns(200)
		SqlDB.SetConnMaxLifetime(15 * time.Minute)
		SqlDB.SetConnMaxIdleTime(5 * time.Minute)

		log.Println("Database connection successfully established")
	})

	return DB
}

// CloseDB safely closes the database connection when it is no longer needed.
// This should be called in a defer statement.
func CloseDB() {
	if SqlDB != nil {
		err := SqlDB.Close()
		if err != nil {
			log.Printf("Error closing the database connection: %v", err)
		} else {
			log.Println("Database connection successfully closed")
		}
	}
}
