package services

import (
	"fmt"

	"golang_starter_kit_2025/facades"

	"gorm.io/gorm"
)

// DatabaseService provides methods to work with multiple databases
type DatabaseService struct{}

// NewDatabaseService creates a new database service instance
func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

// GetDB returns the default database connection
func (s *DatabaseService) GetDB() *gorm.DB {
	return facades.GetDB()
}

// GetMySQL returns the MySQL database connection
func (s *DatabaseService) GetMySQL() (*gorm.DB, error) {
	conn, err := facades.MySQL()
	if err != nil {
		return nil, fmt.Errorf("failed to get MySQL connection: %v", err)
	}
	return conn.DB, nil
}

// GetPostgreSQL returns the PostgreSQL database connection
func (s *DatabaseService) GetPostgreSQL() (*gorm.DB, error) {
	conn, err := facades.PostgreSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to get PostgreSQL connection: %v", err)
	}
	return conn.DB, nil
}

// GetMySQLSecondary returns the secondary MySQL database connection
func (s *DatabaseService) GetMySQLSecondary() (*gorm.DB, error) {
	conn, err := facades.MySQLSecondary()
	if err != nil {
		return nil, fmt.Errorf("failed to get MySQL secondary connection: %v", err)
	}
	return conn.DB, nil
}

// GetConnection returns a specific database connection by name
func (s *DatabaseService) GetConnection(connectionName string) (*gorm.DB, error) {
	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}
	return conn.DB, nil
}

// ExecuteOnMySQL executes a function on MySQL database
func (s *DatabaseService) ExecuteOnMySQL(fn func(*gorm.DB) error) error {
	db, err := s.GetMySQL()
	if err != nil {
		return err
	}
	return fn(db)
}

// ExecuteOnPostgreSQL executes a function on PostgreSQL database
func (s *DatabaseService) ExecuteOnPostgreSQL(fn func(*gorm.DB) error) error {
	db, err := s.GetPostgreSQL()
	if err != nil {
		return err
	}
	return fn(db)
}

// ExecuteOnMySQLSecondary executes a function on secondary MySQL database
func (s *DatabaseService) ExecuteOnMySQLSecondary(fn func(*gorm.DB) error) error {
	db, err := s.GetMySQLSecondary()
	if err != nil {
		return err
	}
	return fn(db)
}

// TransactionOnMySQL executes a transaction on MySQL database
func (s *DatabaseService) TransactionOnMySQL(fn func(*gorm.DB) error) error {
	db, err := s.GetMySQL()
	if err != nil {
		return err
	}
	return db.Transaction(fn)
}

// TransactionOnPostgreSQL executes a transaction on PostgreSQL database
func (s *DatabaseService) TransactionOnPostgreSQL(fn func(*gorm.DB) error) error {
	db, err := s.GetPostgreSQL()
	if err != nil {
		return err
	}
	return db.Transaction(fn)
}

// GetConnectionStats returns statistics for all connections
func (s *DatabaseService) GetConnectionStats() (map[string]interface{}, error) {
	manager := facades.GetManager()
	stats := make(map[string]interface{})

	connections := []string{"mysql", "postgres", "mysql_secondary"}
	for _, connName := range connections {
		// Always try to get the connection (will try to connect if not established)
		conn, err := facades.GetConnection(connName)
		if err == nil && conn != nil && manager.IsConnected(connName) {
			connStats, errStats := manager.GetConnectionStats(connName)
			if errStats == nil {
				stats[connName] = map[string]interface{}{
					"connected":           true,
					"open_connections":    connStats.OpenConnections,
					"in_use":              connStats.InUse,
					"idle":                connStats.Idle,
					"wait_count":          connStats.WaitCount,
					"wait_duration":       connStats.WaitDuration,
					"max_idle_closed":     connStats.MaxIdleClosed,
					"max_lifetime_closed": connStats.MaxLifetimeClosed,
				}
			} else {
				stats[connName] = map[string]interface{}{
					"connected": false,
					"error":     errStats.Error(),
				}
			}
		} else {
			// Try to provide error from GetConnection if available
			errMsg := "Connection not established"
			if err != nil {
				errMsg = err.Error()
			}
			stats[connName] = map[string]interface{}{
				"connected": false,
				"error":     errMsg,
			}
		}
	}

	return stats, nil
}

// CloseConnection closes a specific database connection
func (s *DatabaseService) CloseConnection(connectionName string) error {
	return facades.CloseConnection(connectionName)
}

// CloseAllConnections closes all database connections
func (s *DatabaseService) CloseAllConnections() {
	facades.CloseDB()
}
