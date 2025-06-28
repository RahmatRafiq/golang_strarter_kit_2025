package facades

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"golang_starter_kit_2025/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectionManager manages multiple database connections
type ConnectionManager struct {
	gormConnections  map[string]*gorm.DB
	mongoConnections map[string]*mongo.Client
	sqlConnections   map[string]*sql.DB
	config           *config.DatabasesConfig
	mutex            sync.RWMutex
}

var (
	connectionManager *ConnectionManager
	managerOnce       sync.Once
)

// GetManager returns the singleton instance of ConnectionManager
func GetManager() *ConnectionManager {
	managerOnce.Do(func() {
		connectionManager = &ConnectionManager{
			gormConnections:  make(map[string]*gorm.DB),
			mongoConnections: make(map[string]*mongo.Client),
			sqlConnections:   make(map[string]*sql.DB),
			config:           config.GetDatabaseConfig(),
		}
	})
	return connectionManager
}

// GetDB returns the default GORM database connection
func GetDB() *gorm.DB {
	return GetManager().Connection("")
}

// Connection returns a GORM database connection by name
func (cm *ConnectionManager) Connection(name string) *gorm.DB {
	if name == "" {
		name = cm.config.Default
	}

	cm.mutex.RLock()
	if conn, exists := cm.gormConnections[name]; exists {
		cm.mutex.RUnlock()
		return conn
	}
	cm.mutex.RUnlock()

	// Create new connection if it doesn't exist
	return cm.createGormConnection(name)
}

// MongoDB returns a MongoDB connection by name
func (cm *ConnectionManager) MongoDB(name string) *mongo.Client {
	if name == "" {
		// Find first MongoDB connection
		for connName, connConfig := range cm.config.Connections {
			if connConfig.Driver == "mongodb" {
				name = connName
				break
			}
		}
	}

	cm.mutex.RLock()
	if conn, exists := cm.mongoConnections[name]; exists {
		cm.mutex.RUnlock()
		return conn
	}
	cm.mutex.RUnlock()

	// Create new MongoDB connection if it doesn't exist
	return cm.createMongoConnection(name)
}

// MongoDatabase returns a MongoDB database instance
func (cm *ConnectionManager) MongoDatabase(connectionName string, databaseName string) *mongo.Database {
	client := cm.MongoDB(connectionName)
	if client == nil {
		return nil
	}

	if databaseName == "" {
		connConfig, err := cm.config.GetConnectionConfig(connectionName)
		if err != nil {
			log.Printf("Error getting connection config: %v", err)
			return nil
		}
		databaseName = connConfig.Database
	}

	return client.Database(databaseName)
}

// SqlDB returns raw SQL database connection
func (cm *ConnectionManager) SqlDB(name string) *sql.DB {
	gormDB := cm.Connection(name)
	if gormDB == nil {
		return nil
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Printf("Error getting SQL DB from GORM: %v", err)
		return nil
	}

	return sqlDB
}

