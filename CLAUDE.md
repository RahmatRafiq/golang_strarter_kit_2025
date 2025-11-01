# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Golang Starter Kit 2025** is a production-ready REST API backend starter template built with Go 1.24.3. The project implements clean architecture principles with a layered approach (Controller → Service → Repository → Database), supports multiple database connections simultaneously (MySQL, PostgreSQL, SQL Server, SQLite), and includes comprehensive authentication, RBAC, file management, and database migration/seeding systems.

**Key Statistics:**
- 93 Go files
- 6 repository interfaces
- 12 controllers
- 10 services
- 9 GORM models
- Multi-database support with connection pooling

## Tech Stack

| Category | Technology |
|----------|------------|
| **Framework** | Gin v1.10.0 (HTTP web framework) |
| **ORM** | GORM v1.26.1 with multi-driver support |
| **Databases** | MySQL, PostgreSQL, SQLite, SQL Server, MongoDB (optional) |
| **Authentication** | JWT (golang-jwt/jwt v5.2.2) |
| **Password Hashing** | Argon2id (primary), Bcrypt (legacy support) |
| **Validation** | go-playground/validator v10.26.0 |
| **API Documentation** | Swaggo (Swagger/OpenAPI) |
| **Testing** | Ginkgo v2.20.2 & Gomega v1.34.2 (BDD framework) |
| **CLI** | urfave/cli v2.27.6 |
| **CORS** | gin-contrib/cors v1.7.5 |
| **Development** | Air (hot reload - external tool) |

## Development Commands

### Running the Application

```bash
# Development with hot reload (recommended for active development)
air
# Watches file changes and automatically recompiles/restarts server
# Configure in .air.toml if needed

# Standard development run
go run main.go
# Default server starts on port 9999

# Production build
make build
# Creates ./main executable

# Run compiled binary
./main
```

**Important URLs:**
- API Server: `http://localhost:9999`
- Swagger UI: `http://localhost:9999/swagger/index.html`
- Health Check: `http://localhost:9999/health`
- Multi-DB Health: `http://localhost:9999/health/databases`

### Testing

```bash
# Run all tests with race detector (recommended)
ginkgo -r --race

# Run all tests (standard Go test)
go test ./...

# Test with coverage report
go test -cover ./...

# Test specific package
go test ./app/controllers
go test ./app/services
go test ./app/repositories

# Run specific test file
go test ./app/controllers/auth_controllers_test.go
```

### API Documentation (Swagger)

```bash
# Generate/update Swagger documentation
swag init
# Scans all @-annotations in controllers and generates:
# - docs/docs.go
# - docs/swagger.json
# - docs/swagger.yaml

# ALWAYS run this after:
# 1. Adding new endpoints
# 2. Modifying request/response structs
# 3. Changing API documentation comments
```

**Swagger Annotations Example:**
```go
// @Summary Get user by ID
// @Description Get detailed user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} responses.UserResponse
// @Failure 404 {object} helpers.ResponseParams
// @Router /users/{id} [get]
// @Security Bearer
```

### Database Management

#### Migration Commands

```bash
# Create new migration file
go run main.go make:migration create_products_table
# Creates: app/database/migrations/YYYYMMDDHHMMSS_create_products_table.sql
# File contains UP and DOWN sections

# Run all pending migrations (creates new batch)
go run main.go migrate:all
go run main.go migrate:all --connection=postgres

# Run specific migration file
go run main.go migrate --file=20250426184415_create_roles_table.sql
go run main.go migrate --file=20250426184415_create_roles_table.sql --connection=mysql_secondary

# View migration status
go run main.go migrate:status
go run main.go migrate:status --connection=postgres
# Shows: ✅ Ran (with batch number) | ⏳ Pending

# Rollback last migration batch
go run main.go rollback:batch
go run main.go rollback:batch --connection=postgres

# Rollback N last batches (like Laravel)
go run main.go rollback:batch --step=3

# Rollback specific batch number
go run main.go rollback:batch --batch=2

# Rollback ALL migrations
go run main.go rollback:all

# Rollback specific migration file
go run main.go rollback --file=20250426184415_create_roles_table.sql

# Reset database (rollback all + migrate all)
go run main.go migrate:reset

# Fresh migration (DROP all tables + migrate all)
go run main.go migrate:fresh
# ⚠️ DESTRUCTIVE: Deletes ALL data

# Wipe database completely (drop all tables)
go run main.go db:wipe
# ⚠️ DESTRUCTIVE: Nuclear option
```

**Migration File Structure:**
```sql
-- +++ UP Migration
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    reference VARCHAR(100) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- --- DOWN Migration
DROP TABLE IF EXISTS users;
```

**Migration System Features:**
- Tracks migrations in `migrations` table with batch numbers
- Supports database-specific syntax (MySQL `AUTO_INCREMENT` vs PostgreSQL `SERIAL`)
- Batch system allows granular rollbacks
- Files named: `YYYYMMDDHHMMSS_description.sql`
- Stored in: `app/database/migrations/`
- Managed by: `app/database/migration_manager.go`

#### Seeder Commands

```bash
# Create new seeder
go run main.go make:seeder --name=ProductsSeeder
# Creates: app/database/seeds/YYYYMMDDHHMMSS_ProductsSeeder.go
# With SeedProductsSeeder() and RollbackProductsSeeder() functions

# Run all seeders
go run main.go db:seed
go run main.go db:seed --connection=postgres

# Run specific seeder class
go run main.go db:seed --class=UserSeeder
go run main.go db:seed --class=UserSeeder --connection=postgres

# Rollback last seeder batch
go run main.go rollback:seeder
go run main.go rollback:seeder --connection=postgres

# Rollback specific seeder batch
go run main.go rollback:seeder --batch=2
```

**Seeder File Structure:**
```go
package seeds

import (
    "golang_starter_kit_2025/app/helpers"
    "golang_starter_kit_2025/app/models"
    "log"
    "time"
    "gorm.io/gorm"
)

// Naming: Function name MUST match filename
// File: 20250423230248_UserSeeder.go
// Functions: SeedUserSeeder, RollbackUserSeeder

func SeedUserSeeder(db *gorm.DB) error {
    log.Println("🌱 Seeding UserSeeder...")

    data := models.User{
        Reference: helpers.GenerateReference("USR"),
        Username:  "admin",
        Email:     "admin@example.com",
        Password:  "admin123", // Will be hashed by BeforeCreate hook
        Pin:       "123456",   // Will be hashed by BeforeCreate hook
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    return db.Create(&data).Error
}

func RollbackUserSeeder(db *gorm.DB) error {
    log.Println("🗑️ Rolling back UserSeeder...")

    // Unscoped() performs hard delete (not soft delete)
    return db.Unscoped().
        Where("username = ?", "admin").
        Delete(&models.User{}).
        Error
}
```

**Seeder System Features:**
- Batch tracking system (like migrations)
- Function names MUST match file names
- Supports rollback for data cleanup
- Use `Unscoped()` for hard deletes in rollbacks
- Managed by: `app/database/seeder_manager.go`

**Connection Flag:**
All migration and seeder commands support `--connection` flag:
- `--connection=mysql` (default)
- `--connection=postgres`
- `--connection=mysql_secondary`

## Architecture Deep Dive

### Project Structure (93 Go Files)

