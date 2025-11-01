# AI Context Documentation - Golang Starter Kit 2025

> **Tujuan**: Dokumentasi ini dibuat sebagai "base pemikiran" untuk AI assistant agar dapat memahami arsitektur, struktur, dan pattern yang digunakan dalam codebase ini tanpa perlu screening ulang setiap kali ada request optimasi atau update.

---

## 📋 Table of Contents
1. [Arsitektur Overview](#arsitektur-overview)
2. [Database Architecture](#database-architecture)
3. [Flow Request-Response](#flow-request-response)
4. [Models & Relationships](#models--relationships)
5. [Authentication & Authorization](#authentication--authorization)
6. [API Endpoints Structure](#api-endpoints-structure)
7. [Helper Functions](#helper-functions)
8. [Testing Structure](#testing-structure)
9. [Configuration System](#configuration-system)
10. [Best Practices](#best-practices)

---

## 🏗️ Arsitektur Overview

### Pattern yang Digunakan
- **MVC Pattern**: Model-View-Controller (tanpa View karena API)
- **Repository Pattern**: TIDAK diimplementasikan (logic langsung di Services)
- **Facade Pattern**: Database facade untuk abstraksi koneksi database
- **Manager Pattern**: Database Manager untuk multi-database connections
- **Dependency Injection**: Constructor injection di Controllers

### Struktur Folder Utama
```
app/
├── controllers/     → HTTP request handling (thin controllers)
├── services/        → Business logic (fat services)
├── models/          → GORM models & database entities
├── requests/        → Request validation structs
├── responses/       → Response formatting structs
├── middleware/      → HTTP middleware (auth, logging)
├── helpers/         → Utility functions
├── casts/           → Data transformation objects (JWT claims, tokens)
├── handlers/        → Response handlers
└── database/        → Migration & seeder management
```

### Dependency Flow
```
main.go
  └─> bootstrap/main.go (Init)
       ├─> facades.ConnectDB()                    # Inisialisasi database
       ├─> routes.RegisterRoutes(router)          # Register semua routes
       └─> router.Run()                           # Start server

Request Flow:
HTTP Request → Middleware → Controller → Service → Model/DB → Response
```

---

## 🗄️ Database Architecture

### Multi-Database Support
Project ini menggunakan **Database Manager Pattern** untuk mendukung koneksi ke multiple databases secara bersamaan.

#### File Kunci:
- `config/database_config.go` - Konfigurasi database dari environment variables
- `database/manager.go` - Manager untuk handle multiple connections
- `facades/database.go` - Facade untuk akses database yang lebih mudah

#### Supported Databases:
1. **MySQL** (primary/default) - `DB_CONNECTION=mysql`
2. **PostgreSQL** - Multi-connection support
3. **MySQL Secondary** - Multiple MySQL instances
4. **MongoDB** (optional) - Via mongo_config.go

#### Koneksi Database:
```go
// Mendapatkan koneksi default
db := facades.GetDB()  // atau facades.DB

// Mendapatkan koneksi spesifik
mysql, _ := facades.MySQL()
postgres, _ := facades.PostgreSQL()
mysqlSecondary, _ := facades.MySQLSecondary()

// Atau via manager
manager := facades.GetManager()
conn, _ := manager.GetConnection("postgres")
```

#### Environment Variables:
```bash
# Default Connection
DB_CONNECTION=mysql  # Default database yang digunakan

# MySQL Primary
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang-starter-kit-mysql
MYSQL_USER=root
MYSQL_PASSWORD=xxx
MYSQL_CHARSET=utf8mb4
MYSQL_TIMEZONE=Local
MYSQL_MAX_IDLE_CONNS=10
MYSQL_MAX_OPEN_CONNS=200
MYSQL_CONN_MAX_LIFETIME=15m
MYSQL_CONN_MAX_IDLE_TIME=5m

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=golang-starter-kit-pgs
POSTGRES_USER=root
POSTGRES_PASSWORD=xxx
POSTGRES_SSLMODE=disable
POSTGRES_TIMEZONE=UTC
# ... (sama seperti MySQL)

# MySQL Secondary (optional)
MYSQL_SECONDARY_HOST=localhost
MYSQL_SECONDARY_PORT=3306
# ... (sama seperti MySQL)
```

### Migration System
```bash
# CLI Commands (via cmd/migrate.go)
go run main.go make:migration create_table_name    # Buat migration baru
go run main.go migrate:all                         # Run semua migrations
go run main.go migrate:fresh                       # Drop all & re-migrate
go run main.go rollback:batch                      # Rollback batch terakhir
go run main.go rollback:all                        # Rollback semua migrations

# Migration Files Location
app/database/migrations/*.sql
```

**Format Migration File**: `YYYYMMDDHHMMSS_migration_name.sql`
- Menggunakan raw SQL (bukan GORM AutoMigrate)
- Tracking via `migrations` table di database

### Seeder System
```bash
# CLI Commands (via cmd/seeder.go)
go run main.go make:seeder --name=users           # Buat seeder baru
go run main.go db:seed                            # Run semua seeders
go run main.go rollback:seeder                    # Rollback seeders

# Seeder Files Location
app/database/seeds/*.go
```

**Format Seeder**: Go files with `Up()` and `Down()` methods

---

## 🔄 Flow Request-Response

### 1. Public Endpoints (No Auth)
```
HTTP Request
  → Gin Router (routes/web.go)
  → Controller Handler
  → Service Business Logic
  → GORM Model → Database
  ← JSON Response
```

**Contoh**: Login endpoint
```go
// routes/web.go
route.PUT("/auth/login", authController.Login)

// app/controllers/auth_controllers.go
func (ctrl *AuthController) Login(c *gin.Context) {
    var req requests.LoginRequest
    // Bind & validate request
    token, err := ctrl.Service.Login(req)
    // Return response
}

// app/services/auth_service.go
func (s *AuthService) Login(req requests.LoginRequest) (string, error) {
    // Business logic: verify credentials, generate token
    user := models.User{}
    facades.DB.Where("email = ?", req.Email).First(&user)
    // Generate JWT token
    return jwtService.GenerateToken(user)
}
```

### 2. Protected Endpoints (With Auth)
```
HTTP Request
  → middleware.AuthMiddleware() (verify JWT)
  → Controller Handler
  → Service Business Logic
  → GORM Model → Database
  ← JSON Response
```

**Middleware Chain**:
```go
// routes/web.go
userRoutes := route.Group("/users", middleware.AuthMiddleware())
{
    userRoutes.GET("", userController.List)
    // ...
}
```

**Auth Middleware** (app/middleware/auth_middleware.go):
- Extract token dari header `Authorization: Bearer <token>`
- Validate JWT token menggunakan `jwt_service.go`
- Set user info ke Gin Context
- Return 401 jika invalid

### 3. Response Format Standard
```go
// Success Response (via helpers/response_helper.go)
{
    "status": "success",
    "message": "Data retrieved successfully",
    "data": { ... },
    "meta": {
        "page": 1,
        "limit": 10,
        "total": 100
    }
}

// Error Response
{
    "status": "error",
    "message": "Validation failed",
    "errors": {
        "email": "Email is required",
        "password": "Password must be at least 6 characters"
    }
}
```

---

## 📊 Models & Relationships

### Core Models (app/models/)

#### 1. User Model (`user.go`)
```go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:255"`
    Email     string    `gorm:"size:255;unique"`
    Password  string    `gorm:"size:255"`
    CreatedAt time.Time
    UpdatedAt time.Time

    // Relationships
    Roles []Role `gorm:"many2many:user_has_roles"`
}
```

#### 2. Role Model (`role.go`)
```go
type Role struct {
    ID          uint   `gorm:"primaryKey"`
    Name        string `gorm:"size:255;unique"`
    Description string `gorm:"type:text"`

    // Relationships
    Users       []User       `gorm:"many2many:user_has_roles"`
    Permissions []Permission `gorm:"many2many:role_has_permissions"`
}
```

#### 3. Permission Model (`permission.go`)
```go
type Permission struct {
    ID          uint   `gorm:"primaryKey"`
    Name        string `gorm:"size:255;unique"`
    Description string `gorm:"type:text"`

    // Relationships
    Roles []Role `gorm:"many2many:role_has_permissions"`
}
```

#### 4. Category Model (`category.go`)
```go
type Category struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"size:255"`
    Description string    `gorm:"type:text"`
    CreatedAt   time.Time
    UpdatedAt   time.Time

    // Relationships
    Products []Product `gorm:"foreignKey:CategoryID"`
}
```

#### 5. Product Model (`product.go`)
```go
type Product struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"size:255"`
    Description string    `gorm:"type:text"`
    Price       float64   `gorm:"type:decimal(10,2)"`
    Stock       int       `gorm:"default:0"`
    CategoryID  uint      `gorm:"index"`
    CreatedAt   time.Time
    UpdatedAt   time.Time

    // Relationships
    Category Category `gorm:"foreignKey:CategoryID"`
}
```

#### 6. Test Postgres Model (`test_postgres.go`)
```go
type TestPostgres struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:255"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```
**Note**: Model ini khusus untuk testing PostgreSQL connection

### Pivot Tables (Many-to-Many)
- `user_has_roles` - User ↔ Role
- `role_has_permissions` - Role ↔ Permission

### GORM Scopes
**Location**: `app/models/scopes/pagination.go`

```go
// Usage
db.Scopes(scopes.Paginate(page, limit)).Find(&users)
```

---

## 🔐 Authentication & Authorization

### JWT Authentication Flow

#### 1. Login Process
```go
// User submits email + password
→ AuthController.Login()
  → AuthService.Login()
    → Verify credentials (hash password check via helpers/hash_helper.go)
    → Generate JWT token (via services/jwt_service.go)
    → Return token to client
```

#### 2. JWT Token Structure (`app/casts/jwt_claims.go`)
```go
type JWTClaims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
    jwt.StandardClaims
}
```

#### 3. Token Generation (`app/services/jwt_service.go`)
```go
// Generate token dengan expire time
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
signedToken, _ := token.SignedString([]byte(secretKey))
```

**Environment**:
```bash
JWT_SECRET_KEY=your_jwt_secret_key_here
JWT_EXPIRE_MINUTES=60  # Token expire dalam 60 menit
```

#### 4. Protected Route Access
```
Client sends: Authorization: Bearer <jwt_token>
  → middleware.AuthMiddleware()
    → Extract token from header
    → Validate token signature
    → Parse claims
    → Set user info to context
    → Pass to controller
```

#### 5. Authorization (Role-Based)
```go
// Check user roles
user := models.User{}
facades.DB.Preload("Roles").First(&user, userID)

// Check permissions
role := models.Role{}
facades.DB.Preload("Permissions").First(&role, roleID)
```

### Password Hashing
**Helper**: `app/helpers/hash_helper.go`
```go
// Hash password (bcrypt)
hashedPassword := helpers.HashPassword(password)

// Verify password
isValid := helpers.CheckPasswordHash(password, hashedPassword)
```

---

## 🛣️ API Endpoints Structure

### Routes Registration (`routes/web.go`)

#### Public Routes (No Auth Required)
```go
route.GET("", controller.HelloWorld)                # GET / - Hello World
route.PUT("/auth/login", authController.Login)      # PUT /auth/login - Login
```

#### Protected Routes (Auth Required)

**Auth Routes**:
```go
authRoutes := route.Group("/auth").Use(middleware.AuthMiddleware())
{
    authRoutes.GET("/logout", authController.Logout)     # GET /auth/logout
    authRoutes.GET("/refresh", authController.Refresh)   # GET /auth/refresh
}
```

**User Routes**:
```go
userRoutes := route.Group("/users", middleware.AuthMiddleware())
{
    userRoutes.GET("", userController.List)                    # GET /users
    userRoutes.GET("/:id", userController.Get)                 # GET /users/:id
    userRoutes.PUT("", userController.Put)                     # PUT /users (Create/Update)
    userRoutes.DELETE("/:id", userController.Delete)           # DELETE /users/:id
    userRoutes.POST("/:id/roles", userController.AssignRoles)  # POST /users/:id/roles
    userRoutes.GET("/:id/roles", userController.GetRoles)      # GET /users/:id/roles
}
```

**Role Routes**:
```go
roleRoutes := route.Group("/roles", middleware.AuthMiddleware())
{
    roleRoutes.GET("", roleController.List)                                 # GET /roles
    roleRoutes.PUT("", roleController.Put)                                  # PUT /roles
    roleRoutes.DELETE("/:id", roleController.Delete)                        # DELETE /roles/:id
    roleRoutes.POST("/:id/permissions", roleController.AssignPermissions)   # POST /roles/:id/permissions
    roleRoutes.GET("/:id/permissions", roleController.GetPermissions)       # GET /roles/:id/permissions
}
```

**Permission Routes**:
```go
permissionRoutes := route.Group("/permissions", middleware.AuthMiddleware())
{
    permissionRoutes.GET("", permissionController.List)          # GET /permissions
    permissionRoutes.PUT("", permissionController.Put)           # PUT /permissions
    permissionRoutes.DELETE("/:id", permissionController.Delete) # DELETE /permissions/:id
}
```

**Category Routes**:
```go
categoryRoutes := route.Group("/categories", middleware.AuthMiddleware())
{
    categoryRoutes.GET("/", categoryController.List)          # GET /categories
    categoryRoutes.GET("/:id", categoryController.Get)        # GET /categories/:id
    categoryRoutes.PUT("/", categoryController.Put)           # PUT /categories
    categoryRoutes.DELETE("/:id", categoryController.Delete)  # DELETE /categories/:id
}
```

**Product Routes**:
```go
productRoutes := route.Group("/products", middleware.AuthMiddleware())
{
    productRoutes.GET("/", productController.GetAll)        # GET /products
    productRoutes.GET("/:id", productController.GetByID)    # GET /products/:id
    productRoutes.PUT("/", productController.Put)           # PUT /products
    productRoutes.DELETE("/:id", productController.Delete)  # DELETE /products/:id
}
```

**File Routes** (Public):
```go
fileRoutes := route.Group("/file")
{
    fileRoutes.GET("/:key/:filename", fileController.ServeFile)  # GET /file/:key/:filename
}
```

**Database Health Check Routes** (Public):
```go
route.GET("/health", ...)                    # GET /health - Primary DB health
route.GET("/health/databases", ...)          # GET /health/databases - All DBs health

databaseRoutes := route.Group("/api/database")
{
    databaseRoutes.GET("/status", databaseController.GetConnectionStatus)
    databaseRoutes.GET("/health", databaseController.HealthCheck)
    databaseRoutes.GET("/test", databaseController.TestConnection)
}
```

**Test Routes (PostgreSQL)** - No Auth:
```go
testRoutes := route.Group("/tests")
{
    testRoutes.GET("", testController.List)           # GET /tests
    testRoutes.GET(":id", testController.Get)         # GET /tests/:id
    testRoutes.POST("", testController.Create)        # POST /tests
    testRoutes.PUT(":id", testController.Update)      # PUT /tests/:id
    testRoutes.DELETE(":id", testController.Delete)   # DELETE /tests/:id
}
```

### REST Convention
**PENTING**: Project ini menggunakan convention yang BERBEDA dari REST standard:
- ❌ **TIDAK** menggunakan `POST` untuk create
- ✅ **MENGGUNAKAN** `PUT` untuk create dan update
- ✅ `GET` untuk read/list
- ✅ `DELETE` untuk delete

**Logika Create/Update** (di controller):
```go
// PUT /users atau PUT /categories, dll.
func (ctrl *Controller) Put(c *gin.Context) {
    var req Request
    c.ShouldBindJSON(&req)

    if req.ID == 0 {
        // CREATE - jika ID tidak ada
        service.Create(req)
    } else {
        // UPDATE - jika ID ada
        service.Update(req)
    }
}
```

---

## 🛠️ Helper Functions

### File Locations: `app/helpers/*_helper.go`

#### 1. Response Helper (`response_helper.go`)
```go
// Success response
helpers.SuccessResponse(c, data, "Success message")
// Returns: {"status": "success", "message": "...", "data": {...}}

// Error response
helpers.ErrorResponse(c, statusCode, "Error message", errors)
// Returns: {"status": "error", "message": "...", "errors": {...}}

// Paginated response
helpers.PaginatedResponse(c, data, page, limit, total)
// Returns: {"status": "success", "data": [...], "meta": {...}}
```

#### 2. Hash Helper (`hash_helper.go`)
```go
// Hash password dengan bcrypt
hashedPassword := helpers.HashPassword(password)

// Verify password
isValid := helpers.CheckPasswordHash(password, hashedPassword)
```

#### 3. File Helper (`file_helper.go`)
```go
// Upload file
filePath, err := helpers.UploadFile(c, "file_field_name", "upload/path")

// Delete file
err := helpers.DeleteFile(filePath)

// Get file extension
ext := helpers.GetFileExtension(filename)
```

#### 4. Base64 File Helper (`base64file_helper.go`)
```go
// Decode base64 string to file
filePath, err := helpers.DecodeBase64ToFile(base64String, outputPath)

// Encode file to base64
base64String, err := helpers.EncodeFileToBase64(filePath)
```

#### 5. Environment Helper (`env_helper.go`)
```go
// Get environment variable dengan default value
value := helpers.GetEnv("KEY", "default_value")

// Get as int
intValue := helpers.GetEnvAsInt("PORT", 8080)

// Get as bool
boolValue := helpers.GetEnvAsBool("DEBUG", false)
```

#### 6. Reference Helper (`reference_helper.go`)
```go
// Generate unique reference code
refCode := helpers.GenerateReference("PREFIX")
// Returns: PREFIX-YYYYMMDD-XXXX

// Generate UUID
uuid := helpers.GenerateUUID()
```

#### 7. Path Helper (`path_helper.go`)
```go
// Get absolute path
absPath := helpers.GetAbsolutePath("relative/path")

// Check if file exists
exists := helpers.FileExists(filePath)

// Create directory if not exists
err := helpers.CreateDirIfNotExists(dirPath)
```

#### 8. URL Helper (`url_helper.go`)
```go
// Build URL dengan query params
url := helpers.BuildURL("http://example.com", map[string]string{
    "page": "1",
    "limit": "10",
})

// Parse query string
params := helpers.ParseQueryString(queryString)
```

#### 9. Error Helper (`error_helper.go`)
```go
// Check if error is specific type
isNotFound := helpers.IsNotFoundError(err)
isValidation := helpers.IsValidationError(err)

// Wrap error dengan context
err = helpers.WrapError(err, "additional context")
```

---

## 🧪 Testing Structure

### Test Files Location
```
app/
├── controllers/
│   ├── auth_controllers_test.go
│   └── controllers_test.go
├── helpers/
│   ├── hash_helper_test.go
│   ├── env_helper_test.go
│   ├── reference_helper_test.go
│   ├── response_helper_test.go
│   └── helpers_test.go
└── casts/
    ├── jwt_claims_test.go
    └── token_test.go
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests dengan coverage
go test -cover ./...

# Run tests di specific package
go test ./app/controllers
go test ./app/helpers
go test ./app/casts

# Verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Pattern
```go
package controllers

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFunctionName(t *testing.T) {
    // Arrange
    expected := "expected value"

    // Act
    result := FunctionToTest()

    // Assert
    assert.Equal(t, expected, result)
}
```

---

## ⚙️ Configuration System

### Environment Configuration (.env)
```bash
# Application
APP_NAME="Golang Starter Kit 2025"
APP_ENV=development              # development | production | staging
APP_SCHEME=http                  # http | https
APP_HOST=localhost
APP_PORT=9999                    # Default: 9999 (BUKAN 8080!)
APP_VERSION=1.0.0
APP_DESCRIPTION="Multi-Database Golang Starter Kit"

# Database - See Database Architecture section above
DB_CONNECTION=mysql              # mysql | postgres | sqlite
# ... (lihat section Database Architecture)

# JWT
JWT_SECRET_KEY=your_jwt_secret_key_here
JWT_EXPIRE_MINUTES=60

# File Upload
IMAGE_EXPIRE_MINUTES=2
```

### Config Files
- `config/database_config.go` - Database configuration dari env vars
- `config/mongo_config.go` - MongoDB configuration (optional)

### Bootstrap Process (`bootstrap/main.go`)
```go
func Init() {
    // 1. Load .env file
    godotenv.Load()

    // 2. Connect to database
    facades.ConnectDB()
    defer facades.CloseDB()

    // 3. Setup CLI commands (migration, seeder)
    app := &cli.App{...}

    // 4. If no CLI args, start web server
    if len(os.Args) == 1 {
        router := Router()
        router.Run(":" + helpers.GetEnv("APP_PORT", "8080"))
    }
}

func Router() *gin.Engine {
    route := gin.Default()

    // CORS middleware
    route.Use(cors.New(cors.Config{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"GET", "PUT", "DELETE"},
        AllowHeaders: []string{"*"},
    }))

    // Register routes
    routes.RegisterRoutes(route)

    // Setup Swagger docs
    docs.SwaggerInfo.Title = ...
    route.GET("/swagger/*any", ginSwagger.WrapHandler(...))

    return route
}
```

---

## 📝 Best Practices & Patterns

### 1. Controller Pattern
✅ **Thin Controllers** - Controller hanya handle HTTP request/response
```go
func (ctrl *Controller) Action(c *gin.Context) {
    // 1. Bind & validate request
    var req Request
    if err := c.ShouldBindJSON(&req); err != nil {
        helpers.ErrorResponse(c, 400, "Validation error", err)
        return
    }

    // 2. Call service
    result, err := ctrl.Service.DoSomething(req)
    if err != nil {
        helpers.ErrorResponse(c, 500, err.Error(), nil)
        return
    }

    // 3. Return response
    helpers.SuccessResponse(c, result, "Success")
}
```

### 2. Service Pattern
✅ **Fat Services** - Semua business logic ada di service layer
```go
type Service struct {
    // Dependencies (jika ada)
}

func (s *Service) DoSomething(req Request) (Result, error) {
    // Business logic here
    // Database operations here
    // Validations here

    return result, nil
}
```

### 3. Database Access
✅ **Via Facades**:
```go
// Default connection
facades.DB.Where(...).Find(&models)

// Specific connection
postgres, _ := facades.PostgreSQL()
postgres.DB.Create(&model)
```

❌ **JANGAN** create koneksi baru di setiap function
❌ **JANGAN** pass *gorm.DB sebagai parameter (kecuali untuk transaction)

### 4. Error Handling
✅ **Return errors ke controller**:
```go
func (s *Service) Create(data Data) error {
    if err := facades.DB.Create(&data).Error; err != nil {
        return fmt.Errorf("failed to create: %w", err)
    }
    return nil
}
```

✅ **Handle di controller**:
```go
if err := service.Create(data); err != nil {
    helpers.ErrorResponse(c, 500, err.Error(), nil)
    return
}
```

### 5. Request Validation
✅ **Gunakan struct tags**:
```go
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

### 6. Migration Best Practices
✅ **Gunakan raw SQL** (bukan AutoMigrate)
✅ **Timestamp format**: `YYYYMMDDHHMMSS_description.sql`
✅ **Reversible**: Setiap migration harus bisa di-rollback

### 7. Tidak Ada Repository Layer
❌ **TIDAK** ada folder `repositories/`
✅ Database operations langsung di **Services**

Kenapa? Karena dengan GORM, repository pattern menjadi over-engineering untuk kebanyakan kasus. Service sudah cukup sebagai abstraction layer.

---

## 🔍 Common Patterns dalam Codebase

### Pattern 1: Create or Update (Upsert)
```go
func (s *Service) Put(req Request) error {
    if req.ID == 0 {
        // Create new
        return facades.DB.Create(&model).Error
    } else {
        // Update existing
        return facades.DB.Model(&model).Where("id = ?", req.ID).Updates(req).Error
    }
}
```

### Pattern 2: List dengan Pagination
```go
func (s *Service) List(page, limit int) ([]Model, int64, error) {
    var models []Model
    var total int64

    db := facades.DB.Model(&Model{})
    db.Count(&total)

    err := db.Scopes(scopes.Paginate(page, limit)).Find(&models).Error

    return models, total, err
}
```

### Pattern 3: Delete Soft Delete
```go
// GORM akan otomatis soft delete jika model punya DeletedAt field
type Model struct {
    ID        uint
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Di service
func (s *Service) Delete(id uint) error {
    return facades.DB.Delete(&Model{}, id).Error  // Soft delete
}
```

### Pattern 4: Transaction
```go
func (s *Service) ComplexOperation() error {
    return facades.DB.Transaction(func(tx *gorm.DB) error {
        // Operation 1
        if err := tx.Create(&model1).Error; err != nil {
            return err  // Rollback
        }

        // Operation 2
        if err := tx.Create(&model2).Error; err != nil {
            return err  // Rollback
        }

        return nil  // Commit
    })
}
```

### Pattern 5: Preload Relationships
```go
// Single preload
facades.DB.Preload("Category").Find(&products)

// Multiple preloads
facades.DB.Preload("Roles").Preload("Roles.Permissions").Find(&users)

// Conditional preload
facades.DB.Preload("Products", "stock > ?", 0).Find(&categories)
```

---

## 🎯 Quick Reference untuk AI

### Saat User Minta Tambah Endpoint Baru:
1. **Buat Model** di `app/models/` (jika perlu)
2. **Buat Migration** via `go run main.go make:migration`
3. **Buat Request Validation** di `app/requests/`
4. **Buat Service** di `app/services/` (business logic)
5. **Buat Controller** di `app/controllers/` (HTTP handler)
6. **Register Route** di `routes/web.go`
7. **Update Swagger docs** dengan comment annotations

### Saat User Minta Optimasi:
1. **Cek Query Performance**: Tambah index di migration
2. **Cek N+1 Problem**: Gunakan `Preload()`
3. **Cek Connection Pool**: Sesuaikan MAX_OPEN_CONNS di .env
4. **Cek Middleware**: Pastikan tidak ada middleware yang berat
5. **Cek Response Size**: Gunakan pagination & limit fields

### Saat User Minta Fix Bug:
1. **Identifikasi Layer**: Controller, Service, Model, atau Helper?
2. **Cek Error Handling**: Pastikan error di-handle dengan baik
3. **Cek Validation**: Pastikan request validation benar
4. **Cek Database Query**: Test query di raw SQL dulu
5. **Tambah Test**: Buat unit test untuk prevent regression

### Files yang TIDAK BOLEH diubah sembarangan:
- `database/manager.go` - Core database manager
- `facades/database.go` - Database facade
- `bootstrap/main.go` - Application bootstrap
- `config/database_config.go` - Database configuration

### Files yang SERING diubah:
- `routes/web.go` - Tambah/update routes
- `app/controllers/*` - Tambah/update endpoints
- `app/services/*` - Tambah/update business logic
- `app/models/*` - Tambah/update database models
- `app/database/migrations/*` - Database schema changes

---

## 📌 Important Notes

### ⚠️ CRITICAL - Perbedaan dengan Standard REST:
1. **PUT digunakan untuk CREATE dan UPDATE** (bukan POST untuk create)
2. **Port default: 9999** (BUKAN 8080)
3. **TIDAK ADA Repository Layer** (logic di Services)

### ⚠️ Database Connection:
1. **Jangan create koneksi baru** - Gunakan facades
2. **Multi-database via Manager Pattern** - facades.GetConnection()
3. **Default connection** via DB_CONNECTION env var

### ⚠️ Authentication:
1. **JWT Token di header**: `Authorization: Bearer <token>`
2. **Token expire** diatur via JWT_EXPIRE_MINUTES
3. **Password hashing** menggunakan bcrypt (via hash_helper.go)

### ⚠️ Migration vs GORM AutoMigrate:
1. **Production**: Gunakan SQL migrations (app/database/migrations/)
2. **Development**: Boleh gunakan AutoMigrate (tapi TIDAK recommended)
3. **Reason**: SQL migrations lebih terkontrol dan reversible

---

## 🚀 Version Information

- **Go Version**: 1.21+
- **Gin Framework**: Latest
- **GORM**: v2
- **Database**: MySQL 8.0+, PostgreSQL 13+
- **JWT**: golang-jwt/jwt v4

---

**Last Updated**: 2025-11-01
**Maintained by**: AI Assistant for Golang Starter Kit 2025

---

_Dokumentasi ini akan terus di-update seiring perkembangan codebase._
