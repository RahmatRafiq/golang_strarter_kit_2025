package bootstrap

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/cmd"
	"golang_starter_kit_2025/docs"
	"golang_starter_kit_2025/facades"
	"golang_starter_kit_2025/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/urfave/cli/v2"
)

func Init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Validate required environment variables
	validateRequiredEnvVars()

	// Set Gin mode based on environment
	appEnv := helpers.GetEnv("APP_ENV", "production")
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if appEnv == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.TestMode)
	}

	// Connect to database
	facades.ConnectDB()
	defer facades.CloseDB()

	// CLI commands
	app := &cli.App{
		Name:  "Golang Starter Kit",
		Usage: "CLI tool for managing migrations",
		Commands: []*cli.Command{
			cmd.MakeMigrationCommand,
			cmd.MigrationCommand,
			cmd.RollbackCommand,
			cmd.MigrateAllCommand,
			cmd.MigrateFreshCommand,
			cmd.RollbackAllCommand,
			cmd.RollbackBatchCommand,
			cmd.MakeSeederCommand,
			cmd.DBSeedCommand,
			cmd.RollbackSeederCommand,
		},
	}

	// If CLI command provided, run it and exit
	if len(os.Args) > 1 {
		if err := app.Run(os.Args); err != nil {
			log.Fatal(err)
		}
		return
	}

	if !helpers.CheckSecurityWarning() {
		os.Exit(0)
	}

	r := gin.Default()
	facades.ConnectDB()
	// Start web server
	r := Router()
	appPort := helpers.GetEnv("APP_PORT", "9999")
	fmt.Printf("🚀 Server is running on port %s\n", appPort)
	fmt.Printf("📚 Swagger documentation: http://localhost:%s/swagger/index.html\n", appPort)
	fmt.Printf("💚 Health check: http://localhost:%s/health\n", appPort)
	r.Run(":" + appPort)
}

// validateRequiredEnvVars validates that required environment variables are set
func validateRequiredEnvVars() {
	required := map[string]string{
		"DB_CONNECTION":  "Database connection type",
		"JWT_SECRET_KEY": "JWT secret key",
		"APP_PORT":       "Application port",
	}

	r = Router()
	appPort := helpers.GetEnv("APP_PORT", "8080")

	helpers.PrintServerStartupInfo()

	r.Run(":" + appPort)
	var missing []string
	for key, description := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, fmt.Sprintf("%s (%s)", key, description))
		}
	}

	// Check JWT secret is not placeholder
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "your_jwt_secret_key_here" || jwtSecret == "CHANGE_THIS_TO_RANDOM_STRING_AT_LEAST_32_CHARS" {
		log.Fatal("❌ SECURITY ERROR: JWT_SECRET_KEY is still using placeholder value. Please generate a strong secret key using: openssl rand -base64 48")
	}

	if len(missing) > 0 {
		log.Fatalf("❌ ERROR: Required environment variables are not set:\n  - %s", strings.Join(missing, "\n  - "))
	}
}

func Router() *gin.Engine {
	route := gin.Default()

	// Configure CORS properly
	allowedOrigins := helpers.GetEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	route.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(allowedOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Api-Key", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterRoutes(route)

	appName := helpers.GetEnv("APP_NAME", "My App")
	appVersion := helpers.GetEnv("APP_VERSION", "1.0.0")
	appHost := helpers.GetEnv("APP_HOST", "localhost")
	appPort := helpers.GetEnv("APP_PORT", "8080")
	appScheme := helpers.GetEnv("APP_SCHEME", "http")
	appDescription := helpers.GetEnv("APP_DESCRIPTION", "API untuk Supply Chain Retail")

	docs.SwaggerInfo.Title = appName + " API"
	docs.SwaggerInfo.Description = appDescription
	docs.SwaggerInfo.Version = appVersion
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", appHost, appPort)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{appScheme}

	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return route
}