```
golang_starter_kit_2025/
├── main.go                          # Entry point, delegates to bootstrap
├── go.mod & go.sum                  # Dependencies
├── Makefile                         # Build shortcuts
├── .env.example                     # Environment template
│
├── bootstrap/
│   └── main.go                      # App initialization, router setup, CLI handler
│
├── cmd/                             # CLI command definitions
│   ├── migrate.go                   # 10 migration commands
│   └── seeder.go                    # 3 seeder commands
│
├── config/                          # Configuration loaders
│   ├── database_config.go           # Multi-DB config from env vars
│   └── mongo_config.go              # MongoDB config (optional)
│
├── database/                        # Core database management
│   └── manager.go                   # Connection pool manager (singleton)
│
├── facades/                         # Simplified database access
│   ├── database.go                  # DB facade with helper methods
│   └── mongo.go                     # MongoDB facade (optional)
│
├── routes/
│   └── web.go                       # All route definitions + DI wiring
│
├── app/
│   ├── controllers/                 # HTTP request handlers (12 files)
│   │   ├── auth_controllers.go      # Login, Logout, Refresh
│   │   ├── user_controller.go       # User CRUD + role assignment
│   │   ├── role_controller.go       # Role CRUD + permission assignment
│   │   ├── permission_controller.go # Permission CRUD
│   │   ├── product_controller.go    # Product CRUD with categories
│   │   ├── category_controller.go   # Category CRUD
│   │   ├── file_controller.go       # File upload/download/serve
│   │   ├── database_controller.go   # DB health checks & stats
│   │   └── test_postgres_controller.go # Multi-DB testing
│   │
│   ├── services/                    # Business logic layer (10 files)
│   │   ├── auth_service.go          # Login/logout/refresh logic
│   │   ├── jwt_service.go           # Token generation & validation
│   │   ├── user_services.go         # User business operations
│   │   ├── role_service.go          # Role + permission assignment
│   │   ├── permission_service.go    # Permission operations
│   │   ├── product_service.go       # Product inventory logic
│   │   ├── category_service.go      # Category management
│   │   ├── file_service.go          # File processing logic
│   │   └── database_service.go      # Cross-database operations
│   │
│   ├── repositories/                # Data access layer
│   │   ├── interfaces/              # Repository contracts
│   │   │   ├── user_repository_interface.go
│   │   │   ├── role_repository_interface.go
│   │   │   ├── permission_repository_interface.go
│   │   │   ├── product_repository_interface.go
│   │   │   ├── category_repository_interface.go
│   │   │   └── auth_repository_interface.go
│   │   ├── user_repository.go       # User DB operations
│   │   ├── role_repository.go       # Role DB operations
│   │   ├── permission_repository.go # Permission DB operations
│   │   ├── product_repository.go    # Product DB operations
│   │   ├── category_repository.go   # Category DB operations
│   │   └── auth_repository.go       # Auth-specific queries
│   │
│   ├── models/                      # GORM entities (9 models)
│   │   ├── user.go                  # User + BeforeCreate hook (Argon2 hashing)
│   │   ├── role.go                  # Role with many-to-many relations
│   │   ├── permission.go            # Permission model
│   │   ├── role_has_permission.go   # Pivot table model
│   │   ├── user_has_role.go         # Pivot table model
│   │   ├── product.go               # Product + GORM hooks (AfterFind, AfterCreate, AfterUpdate)
│   │   ├── category.go              # Category with product relation
│   │   ├── test_postgres.go         # Test model for multi-DB
│   │   └── scopes/
│   │       └── pagination.go        # Reusable pagination scopes
│   │
│   ├── requests/                    # Input validation structs (11 files)
│   │   ├── login_request.go         # Email + password validation
│   │   ├── filter_request.go        # Pagination params (page, limit, offset)
│   │   ├── category_request.go      # Category creation/update
│   │   ├── product_request.go       # Product validation rules
│   │   ├── permission_request.go    # Permission validation
│   │   ├── role_request.go          # Role validation
│   │   └── ...                      # Other domain-specific requests
│   │
│   ├── responses/                   # Output formatting structs
│   │
│   ├── middleware/                  # HTTP middleware
│   │   ├── auth_middleware.go       # JWT validation + SKIP_AUTH support
│   │   └── logger_middleware.go     # Request/response logging
│   │
│   ├── helpers/                     # Utility functions (17 files)
│   │   ├── hash_helper.go           # Argon2id & Bcrypt password hashing
│   │   ├── response_helper.go       # Standardized API responses
│   │   ├── file_helper.go           # File operations & validation
│   │   ├── base64file_helper.go     # Base64 encoding/decoding
│   │   ├── reference_helper.go      # Generate unique codes (USR-20250101-ABCD)
│   │   ├── url_helper.go            # URL generation & file URLs
│   │   ├── env_helper.go            # Environment variable helpers
│   │   ├── error_helper.go          # Error formatting
│   │   ├── path_helper.go           # Path utilities
│   │   └── security_warning.go      # Dev mode security warnings
│   │
│   ├── handlers/
│   │   └── response_handler.go      # Response formatting logic
│   │
│   ├── casts/                       # Data transformation objects
│   │   ├── jwt_claims.go            # JWT claims struct (UserID, ExpiredAt)
│   │   └── token.go                 # Token response struct
│   │
│   └── database/                    # Database management logic
│       ├── migration_manager.go     # Migration system implementation
│       ├── seeder_manager.go        # Seeder system implementation
│       ├── migrations/              # SQL migration files
│       │   ├── 20250426184415_create_roles_table.sql
│       │   ├── 20250426184424_create_permissions_table.sql
│       │   ├── 20250426184432_create_users_table.sql
│       │   └── ...
│       └── seeds/                   # Go seeder files
│           └── 20250423230248_UserSeeder.go
│
├── storage/                         # File uploads storage
├── documentation/                   # Project documentation
│   ├── DATABASE.md                  # Migration & seeder guide
│   ├── MULTI_DATABASE.md            # Multi-DB configuration
│   ├── API_REFERENCE.md             # API endpoint documentation
│   └── GETTING_STARTED.md           # Setup guide
│
└── examples/
    └── multi_database_usage.go      # Multi-DB usage examples
```

### Layered Architecture Pattern

**Request Flow:**
```
HTTP Request
    ↓
[Middleware Layer]
    ├─ CORS (gin-contrib/cors)
    ├─ Logger (custom)
    └─ Auth (JWT validation or SKIP_AUTH)
    ↓
[Controller Layer] (app/controllers/)
    ├─ Parse request (JSON binding)
    ├─ Validate input (go-playground/validator)
    ├─ Call service method
    ├─ Format response (helpers.SuccessResponse / ErrorResponse)
    └─ Return HTTP response
    ↓
[Service Layer] (app/services/)
    ├─ Business logic
    ├─ Data transformation
    ├─ Call repository methods
    ├─ Aggregate data from multiple repositories
    └─ Return domain objects or errors
    ↓
[Repository Layer] (app/repositories/)
    ├─ Database queries (GORM)
    ├─ Data mapping
    ├─ Error handling
    └─ Return models or errors
    ↓
[Model Layer] (app/models/)
    ├─ GORM models with tags
    ├─ Relationships (BelongsTo, HasMany, ManyToMany)
    ├─ GORM hooks (BeforeCreate, AfterFind, etc.)
    └─ Database representation
    ↓
[Database]
    └─ MySQL / PostgreSQL / SQLite / SQL Server
```

**Dependency Flow (Dependency Injection):**
```
routes/web.go (Wiring Layer)
    ↓
1. Create Repositories
   userRepo := repositories.NewUserRepository(facades.DB)
   roleRepo := repositories.NewRoleRepository(facades.DB)
   permissionRepo := repositories.NewPermissionRepository(facades.DB)
    ↓
2. Create Services (inject repositories)
   userService := services.NewUserService(userRepo)
   roleService := services.NewRoleService(roleRepo, permissionRepo)
    ↓
3. Create Controllers (inject services)
   userController := controllers.NewUserController(*userService)
   roleController := controllers.NewRoleController(*roleService)
    ↓
4. Register Routes
   userRoutes.GET("/:id", userController.Get)
   userRoutes.PUT("", userController.Put)
```

**Why This Architecture?**
- **Separation of Concerns**: Each layer has single responsibility
- **Testability**: Layers can be tested independently with mocks
- **Maintainability**: Changes in one layer don't cascade
- **Scalability**: Easy to add new features following same pattern
- **Repository Interfaces**: Allow swapping implementations (e.g., different DB, caching layer)

### Multi-Database Manager System

**Core Components:**
1. **config/database_config.go**: Loads configurations from environment variables
2. **database/manager.go**: Singleton manager handling multiple connections
3. **facades/database.go**: Simplified access layer with helper functions

**Manager Features:**
- Singleton pattern (thread-safe with sync.Once)
- Connection pooling per database
- Health monitoring (Ping)
- Statistics tracking (sql.DBStats)
- Lazy connection (connects on first use)
- Graceful shutdown

