package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang_strarter_kit_2025/app/database"

	"github.com/urfave/cli/v2"
)

var MigrationCommand = &cli.Command{
	Name:  "migrate",
	Usage: "Run a specific migration file",
	Flags: []cli.Flag{&cli.StringFlag{Name: "file", Required: true}},
	Action: func(c *cli.Context) error {
		name := c.String("file")
		fmt.Println("🚀 Migrate:", name)
		return database.RunMigration(name)
	},
}

var RollbackCommand = &cli.Command{
	Name:  "rollback",
	Usage: "Rollback a specific migration",
	Flags: []cli.Flag{&cli.StringFlag{Name: "file", Required: true}},
	Action: func(c *cli.Context) error {
		name := c.String("file")
		fmt.Println("🔄 Rollback:", name)
		return database.RollbackMigration(name)
	},
}

var MakeMigrationCommand = &cli.Command{
	Name:  "make:migration",
	Usage: "Create new migration template",
	Action: func(c *cli.Context) error {
		if c.Args().Len() < 1 {
			return fmt.Errorf("nama migration dibutuhkan")
		}
		return CreateMigration(c.Args().First())
	},
}

var MigrateAllCommand = &cli.Command{
	Name:  "migrate:all",
	Usage: "Run all pending migrations",
	Action: func(c *cli.Context) error {
		fmt.Println("🚀 Migrate all")
		return database.RunAllMigrations()
	},
}

var RollbackAllCommand = &cli.Command{
	Name:  "rollback:all",
	Usage: "Rollback all batches",
	Action: func(c *cli.Context) error {
		fmt.Println("🔄 Rollback all")
		return database.RunAllRollbacks()
	},
}

var RollbackBatchCommand = &cli.Command{
	Name:  "rollback:batch",
	Usage: "Rollback specific batch",
	Flags: []cli.Flag{&cli.IntFlag{Name: "batch"}},
	Action: func(c *cli.Context) error {
		b := c.Int("batch")
		if b == 0 {
			return database.RollbackLastBatch()
		}
		return database.RollbackBatch(b)
	},
}

var MigrateFreshCommand = &cli.Command{
	Name:  "migrate:fresh",
	Usage: "Reset and re-run all migrations",
	Action: func(c *cli.Context) error {
		fmt.Println("🔄 Fresh: rollback all then migrate all")
		if err := database.RunAllRollbacks(); err != nil {
			return err
		}
		return database.RunAllMigrations()
	},
}

func CreateMigration(name string) error {
	ts := time.Now().Format("20060102150405")
	fname := fmt.Sprintf("%s_%s.sql", ts, name)
	dir, _ := os.Getwd()
	path := fmt.Sprintf("%s/app/database/migrations/%s", dir, fname)
	up, down := getMigrationTemplate(name)
	content := fmt.Sprintf("%s\n%s\n%s\n%s", upMarker, up, downMarker, down)
	return os.WriteFile(path, []byte(content), 0644)
}

var upMarker = "-- +++ UP Migration"
var downMarker = "-- --- DOWN Migration"

func getMigrationTemplate(name string) (string, string) {
	if strings.HasPrefix(name, "create_") {
		tbl := strings.TrimPrefix(name, "create_")
		tbl = strings.TrimSuffix(tbl, "_table")
		up := fmt.Sprintf(`CREATE TABLE %s (
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL DEFAULT NULL
);`, tbl) // Assuming the column name is derived from the table name
		down := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tbl)
		return up, down
	}

	return "-- up SQL here", "-- down SQL here"
}
