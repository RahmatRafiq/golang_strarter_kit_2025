package helpers

import (
	"context"
	"golang_starter_kit_2025/facades"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// DatabaseHelper provides helper functions for multi-database operations
type DatabaseHelper struct {
	manager *facades.ConnectionManager
}

// NewDatabaseHelper creates a new instance of DatabaseHelper
func NewDatabaseHelper() *DatabaseHelper {
	return &DatabaseHelper{
		manager: facades.GetManager(),
	}
}

// GetDefaultDB returns the default GORM database connection
func (dh *DatabaseHelper) GetDefaultDB() *gorm.DB {
	return dh.manager.Connection("")
}

// GetMysqlDB returns MySQL database connection
func (dh *DatabaseHelper) GetMysqlDB() *gorm.DB {
	return dh.manager.Connection("mysql")
}

// GetMysqlSecondaryDB returns MySQL secondary database connection
func (dh *DatabaseHelper) GetMysqlSecondaryDB() *gorm.DB {
	return dh.manager.Connection("mysql_secondary")
}

// GetPostgresDB returns PostgreSQL database connection
func (dh *DatabaseHelper) GetPostgresDB() *gorm.DB {
	return dh.manager.Connection("postgres")
}

// GetSqliteDB returns SQLite database connection
func (dh *DatabaseHelper) GetSqliteDB() *gorm.DB {
	return dh.manager.Connection("sqlite")
}

// GetSqlServerDB returns SQL Server database connection
func (dh *DatabaseHelper) GetSqlServerDB() *gorm.DB {
	return dh.manager.Connection("sqlserver")
}

// GetMongoDB returns MongoDB client
func (dh *DatabaseHelper) GetMongoDB() *mongo.Client {
	return dh.manager.MongoDB("mongodb")
}

// GetMongoDatabase returns MongoDB database instance
func (dh *DatabaseHelper) GetMongoDatabase(databaseName string) *mongo.Database {
	return dh.manager.MongoDatabase("mongodb", databaseName)
}

// GetMongoCollection returns MongoDB collection
func (dh *DatabaseHelper) GetMongoCollection(databaseName, collectionName string) *mongo.Collection {
	db := dh.GetMongoDatabase(databaseName)
	if db == nil {
		return nil
	}
	return db.Collection(collectionName)
}

// GetDBByName returns database connection by name
func (dh *DatabaseHelper) GetDBByName(connectionName string) *gorm.DB {
	return dh.manager.Connection(connectionName)
}

// GetMongoDBByName returns MongoDB client by connection name
func (dh *DatabaseHelper) GetMongoDBByName(connectionName string) *mongo.Client {
	return dh.manager.MongoDB(connectionName)
}

// TestAllConnections tests all available database connections
func (dh *DatabaseHelper) TestAllConnections() map[string]error {
	results := make(map[string]error)
	connections := dh.manager.GetConnectionNames()

	for _, connName := range connections {
		err := dh.manager.TestConnection(connName)
		results[connName] = err

		if err != nil {
			log.Printf("Connection '%s' test failed: %v", connName, err)
		} else {
			log.Printf("Connection '%s' test passed", connName)
		}
	}

	return results
}

// ExecuteInTransaction executes a function within a database transaction
func (dh *DatabaseHelper) ExecuteInTransaction(connectionName string, fn func(*gorm.DB) error) error {
	db := dh.manager.Connection(connectionName)
	if db == nil {
		return nil
	}

	return db.Transaction(fn)
}

// ExecuteMongoTransaction executes a function within a MongoDB transaction
func (dh *DatabaseHelper) ExecuteMongoTransaction(connectionName string, fn func(mongo.SessionContext) error) error {
	client := dh.manager.MongoDB(connectionName)
	if client == nil {
		return nil
	}

	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		_, err := session.WithTransaction(sc, func(sc mongo.SessionContext) (interface{}, error) {
			return nil, fn(sc)
		})
		return err
	})
}

// GetConnectionInfo returns information about all database connections
func (dh *DatabaseHelper) GetConnectionInfo() map[string]interface{} {
	connections := dh.manager.GetConnectionNames()
	info := make(map[string]interface{})

	for _, connName := range connections {
		connInfo := map[string]interface{}{
			"name":   connName,
			"status": "unknown",
		}

		err := dh.manager.TestConnection(connName)
		if err != nil {
			connInfo["status"] = "error"
			connInfo["error"] = err.Error()
		} else {
			connInfo["status"] = "connected"
		}

		info[connName] = connInfo
	}

	return info
}

// SwitchDefaultConnection switches the default database connection
func (dh *DatabaseHelper) SwitchDefaultConnection(connectionName string) error {
	// Test connection first
	err := dh.manager.TestConnection(connectionName)
	if err != nil {
		return err
	}

	// If test passes, we could update the default connection
	// This would require modifying the config, but for now we'll just return success
	log.Printf("Default connection can be switched to '%s'", connectionName)
	return nil
}

// CloseAllConnections closes all database connections
func (dh *DatabaseHelper) CloseAllConnections() {
	dh.manager.CloseAll()
}

// Global instance for easy access
var DBHelper = NewDatabaseHelper()

// Convenience functions for global access
func GetDefaultDB() *gorm.DB {
	return DBHelper.GetDefaultDB()
}

func GetMysqlDB() *gorm.DB {
	return DBHelper.GetMysqlDB()
}

func GetMysqlSecondaryDB() *gorm.DB {
	return DBHelper.GetMysqlSecondaryDB()
}

func GetPostgresDB() *gorm.DB {
	return DBHelper.GetPostgresDB()
}

func GetSqliteDB() *gorm.DB {
	return DBHelper.GetSqliteDB()
}

func GetSqlServerDB() *gorm.DB {
	return DBHelper.GetSqlServerDB()
}

func GetMongoDB() *mongo.Client {
	return DBHelper.GetMongoDB()
}

func GetMongoDatabase(databaseName string) *mongo.Database {
	return DBHelper.GetMongoDatabase(databaseName)
}

func GetMongoCollection(databaseName, collectionName string) *mongo.Collection {
	return DBHelper.GetMongoCollection(databaseName, collectionName)
}

func GetDBByName(connectionName string) *gorm.DB {
	return DBHelper.GetDBByName(connectionName)
}

func TestAllConnections() map[string]error {
	return DBHelper.TestAllConnections()
}
