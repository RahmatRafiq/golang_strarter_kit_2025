package database

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"golang_starter_kit_2025/facades"

	"gorm.io/gorm"
)

const (
	upMarker   = "-- +++ UP Migration"
	downMarker = "-- --- DOWN Migration"
)

func ensureMigrationsTable(connectionName string) error {
	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var createTableSQL string
	if conn.IsPostgreSQL() {
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS migrations (
				id SERIAL PRIMARY KEY,
				connection_name VARCHAR(50) NOT NULL,
				filename VARCHAR(255) NOT NULL,
				batch INTEGER NOT NULL,
				migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE (connection_name, filename)
			)`
	} else {
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS migrations (
				id INT PRIMARY KEY AUTO_INCREMENT,
				connection_name VARCHAR(50) NOT NULL,
				filename VARCHAR(255) NOT NULL,
				batch INT NOT NULL,
				migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE KEY unique_migration (connection_name, filename)
			)`
	}

	return conn.DB.Exec(createTableSQL).Error
}

func getLastBatch(connectionName string) (int, error) {
	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var res struct{ Batch int }
	if err := conn.DB.Raw("SELECT COALESCE(MAX(batch),0) AS batch FROM migrations WHERE connection_name = ?", connectionName).Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Batch, nil
}

func parseMigrationFile(content string) (upStmts, downStmts []string) {
	parts := strings.Split(content, downMarker)
	upPart := parts[0]
	downPart := ""
	if len(parts) > 1 {
		downPart = parts[1]
	}
	upPart = strings.Replace(upPart, upMarker, "", 1)
	return parseSQLStatements(upPart), parseSQLStatements(downPart)
}

func RunMigration(filename string) error {
	return RunMigrationOnConnection(filename, "")
}

func RunMigrationOnConnection(filename, connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql"
	}

	if err := ensureMigrationsTable(connectionName); err != nil {
		return err
	}

	if err := ensureMigrationLogsTable(connectionName); err != nil {
		return err
	}

	last, err := getLastBatch(connectionName)
	if err != nil {
		return err
	}
	batch := last + 1

	path := fmt.Sprintf("app/database/migrations/%s.sql", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	// Start timing
	startTime := time.Now()

	ups, _ := parseMigrationFile(string(data))
	var execError error

	for _, sql := range ups {
		if err := conn.DB.Exec(sql).Error; err != nil {
			execError = fmt.Errorf("failed to execute migration: %v", err)
			break
		}
	}

	// Calculate execution time
	executionTime := time.Since(startTime).Milliseconds()

	// If execution failed, log and return error
	if execError != nil {
		if logErr := logMigrationExecution(connectionName, filename, batch, executionTime, "failed", execError.Error()); logErr != nil {
			fmt.Printf("Warning: failed to log migration execution: %v\n", logErr)
		}
		return execError
	}

	// Record migration
	if err := conn.DB.Exec(
		"INSERT INTO migrations(connection_name,filename,batch) VALUES(?,?,?)", connectionName, filename, batch,
	).Error; err != nil {
		if logErr := logMigrationExecution(connectionName, filename, batch, executionTime, "failed", err.Error()); logErr != nil {
			fmt.Printf("Warning: failed to log migration execution: %v\n", logErr)
		}
		return fmt.Errorf("failed to record migration: %v", err)
	}

	// Log successful execution
	if err := logMigrationExecution(connectionName, filename, batch, executionTime, "success", ""); err != nil {
		fmt.Printf("Warning: failed to log migration execution: %v\n", err)
	}

	fmt.Printf("Migrated: %s (%s)\n", filename, FormatDuration(executionTime))
	return nil
}

// RollbackMigration rolls back a specific migration on the default connection
func RollbackMigration(filename string) error {
	return RollbackMigrationOnConnection(filename, "")
}

