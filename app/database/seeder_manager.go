package database

import (
	"fmt"
	"time"

	"golang_starter_kit_2025/app/database/seeds"
	"golang_starter_kit_2025/facades"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Seeder struct {
	Name            string
	Run             func(db *gorm.DB) error
	Rollback        func(db *gorm.DB) error
	DependsOn       []string // Dependencies that must run before this seeder
	WithTransaction bool     // Whether to wrap in transaction (default: true)
	Batch           int64
}

var SeederList = []Seeder{
	{
		Name:            "UserSeeder",
		Run:             seeds.SeedUserSeeder,
		Rollback:        seeds.RollbackUserSeeder,
		DependsOn:       []string{}, // No dependencies
		WithTransaction: true,       // Use transaction by default
	},
}

func ensureSeedsTable(connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql"
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	// Use different table creation syntax based on database type
	var createTableSQL string
	if conn.IsPostgreSQL() {
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS seeds (
				id SERIAL PRIMARY KEY,
				connection_name VARCHAR(50) NOT NULL,
				filename VARCHAR(255) NOT NULL,
				batch BIGINT NOT NULL,
				seeded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE (connection_name, filename)
			)`
	} else {
		// MySQL/MariaDB
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS seeds (
				id INT PRIMARY KEY AUTO_INCREMENT,
				connection_name VARCHAR(50) NOT NULL,
				filename VARCHAR(255) NOT NULL,
				batch BIGINT NOT NULL,
				seeded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE KEY unique_seed (connection_name, filename)
			)`
	}

	return conn.DB.Exec(createTableSQL).Error
}

func getLastSeedBatch(connectionName string) (int64, error) {
	if connectionName == "" {
		connectionName = "mysql"
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return 0, fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var res struct{ Batch int64 }
	if err := conn.DB.
		Raw("SELECT COALESCE(MAX(batch),0) AS batch FROM seeds WHERE connection_name = ?", connectionName).
		Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Batch, nil
}

func isSeedApplied(name, connectionName string) (bool, error) {
	if connectionName == "" {
		connectionName = "mysql"
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return false, fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var cnt int64
	if err := conn.DB.
		Raw("SELECT COUNT(*) FROM seeds WHERE connection_name = ? AND filename = ?", connectionName, name).
		Scan(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// RunAllSeeders runs all seeders on the default connection
func RunAllSeeders() error {
	return RunAllSeedersOnConnection("")
}

// RunAllSeedersOnConnection runs all seeders on a specified connection with dependency resolution
func RunAllSeedersOnConnection(connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql"
	}

	if err := ensureSeedsTable(connectionName); err != nil {
		return err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	// Get pending seeders
	var pending []Seeder
	for _, s := range SeederList {
		applied, checkErr := isSeedApplied(s.Name, connectionName)
		if checkErr != nil {
			return checkErr
		}
		if !applied {
			pending = append(pending, s)
		}
	}

	if len(pending) == 0 {
		log.Info().Msg("No pending seeders to run")
		return nil
	}

	// Resolve dependencies using topological sort
	graph := NewDependencyGraph(pending)

	// Validate dependencies first
	if validateErr := graph.ValidateDependencies(); validateErr != nil {
		return fmt.Errorf("dependency validation failed: %v", validateErr)
	}

	// Get correct order
	orderedNames, err := graph.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %v", err)
	}

	// Create ordered seeder list
	orderedSeeders := make([]Seeder, 0, len(orderedNames))
	seederMap := make(map[string]Seeder)
	for _, s := range pending {
		seederMap[s.Name] = s
	}
	for _, name := range orderedNames {
		orderedSeeders = append(orderedSeeders, seederMap[name])
	}

	// Run seeders in order with new batch
	newBatch := time.Now().Unix()
	log.Info().
		Int("count", len(orderedSeeders)).
		Str("connection", connectionName).
		Msg("Running seeders in dependency order")

	for _, s := range orderedSeeders {
		s.Batch = newBatch

		// Show dependency info
		if len(s.DependsOn) > 0 {
			log.Info().
				Str("seeder", s.Name).
				Strs("dependencies", s.DependsOn).
				Msg("Seeding")
		} else {
			log.Info().
				Str("seeder", s.Name).
				Msg("Seeding")
		}

		// Run with or without transaction
		if err := runSeederWithTransaction(s, conn.DB, connectionName); err != nil {
			return fmt.Errorf("failed to run seeder %s: %w", s.Name, err)
		}
	}

	log.Info().
		Int64("batch", newBatch).
		Str("connection", connectionName).
		Msg("Seed batch applied successfully")
	return nil
}

// RollbackSeedBatch rolls back a specific seed batch on default connection
func RollbackSeedBatch(batch int64) error {
	return RollbackSeedBatchOnConnection(batch, "")
}

// RollbackSeedBatchOnConnection rolls back a specific seed batch on a specified connection
func RollbackSeedBatchOnConnection(batch int64, connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql"
	}

	if err := ensureSeedsTable(connectionName); err != nil {
		return err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	var rows []struct{ Filename string }
	if err := conn.DB.
		Raw("SELECT filename FROM seeds WHERE connection_name = ? AND batch = ? ORDER BY id DESC", connectionName, batch).
		Scan(&rows).Error; err != nil {
		return err
	}
	if len(rows) == 0 {
		log.Warn().Int64("batch", batch).Msg("No seeders in batch")
		return nil
	}

	for _, r := range rows {
		log.Info().Str("seeder", r.Filename).Msg("Rolling back seeder")
		for _, s := range SeederList {
			if s.Name == r.Filename {
				if s.Rollback != nil {
					if err := s.Rollback(conn.DB); err != nil {
						return fmt.Errorf("rollback seeder %s failed: %w", s.Name, err)
					}
				}
				break
			}
		}
		if err := conn.DB.
			Exec("DELETE FROM seeds WHERE connection_name = ? AND filename = ? AND batch = ?", connectionName, r.Filename, batch).
			Error; err != nil {
			return err
		}
	}
	log.Info().
		Int64("batch", batch).
		Str("connection", connectionName).
		Msg("Seeder batch rolled back")
	return nil
}

// RollbackLastSeedBatch rolls back the last seed batch on default connection
func RollbackLastSeedBatch() error {
	return RollbackLastSeedBatchOnConnection("")
}

// RollbackLastSeedBatchOnConnection rolls back the last seed batch on a specified connection
func RollbackLastSeedBatchOnConnection(connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql"
	}

	b, err := getLastSeedBatch(connectionName)
	if err != nil {
		return err
	}
	if b == 0 {
		log.Warn().Msg("No seed batch to rollback")
		return nil
	}
	return RollbackSeedBatchOnConnection(b, connectionName)
}

// RunSpecificSeeder runs a specific seeder on default connection
func RunSpecificSeeder(name string) error {
	return RunSpecificSeederOnConnection(name, "")
}

// RunSpecificSeederOnConnection runs a specific seeder on a specified connection
func RunSpecificSeederOnConnection(name, connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql"
	}

	if err := ensureSeedsTable(connectionName); err != nil {
		return err
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	// Find the seeder
	var seeder *Seeder
	for _, s := range SeederList {
		if s.Name == name {
			seeder = &s
			break
		}
	}

	if seeder == nil {
		return fmt.Errorf("seeder '%s' not found", name)
	}

	// Check if already applied
	applied, err := isSeedApplied(name, connectionName)
	if err != nil {
		return err
	}

	if applied {
		log.Warn().Str("seeder", name).Msg("Seeder has already been run")
		return nil
	}

	// Run the seeder
	newBatch := time.Now().Unix()
	log.Info().Str("seeder", name).Msg("Running seeder")

	if err := seeder.Run(conn.DB); err != nil {
		return fmt.Errorf("failed to run seeder %s: %w", name, err)
	}

	if err := conn.DB.
		Exec("INSERT INTO seeds (connection_name, filename, batch) VALUES (?, ?, ?)", connectionName, name, newBatch).
		Error; err != nil {
		return err
	}

	log.Info().
		Str("seeder", name).
		Str("connection", connectionName).
		Msg("Seeder completed successfully")
	return nil
}
