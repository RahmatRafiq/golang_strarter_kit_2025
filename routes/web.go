package routes

import (
	"net/http"

	"golang_starter_kit_2025/app/controllers"
	"golang_starter_kit_2025/app/middleware"
	"golang_starter_kit_2025/app/repositories"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/facades"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine) {
	// Apply structured logging middleware
	route.Use(middleware.ZerologMiddleware())

	userRepo := repositories.NewUserRepository(facades.DB)
	roleRepo := repositories.NewRoleRepository(facades.DB)
	permissionRepo := repositories.NewPermissionRepository(facades.DB)
	productRepo := repositories.NewProductRepository(facades.DB)
	categoryRepo := repositories.NewCategoryRepository(facades.DB)
	authRepo := repositories.NewAuthRepository(facades.DB)

	userService := services.NewUserService(userRepo)
	roleService := services.NewRoleService(roleRepo, permissionRepo)
	permissionService := services.NewPermissionService(permissionRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo, productRepo)
	authService := services.NewAuthService(authRepo)

	testService := services.TestService{}
	testController := controllers.NewTestController(testService)
	testRoutes := route.Group("/tests")
	{
		testRoutes.GET("", testController.List)
		testRoutes.GET(":id", testController.Get)
		testRoutes.POST("", testController.Create)
		testRoutes.PUT(":id", testController.Update)
		testRoutes.DELETE(":id", testController.Delete)
	}

	controller := controllers.Controller{}
	route.GET("", controller.HelloWorld)

	authController := controllers.NewAuthController(*authService)
	route.PUT("/auth/login", authController.Login)
	authRoutes := route.Group("/auth").Use(middleware.AuthMiddleware())
	{
		authRoutes.GET("/logout", authController.Logout)
		authRoutes.GET("/refresh", authController.Refresh)
	}

	categoryController := controllers.NewCategoryController(*categoryService)
	categoryRoutes := route.Group("/categories", middleware.AuthMiddleware())
	{
		categoryRoutes.GET("/", categoryController.List)
		categoryRoutes.GET("/:id", categoryController.Get)
		categoryRoutes.PUT("/", categoryController.Put)
		categoryRoutes.DELETE("/:id", categoryController.Delete)
	}

	productController := controllers.NewProductController(*productService)
	productRoutes := route.Group("/products", middleware.AuthMiddleware())
	{
		productRoutes.GET("/", productController.GetAll)
		productRoutes.GET("/:id", productController.GetByID)
		productRoutes.PUT("/", productController.Put)
		productRoutes.DELETE("/:id", productController.Delete)
	}

	userController := controllers.NewUserController(*userService)
	userRoutes := route.Group("/users", middleware.AuthMiddleware())
	{
		userRoutes.GET("", userController.List)
		userRoutes.GET("/:id", userController.Get)
		userRoutes.PUT("", userController.Put)
		userRoutes.DELETE("/:id", userController.Delete)
		userRoutes.POST("/:id/roles", userController.AssignRoles)
		userRoutes.GET("/:id/roles", userController.GetRoles)
	}

	roleController := controllers.NewRoleController(*roleService)
	roleRoutes := route.Group("/roles", middleware.AuthMiddleware())
	{
		roleRoutes.GET("", roleController.List)
		roleRoutes.PUT("", roleController.Put)
		roleRoutes.DELETE("/:id", roleController.Delete)
		roleRoutes.POST("/:id/permissions", roleController.AssignPermissions)
		roleRoutes.GET("/:id/permissions", roleController.GetPermissions)
	}

	permissionController := controllers.NewPermissionController(*permissionService)
	permissionRoutes := route.Group("/permissions", middleware.AuthMiddleware())
	{
		permissionRoutes.GET("", permissionController.List)
		permissionRoutes.PUT("", permissionController.Put)
		permissionRoutes.DELETE("/:id", permissionController.Delete)
	}

	fileController := controllers.NewFileController()
	fileRoutes := route.Group("/file")
	{
		fileRoutes.GET("/:key/:filename", fileController.ServeFile)
	}

	// Database management routes (protected by AuthMiddleware)
	databaseController := controllers.NewDatabaseController()
	databaseRoutes := route.Group("/api/database")
	{
		databaseRoutes.GET("/status", databaseController.GetConnectionStatus)
		databaseRoutes.GET("/health", databaseController.HealthCheck)
		databaseRoutes.GET("/test", databaseController.TestConnection)
	}

	// Endpoint untuk mengecek kesehatan koneksi facades
	route.GET("/health", func(c *gin.Context) {
		sqlDB, err := facades.DB.DB() // Mengambil facades/sql *DB dari GORM *DB
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get facades connection",
				"error":   err.Error(),
			})
			return
		}

		err = sqlDB.Ping() // Menggunakan sqlDB untuk ping ke facades
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "facades connection failed",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "facades is connected",
			"facades": "supply_chain_retail", // Sesuaikan dengan nama facades Anda
		})
	})

	// Multi-database health check endpoint (public)
	route.GET("/health/databases", func(c *gin.Context) {
		manager := facades.GetManager()
		health := make(map[string]interface{})
		connections := []string{"mysql", "postgres", "mysql_secondary"}

		allHealthy := true
		for _, connName := range connections {
			if manager.IsConnected(connName) {
				stats, err := manager.GetConnectionStats(connName)
				if err == nil {
					health[connName] = gin.H{
						"status": "healthy",
						"stats":  stats,
					}
				} else {
					health[connName] = gin.H{
						"status": "unhealthy",
						"error":  err.Error(),
					}
					allHealthy = false
				}
			} else {
				health[connName] = gin.H{
					"status": "disconnected",
				}
				allHealthy = false
			}
		}

		statusCode := http.StatusOK
		if !allHealthy {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, gin.H{
			"overall_health": allHealthy,
			"connections":    health,
		})
	})
}
