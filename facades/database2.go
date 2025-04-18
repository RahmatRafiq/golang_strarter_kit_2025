package facades

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(c ...*gin.Context) (db *gorm.DB) {
	db, err := gorm.Open(mysql.Open("ubaid:plmoknijb@tcp(localhost:3307)/supply_chain_retail?parseTime=true"), &gorm.Config{})
	if err != nil {
		// Return Gin response with 500 status code and error message
		ginErr := fmt.Errorf("failed to connect to the database: %v", err)
		ginErrResponse := gin.H{
			"error": ginErr.Error(),
		}
		c[0].JSON(http.StatusInternalServerError, ginErrResponse)
		return nil
	}

	// dbq, _ := db.DB()
	// dbq.Close()

	return db
}