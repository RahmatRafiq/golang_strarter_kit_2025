package controllers

import (
	"context"
	"net/http"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// DatabaseController handles multi-database operations
type DatabaseController struct{}

// NewDatabaseController creates a new DatabaseController
func NewDatabaseController() *DatabaseController {
	return &DatabaseController{}
}

// GetConnectionsInfo returns information about all database connections
// @Summary Get database connections info
// @Description Get information about all available database connections and their status
// @Tags Database
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/database/connections [get]
func (dc *DatabaseController) GetConnectionsInfo(c *gin.Context) {
	info := helpers.DBHelper.GetConnectionInfo()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}

// TestConnections tests all database connections
// @Summary Test all database connections
// @Description Test connectivity to all configured databases
// @Tags Database
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/database/test [get]
func (dc *DatabaseController) TestConnections(c *gin.Context) {
	results := helpers.TestAllConnections()

	allPassed := true
	for _, err := range results {
		if err != nil {
			allPassed = false
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    allPassed,
		"results":    results,
		"message":    "Database connection tests completed",
		"all_passed": allPassed,
	})
}

// CreateUserInMultipleDBs creates a user in multiple databases to demonstrate multi-DB usage
// @Summary Create user in multiple databases
// @Description Create a user record in both MySQL and PostgreSQL to demonstrate multi-database operations
// @Tags Database
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 200 {object} map[string]interface{}
// @Router /api/database/create-user-multi [post]
func (dc *DatabaseController) CreateUserInMultipleDBs(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	results := make(map[string]interface{})

	// Create in MySQL (primary)
	mysqlDB := helpers.GetMysqlDB()
	if mysqlDB != nil {
		err := mysqlDB.Create(&user).Error
		results["mysql"] = map[string]interface{}{
			"success": err == nil,
			"error":   getErrorString(err),
		}
	} else {
		results["mysql"] = map[string]interface{}{
			"success": false,
			"error":   "MySQL connection not available",
		}
	}

	// Create in PostgreSQL
	postgresDB := helpers.GetPostgresDB()
	if postgresDB != nil {
		err := postgresDB.Create(&user).Error
		results["postgres"] = map[string]interface{}{
			"success": err == nil,
			"error":   getErrorString(err),
		}
	} else {
		results["postgres"] = map[string]interface{}{
			"success": false,
			"error":   "PostgreSQL connection not available",
		}
	}

	// Create in SQLite
	sqliteDB := helpers.GetSqliteDB()
	if sqliteDB != nil {
		err := sqliteDB.Create(&user).Error
		results["sqlite"] = map[string]interface{}{
			"success": err == nil,
			"error":   getErrorString(err),
		}
	} else {
		results["sqlite"] = map[string]interface{}{
			"success": false,
			"error":   "SQLite connection not available",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User creation attempted in multiple databases",
		"results": results,
		"user":    user,
	})
}

// CreateDocumentInMongoDB creates a document in MongoDB
// @Summary Create document in MongoDB
// @Description Create a document in MongoDB to demonstrate NoSQL operations
// @Tags Database
// @Accept json
// @Produce json
// @Param document body map[string]interface{} true "Document data"
// @Success 200 {object} map[string]interface{}
// @Router /api/database/create-mongo-document [post]
func (dc *DatabaseController) CreateDocumentInMongoDB(c *gin.Context) {
	var document map[string]interface{}
	if err := c.ShouldBindJSON(&document); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Add timestamp
	document["created_at"] = time.Now()
	document["updated_at"] = time.Now()

	// Get MongoDB collection
	collection := helpers.GetMongoCollection("test_db", "documents")
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "MongoDB connection not available",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create document in MongoDB",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Document created successfully in MongoDB",
		"document":  document,
		"insert_id": result.InsertedID,
	})
}

// GetDocumentsFromMongoDB retrieves documents from MongoDB
// @Summary Get documents from MongoDB
// @Description Retrieve documents from MongoDB collection
// @Tags Database
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/database/mongo-documents [get]
func (dc *DatabaseController) GetDocumentsFromMongoDB(c *gin.Context) {
	collection := helpers.GetMongoCollection("test_db", "documents")
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "MongoDB connection not available",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to query MongoDB",
			"error":   err.Error(),
		})
		return
	}
	defer cursor.Close(ctx)

	var documents []map[string]interface{}
	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		documents = append(documents, doc)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Documents retrieved successfully from MongoDB",
		"documents": documents,
		"count":     len(documents),
	})
}

// TransactionExample demonstrates transaction across multiple databases
// @Summary Transaction example
// @Description Demonstrate transaction operations across multiple databases
// @Tags Database
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "Transaction data"
// @Success 200 {object} map[string]interface{}
// @Router /api/database/transaction-example [post]
func (dc *DatabaseController) TransactionExample(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	results := make(map[string]interface{})

	// MySQL Transaction
	err := helpers.DBHelper.ExecuteInTransaction("mysql", func(tx *gorm.DB) error {
		// Simulate some database operations
		// For example, creating a user
		user := models.User{
			Username: data["username"].(string),
			Email:    data["email"].(string),
			Password: "dummy_password", // In real scenario, this should be hashed
		}
		return tx.Create(&user).Error
	})

	results["mysql_transaction"] = map[string]interface{}{
		"success": err == nil,
		"error":   getErrorString(err),
	}

	// MongoDB Transaction
	mongoErr := helpers.DBHelper.ExecuteMongoTransaction("mongodb", func(sc mongo.SessionContext) error {
		collection := helpers.GetMongoCollection("test_db", "transactions")
		if collection == nil {
			return mongo.ErrNoDocuments
		}

		document := bson.M{
			"name":       data["name"],
			"email":      data["email"],
			"username":   data["username"],
			"created_at": time.Now(),
		}

		_, err := collection.InsertOne(sc, document)
		return err
	})

	results["mongodb_transaction"] = map[string]interface{}{
		"success": mongoErr == nil,
		"error":   getErrorString(mongoErr),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transaction operations completed",
		"results": results,
		"data":    data,
	})
}

// SwitchDatabase demonstrates switching between database connections
// @Summary Switch database connection
// @Description Demonstrate switching between different database connections
// @Tags Database
// @Accept json
// @Produce json
// @Param connection query string true "Connection name"
// @Success 200 {object} map[string]interface{}
// @Router /api/database/switch [post]
func (dc *DatabaseController) SwitchDatabase(c *gin.Context) {
	connectionName := c.Query("connection")
	if connectionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Connection name is required",
		})
		return
	}

	// Get the database connection by name
	db := helpers.GetDBByName(connectionName)
	if db == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid connection name or connection not available",
		})
		return
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get database instance",
			"error":   err.Error(),
		})
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database connection test failed",
			"error":   err.Error(),
		})
		return
	}

	// Get connection manager for more info
	manager := facades.GetManager()
	availableConnections := manager.GetConnectionNames()

	c.JSON(http.StatusOK, gin.H{
		"success":               true,
		"message":               "Successfully switched to database connection",
		"current_connection":    connectionName,
		"available_connections": availableConnections,
		"connection_status":     "active",
	})
}

// helper function to convert error to string
func getErrorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
