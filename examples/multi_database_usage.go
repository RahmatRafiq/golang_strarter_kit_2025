//go:build examples

package main

import (
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/facades"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Warn().
			Err(err).
			Msg(".env file not found")
	}

	// Initialize database connections
	facades.ConnectDB()
	defer facades.CloseDB()

	// Example 1: Basic database service usage
	log.Info().
		Msg("Database Service Examples")
	basicDatabaseExample()

	// Example 4: Connection status and health
	log.Info().
		Msg("Connection Health Examples")
	healthCheckExample()
}

func basicDatabaseExample() {
	dbService := services.NewDatabaseService()

	// Test MySQL connection
	log.Info().
		Msg("Testing MySQL connection")
	err := dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
		var result struct {
			Version string
		}
		return db.Raw("SELECT VERSION() as version").Scan(&result).Error
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("database", "mysql").
			Msg("MySQL connection failed")
	} else {
		log.Info().
			Str("database", "mysql").
			Msg("MySQL connection successful")
	}

	// Test PostgreSQL connection
	log.Info().
		Msg("Testing PostgreSQL connection")
	err = dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
		var result struct {
			Version string
		}
		return db.Raw("SELECT version() as version").Scan(&result).Error
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("database", "postgresql").
			Msg("PostgreSQL connection failed")
	} else {
		log.Info().
			Str("database", "postgresql").
			Msg("PostgreSQL connection successful")
	}

	// Test MySQL Secondary connection
	log.Info().
		Msg("Testing MySQL Secondary connection")
	err = dbService.ExecuteOnMySQLSecondary(func(db *gorm.DB) error {
		var result struct {
			Version string
		}
		return db.Raw("SELECT VERSION() as version").Scan(&result).Error
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("database", "mysql_secondary").
			Msg("MySQL Secondary connection failed")
	} else {
		log.Info().
			Str("database", "mysql_secondary").
			Msg("MySQL Secondary connection successful")
	}
}

func healthCheckExample() {
	dbService := services.NewDatabaseService()

	log.Info().
		Msg("Getting connection statistics")
	stats, err := dbService.GetConnectionStats()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to get connection statistics")
		return
	}

	for connName, connStats := range stats {
		if statsMap, ok := connStats.(map[string]interface{}); ok {
			if connected, exists := statsMap["connected"]; exists && connected == true {
				logEvent := log.Info().
					Str("connection", connName).
					Str("status", "connected")

				if openConns, exists := statsMap["open_connections"]; exists {
					logEvent = logEvent.Interface("open_connections", openConns)
				}
				if inUse, exists := statsMap["in_use"]; exists {
					logEvent = logEvent.Interface("in_use", inUse)
				}
				if idle, exists := statsMap["idle"]; exists {
					logEvent = logEvent.Interface("idle", idle)
				}
				logEvent.Msg("Connection status")
			} else {
				logEvent := log.Warn().
					Str("connection", connName).
					Str("status", "disconnected")

				if errMsg, exists := statsMap["error"]; exists {
					logEvent = logEvent.Interface("error", errMsg)
				}
				logEvent.Msg("Connection status")
			}
		}
	}

	// Test individual connections
	log.Info().
		Msg("Testing individual connections")
	connections := []string{"mysql", "postgres", "mysql_secondary"}
	manager := facades.GetManager()

	for _, connName := range connections {
		if manager.IsConnected(connName) {
			log.Info().
				Str("connection", connName).
				Str("health", "healthy").
				Msg("Connection health check")
		} else {
			log.Warn().
				Str("connection", connName).
				Str("health", "unhealthy").
				Msg("Connection health check")
		}
	}
}
