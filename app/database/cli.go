package database

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/urfave/cli/v2"
)

var MigrationCmd = &cli.Command{
	Name:  "make:migration",
	Usage: "Generate a new SQL migration file",
	Action: func(c *cli.Context) error {
		name := c.Args().First()
		if name == "" {
			return fmt.Errorf("nama migration harus disediakan")
		}

		timestamp := time.Now().Format("20060102150405")
		filename := fmt.Sprintf("%s_%s", timestamp, name)

		rootPath, _ := os.Getwd() // <-- ambil direktori project
		migrationPath := filepath.Join(rootPath, "app", "database", "migrations")

		// Pastikan folder migrations ada
		if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
			if err := os.MkdirAll(migrationPath, 0755); err != nil {
				return fmt.Errorf("gagal membuat folder migrations: %v", err)
			}
		}

		// Template isi
		upFile := fmt.Sprintf("%s/%s.up.sql", migrationPath, filename)
		downFile := fmt.Sprintf("%s/%s.down.sql", migrationPath, filename)

		writeTemplate(upFile, "-- +++ Write your UP migration here\n")
		writeTemplate(downFile, "-- --- Write your DOWN migration here\n")

		fmt.Println("Migration created:")
		fmt.Println(" -", upFile)
		fmt.Println(" -", downFile)
		return nil
	},
}

func writeTemplate(path, content string) {
	tmpl := template.Must(template.New("migration").Parse(content))
	f, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("gagal membuat file %s: %v", path, err))
	}
	defer f.Close()
	tmpl.Execute(f, nil)
}
