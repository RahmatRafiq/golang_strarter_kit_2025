package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"golang_starter_kit_2025/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Manager handles multiple database connections
type Manager struct {
	connections map[string]*Connection
	configs     *config.DatabaseConfigs
	mutex       sync.RWMutex
}

// Connection wraps GORM DB with additional metadata
type Connection struct {
	DB     *gorm.DB
	SqlDB  *sql.DB
	Config *config.DatabaseConfig
	Name   string
}

var (
	manager *Manager
	once    sync.Once
)

// GetManager returns the singleton database manager
func GetManager() *Manager {
	once.Do(func() {
		manager = &Manager{
			connections: make(map[string]*Connection),
			configs:     config.GetDatabaseConfigs(),
		}
	})
	return manager
}

// Connect establishes a connection to a specific database
func (m *Manager) Connect(connectionName string) (*Connection, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Return existing connection if available
	if conn, exists := m.connections[connectionName]; exists {
		return conn, nil
	}

	// Get configuration for the connection
	cfg, exists := m.configs.Connections[connectionName]
	if !exists {
		return nil, fmt.Errorf("database configuration '%s' not found", connectionName)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration for '%s': %v", connectionName, err)
	}

	// Create GORM dialector based on database type
	var dialector gorm.Dialector
	switch cfg.Type {
	case config.MySQL:
		dialector = mysql.Open(cfg.BuildDSN())
	case config.PostgreSQL:
		dialector = postgres.Open(cfg.BuildDSN())
	case config.SQLite:
		dialector = sqlite.Open(cfg.BuildDSN())
	case config.SQLServer:
		dialector = sqlserver.Open(cfg.BuildDSN())
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		PrepareStmt: true,
	}

	// Open database connection
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database '%s': %v", connectionName, err)
	}

	// Get underlying SQL DB for connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB object for '%s': %v", connectionName, err)
	}

	// Configure connection pooling
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database '%s': %v", connectionName, err)
	}

	// Create connection object
	connection := &Connection{
		DB:     db,
		SqlDB:  sqlDB,
		Config: cfg,
		Name:   connectionName,
	}

	// Store connection
	m.connections[connectionName] = connection

	log.Printf("✅ Database connection '%s' (%s) established successfully", connectionName, cfg.Type)
	return connection, nil
}

// GetConnection returns an existing connection or creates a new one
func (m *Manager) GetConnection(connectionName string) (*Connection, error) {
	m.mutex.RLock()
	if conn, exists := m.connections[connectionName]; exists {
		m.mutex.RUnlock()
		return conn, nil
	}
	m.mutex.RUnlock()

	return m.Connect(connectionName)
}

// GetDefaultConnection returns the default database connection
func (m *Manager) GetDefaultConnection() (*Connection, error) {
	return m.GetConnection(m.configs.Default)
}

// CloseConnection closes a specific database connection
func (m *Manager) CloseConnection(connectionName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if conn, exists := m.connections[connectionName]; exists {
		if err := conn.SqlDB.Close(); err != nil {
			log.Printf("Error closing connection '%s': %v", connectionName, err)
			return err
		}
		delete(m.connections, connectionName)
		log.Printf("✅ Database connection '%s' closed successfully", connectionName)
	}

	return nil
}

// CloseAllConnections closes all database connections
func (m *Manager) CloseAllConnections() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for name, conn := range m.connections {
		if err := conn.SqlDB.Close(); err != nil {
			log.Printf("Error closing connection '%s': %v", name, err)
		} else {
			log.Printf("✅ Database connection '%s' closed successfully", name)
		}
	}

	m.connections = make(map[string]*Connection)
}

// ListConnections returns the names of all configured connections
func (m *Manager) ListConnections() []string {
	var connections []string
	for name := range m.configs.Connections {
		connections = append(connections, name)
	}
	return connections
}

// IsConnected checks if a connection is established and healthy
func (m *Manager) IsConnected(connectionName string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if conn, exists := m.connections[connectionName]; exists {
		if err := conn.SqlDB.Ping(); err == nil {
			return true
		}
	}
	return false
}

// GetConnectionStats returns statistics for a connection
func (m *Manager) GetConnectionStats(connectionName string) (sql.DBStats, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if conn, exists := m.connections[connectionName]; exists {
		return conn.SqlDB.Stats(), nil
	}

	return sql.DBStats{}, fmt.Errorf("connection '%s' not found", connectionName)
}

// Helper methods for the Connection struct

// Transaction starts a new transaction on this connection
func (c *Connection) Transaction(fn func(*gorm.DB) error) error {
	return c.DB.Transaction(fn)
}

// Migrate runs auto migration for models on this connection
func (c *Connection) Migrate(models ...interface{}) error {
	return c.DB.AutoMigrate(models...)
}

// GetType returns the database type for this connection
func (c *Connection) GetType() config.DatabaseType {
	return c.Config.Type
}

// IsMySQL checks if this connection is MySQL
func (c *Connection) IsMySQL() bool {
	return c.Config.Type == config.MySQL
}

// IsPostgreSQL checks if this connection is PostgreSQL
func (c *Connection) IsPostgreSQL() bool {
	return c.Config.Type == config.PostgreSQL
}
