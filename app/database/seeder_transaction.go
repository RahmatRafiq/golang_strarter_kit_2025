package database

import (
	"fmt"
	"time"

	"golang_starter_kit_2025/facades"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// runSeederWithTransaction runs a seeder with optional transaction wrapping
func runSeederWithTransaction(seeder Seeder, db *gorm.DB, connectionName string) error {
	// Check if seeder wants transaction (default: true if not set)
	useTransaction := true
	if !seeder.WithTransaction {
		useTransaction = false
	}

	if !useTransaction {
		log.Debug().Str("seeder", seeder.Name).Msg("Running without transaction")
		// Run without transaction
		if err := seeder.Run(db); err != nil {
			return err
		}

		// Record in seeds table
		return db.Exec(
			"INSERT INTO seeds (connection_name, filename, batch) VALUES (?, ?, ?)",
			connectionName, seeder.Name, seeder.Batch,
		).Error
	}

	// Run with transaction for atomicity
	log.Debug().Str("seeder", seeder.Name).Msg("Running with transaction (atomic)")
	return db.Transaction(func(tx *gorm.DB) error {
		// Run the seeder
		if err := seeder.Run(tx); err != nil {
			log.Error().
				Str("seeder", seeder.Name).
				Err(err).
				Msg("Seeder failed, rolling back transaction")
			return fmt.Errorf("seeder execution failed: %w", err)
		}

		// Record in seeds table
		if err := tx.Exec(
			"INSERT INTO seeds (connection_name, filename, batch) VALUES (?, ?, ?)",
			connectionName, seeder.Name, seeder.Batch,
		).Error; err != nil {
			return fmt.Errorf("failed to record seeder: %w", err)
		}

		log.Debug().Str("seeder", seeder.Name).Msg("Seeder completed successfully")
		return nil
	})
}

// RunSeederWithDependencies runs a specific seeder and all its dependencies
func RunSeederWithDependencies(seederName, connectionName string) error {
	if connectionName == "" {
		connectionName = "mysql"
	}

	if err := ensureSeedsTable(connectionName); err != nil {
		return err
	}

	// Find the seeder
	var targetSeeder *Seeder
	for i, s := range SeederList {
		if s.Name == seederName {
			targetSeeder = &SeederList[i]
			break
		}
	}

	if targetSeeder == nil {
		return fmt.Errorf("seeder '%s' not found", seederName)
	}

	// Build dependency graph from all seeders
	graph := NewDependencyGraph(SeederList)

	// Get dependency chain for this specific seeder
	chain, err := graph.GetDependencyChain(seederName)
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %v", err)
	}

	conn, err := facades.GetConnection(connectionName)
	if err != nil {
		return fmt.Errorf("failed to get connection '%s': %v", connectionName, err)
	}

	// Filter out already applied seeders
	var toRun []Seeder
	for _, name := range chain {
		applied, err := isSeedApplied(name, connectionName)
		if err != nil {
			return err
		}

		if !applied {
			// Find seeder by name
			for _, s := range SeederList {
				if s.Name == name {
					toRun = append(toRun, s)
					break
				}
			}
		}
	}

	if len(toRun) == 0 {
		log.Info().
			Str("seeder", seederName).
			Msg("Seeder and all dependencies are already applied")
		return nil
	}

	// Run seeders in dependency order
	newBatch := time.Now().Unix()
	log.Info().
		Int("count", len(toRun)).
		Str("target_seeder", seederName).
		Msg("Running seeder(s) including dependencies")

	for _, s := range toRun {
		s.Batch = newBatch

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

		if err := runSeederWithTransaction(s, conn.DB, connectionName); err != nil {
			return fmt.Errorf("failed to run seeder %s: %w", s.Name, err)
		}
	}

	log.Info().
		Int("count", len(toRun)).
		Msg("Successfully ran seeder(s)")
	return nil
}