**Supported Database Types:**
```go
const (
    MySQL      DatabaseType = "mysql"
    PostgreSQL DatabaseType = "postgres"
    SQLite     DatabaseType = "sqlite"
    SQLServer  DatabaseType = "sqlserver"
)
```

**Configured Connections:**
1. **mysql** - Primary MySQL connection (default)
2. **postgres** - PostgreSQL connection
3. **mysql_secondary** - Secondary MySQL instance

**Database Access Patterns:**

```go
// 1. Default database (backward compatibility)
db := facades.DB
users := []models.User{}
db.Find(&users)

// 2. Get specific connection via facade helpers
mysqlConn, _ := facades.MySQL()
postgresConn, _ := facades.PostgreSQL()
secondaryConn, _ := facades.MySQLSecondary()

// 3. Get connection via manager
manager := facades.GetManager()
conn, _ := manager.GetConnection("postgres")
conn.DB.Find(&users)

// 4. Use in repositories (injected connection)
type userRepository struct {
    db *gorm.DB // Can be any connection
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepositoryInterface {
    return &userRepository{db: db}
}

// In routes/web.go, choose which DB to use:
userRepo := repositories.NewUserRepository(facades.DB) // Default MySQL
// OR
postgresConn, _ := facades.PostgreSQL()
userRepo := repositories.NewUserRepository(postgresConn.DB) // PostgreSQL
```

**Connection Pool Configuration:**
Each connection has independent pool settings (via env vars):
- `MAX_IDLE_CONNS`: Maximum idle connections (default: 10)
- `MAX_OPEN_CONNS`: Maximum open connections (default: 200)
- `CONN_MAX_LIFETIME`: Connection max lifetime (default: 15m)
- `CONN_MAX_IDLE_TIME`: Idle connection max time (default: 5m)

**Example Multi-DB Usage:**
```go
// Sync data from MySQL to PostgreSQL
mysqlConn, _ := facades.MySQL()
postgresConn, _ := facades.PostgreSQL()

var users []models.User
mysqlConn.DB.Find(&users) // Read from MySQL

for _, user := range users {
    postgresConn.DB.Create(&user) // Write to PostgreSQL
}
```

**Health Check Endpoint:**
```bash
GET /health/databases
```
Returns connection status and statistics for all configured databases.

## Configuration

### Environment Variables (.env)

**Required Variables:**
```bash
# Application
APP_NAME="Golang Starter Kit 2025"
APP_ENV=development                    # development | production | test
APP_SCHEME=http                        # http | https
APP_HOST=localhost
APP_PORT=9999
APP_VERSION=1.0.0

# Default Database Connection
DB_CONNECTION=mysql                    # mysql | postgres | sqlite | sqlserver

# JWT Configuration
JWT_SECRET_KEY=<STRONG_RANDOM_KEY>     # Generate: openssl rand -base64 48
JWT_EXPIRE_MINUTES=60                  # Token expiry in minutes
# ⚠️ CRITICAL: Must not be placeholder value!
```

**MySQL Configuration (Primary):**
```bash
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=root
MYSQL_PASSWORD=secure_password
MYSQL_CHARSET=utf8mb4
MYSQL_TIMEZONE=Local                   # Local | UTC | Asia/Jakarta
MYSQL_MAX_IDLE_CONNS=10
MYSQL_MAX_OPEN_CONNS=200
MYSQL_CONN_MAX_LIFETIME=15m
MYSQL_CONN_MAX_IDLE_TIME=5m
```

**PostgreSQL Configuration:**
```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=golang_starter_kit_2025_pg
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secure_password
POSTGRES_SSLMODE=disable               # disable | require | verify-full
POSTGRES_TIMEZONE=UTC
POSTGRES_MAX_IDLE_CONNS=10
POSTGRES_MAX_OPEN_CONNS=200
POSTGRES_CONN_MAX_LIFETIME=15m
POSTGRES_CONN_MAX_IDLE_TIME=5m
```

**MySQL Secondary Configuration (Optional):**
```bash
MYSQL_SECONDARY_HOST=localhost
MYSQL_SECONDARY_PORT=3307
MYSQL_SECONDARY_DB=golang_starter_kit_2025_secondary
MYSQL_SECONDARY_USER=root
MYSQL_SECONDARY_PASSWORD=secure_password
MYSQL_SECONDARY_CHARSET=utf8mb4
MYSQL_SECONDARY_TIMEZONE=Local
# ... same pool config vars as primary
```

**Development/Security:**
```bash
# Skip Authentication (DEVELOPMENT ONLY!)
SKIP_AUTH=true                         # Bypasses JWT auth, sets user_id=1
# ⚠️ NEVER use in production!

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
# Comma-separated list of allowed origins
```

**Optional (MongoDB):**
```bash
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DB=golang_starter_kit_2025
MONGO_USER=root
MONGO_PASS=secure_password
MONGO_COLL="${APP_ENV}"
```

### Application Bootstrap Process

**Startup Flow** (`bootstrap/main.go`):
```go
func Init() {
    // 1. Load .env file
    godotenv.Load()

    // 2. Validate required env vars
    validateRequiredEnvVars()
    // Checks: DB_CONNECTION, JWT_SECRET_KEY, APP_PORT
    // Fails if JWT_SECRET_KEY is placeholder value

    // 3. Set Gin mode
    if APP_ENV == "production" {
        gin.SetMode(gin.ReleaseMode)
    } else if APP_ENV == "development" {
        gin.SetMode(gin.DebugMode)
    }

    // 4. Initialize database
    facades.ConnectDB()
    defer facades.CloseDB()
    // Connects to default database defined in DB_CONNECTION

    // 5. Setup CLI app
    app := &cli.App{
        Commands: []*cli.Command{
            cmd.MakeMigrationCommand,
            cmd.MigrateAllCommand,
            // ... all 13 CLI commands
        },
    }

    // 6. Check if CLI command invoked
    if len(os.Args) > 1 {
        app.Run(os.Args)
        return
    }

    // 7. Check security warnings (if development mode)
    if !helpers.CheckSecurityWarning() {
        os.Exit(0)
    }

    // 8. Setup HTTP router
    r := Router() // Configures CORS, routes, Swagger

    // 9. Start server
    r.Run(":" + APP_PORT)
}
```

**Router Setup** (`bootstrap/main.go:Router()`):
```go
func Router() *gin.Engine {
    route := gin.Default()

    // CORS middleware
    allowedOrigins := helpers.GetEnv("ALLOWED_ORIGINS", "http://localhost:3000")
    route.Use(cors.New(cors.Config{
        AllowOrigins:     strings.Split(allowedOrigins, ","),
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Api-Key", "X-Request-ID"},
        ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Register all routes with dependency injection
    routes.RegisterRoutes(route)

    // Swagger configuration
    docs.SwaggerInfo.Title = APP_NAME + " API"
    docs.SwaggerInfo.Description = APP_DESCRIPTION
    docs.SwaggerInfo.Version = APP_VERSION
    docs.SwaggerInfo.Host = APP_HOST + ":" + APP_PORT
    docs.SwaggerInfo.Schemes = []string{APP_SCHEME}

    route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return route
}
```

## API Structure & Endpoints

### Route Registration (`routes/web.go`)

All routes are defined in `routes/web.go` with dependency injection:

