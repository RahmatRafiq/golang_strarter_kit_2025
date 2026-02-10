package bootstrap

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/cmd"
	"golang_starter_kit_2025/docs"
	"golang_starter_kit_2025/facades"
	"golang_starter_kit_2025/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/urfave/cli/v2"
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("No .env file found, using environment variables")
	}

	// Initialize structured logging
	helpers.InitLogger()

	validateRequiredEnvVars()

	appEnv := helpers.GetEnv("APP_ENV", "production")
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if appEnv == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.TestMode)
	}

	facades.ConnectDB()
	defer facades.CloseDB()

	// Initialize Redis (optional - continues even if Redis is unavailable)
	if err := helpers.InitRedis(); err != nil {
		log.Warn().Err(err).Msg("Redis not available, caching disabled")
	}
	defer helpers.CloseRedis()

	// Start health history collector (only if Redis is available)
	if helpers.RedisClient != nil {
		services.StartHealthHistoryCollector()
	}

	app := &cli.App{
		Name:  "Golang Starter Kit",
		Usage: "CLI tool for managing migrations",
		Commands: []*cli.Command{
			// Migration commands
			cmd.MakeMigrationCommand,
			cmd.MigrationCommand,
			cmd.RollbackCommand,
			cmd.MigrateAllCommand,
			cmd.MigrateFreshCommand,
			cmd.MigrateStatusCommand,
			cmd.MigrateResetCommand,
			cmd.RollbackAllCommand,
			cmd.RollbackBatchCommand,
			cmd.MigrateLogsCommand, // NEW: View migration logs

			// Migration lock commands
			cmd.MigrateLockStatusCommand,  // NEW: Check lock status
			cmd.MigrateLockReleaseCommand, // NEW: Force release lock

			// Seeder commands
			cmd.MakeSeederCommand,
			cmd.DBSeedCommand,
			cmd.RollbackSeederCommand,

			// Factory commands
			cmd.MakeFactoryCommand, // NEW: Generate factories

			// Database commands
			cmd.DBWipeCommand,
			cmd.DBStatusCommand,
			cmd.DBConnectionsCommand,
		},
	}

	if len(os.Args) > 1 {
		if err := app.Run(os.Args); err != nil {
			log.Fatal().Err(err).Msg("CLI command failed")
		}
		return
	}

	if !helpers.CheckSecurityWarning() {
		os.Exit(0)
	}

	r := Router()
	appPort := helpers.GetEnv("APP_PORT", "8080")

	helpers.PrintServerStartupInfo()

	r.Run(":" + appPort)
}

func isWeakSecret(secret string) bool {
	// Check for common weak patterns
	weakPatterns := []string{
		"secret", "password", "test", "demo", "admin",
		"12345", "qwerty", "abc123", "changeme",
	}

	lowerSecret := strings.ToLower(secret)
	for _, pattern := range weakPatterns {
		if strings.Contains(lowerSecret, pattern) {
			return true
		}
	}

	// Check if it's too simple (all same character, sequential, etc)
	if len(secret) > 0 {
		allSame := true
		for i := 1; i < len(secret); i++ {
			if secret[i] != secret[0] {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}

	return false
}

func validateRequiredEnvVars() {
	required := map[string]string{
		"DB_CONNECTION":  "Database connection type",
		"JWT_SECRET_KEY": "JWT secret key",
		"APP_PORT":       "Application port",
	}

	var missing []string
	for key, description := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, fmt.Sprintf("%s (%s)", key, description))
		}
	}

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "your_jwt_secret_key_here" || jwtSecret == "CHANGE_THIS_TO_RANDOM_STRING_AT_LEAST_32_CHARS" {
		log.Fatal().
			Msg("SECURITY ERROR: JWT_SECRET_KEY is still using placeholder value. Please generate a strong secret key using: openssl rand -base64 48")
	}

	if len(jwtSecret) < 32 {
		log.Fatal().
			Int("length", len(jwtSecret)).
			Msg("SECURITY ERROR: JWT_SECRET_KEY must be at least 32 characters long for adequate security")
	}

	// Warn if key appears weak
	if isWeakSecret(jwtSecret) {
		log.Warn().Msg("JWT_SECRET_KEY appears to be weak. Consider using: openssl rand -base64 48")
	}

	if len(missing) > 0 {
		log.Fatal().
			Strs("missing_vars", missing).
			Msg("Required environment variables are not set")
	}
}

func Router() *gin.Engine {
	route := gin.Default()

	// Gzip compression middleware (5-10x smaller responses)
	route.Use(gzip.Gzip(gzip.DefaultCompression))

	allowedOrigins := helpers.GetEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	route.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(allowedOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Api-Key", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID", "X-Cache"},
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