// RollbackMigrationOnConnection rolls back a specific migration on a specified connection
func RollbackMigrationOnConnection(filename, connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql" // default connection
	}

	path := fmt.Sprintf("app/database/migrations/%s.sql", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read rollback file: %v", err)
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	// Use transaction for atomic rollback
	_, downs := parseMigrationFile(string(data))

	return conn.DB.Transaction(func(tx *gorm.DB) error {
		for _, sql := range downs {
			if err := tx.Exec(sql).Error; err != nil {
				return fmt.Errorf("failed to rollback: %v", err)
			}
		}

		// Remove from migrations table with connection_name filter
		if err := tx.Exec("DELETE FROM migrations WHERE connection_name = ? AND filename = ?", connectionName, filename).Error; err != nil {
			return fmt.Errorf("failed to delete migration record: %v", err)
		}

		fmt.Printf("Rolled back: %s\n", filename)
		return nil
	})
}

func parseSQLStatements(content string) []string {
	var stmts []string
	for _, s := range strings.Split(content, ";") {
		if t := strings.TrimSpace(s); t != "" {
			stmts = append(stmts, t)
		}
	}
	return stmts
}

// RunAllMigrations runs all pending migrations on the default connection
func RunAllMigrations() error {
	return RunAllMigrationsOnConnection("")
}

// RunAllMigrationsOnConnection runs all pending migrations on a specified connection with lock
func RunAllMigrationsOnConnection(connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql" // default connection
	}

	// Acquire migration lock
	lock := NewMigrationLock(connectionName)
	if err := lock.Acquire(); err != nil {
		return fmt.Errorf("failed to acquire migration lock: %v", err)
	}
	defer func() {
		if err := lock.Release(); err != nil {
			fmt.Printf("Warning: failed to release migration lock: %v\n", err)
		}
	}()

	if err := ensureMigrationsTable(connectionName); err != nil {
		return err
	}

	var lastBatch struct{ Batch int }
	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	if queryErr := conn.DB.Raw(
		"SELECT COALESCE(MAX(batch),0) AS batch FROM migrations WHERE connection_name = ?", connectionName,
	).Scan(&lastBatch).Error; queryErr != nil {
		return queryErr
	}
	batch := lastBatch.Batch + 1

	files, err := os.ReadDir("app/database/migrations")
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}
	var toRun []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			name := strings.TrimSuffix(f.Name(), ".sql")
			var cnt int64
			conn.DB.Raw(
				"SELECT COUNT(*) FROM migrations WHERE connection_name = ? AND filename = ?", connectionName, name,
			).Scan(&cnt)
			if cnt == 0 {
				toRun = append(toRun, name)
			}
		}
	}
	sort.Strings(toRun)

	// Ensure migration_logs table exists
	if err := ensureMigrationLogsTable(connectionName); err != nil {
		return fmt.Errorf("failed to ensure migration_logs table: %v", err)
	}

	for _, name := range toRun {
		fmt.Printf("Migrating: %s\n", name)

		// Start timing
		startTime := time.Now()

		data, err := os.ReadFile(
			fmt.Sprintf("app/database/migrations/%s.sql", name),
		)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", name, err)
		}
		parts := strings.Split(
			string(data), "-- --- DOWN Migration",
		)
		up := strings.Replace(
			parts[0], "-- +++ UP Migration", "", 1,
		)

		var execError error
		for _, stmt := range parseSQLStatements(up) {
			if err := conn.DB.Exec(stmt).Error; err != nil {
				execError = fmt.Errorf("failed to %s: %v", name, err)
				break
			}
		}

		// Calculate execution time
		executionTime := time.Since(startTime).Milliseconds()

		// If execution failed, log and return error
		if execError != nil {
			if logErr := logMigrationExecution(connectionName, name, batch, executionTime, "failed", execError.Error()); logErr != nil {
				fmt.Printf("Warning: failed to log migration execution: %v\n", logErr)
			}
			return execError
		}

		if err := conn.DB.Exec(
			"INSERT INTO migrations(connection_name,filename,batch) VALUES(?,?,?)",
			connectionName, name, batch,
		).Error; err != nil {
			if logErr := logMigrationExecution(connectionName, name, batch, executionTime, "failed", err.Error()); logErr != nil {
				fmt.Printf("Warning: failed to log migration execution: %v\n", logErr)
			}
			return fmt.Errorf("failed to record %s: %v", name, err)
		}

		// Log successful execution
		if err := logMigrationExecution(connectionName, name, batch, executionTime, "success", ""); err != nil {
			fmt.Printf("Warning: failed to log migration execution: %v\n", err)
		}
		fmt.Printf("  ✓ Completed in %s\n", FormatDuration(executionTime))
	}

	fmt.Printf("Batch %d applied.\n", batch)
	return nil
}