```go
func RegisterRoutes(route *gin.Engine) {
    // 1. Create repositories (data access layer)
    userRepo := repositories.NewUserRepository(facades.DB)
    roleRepo := repositories.NewRoleRepository(facades.DB)
    permissionRepo := repositories.NewPermissionRepository(facades.DB)
    productRepo := repositories.NewProductRepository(facades.DB)
    categoryRepo := repositories.NewCategoryRepository(facades.DB)
    authRepo := repositories.NewAuthRepository(facades.DB)

    // 2. Create services (business logic layer)
    userService := services.NewUserService(userRepo)
    roleService := services.NewRoleService(roleRepo, permissionRepo)
    permissionService := services.NewPermissionService(permissionRepo)
    productService := services.NewProductService(productRepo, categoryRepo)
    categoryService := services.NewCategoryService(categoryRepo, productRepo)
    authService := services.NewAuthService(authRepo)

    // 3. Create controllers (HTTP handlers)
    userController := controllers.NewUserController(*userService)
    roleController := controllers.NewRoleController(*roleService)
    // ... etc

    // 4. Define routes

    // Public routes (no auth)
    route.GET("", controller.HelloWorld)
    route.PUT("/auth/login", authController.Login)
    route.GET("/health", healthCheckHandler)
    route.GET("/health/databases", multiDBHealthCheckHandler)

    // Protected routes (require JWT)
    authRoutes := route.Group("/auth").Use(middleware.AuthMiddleware())
    {
        authRoutes.GET("/logout", authController.Logout)
        authRoutes.GET("/refresh", authController.Refresh)
    }

    userRoutes := route.Group("/users", middleware.AuthMiddleware())
    {
        userRoutes.GET("", userController.List)
        userRoutes.GET("/:id", userController.Get)
        userRoutes.PUT("", userController.Put)
        userRoutes.DELETE("/:id", userController.Delete)
        userRoutes.POST("/:id/roles", userController.AssignRoles)
        userRoutes.GET("/:id/roles", userController.GetRoles)
    }

    // ... similar groups for roles, permissions, products, categories, files
}
```

### Authentication System (JWT)

**Components:**
- `app/services/auth_service.go`: Login/logout/refresh logic
- `app/services/jwt_service.go`: Token generation & validation
- `app/middleware/auth_middleware.go`: HTTP middleware
- `app/casts/jwt_claims.go`: JWT claims structure
- `app/casts/token.go`: Token response structure

**JWT Claims Structure:**
```go
type JwtClaims struct {
    UserID    uint  `json:"user_id"`
    ExpiredAt int64 `json:"expired_at"` // Unix timestamp
}

func NewJwtClaims(userID uint, expiredAt int64) jwt.MapClaims {
    return jwt.MapClaims{
        "user_id":    userID,
        "expired_at": expiredAt,
    }
}

func ParseJwtClaims(claims jwt.MapClaims) JwtClaims {
    return JwtClaims{
        UserID:    uint(claims["user_id"].(float64)),
        ExpiredAt: int64(claims["expired_at"].(float64)),
    }
}
```

**Authentication Flow:**

1. **Login** (`PUT /auth/login`):
```go
// Request
{
    "email": "admin@example.com",
    "password": "admin123"
}

// Process (auth_service.go:Login)
1. Find user by email (auth_repository)
2. Compare password using Argon2id (helpers.ComparePasswordArgon2)
3. Generate JWT token with expiry (jwt_service.GenerateToken)
4. Save token to user.jwt_token field
5. Return token + expiry time

// Response
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expired_at": "2025-01-01T12:00:00Z"
}
```

2. **Protected Endpoints** (with `middleware.AuthMiddleware()`):
```go
// Request Headers
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

// Middleware Process (auth_middleware.go)
1. Check SKIP_AUTH env var
   - If true: Set user_id=1, skip validation (DEV ONLY)
2. Extract Authorization header
   - Fail if missing: {"reference": "ERROR-1", "message": "Membutuhkan token"}
3. Verify Bearer prefix
   - Fail if invalid: {"reference": "ERROR-2", "message": "Token tidak valid"}
4. Validate token signature (jwt_service.ValidateToken)
   - Fail if invalid: {"reference": "ERROR-3", "message": "Token tidak valid"}
5. Parse claims (casts.ParseJwtClaims)
6. Check token expiry
   - Fail if expired: {"reference": "ERROR-4", "message": "Token sudah kadaluarsa"}
7. Set context variables:
   c.Set("token", tokenString)
   c.Set("user_id", claims.UserID)
8. Call c.Next() to proceed to controller

// Controller Access
func (ctrl *UserController) Get(c *gin.Context) {
    userID := c.Get("user_id").(uint)
    // Use userID for authorization checks
}
```

3. **Logout** (`GET /auth/logout`):
```go
// Process (auth_service.go:Logout)
1. Extract token from request
2. Parse token to get user_id
3. Find user by ID
4. Clear user.jwt_token field (set to "")
5. Save user
```

4. **Refresh Token** (`GET /auth/refresh`):
```go
// Process (auth_service.go:RefreshToken)
1. Validate current token
2. Extract user_id from claims
3. Generate NEW token with fresh expiry
4. Update user.jwt_token
5. Return new token + expiry
```

**JWT Configuration:**
```bash
JWT_SECRET_KEY=<strong-random-key>     # Used for signing tokens
JWT_EXPIRE_MINUTES=60                  # Token validity duration
```

**Security Features:**
- Tokens signed with HMAC-SHA256
- Expiry validation on every request
- Token stored in database (enables single logout)
- Constant-time comparison for security
- Reference codes for error tracking

**Development Bypass:**
```bash
SKIP_AUTH=true
```
- Bypasses ALL authentication checks
- Automatically sets `user_id=1` in context
- ⚠️ **NEVER use in production!**

### Request Validation

**Validation Structs** (`app/requests/`):
```go
package requests

type ProductRequest struct {
    Name        string   `json:"name" binding:"required"`
    Description string   `json:"description" binding:"required"`
    Price       float64  `json:"price" binding:"required,min=0"`
    CategoryID  uint     `json:"category_id" binding:"required"`
    Stock       int      `json:"stock" binding:"required,min=0"`
    Images      []string `json:"images"`
}

type FilterRequest struct {
    Page   *int `json:"page" form:"page"`
    Limit  *int `json:"limit" form:"limit"`
    Offset *int `json:"offset" form:"offset"`
}
```

**Validation Tags** (go-playground/validator):
- `required`: Field must be present
- `min=N`: Minimum value for numbers
- `max=N`: Maximum value for numbers
- `email`: Valid email format
- `len=N`: Exact length
- `oneof=val1 val2`: Must be one of values

**Controller Usage:**
```go
func (ctrl *ProductController) Put(c *gin.Context) {
    var request requests.ProductRequest

    // Bind JSON and validate
    if err := c.ShouldBindJSON(&request); err != nil {
        helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    // Validation passed, proceed with business logic
    product := models.Product{
        Name:        request.Name,
        Description: request.Description,
        Price:       request.Price,
        // ...
    }

    if err := ctrl.service.Create(&product); err != nil {
        helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    helpers.SuccessResponse(c, "Product created", product)
}
```

### Response Format

**Standardized Responses** (`app/helpers/response_helper.go`):

```go
// Success Response
helpers.SuccessResponse(c, "Operation successful", data)
// Returns:
{
    "message": "Operation successful",
    "data": { ... }
}

// Error Response (simple)
helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid input")
// Returns:
{
    "error": "Invalid input"
}

// Error Response (with reference code)
helpers.ResponseError(c, &helpers.ResponseParams[any]{
    Reference: "ERROR-USER-001",
    Message:   "User not found",
}, http.StatusNotFound)
// Returns:
{
    "reference": "ERROR-USER-001",
    "message": "User not found"
}
```

**ResponseParams Structure:**
```go
type ResponseParams[T any] struct {
    Reference string `json:"reference,omitempty"`
    Message   string `json:"message"`
    Data      T      `json:"data,omitempty"`
}
```

**Reference Code Convention:**
- `ERROR-1`, `ERROR-2`, etc.: Authentication errors
- `ERROR-USER-XXX`: User-related errors
- `ERROR-PRODUCT-XXX`: Product-related errors
- Helps with debugging and logging

### Pagination

**Pagination Scopes** (`app/models/scopes/pagination.go`):

```go
// Option 1: Using FilterRequest
func (r *userRepository) List(page, limit int) ([]models.User, int64, error) {
    var users []models.User
    var total int64

    filter := requests.FilterRequest{
        Page:  &page,
        Limit: &limit,
    }

    // Count total
    r.db.Model(&models.User{}).Count(&total)

    // Get paginated results
    err := r.db.Scopes(scopes.Paginate(filter)).Find(&users).Error

    return users, total, err
}

// Option 2: Simple pagination
err := db.Scopes(scopes.PaginateSimple(page, limit)).Find(&users).Error
```

