package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB holds test database connection and configuration
type TestDB struct {
	DB         *gorm.DB
	Connection string
	DBName     string
}

// SetupTestDB initializes a test database connection
// Similar to Laravel's RefreshDatabase trait
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	// Load .env.testing
	loadTestEnv(t)

	connection := os.Getenv("DB_CONNECTION")
	if connection == "" {
		connection = "mysql"
	}

	var db *gorm.DB
	var err error
	var dbName string

	// Setup database based on connection type
	switch connection {
	case "mysql":
		dbName = os.Getenv("MYSQL_DB")
		db, err = setupMySQLTestDB(dbName)
	case "postgres":
		dbName = os.Getenv("POSTGRES_DB")
		db, err = setupPostgresTestDB(dbName)
	default:
		t.Fatalf("Unsupported DB_CONNECTION: %s", connection)
	}

	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	testDB := &TestDB{
		DB:         db,
		Connection: connection,
		DBName:     dbName,
	}

	// IMPORTANT: Setup facades.DB for services that use it (like JwtService)
	// This is required because some services directly access facades.DB
	setupFacades(t, db)

	// Run migrations (like Laravel's RefreshDatabase)
	RunMigrations(t, testDB)

	return testDB
}

// setupMySQLTestDB creates MySQL test database connection
func setupMySQLTestDB(dbName string) (*gorm.DB, error) {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	charset := os.Getenv("MYSQL_CHARSET")
	timezone := os.Getenv("MYSQL_TIMEZONE")

	// First connect without database to create it
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=%s&parseTime=True&loc=%s",
		user, password, host, port, charset, timezone)

	db, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// Drop and recreate test database (fresh state for each test run)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	if _, dropErr := sqlDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)); dropErr != nil {
		fmt.Printf("Warning: failed to drop test database: %v\n", dropErr)
	}

	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %w", err)
	}
	sqlDB.Close()

	// Connect to the test database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
		user, password, host, port, dbName, charset, timezone)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// setupPostgresTestDB creates PostgreSQL test database connection
func setupPostgresTestDB(dbName string) (*gorm.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	sslmode := os.Getenv("POSTGRES_SSLMODE")
	timezone := os.Getenv("POSTGRES_TIMEZONE")

	// First connect to postgres database to create test database
	dsnWithoutDB := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s TimeZone=%s",
		host, port, user, password, sslmode, timezone)

	db, err := gorm.Open(postgres.Open(dsnWithoutDB), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Drop and recreate test database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	if _, dropErr := sqlDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)); dropErr != nil {
		fmt.Printf("Warning: failed to drop test database: %v\n", dropErr)
	}
	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %w", err)
	}
	sqlDB.Close()

	// Connect to the test database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, password, dbName, sslmode, timezone)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// RunMigrations runs all migrations on test database
// Similar to Laravel's php artisan migrate --env=testing
func RunMigrations(t *testing.T, testDB *TestDB) {
	t.Helper()

	// Use GORM AutoMigrate for quick schema setup in tests
	// This is faster than SQL migrations and good enough for testing
	err := testDB.DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Product{},
		&models.Category{},
		&models.RefreshToken{},
	)

	if err != nil {
		t.Fatalf("Failed to run auto-migrate: %v", err)
	}

	// Create junction tables for many-to-many relationships
	// These are not auto-created by GORM's AutoMigrate
	createJunctionTables(t, testDB.DB)

	t.Log("✓ All migrations completed successfully")
}

// createJunctionTables creates many-to-many junction tables
func createJunctionTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	// user_has_roles table
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_has_roles (
			user_id INT UNSIGNED NOT NULL,
			role_id INT UNSIGNED NOT NULL,
			PRIMARY KEY (user_id, role_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
		)
	`).Error

	if err != nil {
		t.Fatalf("Failed to create user_has_roles table: %v", err)
	}

	// role_has_permissions table
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS role_has_permissions (
			role_id INT UNSIGNED NOT NULL,
			permission_id INT UNSIGNED NOT NULL,
			PRIMARY KEY (role_id, permission_id),
			FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
			FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
		)
	`).Error

	if err != nil {
		t.Fatalf("Failed to create role_has_permissions table: %v", err)
	}

	t.Log("✓ Junction tables created")
}

// CleanupTestDB drops test database and closes connection
// Called in defer after test completes
func CleanupTestDB(t *testing.T, testDB *TestDB) {
	t.Helper()

	if testDB == nil || testDB.DB == nil {
		return
	}

	sqlDB, err := testDB.DB.DB()
	if err != nil {
		t.Logf("Warning: Failed to get sql.DB: %v", err)
		return
	}

	// Close connection
	sqlDB.Close()

	// Note: We don't drop the database here to allow inspection after test failure
	// The database will be dropped on next test run anyway
	t.Logf("Test database %s closed (will be recreated on next run)", testDB.DBName)
}

// TruncateTables removes all data from tables (for test isolation)
// Similar to Laravel's DatabaseTransactions trait
func TruncateTables(t *testing.T, testDB *TestDB) {
	t.Helper()

	tables := []string{
		"user_has_roles",
		"role_has_permissions",
		"products",
		"categories",
		"users",
		"roles",
		"permissions",
	}

	for _, table := range tables {
		if err := testDB.DB.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			t.Logf("Warning: Failed to truncate %s: %v", table, err)
		}
	}

	t.Log("All tables truncated")
}

// setupFacades initializes facades.DB for services that depend on it
func setupFacades(t *testing.T, db *gorm.DB) {
	t.Helper()

	// Set facades.DB to use test database connection
	facades.DB = db

	t.Log("✓ Facades initialized with test database")
}

// loadTestEnv loads .env.testing file and overrides existing env vars
func loadTestEnv(t *testing.T) {
	t.Helper()

	// Get project root directory
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

	envFile := filepath.Join(projectRoot, ".env.testing")

	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		t.Fatalf(".env.testing not found at: %s", envFile)
	}

	// IMPORTANT: Use Overload instead of Load to override existing .env values
	if err := godotenv.Overload(envFile); err != nil {
		t.Fatalf("Failed to load .env.testing: %v", err)
	}

	// Verify APP_ENV is set to testing
	if os.Getenv("APP_ENV") != "testing" {
		log.Fatal("APP_ENV must be 'testing' in .env.testing")
	}

	t.Logf("Loaded .env.testing configuration (DB: %s)", os.Getenv("MYSQL_DB"))
}