// RunAllRollbacks rolls back all migrations on the default connection
func RunAllRollbacks() error {
	return RunAllRollbacksOnConnection("")
}

// RunAllRollbacksOnConnection rolls back all migrations on a specified connection
func RunAllRollbacksOnConnection(connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql" // default connection
	}

	if err := ensureMigrationsTable(connectionName); err != nil {
		return err
	}
	last, err := getLastBatch(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get last batch: %v", err)
	}
	for b := last; b >= 1; b-- {
		if err := RollbackBatchOnConnection(b, connectionName); err != nil {
			return err
		}
	}
	return nil
}

// RollbackBatch rolls back a specific batch on the default connection
func RollbackBatch(batch int) error {
	return RollbackBatchOnConnection(batch, "")
}

// RollbackBatchOnConnection rolls back a specific batch on a specified connection
func RollbackBatchOnConnection(batch int, connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql" // default connection
	}

	if err := ensureMigrationsTable(connectionName); err != nil {
		return err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var rows []struct{ Filename string }
	conn.DB.Raw("SELECT filename FROM migrations WHERE connection_name = ? AND batch = ? ORDER BY id DESC", connectionName, batch).Scan(&rows)
	for _, r := range rows {
		fmt.Printf("Rolling back: %s\n", r.Filename)
		if err := RollbackMigrationOnConnection(r.Filename, connectionName); err != nil {
			return err
		}
	}
	fmt.Printf("Batch %d rolled back.\n", batch)
	return nil
}

// RollbackLastBatch rolls back the last batch on the default connection
func RollbackLastBatch() error {
	return RollbackLastBatchOnConnection("")
}

// RollbackLastBatchOnConnection rolls back the last batch on a specified connection
func RollbackLastBatchOnConnection(connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql" // default connection
	}

	last, err := getLastBatch(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get last batch: %v", err)
	}
	if last == 0 {
		fmt.Printf("No batch to rollback.\n")
		return nil
	}
	return RollbackBatchOnConnection(last, connectionName)
}

// RollbackStepsOnConnection rolls back the last N batches on a specified connection
func RollbackStepsOnConnection(steps int, connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql" // default connection
	}

	if err := ensureMigrationsTable(connectionName); err != nil {
		return err
	}

	last, err := getLastBatch(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get last batch: %v", err)
	}
	if last == 0 {
		fmt.Println("⚠️ No batches to rollback.")
		return nil
	}

	// Calculate target batch (we rollback from last down to last-steps+1)
	targetBatch := last - steps + 1
	if targetBatch < 1 {
		targetBatch = 1
	}

	count := 0
	for b := last; b >= targetBatch && b >= 1; b-- {
		if err := RollbackBatchOnConnection(b, connectionName); err != nil {
			return fmt.Errorf("failed to rollback batch %d: %v", b, err)
		}
		count++
	}

	fmt.Printf("✅ Successfully rolled back %d batch(es).\n", count)
	return nil
}

// FreshMigrations truncates migrations and re-runs all on the default connection
func FreshMigrations() error {
	return FreshMigrationsOnConnection("")
}

// FreshMigrationsOnConnection truncates migrations and re-runs all on a specified connection
func FreshMigrationsOnConnection(connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql" // default connection
	}

	if err := ensureMigrationsTable(connectionName); err != nil {
		return err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	// Delete only records for this connection (don't wipe other connections)
	conn.DB.Exec("DELETE FROM migrations WHERE connection_name = ?", connectionName)

	return RunAllMigrationsOnConnection(connectionName)
}