**Pagination Logic:**
```go
func Paginate(filter FilterRequest) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        page := 1
        if filter.Page != nil {
            page = *filter.Page
        }

        limit := 10 // Default
        if filter.Limit != nil {
            limit = *filter.Limit
        }

        // Constraints
        if limit > 100 {
            limit = 100 // Max limit
        } else if limit <= 0 {
            limit = 10 // Min limit
        }

        offset := 0
        if filter.Offset != nil {
            offset = *filter.Offset
        } else {
            offset = (page - 1) * limit
        }

        return db.Offset(offset).Limit(limit)
    }
}
```

**Pagination Defaults:**
- Default page: 1
- Default limit: 10
- Maximum limit: 100
- Minimum limit: 10

**Client Usage:**
```bash
GET /users?page=2&limit=20
GET /products?page=1&limit=50
GET /products?offset=100&limit=25  # Manual offset
```

## Models & Database Relationships

### GORM Models

**User Model** (`app/models/user.go`):
```go
type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Reference string         `gorm:"type:varchar(100);uniqueIndex" json:"reference"`
    Username  string         `gorm:"type:varchar(100);uniqueIndex" json:"username"`
    Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
    Password  string         `gorm:"type:varchar(255)" json:"password"`
    JwtToken  string         `gorm:"type:varchar(255)" json:"jwt_token" swaggerignore:"true"`
    FcmToken  string         `gorm:"type:varchar(255)" json:"fcm_token" swaggerignore:"true"`
    Pin       string         `gorm:"type:varchar(255)" json:"pin"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" swaggerignore:"true"`

    // Relationships
    Roles []Role `gorm:"many2many:user_has_roles;" json:"roles" swaggerignore:"true"`
}

// Hooks
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    // Generate unique reference
    reference := helpers.GenerateReference("USR") // USR-20250101-ABCD

    // Hash password with Argon2id
    password, err := helpers.HashPasswordArgon2(u.Password, helpers.DefaultParams)
    if err != nil {
        return err
    }

    // Hash PIN with Argon2id
    pin, err := helpers.HashPasswordArgon2(u.Pin, helpers.DefaultParams)
    if err != nil {
        return err
    }

    tx.Statement.SetColumn("reference", reference)
    tx.Statement.SetColumn("password", password)
    tx.Statement.SetColumn("pin", pin)

    return nil
}
```

**Product Model** (`app/models/product.go`):
```go
type Product struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Reference   string    `gorm:"unique" json:"reference"`
    StoreID     uint      `json:"store_id"`
    CategoryID  uint      `json:"category_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Margin      float64   `json:"margin"`
    Stock       int       `json:"stock"`
    Sold        int       `json:"sold"`
    Images      []string  `json:"images" gorm:"serializer:json"` // JSON array
    ReceivedAt  time.Time `json:"received_at"`

    // Relationships
    Category *Category `json:"category"`

    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
}

// Hooks
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
    p.Reference = helpers.GenerateReference("PRD") // PRD-20250101-ABCD
    return nil
}

func (p *Product) AfterFind(tx *gorm.DB) (err error) {
    // Convert stored image paths to full URLs
    if p.Images != nil && len(p.Images) > 0 {
        for i, image := range p.Images {
            p.Images[i] = helpers.GetFileURL(image, "member_lands")
        }
    }
    return nil
}

// AfterCreate and AfterUpdate have similar image URL logic
```

**Role Model** (`app/models/role.go`):
```go
type Role struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Name  string `json:"name"`
    Group string `json:"group"`

    // Relationships
    Users       []User       `gorm:"many2many:user_has_roles;" json:"users"`
    Permissions []Permission `gorm:"many2many:role_has_permissions;" json:"permissions"`
}
```

**Permission Model** (`app/models/permission.go`):
```go
type Permission struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Name  string `json:"name"`
    Group string `json:"group"`
}
```

**Category Model** (`app/models/category.go`):
```go
type Category struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
    Category string `json:"category"`

    // Relationships
    Products *[]Product `gorm:"foreignKey:CategoryID" json:"products,omitempty" swaggerignore:"true"`

    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
}
```

### Database Relationships

**RBAC Relationships:**
```
Users ←─────many2many─────→ Roles ←─────many2many─────→ Permissions
       (user_has_roles)              (role_has_permissions)
```

**E-commerce Relationships:**
```
Categories ←─────one2many─────→ Products
           (foreignKey: CategoryID)
```

**Pivot Tables:**
```go
// user_has_roles (app/models/user_has_role.go)
type UserHasRole struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `json:"user_id"`
    RoleID    uint      `json:"role_id"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// role_has_permissions (app/models/role_has_permission.go)
type RoleHasPermission struct {
    ID           uint `gorm:"primaryKey"`
    RoleID       uint `json:"role_id"`
    PermissionID uint `json:"permission_id"`
}
```

**Querying Relationships:**
```go
// Preload roles with user
var user models.User
db.Preload("Roles").First(&user, userID)

// Preload category with product
var product models.Product
db.Preload("Category").First(&product, productID)

// Preload permissions with role
var role models.Role
db.Preload("Permissions").First(&role, roleID)

// Preload products with category
var category models.Category
db.Preload("Products").First(&category, categoryID)
```

### GORM Hooks

**Available Hooks:**
- `BeforeCreate`: Before inserting record
- `AfterCreate`: After inserting record
- `BeforeUpdate`: Before updating record
- `AfterUpdate`: After updating record
- `BeforeDelete`: Before deleting record
- `AfterDelete`: After deleting record
- `AfterFind`: After querying record(s)

**Use Cases:**
- `BeforeCreate`: Auto-generate references, hash passwords
- `AfterFind`: Transform data (URLs, decryption)
- `AfterCreate/Update`: Update related records, trigger events
- `BeforeDelete`: Cascade delete related records

## Helpers & Utilities

### Password Hashing (`app/helpers/hash_helper.go`)

**Argon2id (Primary Method):**
```go
// Hash password
hashedPassword, err := helpers.HashPasswordArgon2(password, helpers.DefaultParams)

// DefaultParams
&Argon2Params{
    Memory:      64 * 1024, // 64 MB
    Iterations:  3,
    Parallelism: 2,
    SaltLength:  16,
    KeyLength:   32,
}

// Verify password
match, err := helpers.ComparePasswordArgon2(inputPassword, hashedPassword)

// Hash format:
// $argon2id$v=19$m=65536,t=3,p=2$<base64-salt>$<base64-hash>
```

**Why Argon2id?**
- Winner of Password Hashing Competition (2015)
- Resistant to GPU/ASIC attacks
- Memory-hard algorithm
- Configurable parameters
- Better than Bcrypt for modern applications

**Bcrypt (Legacy Support):**
```go
// Still available for backward compatibility
hashedPassword := helpers.HashPasswordBcrypt(password)
match := helpers.CheckPasswordHash(password, hashedPassword)
```

### Reference Code Generation (`app/helpers/reference_helper.go`)

```go
// Generate unique reference
ref := helpers.GenerateReference("USR")
// Output: USR-20250101-A1B2C3D4

ref := helpers.GenerateReference("PRD")
// Output: PRD-20250102-X9Y8Z7W6

// Format: PREFIX-YYYYMMDD-RANDOM8CHAR
```

**Used for:**
- User references (USR)
- Product references (PRD)
- Order references (ORD)
- Any entity needing human-readable unique ID

### File Helpers (`app/helpers/file_helper.go`)

```go
// Get full file URL
url := helpers.GetFileURL("image.jpg", "products")
// Output: http://localhost:9999/file/products/image.jpg

// Validate file type
isValid := helpers.ValidateFileType("image.png", []string{".jpg", ".png", ".gif"})

// Get file extension
ext := helpers.GetFileExtension("document.pdf")
// Output: ".pdf"
```

### Base64 File Helpers (`app/helpers/base64file_helper.go`)

```go
// Decode base64 to file
filename, err := helpers.DecodeBase64ToFile(base64String, "uploads/", "image.jpg")

// Encode file to base64
base64String, err := helpers.EncodeFileToBase64("uploads/image.jpg")
```

### Environment Helpers (`app/helpers/env_helper.go`)

```go
// Get env var with default
appName := helpers.GetEnv("APP_NAME", "Default App")

