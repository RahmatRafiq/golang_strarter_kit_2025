package cmd

import (
	"fmt"
	"log"

	"golang_strarter_kit_2025/app/database"

	"github.com/urfave/cli/v2"
)

// MigrationCommand adalah command CLI untuk membuat migrasi dan menjalankan migrasi
// app/cmd/migration.go

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

		err := database.CreateMigrationFile(name)
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
