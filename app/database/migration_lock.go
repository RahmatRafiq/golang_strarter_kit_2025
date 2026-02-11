package database

import (
	"fmt"
	"os"
	"time"

	"golang_starter_kit_2025/facades"

	"github.com/rs/zerolog/log"
)

// MigrationLock manages migration locking to prevent concurrent migrations
type MigrationLock struct {
	connectionName string
	acquired       bool
}

// NewMigrationLock creates a new migration lock for a connection
func NewMigrationLock(connectionName string) *MigrationLock {
	return &MigrationLock{
		connectionName: connectionName,
		acquired:       false,
	}
}

// ensureMigrationLocksTable creates the migration_locks table if it doesn't exist
func ensureMigrationLocksTable(connectionName string) error {
	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	// Use different table creation syntax based on database type
	var createTableSQL string
	if conn.IsPostgreSQL() {
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS migration_locks (
				id SERIAL PRIMARY KEY,
				connection_name VARCHAR(50) NOT NULL UNIQUE,
				locked_by VARCHAR(255) NOT NULL,
				locked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`
	} else {
		// MySQL/MariaDB
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS migration_locks (
				id INT PRIMARY KEY AUTO_INCREMENT,
				connection_name VARCHAR(50) NOT NULL UNIQUE,
				locked_by VARCHAR(255) NOT NULL,
				locked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`
	}

	return conn.DB.Exec(createTableSQL).Error
}

// Acquire attempts to acquire the migration lock
func (ml *MigrationLock) Acquire() error {
	// Ensure table exists
	if err := ensureMigrationLocksTable(ml.connectionName); err != nil {
		return fmt.Errorf("failed to ensure migration_locks table: %v", err)
	}

	conn, err := facades.GetConnection(ml.connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", ml.connectionName, err)
	}

	// Get hostname for lock tracking
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Add PID for better identification
	lockedBy := fmt.Sprintf("%s (PID: %d)", hostname, os.Getpid())

	// Try to acquire lock
	result := conn.DB.Exec(`
		INSERT INTO migration_locks (connection_name, locked_by) 
		VALUES (?, ?)
	`, ml.connectionName, lockedBy)

	if result.Error != nil {
		// Check if lock already exists
		var existingLock struct {
			LockedAt time.Time
			LockedBy string
		}

		conn.DB.Raw(`
			SELECT locked_by, locked_at 
			FROM migration_locks 
			WHERE connection_name = ?
		`, ml.connectionName).Scan(&existingLock)

		return fmt.Errorf(
			"migration already locked by '%s' since %s. "+
				"If this is a stale lock, run: DELETE FROM migration_locks WHERE connection_name = '%s'",
			existingLock.LockedBy,
			existingLock.LockedAt.Format("2006-01-02 15:04:05"),
			ml.connectionName,
		)
	}

	ml.acquired = true
	log.Info().Str("connection", ml.connectionName).Msg("Migration lock acquired")
	return nil
}

// Release releases the migration lock
func (ml *MigrationLock) Release() error {
	if !ml.acquired {
		return nil // Nothing to release
	}

	conn, err := facades.GetConnection(ml.connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", ml.connectionName, err)
	}

	err = conn.DB.Exec(
		"DELETE FROM migration_locks WHERE connection_name = ?",
		ml.connectionName,
	).Error

	if err != nil {
		return fmt.Errorf("failed to release lock: %v", err)
	}

	ml.acquired = false
	log.Info().Str("connection", ml.connectionName).Msg("Migration lock released")
	return nil
}

// ForceRelease forcefully releases a migration lock (use with caution)
func ForceReleaseLock(connectionName string) error {
	if err := ensureMigrationLocksTable(connectionName); err != nil {
		return err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	result := conn.DB.Exec(
		"DELETE FROM migration_locks WHERE connection_name = ?",
		connectionName,
	)

	if result.Error != nil {
		return fmt.Errorf("failed to force release lock: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Warn().Str("connection", connectionName).Msg("No lock found")
	} else {
		log.Info().Str("connection", connectionName).Msg("Forcefully released lock")
	}

	return nil
}

// CheckLockStatus checks if a migration lock exists for a connection
func CheckLockStatus(connectionName string) (bool, string, error) {
	if err := ensureMigrationLocksTable(connectionName); err != nil {
		return false, "", err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return false, "", fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var lock struct {
		LockedAt time.Time
		LockedBy string
	}

	result := conn.DB.Raw(`
		SELECT locked_by, locked_at 
		FROM migration_locks 
		WHERE connection_name = ?
	`, connectionName).Scan(&lock)

	if result.Error != nil {
		return false, "", result.Error
	}

	if result.RowsAffected == 0 {
		return false, "", nil
	}

	info := fmt.Sprintf("Locked by '%s' since %s",
		lock.LockedBy,
		lock.LockedAt.Format("2006-01-02 15:04:05"),
	)

	return true, info, nil
}