// Get env var as int
port := helpers.GetEnvInt("APP_PORT", 8080)

// Get env var as bool
debug := helpers.GetEnvBool("DEBUG", false)
```

### Response Helpers (`app/helpers/response_helper.go`)

```go
// Success response
helpers.SuccessResponse(c, "Operation successful", userData)

// Error response
helpers.ErrorResponse(c, http.StatusBadRequest, "Invalid input")

// Error with reference code
helpers.ResponseError(c, &helpers.ResponseParams[any]{
    Reference: "ERROR-001",
    Message:   "User not found",
}, http.StatusNotFound)

// Success with reference code
helpers.ResponseSuccess(c, &helpers.ResponseParams[models.User]{
    Reference: "SUCCESS-001",
    Message:   "User created",
    Data:      user,
}, http.StatusCreated)
```

## Role-Based Access Control (RBAC)

### RBAC Architecture

**Three-Layer Model:**
```
Users → Roles → Permissions
```

**Models:**
- `User`: Application users with authentication
- `Role`: Groups of permissions (e.g., "admin", "editor", "viewer")
- `Permission`: Granular access rights (e.g., "users.create", "products.delete")

**Relationships:**
- Users ↔ Roles: Many-to-Many (one user can have multiple roles)
- Roles ↔ Permissions: Many-to-Many (one role can have multiple permissions)

**Indirect Relationship:**
Users get permissions through their assigned roles.

### RBAC Endpoints

**User-Role Assignment:**
```bash
# Assign roles to user
POST /users/:id/roles
Body: {
    "role_ids": [1, 2, 3]
}

# Get user's roles
GET /users/:id/roles
Response: [
    {"id": 1, "name": "admin", "group": "system"},
    {"id": 2, "name": "editor", "group": "content"}
]
```

**Role-Permission Assignment:**
```bash
# Assign permissions to role
POST /roles/:id/permissions
Body: {
    "permission_ids": [1, 2, 3, 4, 5]
}

# Get role's permissions
GET /roles/:id/permissions
Response: [
    {"id": 1, "name": "users.create", "group": "users"},
    {"id": 2, "name": "users.read", "group": "users"}
]
```

**Role Management:**
```bash
# List all roles
GET /roles

# Create/Update role
PUT /roles
Body: {
    "name": "moderator",
    "group": "content"
}

# Delete role
DELETE /roles/:id
```

**Permission Management:**
```bash
# List all permissions
GET /permissions

# Create/Update permission
PUT /permissions
Body: {
    "name": "products.create",
    "group": "products"
}

# Delete permission
DELETE /permissions/:id
```

### RBAC Implementation in Services

**Role Service** (`app/services/role_service.go`):
```go
type RoleService struct {
    roleRepo       interfaces.RoleRepositoryInterface
    permissionRepo interfaces.PermissionRepositoryInterface
}

func (s *RoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
    // Validate role exists
    role, err := s.roleRepo.FindByID(roleID)
    if err != nil {
        return errors.New("role not found")
    }

    // Validate all permissions exist
    for _, permID := range permissionIDs {
        if _, err := s.permissionRepo.FindByID(permID); err != nil {
            return fmt.Errorf("permission %d not found", permID)
        }
    }

    // Assign permissions (handled by repository)
    return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}
```

**User Service** (`app/services/user_services.go`):
```go
func (s *UserService) AssignRoles(userID uint, roleIDs []uint) error {
    // Validate user exists
    _, err := s.repo.FindByID(userID)
    if err != nil {
        return errors.New("user not found")
    }

    // Assign roles (repository handles pivot table)
    return s.repo.AssignRoles(userID, roleIDs)
}

func (s *UserService) GetRoles(userID uint) ([]models.Role, error) {
    return s.repo.GetRoles(userID)
}
```

### RBAC in Database

**Pivot Tables:**
```sql
-- user_has_roles
CREATE TABLE user_has_roles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_role (user_id, role_id)
);

-- role_has_permissions
CREATE TABLE role_has_permissions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    UNIQUE KEY unique_role_permission (role_id, permission_id)
);
```

**Cascade Deletes:**
- Delete user → removes all user-role assignments
- Delete role → removes all user-role and role-permission assignments
- Delete permission → removes all role-permission assignments

### Permission Naming Convention

**Format:** `resource.action`

**Examples:**
```
users.create
users.read
users.update
users.delete

products.create
products.read
products.update
products.delete
products.publish

orders.create
orders.read
orders.update
orders.cancel
orders.refund
```

**Grouping:**
Permissions have `group` field for organization:
- Group: "users" → users.*, roles.*, permissions.*
- Group: "products" → products.*, categories.*
- Group: "orders" → orders.*, invoices.*

## File Upload & Storage

### File System

**Storage Location:**
```
storage/
├── products/         # Product images
├── users/            # User avatars
├── documents/        # PDF, docs
└── temp/             # Temporary uploads
```

### File Endpoints

```bash
# Serve file (public)
GET /file/:key/:filename
# Example: GET /file/products/image123.jpg
# Returns: Image binary with appropriate Content-Type

# Upload file (handled in controllers)
# Usually in product/user creation/update endpoints
```

### File Upload Implementation

**Controller Example** (`app/controllers/product_controller.go`):
```go
func (ctrl *ProductController) Put(c *gin.Context) {
    var request requests.ProductRequest
    c.ShouldBindJSON(&request)

    // Process images (base64 or file paths)
    images := []string{}
    for _, imageData := range request.ImagesBase64 {
        // Decode base64 and save
        filename, err := helpers.DecodeBase64ToFile(
            imageData,
            "storage/products/",
            helpers.GenerateRandomFilename(".jpg"),
        )
        if err == nil {
            images = append(images, filename)
        }
    }

    product := models.Product{
        Name:   request.Name,
        Images: images, // Stored as JSON array in DB
        // ...
    }

    ctrl.service.Create(&product)
}
```

**File Service** (`app/services/file_service.go`):
```go
type FileService struct{}

func (s *FileService) SaveFile(fileData []byte, path string) (string, error) {
    // Generate unique filename
    filename := helpers.GenerateRandomFilename(".jpg")
    fullPath := filepath.Join(path, filename)

    // Save file
    err := ioutil.WriteFile(fullPath, fileData, 0644)
    return filename, err
}

func (s *FileService) DeleteFile(path string) error {
    return os.Remove(path)
}
```

**File Controller** (`app/controllers/file_controller.go`):
```go
func (ctrl *FileController) ServeFile(c *gin.Context) {
    key := c.Param("key")         // Folder: "products", "users"
    filename := c.Param("filename") // File: "image123.jpg"

    filePath := filepath.Join("storage", key, filename)

    // Check file exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        c.JSON(404, gin.H{"error": "File not found"})
        return
    }

    // Serve file
    c.File(filePath)
}
```

### Image URL Generation

**In Models (AfterFind hook):**
```go
func (p *Product) AfterFind(tx *gorm.DB) (err error) {
    // Convert stored paths to full URLs
    for i, image := range p.Images {
        p.Images[i] = helpers.GetFileURL(image, "products")
        // Converts: "image123.jpg"
        // To: "http://localhost:9999/file/products/image123.jpg"
    }
    return nil
}
```

**Manual URL Generation:**
```go
url := helpers.GetFileURL("avatar.jpg", "users")
// Output: http://localhost:9999/file/users/avatar.jpg
```

## Code Patterns & Best Practices

### Adding New Features (Complete Guide)

**Example: Adding "Orders" Feature**

**1. Create Migration:**
```bash
go run main.go make:migration create_orders_table
```

Edit `app/database/migrations/YYYYMMDDHHMMSS_create_orders_table.sql`:
```sql
-- +++ UP Migration
CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    reference VARCHAR(100) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- --- DOWN Migration
DROP TABLE IF EXISTS orders;
```

**2. Run Migration:**
```bash
go run main.go migrate:all
```

**3. Create Model** (`app/models/order.go`):
```go
package models

import (
    "time"
    "golang_starter_kit_2025/app/helpers"
    "gorm.io/gorm"
)

