package main

import (
	"fmt"
	"log"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/facades"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Initialize database connections
	facades.ConnectDB()
	defer facades.CloseDB()

	// Example 1: Basic database service usage
	fmt.Println("=== Database Service Examples ===")
	basicDatabaseExample()

	// Example 2: Repository pattern usage
	fmt.Println("\n=== Repository Pattern Examples ===")
	repositoryExample()

	// Example 3: Cross-database synchronization
	fmt.Println("\n=== Data Synchronization Examples ===")
	syncExample()

	// Example 4: Connection status and health
	fmt.Println("\n=== Connection Health Examples ===")
	healthCheckExample()
}

func basicDatabaseExample() {
	dbService := services.NewDatabaseService()

	// Test MySQL connection
	fmt.Println("1. Testing MySQL connection...")
	err := dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
		var result struct {
			Version string
		}
		return db.Raw("SELECT VERSION() as version").Scan(&result).Error
	})
	if err != nil {
		fmt.Printf("   ❌ MySQL Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ MySQL connection successful\n")
	}

	// Test PostgreSQL connection
	fmt.Println("2. Testing PostgreSQL connection...")
	err = dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
		var result struct {
			Version string
		}
		return db.Raw("SELECT version() as version").Scan(&result).Error
	})
	if err != nil {
		fmt.Printf("   ❌ PostgreSQL Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ PostgreSQL connection successful\n")
	}

	// Test MySQL Secondary connection
	fmt.Println("3. Testing MySQL Secondary connection...")
	err = dbService.ExecuteOnMySQLSecondary(func(db *gorm.DB) error {
		var result struct {
			Version string
		}
		return db.Raw("SELECT VERSION() as version").Scan(&result).Error
	})
	if err != nil {
		fmt.Printf("   ❌ MySQL Secondary Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ MySQL Secondary connection successful\n")
	}
}

func repositoryExample() {
	userRepo := repositories.NewUserRepository()

	// Example user
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Create user in MySQL
	fmt.Println("1. Creating user in MySQL...")
	err := userRepo.CreateOnMySQL(user)
	if err != nil {
		fmt.Printf("   ❌ MySQL Create Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ User created in MySQL with ID: %d\n", user.ID)
	}

	// Create user in PostgreSQL
	fmt.Println("2. Creating user in PostgreSQL...")
	user2 := &models.User{
		Username: "testuserpg",
		Email:    "testpg@example.com",
		Password: "password123",
	}
	err = userRepo.CreateOnPostgreSQL(user2)
	if err != nil {
		fmt.Printf("   ❌ PostgreSQL Create Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ User created in PostgreSQL with ID: %d\n", user2.ID)
	}

	// Get user from MySQL
	fmt.Println("3. Retrieving user from MySQL...")
	if user.ID > 0 {
		retrievedUser, err := userRepo.GetFromMySQL(user.ID)
		if err != nil {
			fmt.Printf("   ❌ MySQL Retrieve Error: %v\n", err)
		} else {
			fmt.Printf("   ✅ Retrieved user from MySQL: %s (%s)\n", retrievedUser.Username, retrievedUser.Email)
		}
	}
}

func syncExample() {
	dbService := services.NewDatabaseService()

	fmt.Println("1. Testing data synchronization between MySQL and PostgreSQL...")
	err := dbService.SyncData(func(mysql, postgres *gorm.DB) error {
		// Count users in each database
		var mysqlCount int64
		var postgresCount int64

		mysql.Model(&models.User{}).Count(&mysqlCount)
		postgres.Model(&models.User{}).Count(&postgresCount)

		fmt.Printf("   MySQL users: %d\n", mysqlCount)
		fmt.Printf("   PostgreSQL users: %d\n", postgresCount)

		return nil
	})

	if err != nil {
		fmt.Printf("   ❌ Sync Error: %v\n", err)
	} else {
		fmt.Printf("   ✅ Sync operation completed\n")
	}
}

func healthCheckExample() {
	dbService := services.NewDatabaseService()

	fmt.Println("1. Getting connection statistics...")
	stats, err := dbService.GetConnectionStats()
	if err != nil {
		fmt.Printf("   ❌ Stats Error: %v\n", err)
		return
	}

	for connName, connStats := range stats {
		fmt.Printf("   Connection: %s\n", connName)
		if statsMap, ok := connStats.(map[string]interface{}); ok {
			if connected, exists := statsMap["connected"]; exists && connected == true {
				fmt.Printf("     ✅ Status: Connected\n")
				if openConns, exists := statsMap["open_connections"]; exists {
					fmt.Printf("     📊 Open Connections: %v\n", openConns)
				}
				if inUse, exists := statsMap["in_use"]; exists {
					fmt.Printf("     🔄 In Use: %v\n", inUse)
				}
				if idle, exists := statsMap["idle"]; exists {
					fmt.Printf("     💤 Idle: %v\n", idle)
				}
			} else {
				fmt.Printf("     ❌ Status: Disconnected\n")
				if errMsg, exists := statsMap["error"]; exists {
					fmt.Printf("     🔥 Error: %v\n", errMsg)
				}
			}
		}
		fmt.Println()
	}

	// Test individual connections
	fmt.Println("2. Testing individual connections...")
	connections := []string{"mysql", "postgres", "mysql_secondary"}
	manager := facades.GetManager()

	for _, connName := range connections {
		fmt.Printf("   Testing %s: ", connName)
		if manager.IsConnected(connName) {
			fmt.Println("✅ Healthy")
		} else {
			fmt.Println("❌ Unhealthy")
		}
	}
}