// ShowMigrationStatus displays the status of all migrations with execution times
func ShowMigrationStatus(connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql"
	}

	if err := ensureMigrationsTable(connectionName); err != nil {
		return err
	}

	// Ensure migration_logs table exists
	if err := ensureMigrationLogsTable(connectionName); err != nil {
		return fmt.Errorf("failed to ensure migration_logs table: %v", err)
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	// Get all migration files
	files, err := os.ReadDir("app/database/migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations folder: %v", err)
	}

	// Get applied migrations with execution time
	var applied []struct {
		Filename        string
		Batch           int
		ExecutionTimeMs int64
	}
	conn.DB.Raw(`
		SELECT m.filename, m.batch, COALESCE(ml.execution_time_ms, 0) as execution_time_ms
		FROM migrations m
		LEFT JOIN (
			SELECT filename, execution_time_ms 
			FROM migration_logs 
			WHERE connection_name = ? AND status = 'success'
			ORDER BY executed_at DESC
		) ml ON m.filename = ml.filename
		WHERE m.connection_name = ?
		ORDER BY m.batch, m.filename
	`, connectionName, connectionName).Scan(&applied)

	appliedMap := make(map[string]struct {
		Batch           int
		ExecutionTimeMs int64
	})
	for _, a := range applied {
		appliedMap[a.Filename] = struct {
			Batch           int
			ExecutionTimeMs int64
		}{a.Batch, a.ExecutionTimeMs}
	}

	fmt.Println()
	fmt.Println("Migration Status on connection:", connectionName)
	fmt.Println(strings.Repeat("=", 95))
	fmt.Printf("%-50s %-10s %-15s %-15s\n", "Migration", "Batch", "Status", "Execution Time")
	fmt.Println(strings.Repeat("-", 95))

	var pending, ran int
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		name := strings.TrimSuffix(f.Name(), ".sql")
		if info, ok := appliedMap[name]; ok {
			timeStr := "-"
			if info.ExecutionTimeMs > 0 {
				timeStr = FormatDuration(info.ExecutionTimeMs)
			}
			fmt.Printf("%-50s %-10d %-15s %-15s\n",
				truncateString(name, 50),
				info.Batch,
				"✅ Ran",
				timeStr)
			ran++
		} else {
			fmt.Printf("%-50s %-10s %-15s %-15s\n",
				truncateString(name, 50),
				"-",
				"⏳ Pending",
				"-")
			pending++
		}
	}

	fmt.Println(strings.Repeat("=", 95))
	fmt.Printf("Total: %d migrations (%d ran, %d pending)\n", ran+pending, ran, pending)

	// Show statistics if there are ran migrations
	if ran > 0 {
		stats, err := GetMigrationStats(connectionName)
		if err != nil {
			fmt.Printf("Warning: failed to get migration statistics: %v\n", err)
		} else if stats != nil {
			fmt.Println()
			fmt.Println("Execution Statistics:")

			// Use comma-ok idiom for type assertions
			if avgTime, ok := stats["avg_time_ms"].(float64); ok {
				fmt.Printf("  Average Time: %s\n", FormatDuration(int64(avgTime)))
			}
			if minTime, ok := stats["min_time_ms"].(int64); ok {
				fmt.Printf("  Fastest: %s\n", FormatDuration(minTime))
			}
			if maxTime, ok := stats["max_time_ms"].(int64); ok {
				fmt.Printf("  Slowest: %s\n", FormatDuration(maxTime))
			}
		}
	}

	fmt.Println()
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// WipeDatabase drops all tables from the database
func WipeDatabase(connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql"
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var tables []string

	if conn.IsPostgreSQL() {
		// PostgreSQL: Get all tables from public schema
		var rows []struct{ TableName string }
		conn.DB.Raw(`
			SELECT tablename as table_name
			FROM pg_tables
			WHERE schemaname = 'public'
		`).Scan(&rows)

		for _, row := range rows {
			tables = append(tables, row.TableName)
		}

		// Drop tables with CASCADE
		for _, table := range tables {
			fmt.Printf("Dropping table: %s\n", table)
			if err := conn.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
				return fmt.Errorf("failed to drop table %s: %v", table, err)
			}
		}
	} else {
		// MySQL/MariaDB: Disable foreign key checks temporarily
		conn.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
		defer conn.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

		// Get all tables
		var rows []struct{ TableName string }
		conn.DB.Raw("SHOW TABLES").Scan(&rows)

		for _, row := range rows {
			tables = append(tables, row.TableName)
		}

		// Drop each table
		for _, table := range tables {
			fmt.Printf("Dropping table: %s\n", table)
			if err := conn.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)).Error; err != nil {
				return fmt.Errorf("failed to drop table %s: %v", table, err)
			}
		}
	}

	fmt.Printf("✅ Successfully dropped %d tables.\n", len(tables))
	return nil
}
