package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"golang_starter_kit_2025/facades"
)

const (
	upMarker   = "-- +++ UP Migration"
	downMarker = "-- --- DOWN Migration"
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
	if err := facades.DB.Raw("SELECT COALESCE(MAX(batch),0) AS batch FROM migrations").Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Batch, nil
}

func isMigrationApplied(filename string) (bool, error) {
	var cnt int64
	if err := facades.DB.Raw("SELECT COUNT(*) FROM migrations WHERE filename = ?", filename).Scan(&cnt).Error; err != nil {
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

func RunMigration(filename string) error {

	if err := ensureMigrationsTable(); err != nil {
		return err
	}

	last, err := getLastBatch()
	if err != nil {
		return err
	}
	batch := last + 1

	path := fmt.Sprintf("app/database/migrations/%s.sql", filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("gagal membaca file migrasi: %v", err)
	}

	ups, _ := parseMigrationFile(string(data))
	for _, sql := range ups {
		if err := facades.DB.Exec(sql).Error; err != nil {
			return fmt.Errorf("gagal menjalankan migrasi: %v", err)
		}
	}

	if err := facades.DB.Exec(
		"INSERT INTO migrations(filename,batch) VALUES(?,?)", filename, batch,
	).Error; err != nil {
		return fmt.Errorf("gagal mencatat migrasi: %v", err)
	}

	return nil
}

func RollbackMigration(filename string) error {
	path := fmt.Sprintf("app/database/migrations/%s.sql", filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("gagal membaca file rollback: %v", err)
	}
	_, downs := parseMigrationFile(string(data))
	for _, sql := range downs {
		if err := facades.DB.Exec(sql).Error; err != nil {
			return fmt.Errorf("gagal rollback: %v", err)
		}
	}
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

func RunAllMigrations() error {

	if err := facades.DB.Exec(`
        CREATE TABLE IF NOT EXISTS migrations (
            id INT PRIMARY KEY AUTO_INCREMENT,
            filename VARCHAR(255) NOT NULL,
            batch INT NOT NULL,
            migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`).Error; err != nil {
		return err
	}

	var lastBatch struct{ Batch int }
	if err := facades.DB.Raw(
		"SELECT COALESCE(MAX(batch),0) AS batch FROM migrations",
	).Scan(&lastBatch).Error; err != nil {
		return err
	}
	batch := lastBatch.Batch + 1

	files, err := ioutil.ReadDir("app/database/migrations")
	if err != nil {
		return fmt.Errorf("gagal baca folder: %v", err)
	}
	var toRun []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			name := strings.TrimSuffix(f.Name(), ".sql")
			var cnt int64
			facades.DB.Raw(
				"SELECT COUNT(*) FROM migrations WHERE filename = ?", name,
			).Scan(&cnt)
			if cnt == 0 {
				toRun = append(toRun, name)
			}
		}
	}
	sort.Strings(toRun)

	for _, name := range toRun {
		log.Println("🚀 Running", name)

		data, err := ioutil.ReadFile(
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
			if err := facades.DB.Exec(stmt).Error; err != nil {
				return fmt.Errorf("gagal %s: %v", name, err)
			}
		}

		if err := facades.DB.Exec(
			"INSERT INTO migrations(filename,batch) VALUES(?,?)",
			name, batch,
		).Error; err != nil {
			return fmt.Errorf("gagal mencatat %s: %v", name, err)
		}
	}

	log.Printf("✅ Batch %d applied", batch)
	return nil
}

func RunAllRollbacks() error {
	if err := ensureMigrationsTable(); err != nil {
		return err
	}
	last, _ := getLastBatch()
	for b := last; b >= 1; b-- {
		if err := RollbackBatch(b); err != nil {
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
	facades.DB.Raw("SELECT filename FROM migrations WHERE batch=? ORDER BY id DESC", batch).Scan(&rows)
	for _, r := range rows {
		log.Println("🔄 Rollback", r.Filename)
		if err := RollbackMigration(r.Filename); err != nil {
			return err
		}
		facades.DB.Exec("DELETE FROM migrations WHERE filename=?", r.Filename)
	}
	log.Printf("✅ Batch %d rolled back", batch)
	return nil
}

func RollbackLastBatch() error {
	last, _ := getLastBatch()
	if last == 0 {
		log.Println("⚠️ No batch to rollback")
		return nil
	}
	return RollbackBatch(last)
}

func FreshMigrations() error {
	if err := ensureMigrationsTable(); err != nil {
		return err
	}
	facades.DB.Exec("TRUNCATE migrations")
	return RunAllMigrations()
}