type Order struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Reference   string         `gorm:"unique" json:"reference"`
    UserID      uint           `gorm:"not null" json:"user_id"`
    TotalAmount float64        `gorm:"type:decimal(10,2)" json:"total_amount"`
    Status      string         `gorm:"type:varchar(50)" json:"status"`
    CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

    // Relationships
    User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
    o.Reference = helpers.GenerateReference("ORD")
    return nil
}
```

**4. Create Repository Interface** (`app/repositories/interfaces/order_repository_interface.go`):
```go
package interfaces

import "golang_starter_kit_2025/app/models"

type OrderRepositoryInterface interface {
    Create(order *models.Order) error
    Update(order *models.Order) error
    Delete(id uint) error
    FindByID(id uint) (*models.Order, error)
    FindByUserID(userID uint, page, limit int) ([]models.Order, int64, error)
    List(page, limit int) ([]models.Order, int64, error)
    FindByStatus(status string, page, limit int) ([]models.Order, int64, error)
    UpdateStatus(orderID uint, status string) error
}
```

**5. Implement Repository** (`app/repositories/order_repository.go`):
```go
package repositories

import (
    "errors"
    "golang_starter_kit_2025/app/models"
    "golang_starter_kit_2025/app/models/scopes"
    "golang_starter_kit_2025/app/repositories/interfaces"
    "gorm.io/gorm"
)

type orderRepository struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) interfaces.OrderRepositoryInterface {
    return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
    return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uint) (*models.Order, error) {
    var order models.Order
    err := r.db.Preload("User").First(&order, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("order not found")
        }
        return nil, err
    }
    return &order, nil
}

func (r *orderRepository) FindByUserID(userID uint, page, limit int) ([]models.Order, int64, error) {
    var orders []models.Order
    var total int64

    r.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total)
    err := r.db.Where("user_id = ?", userID).
        Scopes(scopes.PaginateSimple(page, limit)).
        Preload("User").
        Find(&orders).Error

    return orders, total, err
}

// ... implement other interface methods
```

**6. Create Service** (`app/services/order_service.go`):
```go
package services

import (
    "errors"
    "golang_starter_kit_2025/app/models"
    "golang_starter_kit_2025/app/repositories/interfaces"
)

type OrderService struct {
    orderRepo interfaces.OrderRepositoryInterface
    userRepo  interfaces.UserRepositoryInterface
}

func NewOrderService(
    orderRepo interfaces.OrderRepositoryInterface,
    userRepo interfaces.UserRepositoryInterface,
) *OrderService {
    return &OrderService{
        orderRepo: orderRepo,
        userRepo:  userRepo,
    }
}

func (s *OrderService) CreateOrder(userID uint, totalAmount float64) (*models.Order, error) {
    // Validate user exists
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    // Create order
    order := &models.Order{
        UserID:      user.ID,
        TotalAmount: totalAmount,
        Status:      "pending",
    }

    if err := s.orderRepo.Create(order); err != nil {
        return nil, err
    }

    return order, nil
}

func (s *OrderService) GetUserOrders(userID uint, page, limit int) ([]models.Order, int64, error) {
    return s.orderRepo.FindByUserID(userID, page, limit)
}

func (s *OrderService) UpdateOrderStatus(orderID uint, status string) error {
    // Validate status
    validStatuses := []string{"pending", "processing", "completed", "cancelled"}
    isValid := false
    for _, s := range validStatuses {
        if s == status {
            isValid = true
            break
        }
    }
    if !isValid {
        return errors.New("invalid order status")
    }

    return s.orderRepo.UpdateStatus(orderID, status)
}
```

**7. Create Request Validator** (`app/requests/order_request.go`):
```go
package requests

type CreateOrderRequest struct {
    UserID      uint    `json:"user_id" binding:"required"`
    TotalAmount float64 `json:"total_amount" binding:"required,min=0"`
}

type UpdateOrderStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=pending processing completed cancelled"`
}
```

**8. Create Controller** (`app/controllers/order_controller.go`):
```go
package controllers

import (
    "net/http"
    "strconv"
    "golang_starter_kit_2025/app/helpers"
    "golang_starter_kit_2025/app/requests"
    "golang_starter_kit_2025/app/services"
    "github.com/gin-gonic/gin"
)

type OrderController struct {
    service *services.OrderService
}

func NewOrderController(service services.OrderService) *OrderController {
    return &OrderController{service: &service}
}

// @Summary Create new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body requests.CreateOrderRequest true "Order data"
// @Success 201 {object} models.Order
// @Failure 400 {object} helpers.ResponseParams
// @Router /orders [post]
// @Security Bearer
func (ctrl *OrderController) Create(c *gin.Context) {
    var request requests.CreateOrderRequest

    if err := c.ShouldBindJSON(&request); err != nil {
        helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    order, err := ctrl.service.CreateOrder(request.UserID, request.TotalAmount)
    if err != nil {
        helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    helpers.SuccessResponse(c, "Order created successfully", order)
}

// @Summary Get user's orders
// @Tags orders
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} models.Order
// @Router /users/{user_id}/orders [get]
// @Security Bearer
func (ctrl *OrderController) GetUserOrders(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 32)
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    orders, total, err := ctrl.service.GetUserOrders(uint(userID), page, limit)
    if err != nil {
        helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    helpers.SuccessResponse(c, "Orders retrieved", gin.H{
        "orders": orders,
        "total":  total,
        "page":   page,
        "limit":  limit,
    })
}
```

**9. Register Routes** (`routes/web.go`):
```go
func RegisterRoutes(route *gin.Engine) {
    // ... existing repo creation
    orderRepo := repositories.NewOrderRepository(facades.DB)

    // ... existing service creation
    orderService := services.NewOrderService(orderRepo, userRepo)

    // ... existing controller creation
    orderController := controllers.NewOrderController(*orderService)

    // Register routes
    orderRoutes := route.Group("/orders", middleware.AuthMiddleware())
    {
        orderRoutes.POST("", orderController.Create)
        orderRoutes.GET("/:id", orderController.GetByID)
        orderRoutes.GET("/user/:user_id", orderController.GetUserOrders)
        orderRoutes.PATCH("/:id/status", orderController.UpdateStatus)
    }
}
```

**10. Update Swagger:**
```bash
swag init
```

**11. Create Seeder (Optional):**
```bash
go run main.go make:seeder --name=OrdersSeeder
```

**12. Test:**
```bash
# Create test file: app/controllers/order_controller_test.go
# Run tests
go test ./app/controllers/order_controller_test.go
```

### Repository Pattern Best Practices

**Interface First:**
- Always define interface before implementation
- Put interfaces in `app/repositories/interfaces/`
- One interface per domain model

**Repository Responsibilities:**
- ONLY database operations (CRUD)
- NO business logic
- NO validation (except basic DB constraints)
- Return domain models or errors

**Error Handling:**
```go
// Good: Specific error messages
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, errors.New("user not found")
}

// Bad: Exposing GORM errors
return nil, err
```

**Preloading Relationships:**
```go
// Good: Explicit preload control
func (r *orderRepository) FindByIDWithRelations(id uint) (*models.Order, error) {
    var order models.Order
    err := r.db.Preload("User").Preload("Items").First(&order, id).Error
    return &order, err
}

