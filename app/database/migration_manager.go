package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"golang_strarter_kit_2025/facades"
)

func ensureMigrationsTable() error {
	return facades.DB.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INT PRIMARY KEY AUTO_INCREMENT,
			filename VARCHAR(255) NOT NULL,
			batch INT NOT NULL,
			migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

func getLastBatch() (int, error) {
	var res struct{ Batch int }
	if err := facades.DB.
		Raw("SELECT COALESCE(MAX(batch),0) AS batch FROM migrations").
		Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Batch, nil
}

func isMigrationApplied(filename string) (bool, error) {
	var cnt int64
	if err := facades.DB.
		Raw("SELECT COUNT(*) FROM migrations WHERE filename = ?", filename).
		Scan(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func RunMigration(filename string) error {
	upFilePath := fmt.Sprintf("app/database/migrations/%s.up.sql", filename)
	content, err := ioutil.ReadFile(upFilePath)
	if err != nil {
		return fmt.Errorf("gagal membaca file migrasi: %v", err)
	}

	statements := parseSQLStatements(string(content))

	for _, stmt := range statements {
		if err := facades.DB.Exec(stmt).Error; err != nil {
			return fmt.Errorf("gagal menjalankan migrasi: %v", err)
		}
	}

	return nil
}

func RollbackMigration(filename string) error {
	downFilePath := fmt.Sprintf("app/database/migrations/%s.down.sql", filename)
	content, err := ioutil.ReadFile(downFilePath)
	if err != nil {
		return fmt.Errorf("gagal membaca file rollback migrasi: %v", err)
	}

	statements := parseSQLStatements(string(content))

	for _, stmt := range statements {
		if err := facades.DB.Exec(stmt).Error; err != nil {
			return fmt.Errorf("gagal menjalankan rollback migrasi: %v", err)
		}
	}

	return nil
}

func parseSQLStatements(content string) []string {
	lines := strings.Split(content, "\n")
	cleanedLines := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if idx := strings.Index(line, "--"); idx != -1 {
			line = line[:idx]
		}
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}

		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	cleanedContent := strings.Join(cleanedLines, " ")
	rawStatements := strings.Split(cleanedContent, ";")

	finalStatements := []string{}
	for _, stmt := range rawStatements {
		s := strings.TrimSpace(stmt)
		if s != "" {
			finalStatements = append(finalStatements, s)
		}
	}

	return finalStatements
}

func RunAllMigrations() error {
	if err := ensureMigrationsTable(); err != nil {
		return err
	}
	files, err := ioutil.ReadDir("app/database/migrations/")
	if err != nil {
		return fmt.Errorf("gagal membaca folder: %v", err)
	}
	last, err := getLastBatch()
	if err != nil {
		return err
	}
	batch := last + 1
	var toRun []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			name := strings.TrimSuffix(f.Name(), ".up.sql")
			applied, _ := isMigrationApplied(name)
			if !applied {
				toRun = append(toRun, name)
			}
		}
	}
	sort.Strings(toRun)
	for _, name := range toRun {
		log.Println("🚀 Running:", name)
		if err := RunMigration(name); err != nil {
			return err
		}
		if err := facades.DB.Exec(
			"INSERT INTO migrations (filename, batch) VALUES (?, ?)",
			name, batch,
		).Error; err != nil {
			return err
		}
	}
	log.Printf("✅ Batch %d applied.", batch)
	return nil
}

func RunAllRollbacks() error {
	if err := ensureMigrationsTable(); err != nil {
		return err
	}
	last, err := getLastBatch()
	if err != nil {
		return err
	}
	if last == 0 {
		log.Println("⚠️  No batches to rollback.")
		return nil
	}
	for batch := last; batch >= 1; batch-- {
		log.Printf("🔄 Rolling back batch %d...\n", batch)
		if err := RollbackBatch(batch); err != nil {
			return err
		}
	}
	return nil
}

func RollbackBatch(batch int) error {
	if err := ensureMigrationsTable(); err != nil {
		return err
	}
	var rows []struct{ Filename string }
	facades.DB.
		Raw("SELECT filename FROM migrations WHERE batch = ? ORDER BY id DESC", batch).
		Scan(&rows)
	if len(rows) == 0 {
		log.Printf("⚠️  No migrations in batch %d\n", batch)
		return nil
	}
	for _, r := range rows {
		log.Println("🔄 Rolling back:", r.Filename)
		if err := RollbackMigration(r.Filename); err != nil {
			return err
		}
		facades.DB.Exec("DELETE FROM migrations WHERE filename = ?", r.Filename)
	}
	log.Printf("✅ Batch %d rolled back.\n", batch)
	return nil
}

func RollbackLastBatch() error {
	last, err := getLastBatch()
	if err != nil {
		return err
	}
	if last == 0 {
		log.Println("⚠️  No batch to rollback.")
		return nil
	}
	return RollbackBatch(last)
}

func FreshMigrations() error {
	if err := ensureMigrationsTable(); err != nil {
		return err
	}
	if err := facades.DB.Exec("TRUNCATE TABLE migrations").Error; err != nil {
		return err
	}
	files, err := ioutil.ReadDir("app/database/migrations/")
	if err != nil {
		return fmt.Errorf("gagal membaca folder: %v", err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".up.sql") {
			name := strings.TrimSuffix(f.Name(), ".up.sql")
			if err := RollbackMigration(name); err != nil {
				return err
			}
			facades.DB.Exec("DELETE FROM migrations WHERE filename = ?", name)
		}
	}
	log.Println("✅ All migrations have been rolled back.")
	return nil
}
