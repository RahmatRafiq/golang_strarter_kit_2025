package database

import (
	"fmt"
	"log"
	"time"

	"golang_starter_kit_2025/facades"

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
		log.Printf("  ↳ Running without transaction")
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
	log.Printf("  ↳ Running with transaction (atomic)")
	return db.Transaction(func(tx *gorm.DB) error {
		// Run the seeder
		if err := seeder.Run(tx); err != nil {
			log.Printf("  ✗ Seeder failed, rolling back transaction")
			return fmt.Errorf("seeder execution failed: %w", err)
		}

		// Record in seeds table
		if err := tx.Exec(
			"INSERT INTO seeds (connection_name, filename, batch) VALUES (?, ?, ?)",
			connectionName, seeder.Name, seeder.Batch,
		).Error; err != nil {
			return fmt.Errorf("failed to record seeder: %w", err)
		}

		log.Printf("  ✓ Seeder completed successfully")
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
		log.Printf("✅ Seeder '%s' and all dependencies are already applied", seederName)
		return nil
	}

	// Run seeders in dependency order
	newBatch := time.Now().Unix()
	log.Printf("🚀 Running %d seeder(s) for '%s' (including dependencies)\n", len(toRun), seederName)

	for _, s := range toRun {
		s.Batch = newBatch

		if len(s.DependsOn) > 0 {
			log.Printf("🌱 Seeding: %s (depends on: %v)", s.Name, s.DependsOn)
		} else {
			log.Printf("🌱 Seeding: %s", s.Name)
		}

		if err := runSeederWithTransaction(s, conn.DB, connectionName); err != nil {
			return fmt.Errorf("failed to run seeder %s: %w", s.Name, err)
		}
	}

	log.Printf("✅ Successfully ran %d seeder(s)\n", len(toRun))
	return nil
}