// Separate method for different preload scenarios
func (r *orderRepository) FindByID(id uint) (*models.Order, error) {
    var order models.Order
    err := r.db.First(&order, id).Error
    return &order, err
}
```

### Service Layer Best Practices

**Business Logic:**
- All business rules in services
- Validation before DB operations
- Coordinate multiple repositories
- Return domain errors, not HTTP errors

**Transaction Handling:**
```go
func (s *OrderService) CreateOrderWithItems(orderData OrderData) error {
    return facades.DB.Transaction(func(tx *gorm.DB) error {
        // Create repositories with transaction DB
        orderRepo := repositories.NewOrderRepository(tx)
        itemRepo := repositories.NewOrderItemRepository(tx)

        // Create order
        order := &models.Order{...}
        if err := orderRepo.Create(order); err != nil {
            return err // Automatic rollback
        }

        // Create order items
        for _, itemData := range orderData.Items {
            item := &models.OrderItem{OrderID: order.ID, ...}
            if err := itemRepo.Create(item); err != nil {
                return err // Automatic rollback
            }
        }

        return nil // Commit
    })
}
```

**Aggregating Data:**
```go
func (s *ProductService) GetProductWithFullDetails(id uint) (*ProductDetails, error) {
    // Use multiple repositories
    product, err := s.productRepo.FindByID(id)
    if err != nil {
        return nil, err
    }

    category, _ := s.categoryRepo.FindByID(product.CategoryID)
    reviews, _ := s.reviewRepo.FindByProductID(product.ID, 1, 5)

    return &ProductDetails{
        Product:  product,
        Category: category,
        Reviews:  reviews,
    }, nil
}
```

### Controller Best Practices

**Thin Controllers:**
```go
// Good: Delegate to service
func (ctrl *UserController) Create(c *gin.Context) {
    var request requests.CreateUserRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    user, err := ctrl.service.CreateUser(request)
    if err != nil {
        helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    helpers.SuccessResponse(c, "User created", user)
}

// Bad: Business logic in controller
func (ctrl *UserController) Create(c *gin.Context) {
    // ... validation
    // Hashing password here
    hashedPassword, _ := bcrypt.GenerateFromPassword(...)
    // Direct DB access
    db.Create(&User{Password: hashedPassword})
    // ...
}
```

**Error Mapping:**
```go
func (ctrl *UserController) Get(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

    user, err := ctrl.service.FindByID(uint(id))
    if err != nil {
        // Map service errors to HTTP errors
        if err.Error() == "user not found" {
            helpers.ErrorResponse(c, http.StatusNotFound, err.Error())
            return
        }
        helpers.ErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    helpers.SuccessResponse(c, "User found", user)
}
```

## Testing

### Test Structure

**Using Ginkgo/Gomega (BDD):**
```go
package controllers_test

import (
    "testing"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestControllers(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Controllers Suite")
}

var _ = Describe("UserController", func() {
    Context("When creating a user", func() {
        It("should create user successfully", func() {
            // Test implementation
            user := CreateUser("test@example.com", "password")
            Expect(user.Email).To(Equal("test@example.com"))
        })

        It("should fail with invalid email", func() {
            _, err := CreateUser("invalid-email", "password")
            Expect(err).To(HaveOccurred())
        })
    })
})
```

**Running Tests:**
```bash
# All tests with race detector
ginkgo -r --race

# With coverage
ginkgo -r --cover

# Specific package
ginkgo ./app/controllers

# Watch mode (re-run on file changes)
ginkgo watch -r
```

## Health Checks

### Default Database Health

```bash
GET /health
```

**Response:**
```json
{
    "message": "database is connected",
    "database": "supply_chain_retail"
}
```

### Multi-Database Health

```bash
GET /health/databases
```

**Response:**
```json
{
    "overall_health": true,
    "connections": {
        "mysql": {
            "status": "healthy",
            "stats": {
                "max_open_connections": 200,
                "open_connections": 5,
                "in_use": 2,
                "idle": 3
            }
        },
        "postgres": {
            "status": "healthy",
            "stats": {
                "max_open_connections": 200,
                "open_connections": 3,
                "in_use": 1,
                "idle": 2
            }
        },
        "mysql_secondary": {
            "status": "disconnected"
        }
    }
}
```

**HTTP Status Codes:**
- 200: All connections healthy
- 503: One or more connections unhealthy/disconnected

## Important Notes & Best Practices

### Security

1. **JWT Secret Key**
   - NEVER use placeholder values in production
   - Generate strong key: `openssl rand -base64 48`
   - App validates on startup and fails if placeholder detected

2. **SKIP_AUTH Flag**
   - ONLY for development
   - Bypasses ALL authentication
   - Sets user_id=1 automatically
   - ⚠️ **NEVER enable in production!**

3. **Password Hashing**
   - Primary: Argon2id (memory-hard, GPU-resistant)
   - Always use `helpers.HashPasswordArgon2()` in services
   - User model auto-hashes in `BeforeCreate` hook
   - Never store plain passwords

4. **Environment Variables**
   - Never commit `.env` file
   - Use `.env.example` as template
   - Required vars validated on startup

### Database

1. **Migrations**
   - Never modify executed migrations
   - Create new migration for schema changes
   - Always provide DOWN migration for rollback
   - Test migrations on dev database first
   - Use `--connection` flag for multi-DB setups

2. **Seeders**
   - Filename MUST match function names
   - File: `UserSeeder.go` → Functions: `SeedUserSeeder`, `RollbackUserSeeder`
   - Use `Unscoped()` for hard deletes in rollbacks
   - Seeders tracked in batches like migrations

3. **Repositories**
   - All database access through repository layer
   - Use interfaces for testability
   - Return specific errors, not GORM errors
   - Inject `*gorm.DB` for multi-DB support

### Code Organization

1. **Dependency Injection**
   - Wire all dependencies in `routes/web.go`
   - Constructor-based injection
   - Repositories → Services → Controllers

2. **Layering**
   - Controllers: HTTP handling only
   - Services: Business logic
   - Repositories: Database operations
   - Models: Data structures + GORM hooks

3. **Error Handling**
   - Use reference codes for error tracking
   - Map service errors to HTTP status codes in controllers
   - Provide meaningful error messages

4. **Pagination**
   - Default: page=1, limit=10
   - Maximum limit: 100
   - Use `scopes.Paginate()` or `scopes.PaginateSimple()`

### Development Workflow

1. **API Changes**
   - Update controller methods
   - Update Swagger annotations
   - Run `swag init` to regenerate docs

2. **Adding Features**
   - Follow 12-step pattern (see "Adding New Features")
   - Create interface → implementation → tests
   - Always run tests before committing

3. **Hot Reload**
   - Use `air` for development
   - Watches file changes and auto-restarts
   - Configure in `.air.toml` if needed

### Performance

1. **Connection Pooling**
   - Configure per environment load
   - Default: 10 idle, 200 max open
   - Monitor with health endpoints

2. **Preloading**
   - Only preload when needed
   - Separate repository methods for different preload scenarios
   - Avoid N+1 queries

3. **Pagination**
   - Always paginate list endpoints
   - Enforce maximum limit (100)
   - Return total count for client-side pagination UI

### Common Pitfalls

1. **GORM Hooks**
   - Hooks run AFTER validation
   - Use `tx.Statement.SetColumn()` in hooks, not direct field assignment
   - BeforeCreate runs before INSERT, not before validation

2. **JSON Arrays in MySQL**
   - Use `gorm:"serializer:json"` tag
   - Store as JSON string in database
   - Auto-serializes/deserializes

3. **Soft Deletes**
   - Models with `gorm.DeletedAt` soft delete by default
   - Use `Unscoped()` to query/delete permanently
   - Affects cascade deletes (use `ON DELETE CASCADE` in SQL)

4. **Reference Codes**
   - Auto-generated in BeforeCreate hooks
   - Format: PREFIX-YYYYMMDD-RANDOM8
   - Use `helpers.GenerateReference(prefix)`

5. **File URLs**
   - Stored as relative paths in DB
   - Converted to full URLs in AfterFind hooks
   - Use `helpers.GetFileURL()` for manual conversion

### Project-Specific Conventions

1. **Naming**
   - Models: Singular (User, Product, Order)
   - Tables: Plural (users, products, orders)
   - Repositories: `{Model}Repository` (UserRepository)
   - Services: `{Model}Service` (UserService)
   - Controllers: `{Model}Controller` (UserController)

2. **HTTP Methods**
   - GET: Retrieve resources
   - POST: Create resources (specific endpoint)
   - PUT: Create or Update (upsert pattern used in this project)
   - PATCH: Partial update
   - DELETE: Remove resources

3. **Response Format**
   - Success: `{"message": "...", "data": {...}}`
   - Error: `{"reference": "ERROR-XXX", "message": "..."}`
   - Paginated: `{"data": [...], "total": 100, "page": 1, "limit": 10}`

---

**Project Maintained by:** [Dzyfhuba](https://github.com/Dzyfhuba) & [RahmatRafiq](https://github.com/RahmatRafiq)

**License:** MIT
