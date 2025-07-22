package controllers

import (
	"net/http"

	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
)

// DatabaseController handles database management endpoints
type DatabaseController struct {
	dbService *services.DatabaseService
}

// NewDatabaseController creates a new database controller
func NewDatabaseController() *DatabaseController {
	return &DatabaseController{
		dbService: services.NewDatabaseService(),
	}
}

// GetConnectionStatus returns the status of all database connections
// @Summary Get database connection status
// @Description Get the status and statistics of all configured database connections
// @Tags Database
// @Produce json
// @Success 200 {object} map[string]interface{} "Connection status and statistics"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/database/status [get]
func (ctrl *DatabaseController) GetConnectionStatus(c *gin.Context) {
	stats, err := ctrl.dbService.GetConnectionStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to get connection statistics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Database connection status retrieved successfully",
		"data":    stats,
	})
}

// HealthCheck performs a health check on all database connections
// @Summary Database health check
// @Description Perform a health check on all configured database connections
// @Tags Database
// @Produce json
// @Success 200 {object} map[string]interface{} "Health check results"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/database/health [get]
func (ctrl *DatabaseController) HealthCheck(c *gin.Context) {
	health := make(map[string]interface{})
	connections := []string{"mysql", "postgres", "mysql_secondary"}

	for _, connName := range connections {
		_, err := ctrl.dbService.GetConnection(connName)
		if err != nil {
			health[connName] = gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			}
		} else {
			health[connName] = gin.H{
				"status": "healthy",
			}
		}
	}

	// Determine overall health
	allHealthy := true
	for _, conn := range health {
		if connMap, ok := conn.(gin.H); ok {
			if status, exists := connMap["status"]; exists && status != "healthy" {
				allHealthy = false
				break
			}
		}
	}

	statusCode := http.StatusOK
	if !allHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"error":       !allHealthy,
		"message":     "Database health check completed",
		"overall":     map[string]interface{}{"healthy": allHealthy},
		"connections": health,
	})
}

// TestConnection tests a specific database connection
// @Summary Test database connection
// @Description Test a specific database connection by name
// @Tags Database
// @Param connection query string true "Connection name (mysql, postgres, mysql_secondary)"
// @Produce json
// @Success 200 {object} map[string]interface{} "Connection test result"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/database/test [get]
func (ctrl *DatabaseController) TestConnection(c *gin.Context) {
	connectionName := c.Query("connection")
	if connectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Connection name is required",
		})
		return
	}

	// Test the connection
	db, err := ctrl.dbService.GetConnection(connectionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      true,
			"message":    "Failed to connect to database",
			"connection": connectionName,
			"details":    err.Error(),
		})
		return
	}

	// Try to execute a simple query
	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      true,
			"message":    "Failed to get SQL DB instance",
			"connection": connectionName,
			"details":    err.Error(),
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      true,
			"message":    "Failed to ping database",
			"connection": connectionName,
			"details":    err.Error(),
		})
		return
	}

	stats := sqlDB.Stats()
	c.JSON(http.StatusOK, gin.H{
		"error":      false,
		"message":    "Database connection test successful",
		"connection": connectionName,
		"stats": gin.H{
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
		},
	})
}
