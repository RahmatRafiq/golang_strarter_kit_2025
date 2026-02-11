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

var MigrateResetCommand = &cli.Command{
	Name:  "migrate:reset",
	Usage: "Rollback all migrations (alias for rollback:all)",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		fmt.Printf("🔄 Resetting all migrations on connection %s\n", connection)
		return database.RunAllRollbacksOnConnection(connection)
	},
}

var RollbackBatchCommand = &cli.Command{
	Name:  "rollback:batch",
	Usage: "Rollback specific batch or last N steps",
	Flags: []cli.Flag{
		&cli.IntFlag{Name: "batch", Usage: "Specific batch number to rollback"},
		&cli.IntFlag{Name: "step", Usage: "Number of batches to rollback (from latest)"},
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		b := c.Int("batch")
		step := c.Int("step")
		connection := c.String("connection")

		// If step is provided, rollback N batches
		if step > 0 {
			fmt.Printf("🔄 Rolling back last %d batch(es) on connection %s\n", step, connection)
			return database.RollbackStepsOnConnection(step, connection)
		}

		// If batch is provided, rollback specific batch
		if b > 0 {
			return database.RollbackBatchOnConnection(b, connection)
		}

		// Default: rollback last batch
		return database.RollbackLastBatchOnConnection(connection)
	},
}

var MigrateFreshCommand = &cli.Command{
	Name:  "migrate:fresh",
	Usage: "Reset and re-run all migrations",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
		&cli.BoolFlag{Name: "seed", Usage: "Run seeders after migrations"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		seed := c.Bool("seed")

		fmt.Printf("🔄 Fresh: rollback all then migrate all on connection %s\n", connection)
		if err := database.RunAllRollbacksOnConnection(connection); err != nil {
			return err
		}

		if err := database.RunAllMigrationsOnConnection(connection); err != nil {
			return err
		}

		if seed {
			fmt.Printf("🌱 Running seeders on connection %s...\n", connection)
			if err := database.RunAllSeedersOnConnection(connection); err != nil {
				return fmt.Errorf("failed to run seeders: %v", err)
			}
			fmt.Println("✅ Seeders completed successfully!")
		}

		return nil
	},
}

var DBWipeCommand = &cli.Command{
	Name:  "db:wipe",
	Usage: "Drop all tables from the database",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
		&cli.BoolFlag{Name: "force", Usage: "Force execution without confirmation"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		force := c.Bool("force")

		if !force {
			fmt.Println("⚠️  WARNING: This will DROP ALL TABLES in the database!")
			fmt.Printf("Connection: %s\n\n", connection)
			fmt.Print("Are you sure you want to continue? (type 'yes' to confirm): ")

			var confirmation string
			if _, err := fmt.Scanln(&confirmation); err != nil {
				fmt.Printf("❌ Error reading input: %v\n", err)
				return err
			}

			if confirmation != "yes" {
				fmt.Println("❌ Operation canceled.")
				return nil
			}
		}

		fmt.Printf("🗑️  Dropping all tables on connection %s\n", connection)
		return database.WipeDatabase(connection)
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

var MigrateLockStatusCommand = &cli.Command{
	Name:  "migrate:lock:status",
	Usage: "Check migration lock status",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to check"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		locked, info, err := database.CheckLockStatus(connection)
		if err != nil {
			return fmt.Errorf("failed to check lock status: %v", err)
		}

		if locked {
			fmt.Printf("🔒 Migration LOCKED for connection '%s'\n", connection)
			fmt.Printf("   %s\n", info)
		} else {
			fmt.Printf("🔓 Migration NOT locked for connection '%s'\n", connection)
		}
		return nil
	},
}

var MigrateLockReleaseCommand = &cli.Command{
	Name:  "migrate:lock:release",
	Usage: "Force release migration lock (use with caution)",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection"},
		&cli.BoolFlag{Name: "force", Usage: "Force release without confirmation"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		force := c.Bool("force")

		if !force {
			fmt.Println("⚠️  WARNING: Force releasing a lock can cause issues if migrations are actually running!")
			fmt.Printf("Connection: %s\n\n", connection)
			fmt.Print("Are you sure? (type 'yes' to confirm): ")

			var confirmation string
			if _, err := fmt.Scanln(&confirmation); err != nil {
				fmt.Printf("❌ Error reading input: %v\n", err)
				return err
			}

			if confirmation != "yes" {
				fmt.Println("❌ Operation canceled.")
				return nil
			}
		}

		return database.ForceReleaseLock(connection)
	},
}

var MigrateLogsCommand = &cli.Command{
	Name:  "migrate:logs",
	Usage: "Show migration execution logs",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection"},
		&cli.IntFlag{Name: "limit", Value: 20, Usage: "Number of logs to show"},
		&cli.BoolFlag{Name: "failed", Usage: "Show only failed migrations"},
		&cli.BoolFlag{Name: "slow", Usage: "Show slow migrations (>1s)"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		limit := c.Int("limit")
		showFailed := c.Bool("failed")
		showSlow := c.Bool("slow")

		var logs []database.MigrationLog
		var err error

		if showFailed {
			logs, err = database.GetFailedMigrations(connection, limit)
			fmt.Printf("📊 Failed Migrations (connection: %s)\n\n", connection)
		} else if showSlow {
			logs, err = database.GetSlowMigrations(connection, 1000, limit) // >1s
			fmt.Printf("🐌 Slow Migrations (>1s) (connection: %s)\n\n", connection)
		} else {
			logs, err = database.GetMigrationLogs(connection, limit)
			fmt.Printf("📊 Recent Migration Logs (connection: %s)\n\n", connection)
		}

		if err != nil {
			return fmt.Errorf("failed to get logs: %v", err)
		}

		if len(logs) == 0 {
			fmt.Println("No logs found.")
			return nil
		}

		fmt.Println(strings.Repeat("=", 100))
		fmt.Printf("%-40s %-8s %-12s %-12s %-20s\n", "Migration", "Batch", "Status", "Time", "Executed At")
		fmt.Println(strings.Repeat("-", 100))

		for _, log := range logs {
			status := log.Status
			if status == "success" {
				status = "✅ Success"
			} else {
				status = "❌ Failed"
			}

			fmt.Printf("%-40s %-8d %-12s %-12s %-20s\n",
				truncate(log.Filename, 40),
				log.Batch,
				status,
				database.FormatDuration(log.ExecutionTimeMs),
				log.ExecutedAt.Format("2006-01-02 15:04:05"),
			)

			if log.ErrorMessage != "" && showFailed {
				fmt.Printf("   Error: %s\n", truncate(log.ErrorMessage, 90))
			}
		}

		fmt.Println(strings.Repeat("=", 100))
		return nil
	},
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

var MigrateStatusCommand = &cli.Command{
	Name:  "migrate:status",
	Usage: "Show the status of each migration",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "connection", Value: "mysql", Usage: "Database connection to use (mysql, postgres, mysql_secondary)"},
	},
	Action: func(c *cli.Context) error {
		connection := c.String("connection")
		return database.ShowMigrationStatus(connection)
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
			stats, err := manager.GetConnectionStats(connection)
			if err != nil {
				fmt.Printf("⚠️  Connection '%s' is connected but stats unavailable: %v\n", connection, err)
				return err
			}
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
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	path := fmt.Sprintf("%s/app/database/migrations/%s", dir, fname)
	up, down := getMigrationTemplate(name)
	content := fmt.Sprintf("%s\n%s\n%s\n%s", upMarker, up, downMarker, down)
	return os.WriteFile(path, []byte(content), 0600)
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
