package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"golang_starter_kit_2025/app/database"

	"github.com/urfave/cli/v2"
)

var MakeSeederCommand = &cli.Command{
	Name:  "make:seeder",
	Usage: "Generate a new Go seeder skeleton file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Nama seeder (tanpa ekstensi). Contoh: --name=users_seeder",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		name := c.String("name")
		if name == "" {
			return fmt.Errorf("nama seeder harus disediakan. Contoh: make:seeder --name=users_seeder")
		}

		timestamp := time.Now().Format("20060102150405")
		fileName := fmt.Sprintf("%s_%s.go", timestamp, name)
		dir := path.Join("app", "database", "seeds")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("gagal membuat direktori seeds: %w", err)
		}
		filePath := path.Join(dir, fileName)

		structName := strings.Title(strings.ReplaceAll(name, "_", " "))
		structName = strings.ReplaceAll(structName, " ", "")
		modelName := strings.TrimSuffix(structName, "Seeder")
		content := fmt.Sprintf(`package seeds

import (
	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"log"
	"time"

	"gorm.io/gorm"
)

func Seed%[1]s(db *gorm.DB) error {
	log.Println("🌱 Seeding %[1]s...")

	data := models.%[2]s{
		Reference: helpers.GenerateReference("USR"),
		// Tambahkan field sesuai model
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func Rollback%[1]s(db *gorm.DB) error {
	log.Println("🗑️ Rolling back %[1]s…")
	return db.Unscoped().
		Where("reference LIKE ?", "USR%%").
		// Delete(&models.%[2]s{}).
		Error
}
`, structName, modelName)

		if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
			log.Fatal("❌ Gagal membuat file seeder:", err)
		}

		fmt.Println("✅ File seeder berhasil dibuat:", filePath)
		return nil
	},
}

var DBSeedCommand = &cli.Command{
	Name:  "db:seed",
	Usage: "Run all Go-based seeders",
	Action: func(c *cli.Context) error {
		fmt.Println("🌱 Menjalankan semua seeder Go...")
		if err := database.RunAllSeeders(); err != nil {
			log.Fatal("❌ Gagal menjalankan seeder:", err)
		}
		fmt.Println("✅ Semua seeder berhasil dijalankan!")
		return nil
	},
}
var RollbackSeederCommand = &cli.Command{
	Name:  "rollback:seeder",
	Usage: "Rollback seeders for a specific batch (default last)",
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:    "batch",
			Aliases: []string{"b"},
			Usage:   "Batch number to rollback",
		},
	},
	Action: func(c *cli.Context) error {
		b := c.Int64("batch")
		if b == 0 {
			log.Println("🔄 Rolling back last seed batch...")
			return database.RollbackLastSeedBatch()
		}
		log.Printf("🔄 Rolling back seed batch %d...\n", b)
		return database.RollbackSeedBatch(b)
	},
}
