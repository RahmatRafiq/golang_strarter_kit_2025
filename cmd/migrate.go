package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang_starter_kit_2025/app/database"
	"golang_starter_kit_2025/facades"

	"github.com/urfave/cli/v2"
)

var MigrationCommand = &cli.Command{
	Name:  "migrate",
	Usage: "Run a specific migration file",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "file", Required: true},
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		name := c.String("file")
		connection := c.String("connection")
		fmt.Printf("🚀 Migrate: %s on connection %s\n", name, connection)
		return database.RunMigrationOnConnection(name, connection)
	},
}

var RollbackCommand = &cli.Command{
	Name:  "rollback",
	Usage: "Rollback a specific migration",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "file", Required: true},
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		name := c.String("file")
		connection := c.String("connection")
		fmt.Printf("🔄 Rollback: %s on connection %s\n", name, connection)
		return database.RollbackMigrationOnConnection(name, connection)
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
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		fmt.Printf("🚀 Migrate all on connection %s\n", connection)
		return database.RunAllMigrationsOnConnection(connection)
	},
}

var RollbackAllCommand = &cli.Command{
	Name:  "rollback:all",
	Usage: "Rollback all batches",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		fmt.Printf("🔄 Rollback all on connection %s\n", connection)
		return database.RunAllRollbacksOnConnection(connection)
	},
}

var RollbackBatchCommand = &cli.Command{
	Name:  "rollback:batch",
	Usage: "Rollback specific batch",
	Flags: []cli.Flag{
		&cli.IntFlag{Name: "batch"},
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		b := c.Int("batch")
		connection := c.String("connection")
		if b == 0 {
			return database.RollbackLastBatchOnConnection(connection)
		}
		return database.RollbackBatchOnConnection(b, connection)
	},
}

var MigrateFreshCommand = &cli.Command{
	Name:  "migrate:fresh",
	Usage: "Reset and re-run all migrations",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		fmt.Printf("🔄 Fresh: rollback all then migrate all on connection %s\n", connection)
		if err := database.RunAllRollbacksOnConnection(connection); err != nil {
			return err
		}
		return database.RunAllMigrationsOnConnection(connection)
	},
}

var DBConnectionsCommand = &cli.Command{
	Name:  "db:connections",
	Usage: "List all available database connections",
	Action: func(c *cli.Context) error {
		fmt.Println("📊 Available Database Connections:")
		connections := []string{"mysql", "postgres", "mysql_secondary"}
		for _, conn := range connections {
			fmt.Printf("  - %s\n", conn)
		}
		return nil
	},
}

var DBStatusCommand = &cli.Command{
	Name:  "db:status",
	Usage: "Check database connection status",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to check (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		fmt.Printf("🔍 Checking connection status for: %s\n", connection)

		// Initialize database manager
		manager := facades.GetManager()

		// Try to connect
		conn, err := manager.GetConnection(connection)
		if err != nil {
			fmt.Printf("❌ Connection '%s' failed: %v\n", connection, err)
			return err
		}

		// Check if connected
		if manager.IsConnected(connection) {
			stats, _ := manager.GetConnectionStats(connection)
			fmt.Printf("✅ Connection '%s' is healthy\n", connection)
			fmt.Printf("   Database Type: %s\n", conn.GetType())
			fmt.Printf("   Open Connections: %d\n", stats.OpenConnections)
			fmt.Printf("   In Use: %d\n", stats.InUse)
			fmt.Printf("   Idle: %d\n", stats.Idle)
		} else {
			fmt.Printf("❌ Connection '%s' is not healthy\n", connection)
		}

		return nil
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
);`, tbl)
		down := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tbl)
		return up, down
	}

	if strings.HasPrefix(name, "alter_") {
		tbl := strings.TrimPrefix(name, "alter_")
		tbl = strings.TrimSuffix(tbl, "_table")
		up := fmt.Sprintf(`ALTER TABLE %s 
-- ADD COLUMN new_column_name DATA_TYPE;
`, tbl)
		down := fmt.Sprintf(`ALTER TABLE %s 
-- DROP COLUMN new_column_name;
`, tbl)
		return up, down
	}

	// Default fallback
	return "-- up SQL here", "-- down SQL here"
}
