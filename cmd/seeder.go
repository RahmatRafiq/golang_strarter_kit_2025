package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
	"unicode"

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

		// Convert snake_case to PascalCase
		words := strings.Split(name, "_")
		var parts []string
		for _, word := range words {
			if len(word) > 0 {
				runes := []rune(word)
				runes[0] = unicode.ToUpper(runes[0])
				parts = append(parts, string(runes))
			}
		}
		structName := strings.Join(parts, "")
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

		if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
			log.Fatal("❌ Gagal membuat file seeder:", err)
		}

		fmt.Println("✅ File seeder berhasil dibuat:", filePath)
		return nil
	},
}

var DBSeedCommand = &cli.Command{
	Name:  "db:seed",
	Usage: "Run all Go-based seeders",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
		&cli.StringFlag{Name: "class", Usage: "Specific seeder class to run"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		class := c.String("class")

		if class != "" {
			fmt.Printf("🌱 Running specific seeder: %s on connection %s\n", class, connection)
			return database.RunSpecificSeederOnConnection(class, connection)
		}

		fmt.Printf("🌱 Running all seeders on connection %s\n", connection)
		if err := database.RunAllSeedersOnConnection(connection); err != nil {
			log.Fatal("❌ Failed to run seeders:", err)
		}
		fmt.Println("✅ All seeders completed successfully!")
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
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		b := c.Int64("batch")
		connection := c.String("connection")

		if b == 0 {
			log.Printf("🔄 Rolling back last seed batch on connection %s\n", connection)
			return database.RollbackLastSeedBatchOnConnection(connection)
		}
		log.Printf("🔄 Rolling back seed batch %d on connection %s\n", b, connection)
		return database.RollbackSeedBatchOnConnection(b, connection)
	},
}