// createGormConnection creates a new GORM database connection
func (cm *ConnectionManager) createGormConnection(name string) *gorm.DB {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Double-check if connection was created while waiting for lock
	if conn, exists := cm.gormConnections[name]; exists {
		return conn
	}

	connConfig, err := cm.config.GetConnectionConfig(name)
	if err != nil {
		log.Printf("Error getting connection config for '%s': %v", name, err)
		return nil
	}

	if connConfig.Driver == "mongodb" {
		log.Printf("Cannot create GORM connection for MongoDB driver. Use MongoDB() method instead.")
		return nil
	}

	var dialector gorm.Dialector
	dsn := connConfig.BuildDSN()

	switch connConfig.Driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	case "sqlserver":
		dialector = sqlserver.Open(dsn)
	default:
		log.Printf("Unsupported database driver: %s", connConfig.Driver)
		return nil
	}

	// Configure GORM with custom logger
	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(log.Writer(), "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: 500 * time.Millisecond,
				LogLevel:      logger.Warn,
				Colorful:      true,
			},
		),
		PrepareStmt: true,
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		log.Printf("Error connecting to database '%s': %v", name, err)
		return nil
	}

	// Configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting SQL DB instance for '%s': %v", name, err)
		return nil
	}

	sqlDB.SetMaxIdleConns(connConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(connConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(connConfig.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connConfig.ConnMaxIdleTime)

	// Store connections
	cm.gormConnections[name] = db
	cm.sqlConnections[name] = sqlDB

	log.Printf("Database connection '%s' (%s) established successfully", name, connConfig.Driver)
	return db
}

// createMongoConnection creates a new MongoDB connection
func (cm *ConnectionManager) createMongoConnection(name string) *mongo.Client {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Double-check if connection was created while waiting for lock
	if conn, exists := cm.mongoConnections[name]; exists {
		return conn
	}

	connConfig, err := cm.config.GetConnectionConfig(name)
	if err != nil {
		log.Printf("Error getting connection config for '%s': %v", name, err)
		return nil
	}

	if connConfig.Driver != "mongodb" {
		log.Printf("Connection '%s' is not a MongoDB connection", name)
		return nil
	}

	uri := connConfig.BuildMongoURI()

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	if connConfig.MongoOptions != nil {
		if connConfig.MongoOptions.MaxPoolSize > 0 {
			clientOptions.SetMaxPoolSize(uint64(connConfig.MongoOptions.MaxPoolSize))
		}
		if connConfig.MongoOptions.MinPoolSize > 0 {
			clientOptions.SetMinPoolSize(uint64(connConfig.MongoOptions.MinPoolSize))
		}
		if connConfig.MongoOptions.MaxConnIdleTime > 0 {
			clientOptions.SetMaxConnIdleTime(time.Duration(connConfig.MongoOptions.MaxConnIdleTime) * time.Second)
		}
		if connConfig.MongoOptions.ServerSelectionTimeout > 0 {
			clientOptions.SetServerSelectionTimeout(time.Duration(connConfig.MongoOptions.ServerSelectionTimeout) * time.Second)
		}
	}

	// Create MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Printf("Error connecting to MongoDB '%s': %v", name, err)
		return nil
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Error pinging MongoDB '%s': %v", name, err)
		return nil
	}

	cm.mongoConnections[name] = client
	log.Printf("MongoDB connection '%s' established successfully", name)
	return client
}

// CloseAll closes all database connections
func (cm *ConnectionManager) CloseAll() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Close GORM connections
	for name, db := range cm.gormConnections {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing GORM connection '%s': %v", name, err)
			} else {
				log.Printf("GORM connection '%s' closed successfully", name)
			}
		}
	}

	// Close MongoDB connections
	for name, client := range cm.mongoConnections {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error closing MongoDB connection '%s': %v", name, err)
		} else {
			log.Printf("MongoDB connection '%s' closed successfully", name)
		}
		cancel()
	}

	// Clear connection maps
	cm.gormConnections = make(map[string]*gorm.DB)
	cm.mongoConnections = make(map[string]*mongo.Client)
	cm.sqlConnections = make(map[string]*sql.DB)
}

// GetConnectionNames returns all available connection names
func (cm *ConnectionManager) GetConnectionNames() []string {
	return cm.config.ListConnections()
}

// AddConnection adds a new database connection at runtime
func (cm *ConnectionManager) AddConnection(name string, connConfig config.DatabaseConfig) {
	cm.config.AddConnection(name, connConfig)
}

// RemoveConnection removes a database connection
func (cm *ConnectionManager) RemoveConnection(name string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Close existing connections
	if db, exists := cm.gormConnections[name]; exists {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
		delete(cm.gormConnections, name)
	}

	if client, exists := cm.mongoConnections[name]; exists {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		client.Disconnect(ctx)
		cancel()
		delete(cm.mongoConnections, name)
	}

	// Remove from config
	cm.config.RemoveConnection(name)
}

// TestConnection tests a database connection
func (cm *ConnectionManager) TestConnection(name string) error {
	connConfig, err := cm.config.GetConnectionConfig(name)
	if err != nil {
		return err
	}

	if connConfig.Driver == "mongodb" {
		client := cm.MongoDB(name)
		if client == nil {
			return fmt.Errorf("failed to connect to MongoDB '%s'", name)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return client.Ping(ctx, nil)
	} else {
		db := cm.Connection(name)
		if db == nil {
			return fmt.Errorf("failed to connect to database '%s'", name)
		}

		sqlDB, err := db.DB()
		if err != nil {
			return err
		}

		return sqlDB.Ping()
	}
}
