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

func ensureSeedsTable() error {
	return facades.DB.Exec(`
		CREATE TABLE IF NOT EXISTS seeds (
			id INT PRIMARY KEY AUTO_INCREMENT,
			filename VARCHAR(255) NOT NULL,
			batch BIGINT NOT NULL,
			seeded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

func getLastSeedBatch() (int64, error) {
	var res struct{ Batch int64 }
	if err := facades.DB.
		Raw("SELECT COALESCE(MAX(batch),0) AS batch FROM seeds").
		Scan(&res).Error; err != nil {
		return 0, err
	}
	return res.Batch, nil
}

func isSeedApplied(name string) (bool, error) {
	var cnt int64
	if err := facades.DB.
		Raw("SELECT COUNT(*) FROM seeds WHERE filename = ?", name).
		Scan(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func RunAllSeeders() error {
	if err := ensureSeedsTable(); err != nil {
		return err
	}
	_, err := getLastSeedBatch()
	if err != nil {
		return err
	}
	newBatch := time.Now().Unix()

	var pending []Seeder
	for _, s := range SeederList {
		applied, err := isSeedApplied(s.Name)
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
		if err := s.Run(facades.DB); err != nil {
			return fmt.Errorf("failed to run seeder %s: %w", s.Name, err)
		}
		if err := facades.DB.
			Exec("INSERT INTO seeds (filename, batch) VALUES (?, ?)", s.Name, s.Batch).
			Error; err != nil {
			return err
		}
	}
	log.Printf("✅ Seed batch %d applied.\n", newBatch)
	return nil
}
func RollbackSeedBatch(batch int64) error {
	if err := ensureSeedsTable(); err != nil {
		return err
	}

	var rows []struct{ Filename string }
	if err := facades.DB.
		Raw("SELECT filename FROM seeds WHERE batch = ? ORDER BY id DESC", batch).
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
					if err := s.Rollback(facades.DB); err != nil {
						return fmt.Errorf("rollback seeder %s failed: %w", s.Name, err)
					}
				}
				break
			}
		}
		if err := facades.DB.
			Exec("DELETE FROM seeds WHERE filename = ? AND batch = ?", r.Filename, batch).
			Error; err != nil {
			return err
		}
	}
	log.Printf("✅ Seeder batch %d rolled back.\n", batch)
	return nil
}

func RollbackLastSeedBatch() error {
	b, err := getLastSeedBatch()
	if err != nil {
		return err
	}
	if b == 0 {
		log.Println("⚠️ No seed batch to rollback.")
		return nil
	}
	return RollbackSeedBatch(b)
}
