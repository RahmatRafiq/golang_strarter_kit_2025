package bootstrap

import (
	"fmt"
	"log"
	"os"

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
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	facades.ConnectDB()
	defer facades.CloseDB()

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

	if len(os.Args) > 1 {
		if err := app.Run(os.Args); err != nil {
			log.Fatal(err)
		}
		return
	}

	r := gin.Default()
	facades.ConnectDB()

	defer facades.CloseDB()

	r = Router()
	appPort := helpers.GetEnv("APP_PORT", "8080")
	fmt.Printf("Server is running on port %s\n", appPort)
	r.Run(":" + appPort)
}

func Router() *gin.Engine {
	route := gin.Default()

	route.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "PUT", "DELETE"},
		AllowHeaders: []string{"*"},
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
