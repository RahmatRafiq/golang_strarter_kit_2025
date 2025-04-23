package cmd

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"

	"golang_strarter_kit_2025/app/database"
)

var MakeSeederCommand = &cli.Command{
	Name:  "make:seeder",
	Usage: "Generate a new seeder file",
	Action: func(c *cli.Context) error {
		name := c.Args().First()
		if name == "" {
			return fmt.Errorf("nama seeder harus disediakan. Contoh: make:seeder users_seeder")
		}

		if err := database.CreateSeederFile(name); err != nil {
			log.Fatal("❌ Gagal membuat file seeder:", err)
		}

		fmt.Println("✅ File seeder berhasil dibuat.")
		return nil
	},
}

var DBSeedCommand = &cli.Command{
	Name:  "db:seed",
	Usage: "Run all seeders",
	Action: func(c *cli.Context) error {
		fmt.Println("🌱 Menjalankan semua seeder...")
		if err := database.RunAllSeeders(); err != nil {
			log.Fatal("❌ Gagal menjalankan seeder:", err)
		}
		fmt.Println("✅ Semua seeder berhasil dijalankan!")
		return nil
	},
}
