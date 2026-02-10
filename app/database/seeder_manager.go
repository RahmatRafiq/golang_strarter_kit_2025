package database

import (
	"fmt"
	"log"
	"sort"
	"time"

	"golang_starter_kit_2025/app/database/seeds"
	"golang_starter_kit_2025/facades"

	"gorm.io/gorm"
)

type Seeder struct {
	Name     string
	Run      func(db *gorm.DB) error
	Rollback func(db *gorm.DB) error
	Batch    int64
}

var SeederList = []Seeder{
	{Name: "UserSeeder",
		Run:      seeds.SeedUserSeeder,
		Rollback: seeds.RollbackUserSeeder,
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

// RunAllSeedersOnConnection runs all seeders on a specified connection
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

	_, err = getLastSeedBatch(connectionName)
	if err != nil {
		return err
	}
	newBatch := time.Now().Unix()

	var pending []Seeder
	for _, s := range SeederList {
		applied, err := isSeedApplied(s.Name, connectionName)
		if err != nil {
			return err
		}
		if !applied {
			s.Batch = newBatch
			pending = append(pending, s)
		}
	}
	sort.Slice(pending, func(i, j int) bool { return pending[i].Name < pending[j].Name })

	for _, s := range pending {
		log.Println("🌱 Seeding:", s.Name)
		if err := s.Run(conn.DB); err != nil {
			return fmt.Errorf("failed to run seeder %s: %w", s.Name, err)
		}
		if err := conn.DB.
			Exec("INSERT INTO seeds (connection_name, filename, batch) VALUES (?, ?, ?)", connectionName, s.Name, s.Batch).
			Error; err != nil {
			return err
		}
	}
	log.Printf("✅ Seed batch %d applied on connection %s.\n", newBatch, connectionName)
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
		log.Printf("⚠️ No seeders in batch %d\n", batch)
		return nil
	}

	for _, r := range rows {
		log.Println("🔄 Rolling back seeder:", r.Filename)
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
	log.Printf("✅ Seeder batch %d rolled back on connection %s.\n", batch, connectionName)
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
		log.Println("⚠️ No seed batch to rollback.")
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
		log.Printf("⚠️ Seeder '%s' has already been run\n", name)
		return nil
	}

	// Run the seeder
	newBatch := time.Now().Unix()
	log.Printf("🌱 Running seeder: %s\n", name)

	if err := seeder.Run(conn.DB); err != nil {
		return fmt.Errorf("failed to run seeder %s: %w", name, err)
	}

	if err := conn.DB.
		Exec("INSERT INTO seeds (connection_name, filename, batch) VALUES (?, ?, ?)", connectionName, name, newBatch).
		Error; err != nil {
		return err
	}

	log.Printf("✅ Seeder '%s' completed successfully on connection %s\n", name, connectionName)
	return nil
}
