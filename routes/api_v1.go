package routes

import (
	"golang_starter_kit_2025/app/controllers"
	"golang_starter_kit_2025/app/middleware"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
)

// RegisterV1Routes registers all API v1 routes
// API versioning allows for breaking changes without affecting existing clients
func RegisterV1Routes(
	route *gin.Engine,
	userService *services.UserService,
	roleService *services.RoleService,
	permissionService *services.PermissionService,
	productService *services.ProductService,
	categoryService *services.CategoryService,
	authService *services.AuthService,
) {
	// API v1 group
	v1 := route.Group("/api/v1")

	// ====================
	// PUBLIC ROUTES (No Auth Required)
	// ====================

	// Auth routes with stricter rate limiting (5 attempts per 15min)
	authController := controllers.NewAuthController(*authService)
	authPublicRoutes := v1.Group("/auth")
	authPublicRoutes.Use(middleware.AuthRateLimiter())
	{
		authPublicRoutes.POST("/login", authController.Login)     // Changed from PUT to POST (RESTful)
		authPublicRoutes.POST("/refresh", authController.Refresh) // Refresh token
	}

	// ====================
	// PROTECTED ROUTES (Auth Required)
	// ====================

	// Auth protected routes
	authRoutes := v1.Group("/auth")
	authRoutes.Use(middleware.AuthMiddleware())
	{
		authRoutes.POST("/logout", authController.Logout)
	}

	// User routes
	userController := controllers.NewUserController(*userService)
	userRoutes := v1.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		userRoutes.GET("", userController.List)
		userRoutes.GET("/:id", userController.Get)
		userRoutes.POST("", userController.Put)         // Changed from PUT to POST for create
		userRoutes.PUT("/:id", userController.Put)      // PUT for update
		userRoutes.DELETE("/:id", userController.Delete)
		userRoutes.POST("/:id/roles", userController.AssignRoles)
		userRoutes.GET("/:id/roles", userController.GetRoles)
	}

	// Role routes
	roleController := controllers.NewRoleController(*roleService)
	roleRoutes := v1.Group("/roles")
	roleRoutes.Use(middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		roleRoutes.GET("", roleController.List)
		roleRoutes.POST("", roleController.Put)         // Changed from PUT to POST for create
		roleRoutes.PUT("/:id", roleController.Put)      // PUT for update
		roleRoutes.DELETE("/:id", roleController.Delete)
		roleRoutes.POST("/:id/permissions", roleController.AssignPermissions)
		roleRoutes.GET("/:id/permissions", roleController.GetPermissions)
	}

	// Permission routes
	permissionController := controllers.NewPermissionController(*permissionService)
	permissionRoutes := v1.Group("/permissions")
	permissionRoutes.Use(middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		permissionRoutes.GET("", permissionController.List)
		permissionRoutes.POST("", permissionController.Put)         // Changed from PUT to POST for create
		permissionRoutes.PUT("/:id", permissionController.Put)      // PUT for update
		permissionRoutes.DELETE("/:id", permissionController.Delete)
	}

	// Product routes
	productController := controllers.NewProductController(*productService)
	productRoutes := v1.Group("/products")
	productRoutes.Use(middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		productRoutes.GET("", productController.GetAll)
		productRoutes.GET("/:id", productController.GetByID)
		productRoutes.POST("", productController.Put)         // Changed from PUT to POST for create
		productRoutes.PUT("/:id", productController.Put)      // PUT for update
		productRoutes.DELETE("/:id", productController.Delete)
	}

	// Category routes
	categoryController := controllers.NewCategoryController(*categoryService)
	categoryRoutes := v1.Group("/categories")
	categoryRoutes.Use(middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		categoryRoutes.GET("", categoryController.List)
		categoryRoutes.GET("/:id", categoryController.Get)
		categoryRoutes.POST("", categoryController.Put)         // Changed from PUT to POST for create
		categoryRoutes.PUT("/:id", categoryController.Put)      // PUT for update
		categoryRoutes.DELETE("/:id", categoryController.Delete)
	}

	// Database management routes
	databaseController := controllers.NewDatabaseController()
	databaseRoutes := v1.Group("/database")
	databaseRoutes.Use(middleware.AuthMiddleware())
	{
		databaseRoutes.GET("/status", databaseController.GetConnectionStatus)
		databaseRoutes.GET("/health", databaseController.HealthCheck)
		databaseRoutes.GET("/test", databaseController.TestConnection)
	}
}
