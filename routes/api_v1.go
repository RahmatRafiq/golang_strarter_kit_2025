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
	passwordResetService *services.PasswordResetService,
	emailVerificationService *services.EmailVerificationService,
	oauthService *services.OAuthService,
	storageService *services.StorageService,
) {
	// API v1 group
	v1 := route.Group("/api/v1")

	// ====================
	// PUBLIC ROUTES (No Auth Required)
	// ====================

	// Auth routes with stricter rate limiting (5 attempts per 15min)
	authController := controllers.NewAuthController(*authService)
	passwordResetController := controllers.NewPasswordResetController(passwordResetService)
	emailVerificationController := controllers.NewEmailVerificationController(emailVerificationService)
	authPublicRoutes := v1.Group("/auth")
	authPublicRoutes.Use(middleware.AuthRateLimiter())
	{
		authPublicRoutes.POST("/login", authController.Login)
		authPublicRoutes.POST("/refresh", authController.Refresh)
		authPublicRoutes.POST("/forgot-password", passwordResetController.ForgotPassword)
		authPublicRoutes.POST("/reset-password", passwordResetController.ResetPassword)
		authPublicRoutes.POST("/verify-email", emailVerificationController.VerifyEmail)
	}

	oauthController := controllers.NewOAuthController(oauthService)
	oauthRoutes := v1.Group("/auth")
	{
		oauthRoutes.GET("/google", oauthController.GoogleLogin)
		oauthRoutes.GET("/google/callback", oauthController.GoogleCallback)
		oauthRoutes.GET("/github", oauthController.GitHubLogin)
		oauthRoutes.GET("/github/callback", oauthController.GitHubCallback)
	}

	// ====================
	// PROTECTED ROUTES (Auth Required)
	// ====================

	// Auth protected routes
	authRoutes := v1.Group("/auth")
	authRoutes.Use(middleware.AuthMiddleware())
	{
		authRoutes.POST("/logout", authController.Logout)
		authRoutes.POST("/resend-verification", emailVerificationController.SendVerification)
	}

	// User routes
	userController := controllers.NewUserController(*userService)
	userRoutes := v1.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		userRoutes.GET("", userController.List)
		userRoutes.GET("/:id", userController.Get)
		userRoutes.POST("", userController.Create)
		userRoutes.PUT("/:id", userController.Update)
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
		roleRoutes.POST("", roleController.Create)
		roleRoutes.PUT("/:id", roleController.Update)
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
		permissionRoutes.POST("", permissionController.Create)
		permissionRoutes.PUT("/:id", permissionController.Update)
		permissionRoutes.DELETE("/:id", permissionController.Delete)
	}

	// Product routes
	productController := controllers.NewProductController(*productService)
	productRoutes := v1.Group("/products")
	productRoutes.Use(middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		productRoutes.GET("", productController.List)
		productRoutes.GET("/:id", productController.Get)
		productRoutes.POST("", productController.Create)
		productRoutes.PUT("/:id", productController.Update)
		productRoutes.DELETE("/:id", productController.Delete)
	}

	// Category routes
	categoryController := controllers.NewCategoryController(*categoryService)
	categoryRoutes := v1.Group("/categories")
	categoryRoutes.Use(middleware.AuthMiddleware(), middleware.UserRateLimiter())
	{
		categoryRoutes.GET("", categoryController.List)
		categoryRoutes.GET("/:id", categoryController.Get)
		categoryRoutes.POST("", categoryController.Create)
		categoryRoutes.PUT("/:id", categoryController.Update)
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

	// Storage routes (file upload/download)
	if storageService != nil {
		storageController := controllers.NewStorageController(storageService)
		storageRoutes := v1.Group("/storage")
		storageRoutes.Use(middleware.AuthMiddleware())
		{
			storageRoutes.POST("/upload", storageController.UploadFile)
			storageRoutes.GET("/url/:filename", storageController.GetFileURL)
			storageRoutes.DELETE("/:filename", storageController.DeleteFile)
		}
	}
}
