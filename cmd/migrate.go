package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang_strarter_kit_2025/app/database"

	"github.com/urfave/cli/v2"
)

var MigrationCommand = &cli.Command{
	Name:  "migrate",
	Usage: "Run migration for given file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "file",
			Usage:    "Nama file migration tanpa ekstensi",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		filename := c.String("file")
		fmt.Println("🚀 Menjalankan migration untuk:", filename)

		if err := database.RunMigration(filename); err != nil {
			log.Fatal("❌ Migration gagal:", err)
		}

		fmt.Println("✅ Migration berhasil dijalankan!")
		return nil
	},
}

var RollbackCommand = &cli.Command{
	Name:  "rollback",
	Usage: "Rollback migration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "file",
			Usage:    "Nama file migration tanpa ekstensi",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		filename := c.String("file")
		fmt.Println("🔄 Menjalankan rollback untuk:", filename)

		if err := database.RollbackMigration(filename); err != nil {
			log.Fatal("❌ Rollback gagal:", err)
		}

		fmt.Println("✅ Rollback berhasil!")
		return nil
	},
}

var MakeMigrationCommand = &cli.Command{
	Name:  "make:migration",
	Usage: "Generate a new migration file",
	Action: func(c *cli.Context) error {
		name := c.Args().First()
		if name == "" {
			return fmt.Errorf("nama migrasi harus disediakan. Contoh: make:migration create_users_table")
		}

		err := CreateMigrationFile(name)
		if err != nil {
			log.Fatal("❌ Gagal membuat file migrasi:", err)
		}

		fmt.Println("✅ File migrasi berhasil dibuat.")
		return nil
	},
}
var MigrateAllCommand = &cli.Command{
	Name:  "migrate:all",
	Usage: "Run all migrations",
	Action: func(c *cli.Context) error {
		fmt.Println("🚀 Menjalankan semua migrasi...")

		if err := database.RunAllMigrations(); err != nil {
			log.Fatal("❌ Gagal menjalankan semua migrasi:", err)
		}

		fmt.Println("✅ Semua migrasi berhasil dijalankan!")
		return nil
	},
}

var RollbackAllCommand = &cli.Command{
	Name:  "rollback:all",
	Usage: "Rollback all migrations",
	Action: func(c *cli.Context) error {
		fmt.Println("🔄 Menjalankan rollback untuk semua migrasi...")

		if err := database.RunAllRollbacks(); err != nil {
			log.Fatal("❌ Gagal menjalankan rollback untuk semua migrasi:", err)
		}

		fmt.Println("✅ Semua rollback berhasil!")
		return nil
	},
}

var RollbackBatchCommand = &cli.Command{
	Name:  "rollback:batch",
	Usage: "Rollback migrations for a specific batch (default last)",
	Flags: []cli.Flag{
		&cli.IntFlag{Name: "batch", Usage: "Batch number to rollback"},
	},
	Action: func(c *cli.Context) error {
		b := c.Int("batch")
		if b == 0 {
			fmt.Println("🔄 Rolling back last batch...")
			return database.RollbackLastBatch()
		}
		fmt.Printf("🔄 Rolling back batch %d...\n", b)
		return database.RollbackBatch(b)
	},
}

var CreateMigrationFile = func(name string) error {
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

	upTemplate, downTemplate := getMigrationTemplate(name)

	if err := writeTemplate(upFile, upTemplate); err != nil {
		return err
	}
	if err := writeTemplate(downFile, downTemplate); err != nil {
		return err
	}

	fmt.Println("Migration files created:")
	fmt.Println(" -", upFile)
	fmt.Println(" -", downFile)

	return nil
}

var getMigrationTemplate = func(name string) (string, string) {
	switch {
	case strings.HasPrefix(name, "create_"):
		table := extractTableName(name, "create_")
		up := fmt.Sprintf(`-- +++ UP Migration
-- Create table for %s:

CREATE TABLE %s (
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
	reference VARCHAR(255) UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL
);
`, table, table)

		down := fmt.Sprintf(`-- --- DOWN Migration
-- Drop the %s table

DROP TABLE IF EXISTS %s;
`, table, table)

		return up, down

	case strings.HasPrefix(name, "alter_"):
		table := extractTableName(name, "alter_")
		up := fmt.Sprintf(`-- +++ UP Migration
-- Alter table %s to add a new column

ALTER TABLE %s ADD COLUMN new_column VARCHAR(255);
`, table, table)

		down := fmt.Sprintf(`-- --- DOWN Migration
-- Revert ALTER TABLE by dropping the new_column

ALTER TABLE %s DROP COLUMN new_column;
`, table)

		return up, down

	default:
		return "-- +++ UP Migration\n", "-- --- DOWN Migration\n"
	}
}

var extractTableName = func(name string, prefix string) string {
	trimmed := strings.TrimPrefix(name, prefix)
	trimmed = strings.TrimSuffix(trimmed, "_table")
	return trimmed
}

var writeTemplate = func(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("gagal membuat file %s: %v", filePath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("gagal menulis ke file %s: %v", filePath, err)
	}

	return nil
}

var MigrateFreshCommand = &cli.Command{
	Name:  "migrate:fresh",
	Usage: "Rollback all migrations and re-run them",
	Action: func(c *cli.Context) error {
		fmt.Println("🔄 Rolling back all migrations...")
		if err := database.RunAllRollbacks(); err != nil {
			log.Fatal("❌ Failed to rollback all migrations:", err)
		}

		fmt.Println("🚀 Re-running all migrations...")
		if err := database.RunAllMigrations(); err != nil {
			log.Fatal("❌ Failed to re-run all migrations:", err)
		}

		fmt.Println("✅ Fresh migration completed successfully!")
		return nil
	},
}
