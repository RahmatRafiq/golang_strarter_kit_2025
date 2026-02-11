package database

import (
	"fmt"
	"time"

	"golang_starter_kit_2025/facades"
)

// MigrationLog represents a migration execution log entry
type MigrationLog struct {
	ID              int64
	ConnectionName  string
	Filename        string
	Batch           int
	ExecutionTimeMs int64
	Status          string // 'success' or 'failed'
	ErrorMessage    string
	ExecutedAt      time.Time
}

// ensureMigrationLogsTable creates the migration_logs table if it doesn't exist
func ensureMigrationLogsTable(connectionName string) error {
	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var createTableSQL string
	if conn.IsPostgreSQL() {
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS migration_logs (
				id SERIAL PRIMARY KEY,
				connection_name VARCHAR(50) NOT NULL,
				filename VARCHAR(255) NOT NULL,
				batch INT NOT NULL,
				execution_time_ms BIGINT NOT NULL,
				status VARCHAR(20) NOT NULL,
				error_message TEXT,
				executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`
	} else {
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS migration_logs (
				id INT PRIMARY KEY AUTO_INCREMENT,
				connection_name VARCHAR(50) NOT NULL,
				filename VARCHAR(255) NOT NULL,
				batch INT NOT NULL,
				execution_time_ms BIGINT NOT NULL,
				status VARCHAR(20) NOT NULL,
				error_message TEXT,
				executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				INDEX idx_connection_filename (connection_name, filename),
				INDEX idx_status (status)
			)`
	}

	return conn.DB.Exec(createTableSQL).Error
}

// logMigrationExecution logs a migration execution
func logMigrationExecution(connectionName, filename string, batch int, executionTimeMs int64, status string, errorMsg string) error {
	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return err
	}

	return conn.DB.Exec(`
		INSERT INTO migration_logs 
		(connection_name, filename, batch, execution_time_ms, status, error_message) 
		VALUES (?, ?, ?, ?, ?, ?)
	`, connectionName, filename, batch, executionTimeMs, status, errorMsg).Error
}

// GetMigrationLogs retrieves migration logs for a connection
func GetMigrationLogs(connectionName string, limit int) ([]MigrationLog, error) {
	if err := ensureMigrationLogsTable(connectionName); err != nil {
		return nil, err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var logs []MigrationLog
	query := `
		SELECT id, connection_name, filename, batch, execution_time_ms, status, 
		       COALESCE(error_message, '') as error_message, executed_at
		FROM migration_logs
		WHERE connection_name = ?
		ORDER BY executed_at DESC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	err = conn.DB.Raw(query, connectionName).Scan(&logs).Error
	return logs, err
}

// GetSlowMigrations retrieves migrations that took longer than threshold (in milliseconds)
func GetSlowMigrations(connectionName string, thresholdMs int64, limit int) ([]MigrationLog, error) {
	if err := ensureMigrationLogsTable(connectionName); err != nil {
		return nil, err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var logs []MigrationLog
	query := `
		SELECT id, connection_name, filename, batch, execution_time_ms, status,
		       COALESCE(error_message, '') as error_message, executed_at
		FROM migration_logs
		WHERE connection_name = ? AND execution_time_ms > ?
		ORDER BY execution_time_ms DESC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	err = conn.DB.Raw(query, connectionName, thresholdMs).Scan(&logs).Error
	return logs, err
}

// GetFailedMigrations retrieves failed migration attempts
func GetFailedMigrations(connectionName string, limit int) ([]MigrationLog, error) {
	if err := ensureMigrationLogsTable(connectionName); err != nil {
		return nil, err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var logs []MigrationLog
	query := `
		SELECT id, connection_name, filename, batch, execution_time_ms, status,
		       COALESCE(error_message, '') as error_message, executed_at
		FROM migration_logs
		WHERE connection_name = ? AND status = 'failed'
		ORDER BY executed_at DESC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	err = conn.DB.Raw(query, connectionName).Scan(&logs).Error
	return logs, err
}

// GetMigrationStats retrieves statistics about migrations
func GetMigrationStats(connectionName string) (map[string]interface{}, error) {
	if err := ensureMigrationLogsTable(connectionName); err != nil {
		return nil, err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var stats struct {
		TotalExecutions int64
		SuccessCount    int64
		FailedCount     int64
		AvgTimeMs       float64
		MinTimeMs       int64
		MaxTimeMs       int64
	}

	err = conn.DB.Raw(`
		SELECT 
			COUNT(*) as total_executions,
			SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success_count,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed_count,
			AVG(execution_time_ms) as avg_time_ms,
			MIN(execution_time_ms) as min_time_ms,
			MAX(execution_time_ms) as max_time_ms
		FROM migration_logs
		WHERE connection_name = ?
	`, connectionName).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_executions": stats.TotalExecutions,
		"success_count":    stats.SuccessCount,
		"failed_count":     stats.FailedCount,
		"avg_time_ms":      stats.AvgTimeMs,
		"min_time_ms":      stats.MinTimeMs,
		"max_time_ms":      stats.MaxTimeMs,
	}, nil
}

// FormatDuration formats milliseconds to human-readable duration
func FormatDuration(ms int64) string {
	duration := time.Duration(ms) * time.Millisecond

	if duration < time.Second {
		return fmt.Sprintf("%dms", ms)
	}

	seconds := duration.Seconds()
	if seconds < 60 {
		return fmt.Sprintf("%.2fs", seconds)
	}

	minutes := int(seconds / 60)
	remainingSeconds := seconds - float64(minutes*60)
	return fmt.Sprintf("%dm %.1fs", minutes, remainingSeconds)
}
