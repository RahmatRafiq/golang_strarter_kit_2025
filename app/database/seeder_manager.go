package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"golang_strarter_kit_2025/facades"
)

func ensureSeedsTable() error {
	return facades.DB.Exec(`
		CREATE TABLE IF NOT EXISTS seeds (
			id INT PRIMARY KEY AUTO_INCREMENT,
			filename VARCHAR(255) NOT NULL,
			batch INT NOT NULL,
			seeded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

func CreateSeederFile(name string) error {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)

	rootPath, _ := os.Getwd()
	seedsPath := fmt.Sprintf("%s/app/database/seeds", rootPath)

	if _, err := os.Stat(seedsPath); os.IsNotExist(err) {
		if err := os.MkdirAll(seedsPath, 0755); err != nil {
			return fmt.Errorf("gagal membuat folder seeds: %v", err)
		}
	}

	filePath := fmt.Sprintf("%s/%s", seedsPath, filename)
	content := `-- Contoh:
-- INSERT INTO users (name, email) VALUES ('Admin', 'admin@example.com');
`
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

func isSeederApplied(filename string) (bool, error) {
	var cnt int64
	if err := facades.DB.Raw("SELECT COUNT(*) FROM seeds WHERE filename = ?", filename).Scan(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func RunSeeder(filename string) error {
	filePath := fmt.Sprintf("app/database/seeds/%s.sql", filename)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("gagal membaca file seeder: %v", err)
	}

	statements := parseSQLStatements(string(content))
	for _, stmt := range statements {
		if err := facades.DB.Exec(stmt).Error; err != nil {
			return fmt.Errorf("gagal menjalankan seeder: %v", err)
		}
	}

	return nil
}

func RunAllSeeders() error {
	if err := ensureSeedsTable(); err != nil {
		return err
	}
	files, err := ioutil.ReadDir("app/database/seeds/")
	if err != nil {
		return fmt.Errorf("gagal membaca folder seeder: %v", err)
	}

	batch := time.Now().Unix() // pakai timestamp unik
	var toRun []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			name := strings.TrimSuffix(f.Name(), ".sql")
			applied, _ := isSeederApplied(name)
			if !applied {
				toRun = append(toRun, name)
			}
		}
	}
	sort.Strings(toRun)

	for _, name := range toRun {
		log.Println("🌱 Seeding:", name)
		if err := RunSeeder(name); err != nil {
			return err
		}
		if err := facades.DB.Exec("INSERT INTO seeds (filename, batch) VALUES (?, ?)", name, batch).Error; err != nil {
			return err
		}
	}
	log.Printf("✅ Seeder batch %d applied.\n", batch)
	return nil
}
