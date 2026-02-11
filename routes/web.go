package routes

import (
	"net/http"
	"time"

	"golang_starter_kit_2025/app/controllers"
	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/middleware"
	"golang_starter_kit_2025/app/repositories"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/facades"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterRoutes(route *gin.Engine) {
	// Apply observability middleware (must be first for accurate metrics)
	route.Use(middleware.ZerologMiddleware())  // Structured logging
	route.Use(middleware.MetricsMiddleware())  // Prometheus metrics

	// Apply performance middleware
	route.Use(middleware.ContextTimeoutMiddleware())  // Request timeout
	cacheTTL := time.Duration(helpers.GetEnvInt("CACHE_TTL_MINUTES", 5)) * time.Minute
	route.Use(middleware.CacheMiddleware(cacheTTL))   // HTTP caching

	// Apply global security middleware
	route.Use(middleware.GlobalRateLimiter())   // Rate limiting
	route.Use(middleware.SQLInjectionCheck())   // SQL injection protection
	route.Use(middleware.XSSProtection())       // XSS protection
	route.Use(middleware.SanitizeInput())       // Input sanitization

	// Start database metrics collector
	helpers.CollectDBMetrics()
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

	// Prometheus metrics endpoint
	route.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// ========================================
	// API v1 Routes (Recommended)
	// ========================================
	RegisterV1Routes(route, userService, roleService, permissionService, productService, categoryService, authService)

	// ========================================
	// Legacy Routes (Deprecated - For Backward Compatibility)
	// TODO: Remove these in v2.0.0
	// ========================================
	authController := controllers.NewAuthController(*authService)

	// Auth routes with stricter rate limiting (5 attempts per 15min)
	authPublicRoutes := route.Group("/auth")
	authPublicRoutes.Use(middleware.AuthRateLimiter())
	{
		authPublicRoutes.PUT("/login", authController.Login) // DEPRECATED: Use POST /api/v1/auth/login
	}

	// Refresh endpoint with moderate rate limiting
	route.POST("/auth/refresh", authController.Refresh) // DEPRECATED: Use /api/v1/auth/refresh

	// Protected auth routes
	authRoutes := route.Group("/auth").Use(middleware.AuthMiddleware())
	{
		authRoutes.GET("/logout", authController.Logout) // DEPRECATED: Use POST /api/v1/auth/logout
	}

	categoryController := controllers.NewCategoryController(*categoryService)
	categoryRoutes := route.Group("/categories", middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		categoryRoutes.GET("/", categoryController.List)   // DEPRECATED: Use /api/v1/categories
		categoryRoutes.GET("/:id", categoryController.Get)
		categoryRoutes.PUT("/", categoryController.Put)
		categoryRoutes.DELETE("/:id", categoryController.Delete)
	}

	productController := controllers.NewProductController(*productService)
	productRoutes := route.Group("/products", middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		productRoutes.GET("/", productController.GetAll)   // DEPRECATED: Use /api/v1/products
		productRoutes.GET("/:id", productController.GetByID)
		productRoutes.PUT("/", productController.Put)
		productRoutes.DELETE("/:id", productController.Delete)
	}

	userController := controllers.NewUserController(*userService)
	userRoutes := route.Group("/users", middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		userRoutes.GET("", userController.List)           // DEPRECATED: Use /api/v1/users
		userRoutes.GET("/:id", userController.Get)
		userRoutes.PUT("", userController.Put)
		userRoutes.DELETE("/:id", userController.Delete)
		userRoutes.POST("/:id/roles", userController.AssignRoles)
		userRoutes.GET("/:id/roles", userController.GetRoles)
	}

	roleController := controllers.NewRoleController(*roleService)
	roleRoutes := route.Group("/roles", middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		roleRoutes.GET("", roleController.List)           // DEPRECATED: Use /api/v1/roles
		roleRoutes.PUT("", roleController.Put)
		roleRoutes.DELETE("/:id", roleController.Delete)
		roleRoutes.POST("/:id/permissions", roleController.AssignPermissions)
		roleRoutes.GET("/:id/permissions", roleController.GetPermissions)
	}

	permissionController := controllers.NewPermissionController(*permissionService)
	permissionRoutes := route.Group("/permissions", middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		permissionRoutes.GET("", permissionController.List) // DEPRECATED: Use /api/v1/permissions
		permissionRoutes.PUT("", permissionController.Put)
		permissionRoutes.DELETE("/:id", permissionController.Delete)
	}

	// File serving (not versioned - utility endpoint)
	fileController := controllers.NewFileController()
	fileRoutes := route.Group("/file")
	{
		fileRoutes.GET("/:key/:filename", fileController.ServeFile)
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
