package facades

// This file is kept for backward compatibility
// All new functionality is in database_manager.go

import (
	"database/sql"

	"gorm.io/gorm"
)

var (
	// DB is the global GORM database instance (for backward compatibility)
	DB *gorm.DB
	// SqlDB is the global raw SQL DB instance (for backward compatibility)
	SqlDB *sql.DB
)

// ConnectDB initializes database connections (for backward compatibility)
func ConnectDB(envFiles ...string) *gorm.DB {
	manager := GetManager()
	db := manager.Connection("")

	// Set global variables for backward compatibility
	DB = db
	if db != nil {
		SqlDB, _ = db.DB()
	}

	return db
}

// CloseDB closes all database connections (for backward compatibility)
func CloseDB() {
	GetManager().CloseAll()
}
