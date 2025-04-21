package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"golang_strarter_kit_2025/facades"
)

// CreateMigrationFile membuat file migrasi baru
func CreateMigrationFile(name string) error {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s", timestamp, name)

	rootPath, _ := os.Getwd()
	migrationPath := fmt.Sprintf("%s/app/database/migrations", rootPath)

	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		if err := os.MkdirAll(migrationPath, 0755); err != nil {
			return fmt.Errorf("gagal membuat folder migrations: %v", err)
		}
	}

	upFile := fmt.Sprintf("%s/%s.up.sql", migrationPath, filename)
	downFile := fmt.Sprintf("%s/%s.down.sql", migrationPath, filename)

	writeTemplate(upFile, "-- +++ Write your UP migration here\n")
	writeTemplate(downFile, "-- --- Write your DOWN migration here\n")

	fmt.Println("Migration files created:")
	fmt.Println(" -", upFile)
	fmt.Println(" -", downFile)

	return nil
}

func writeTemplate(path, content string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("gagal membuat file %s: %v", path, err)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		log.Fatalf("gagal menulis ke file %s: %v", path, err)
	}
}

// RunMigration menjalankan file migrasi .up.sql
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

// RollbackMigration menjalankan file migrasi .down.sql
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

// parseSQLStatements memisahkan dan membersihkan statement SQL
func parseSQLStatements(content string) []string {
	rawStatements := strings.Split(content, ";")
	cleaned := []string{}

	for _, stmt := range rawStatements {
		s := strings.TrimSpace(stmt)
		if s == "" || strings.HasPrefix(s, "--") || strings.HasPrefix(s, "#") {
			continue
		}
		cleaned = append(cleaned, s)
	}

	return cleaned
}

// RunAllMigrations menjalankan semua migrasi yang ada di folder migrations
func RunAllMigrations() error {
	migrationPath := "app/database/migrations/"
	files, err := ioutil.ReadDir(migrationPath)
	if err != nil {
		return fmt.Errorf("gagal membaca folder migrasi: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			filename := strings.TrimSuffix(file.Name(), ".up.sql")
			if err := RunMigration(filename); err != nil {
				return fmt.Errorf("gagal menjalankan migrasi %s: %v", filename, err)
			}
		}
	}

	return nil
}
