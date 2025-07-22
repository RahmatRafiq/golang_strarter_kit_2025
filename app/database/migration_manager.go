package database

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang_starter_kit_2025/facades"
)

const (
	upMarker   = "-- +++ UP Migration"
	downMarker = "-- --- DOWN Migration"
)

// ensureMigrationsTable creates migrations table for a specific connection
func ensureMigrationsTable(connectionName string) error {
	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	// Use different table creation syntax based on database type
	var createTableSQL string
	if conn.IsPostgreSQL() {
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS migrations (
				id SERIAL PRIMARY KEY,
				filename VARCHAR(255) NOT NULL,
				batch INTEGER NOT NULL,
				migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`
	} else {
		// MySQL/MariaDB
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS migrations (
				id INT PRIMARY KEY AUTO_INCREMENT,
				filename VARCHAR(255) NOT NULL,
				batch INT NOT NULL,
				migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
	if err := conn.DB.Raw("SELECT COALESCE(MAX(batch),0) AS batch FROM migrations").Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Batch, nil
}

func isMigrationApplied(filename, connectionName string) (bool, error) {
	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return false, fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var cnt int64
	if err := conn.DB.Raw("SELECT COUNT(*) FROM migrations WHERE filename = ?", filename).Scan(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
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

// RunMigration runs a specific migration on the default connection
func RunMigration(filename string) error {
	return RunMigrationOnConnection(filename, "")
}

// RunMigrationOnConnection runs a specific migration on a specified connection
func RunMigrationOnConnection(filename, connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql" // default connection
	}

	if err := ensureMigrationsTable(connectionName); err != nil {
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
		return fmt.Errorf("gagal membaca file migrasi: %v", err)
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	ups, _ := parseMigrationFile(string(data))
	for _, sql := range ups {
		if err := conn.DB.Exec(sql).Error; err != nil {
			return fmt.Errorf("gagal menjalankan migrasi: %v", err)
		}
	}

	if err := conn.DB.Exec(
		"INSERT INTO migrations(filename,batch) VALUES(?,?)", filename, batch,
	).Error; err != nil {
		return fmt.Errorf("gagal mencatat migrasi: %v", err)
	}

	fmt.Printf("Migrated: %s\n", filename)
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
		return fmt.Errorf("gagal membaca file rollback: %v", err)
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	_, downs := parseMigrationFile(string(data))
	for _, sql := range downs {
		if err := conn.DB.Exec(sql).Error; err != nil {
			return fmt.Errorf("gagal rollback: %v", err)
		}
	}

	// Remove from migrations table
	if err := conn.DB.Exec("DELETE FROM migrations WHERE filename=?", filename).Error; err != nil {
		return fmt.Errorf("gagal menghapus record migrasi: %v", err)
	}

	fmt.Printf("Rolled back: %s\n", filename)
	return nil
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

// RunAllMigrationsOnConnection runs all pending migrations on a specified connection
func RunAllMigrationsOnConnection(connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql" // default connection
	}

	if err := ensureMigrationsTable(connectionName); err != nil {
		return err
	}

	var lastBatch struct{ Batch int }
	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	if err := conn.DB.Raw(
		"SELECT COALESCE(MAX(batch),0) AS batch FROM migrations",
	).Scan(&lastBatch).Error; err != nil {
		return err
	}
	batch := lastBatch.Batch + 1

	files, err := os.ReadDir("app/database/migrations")
	if err != nil {
		return fmt.Errorf("gagal baca folder: %v", err)
	}
	var toRun []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			name := strings.TrimSuffix(f.Name(), ".sql")
			var cnt int64
			conn.DB.Raw(
				"SELECT COUNT(*) FROM migrations WHERE filename = ?", name,
			).Scan(&cnt)
			if cnt == 0 {
				toRun = append(toRun, name)
			}
		}
	}
	sort.Strings(toRun)

	for _, name := range toRun {
		fmt.Printf("Migrating: %s\n", name)

		data, err := os.ReadFile(
			fmt.Sprintf("app/database/migrations/%s.sql", name),
		)
		if err != nil {
			return fmt.Errorf("gagal membaca %s: %v", name, err)
		}
		parts := strings.Split(
			string(data), "-- --- DOWN Migration",
		)
		up := strings.Replace(
			parts[0], "-- +++ UP Migration", "", 1,
		)

		for _, stmt := range parseSQLStatements(up) {
			if err := conn.DB.Exec(stmt).Error; err != nil {
				return fmt.Errorf("gagal %s: %v", name, err)
			}
		}

		if err := conn.DB.Exec(
			"INSERT INTO migrations(filename,batch) VALUES(?,?)",
			name, batch,
		).Error; err != nil {
			return fmt.Errorf("gagal mencatat %s: %v", name, err)
		}
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
	last, _ := getLastBatch(connectionName)
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
	conn.DB.Raw("SELECT filename FROM migrations WHERE batch=? ORDER BY id DESC", batch).Scan(&rows)
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

	last, _ := getLastBatch(connectionName)
	if last == 0 {
		fmt.Printf("No batch to rollback.\n")
		return nil
	}
	return RollbackBatchOnConnection(last, connectionName)
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

	// Use different truncate syntax for PostgreSQL
	if conn.IsPostgreSQL() {
		conn.DB.Exec("TRUNCATE migrations RESTART IDENTITY")
	} else {
		conn.DB.Exec("TRUNCATE migrations")
	}

	return RunAllMigrationsOnConnection(connectionName)
}
