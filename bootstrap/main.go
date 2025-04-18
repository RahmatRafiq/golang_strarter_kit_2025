package bootstrap

import (
	"fmt"
	"log"

	"golang_strarter_kit_2025/app/helpers"
	"golang_strarter_kit_2025/docs"
	"golang_strarter_kit_2025/facades"
	"golang_strarter_kit_2025/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() {
	// Memuat variabel lingkungan dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Inisialisasi koneksi facades
	facades.ConnectDB()

	defer facades.CloseDB()

	r := Router()
	fmt.Println("Server is running on port 8080")
	r.Run(":8080")
}

func Router() *gin.Engine {
	// Menginisialisasi Gin Router
	route := gin.Default()

	route.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "PUT", "DELETE"},
		AllowHeaders: []string{"*"},
	}))

	// Daftarkan routes
	routes.RegisterRoutes(route)

	docs.SwaggerInfo.Title = "Supply Chain Retail API"
	docs.SwaggerInfo.Description = "API untuk Supply Chain Retail"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = helpers.GetEnv("SWAGGER_HOST", "localhost:8080")
	docs.SwaggerInfo.BasePath = "/"

	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return route
}
