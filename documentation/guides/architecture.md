# Architecture Guide

## Overview

Aplikasi ini menggunakan arsitektur modular yang terinspirasi dari Laravel, dengan fokus pada separation of concerns, maintainability, dan scalability. Arsitektur ini mendukung multiple database connections dan pattern repository untuk consistency.

## Architectural Patterns

### 1. Layered Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
├─────────────────────────────────────────────────────────────┤
│  Controllers  │  Middleware  │  Handlers  │   Routes        │
├─────────────────────────────────────────────────────────────┤
│                    Business Logic Layer                     │
├─────────────────────────────────────────────────────────────┤
│  Services     │  Validators  │  Helpers   │   Casts         │
├─────────────────────────────────────────────────────────────┤
│                    Data Access Layer                        │
├─────────────────────────────────────────────────────────────┤
│  Repositories │  Models      │  Database  │   Migrations    │
├─────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                     │
├─────────────────────────────────────────────────────────────┤
│  Database     │  Config      │  Facades   │   CLI Commands  │
│  Manager      │              │            │                 │
└─────────────────────────────────────────────────────────────┘
```

### 2. Repository Pattern

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Controller    │────│    Service      │────│   Repository    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │              ┌─────────────────┐
         │                       │              │   Interface     │
         │                       │              │   Contract      │
         │                       │              └─────────────────┘
         │                       │                       │
         │              ┌─────────────────┐              │
         │              │   Business      │              │
         │              │   Logic         │              │
         │              └─────────────────┘              │
         │                                               │
         └───────────────────────────────────────────────┼────────┐
                                                         │        │
                                              ┌─────────────────┐ │
                                              │     MySQL       │ │
                                              │   Repository    │ │
                                              └─────────────────┘ │
                                                         │        │
                                              ┌─────────────────┐ │
                                              │   PostgreSQL    │ │
                                              │   Repository    │ │
                                              └─────────────────┘ │
                                                                  │
                                              ┌─────────────────┐ │
                                              │     Cache       │ │
                                              │   Repository    │ │
                                              └─────────────────┘─┘
```

## Directory Structure

```
golang_starter_kit_2025/
├── app/
│   ├── controllers/          # HTTP request handlers
│   ├── services/            # Business logic layer
│   ├── repositories/        # Data access layer
│   ├── models/              # Data models & entities
│   ├── middleware/          # HTTP middleware
│   ├── requests/            # Request validation
│   ├── responses/           # Response formatting
│   ├── helpers/             # Utility functions
│   ├── casts/              # Data casting & transformation
│   └── database/           # Database management
│       ├── migrations/     # Database migrations
│       └── seeds/          # Database seeders
├── bootstrap/              # Application initialization
├── config/                 # Configuration management
├── database/              # Database manager & facades
├── docs/                  # Documentation
├── facades/               # Facade pattern implementations
├── interfaces/            # Interface definitions
├── routes/                # Route definitions
└── services/              # External service integrations
```

## Component Responsibilities

### 1. Controllers Layer

**Purpose**: Handle HTTP requests and responses

```go
// app/controllers/user_controller.go
type UserController struct {
    userService *services.UserService
}

func (c *UserController) GetUser(ctx *gin.Context) {
    id := ctx.Param("id")
    
    user, err := c.userService.GetUserByID(id)
    if err != nil {
        responses.ErrorResponse(ctx, http.StatusNotFound, "User not found")
        return
    }
    
    responses.SuccessResponse(ctx, "User retrieved successfully", user)
}
```

**Responsibilities**:
- ✅ Parse and validate HTTP requests
- ✅ Call appropriate services
- ✅ Format HTTP responses
- ❌ Business logic
- ❌ Database operations
- ❌ Complex data manipulation

### 2. Services Layer

**Purpose**: Implement business logic and orchestrate operations

```go
// app/services/user_service.go
type UserService struct {
    userRepo      repositories.UserRepositoryInterface
    analyticsRepo repositories.AnalyticsRepositoryInterface
    jwtService    *JWTService
}

func (s *UserService) CreateUser(req *requests.CreateUserRequest) (*models.User, error) {
    // Business logic
    hashedPassword, err := helpers.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
    }
    
    // Create user in primary database
    createdUser, err := s.userRepo.Create(user)
    if err != nil {
        return nil, err
    }
    
    // Log to analytics database
    s.analyticsRepo.LogUserCreation(createdUser.ID)
    
    return createdUser, nil
}
```

**Responsibilities**:
- ✅ Business rules implementation
- ✅ Data validation and transformation
- ✅ Orchestrate multiple repository calls
- ✅ Transaction management
- ❌ HTTP request/response handling
- ❌ Direct database queries

### 3. Repository Layer

**Purpose**: Abstract data access operations

```go
// interfaces/user_repository_interface.go
type UserRepositoryInterface interface {
    Create(user *models.User) (*models.User, error)
    GetByID(id uint) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
    GetWithPagination(page, limit int) ([]models.User, int64, error)
}

// app/repositories/mysql_user_repository.go
type MySQLUserRepository struct {
    db *database.Connection
}

func (r *MySQLUserRepository) Create(user *models.User) (*models.User, error) {
    result := r.db.DB.Create(user)
    if result.Error != nil {
        return nil, result.Error
    }
    return user, nil
}

func (r *MySQLUserRepository) GetByID(id uint) (*models.User, error) {
    var user models.User
    result := r.db.DB.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
```

**Responsibilities**:
- ✅ Database operations (CRUD)
- ✅ Query optimization
- ✅ Data mapping
- ✅ Connection management
- ❌ Business logic
- ❌ Request validation
- ❌ Response formatting

### 4. Models Layer

**Purpose**: Define data structures and relationships

```go
// app/models/user.go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
    
    // Relationships
    Orders []Order `json:"orders,omitempty" gorm:"foreignKey:UserID"`
    Roles  []Role  `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

// Database-specific methods
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // Auto-generate UUID or other preprocessing
    return nil
}

func (u *User) TableName() string {
    return "users"
}
```

**Responsibilities**:
- ✅ Data structure definition
- ✅ Database relationships
- ✅ Validation rules
- ✅ Serialization/Deserialization
- ❌ Business logic
- ❌ Database queries

## Database Manager Architecture

### 1. Connection Manager

```go
// database/manager.go
type Manager struct {
    connections map[string]*Connection
    configs     *config.DatabaseConfigs
    mutex       sync.RWMutex
}

type Connection struct {
    DB     *gorm.DB
    SqlDB  *sql.DB
    Config *config.DatabaseConfig
    Name   string
}
```

### 2. Facade Pattern

```go
// facades/database.go
func MySQL() (*database.Connection, error) {
    return GetConnection("mysql")
}

func PostgreSQL() (*database.Connection, error) {
    return GetConnection("postgres")
}

func GetConnection(name string) (*database.Connection, error) {
    if manager == nil {
        ConnectDB()
    }
    return manager.GetConnection(name)
}
```

## Request/Response Flow

### 1. Complete Request Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │    │    Route    │    │ Middleware  │    │ Controller  │
│   Request   │───▶│   Handler   │───▶│   Stack     │───▶│   Method    │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
                                                                  │
┌─────────────┐    ┌─────────────┐    ┌─────────────┐           │
│  Response   │◀───│   Format    │◀───│   Service   │◀──────────┘
│   to Client │    │  Response   │    │   Layer     │
└─────────────┘    └─────────────┘    └─────────────┘
                                              │
                           ┌─────────────────────────────────┐
                           │                                 │
                           ▼                                 ▼
                  ┌─────────────┐                  ┌─────────────┐
                  │ Repository  │                  │ Repository  │
                  │   MySQL     │                  │ PostgreSQL  │
                  └─────────────┘                  └─────────────┘
```

### 2. Example Implementation

```go
// routes/user_routes.go
func RegisterUserRoutes(r *gin.Engine) {
    userGroup := r.Group("/api/users")
    userGroup.Use(middleware.AuthMiddleware())
    
    userController := controllers.NewUserController()
    
    userGroup.GET("/:id", userController.GetUser)
    userGroup.POST("/", userController.CreateUser)
    userGroup.PUT("/:id", userController.UpdateUser)
    userGroup.DELETE("/:id", userController.DeleteUser)
}

// controllers/user_controller.go
func (c *UserController) CreateUser(ctx *gin.Context) {
    var req requests.CreateUserRequest
    
    if err := ctx.ShouldBindJSON(&req); err != nil {
        responses.ValidationErrorResponse(ctx, err)
        return
    }
    
    user, err := c.userService.CreateUser(&req)
    if err != nil {
        responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
        return
    }
    
    responses.SuccessResponse(ctx, "User created successfully", user)
}
```

## Dependency Injection

### 1. Constructor Injection

```go
// services/user_service.go
type UserService struct {
    userRepo      repositories.UserRepositoryInterface
    roleRepo      repositories.RoleRepositoryInterface
    emailService  *EmailService
    cacheService  *CacheService
}

func NewUserService(
    userRepo repositories.UserRepositoryInterface,
    roleRepo repositories.RoleRepositoryInterface,
    emailService *EmailService,
    cacheService *CacheService,
) *UserService {
    return &UserService{
        userRepo:     userRepo,
        roleRepo:     roleRepo,
        emailService: emailService,
        cacheService: cacheService,
    }
}
```

### 2. Service Container (Optional)

```go
// bootstrap/container.go
type Container struct {
    services map[string]interface{}
}

func (c *Container) Register(name string, service interface{}) {
    c.services[name] = service
}

func (c *Container) Resolve(name string) interface{} {
    return c.services[name]
}

// Usage
container := &Container{services: make(map[string]interface{})}
container.Register("userService", NewUserService(...))
```

## Error Handling Strategy

### 1. Layered Error Handling

```go
// helpers/errors.go
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
    return e.Message
}

// Custom error types
var (
    ErrUserNotFound     = &AppError{Code: "USER_NOT_FOUND", Message: "User not found"}
    ErrInvalidPassword  = &AppError{Code: "INVALID_PASSWORD", Message: "Invalid password"}
    ErrDatabaseError    = &AppError{Code: "DATABASE_ERROR", Message: "Database operation failed"}
)

// Repository layer
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
    var user models.User
    result := r.db.DB.First(&user, id)
    
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, ErrUserNotFound
    }
    
    if result.Error != nil {
        return nil, fmt.Errorf("%w: %v", ErrDatabaseError, result.Error)
    }
    
    return &user, nil
}
```

### 2. Global Error Handler

```go
// middleware/error_handler.go
func ErrorHandlerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            var appErr *AppError
            if errors.As(err, &appErr) {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": appErr,
                })
                return
            }
            
            // Handle other error types
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal server error",
            })
        }
    }
}
```

## Testing Architecture

### 1. Repository Testing

```go
// repositories/user_repository_test.go
func TestUserRepository_Create(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    repo := &MySQLUserRepository{db: db}
    
    // Test data
    user := &models.User{
        Name:  "Test User",
        Email: "test@example.com",
    }
    
    // Execute
    result, err := repo.Create(user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, result.ID)
    assert.Equal(t, "Test User", result.Name)
}
```

### 2. Service Testing with Mocks

```go
// services/user_service_test.go
func TestUserService_CreateUser(t *testing.T) {
    // Setup mocks
    mockRepo := &mocks.UserRepositoryInterface{}
    service := NewUserService(mockRepo, nil, nil, nil)
    
    // Setup expectations
    mockRepo.On("Create", mock.AnythingOfType("*models.User")).
        Return(&models.User{ID: 1, Name: "Test User"}, nil)
    
    // Execute
    req := &requests.CreateUserRequest{
        Name:  "Test User",
        Email: "test@example.com",
    }
    
    user, err := service.CreateUser(req)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, uint(1), user.ID)
    mockRepo.AssertExpectations(t)
}
```

## Performance Considerations

### 1. Database Query Optimization

```go
// Good: Use eager loading for relationships
func (r *UserRepository) GetWithOrders(id uint) (*models.User, error) {
    var user models.User
    result := r.db.DB.Preload("Orders").First(&user, id)
    return &user, result.Error
}

// Good: Use pagination for large datasets
func (r *UserRepository) GetWithPagination(page, limit int) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    offset := (page - 1) * limit
    
    r.db.DB.Model(&models.User{}).Count(&total)
    result := r.db.DB.Offset(offset).Limit(limit).Find(&users)
    
    return users, total, result.Error
}
```

### 2. Caching Strategy

```go
// services/user_service.go
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("user:%d", id)
    if cached := s.cacheService.Get(cacheKey); cached != nil {
        return cached.(*models.User), nil
    }
    
    // Fallback to database
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // Cache for future requests
    s.cacheService.Set(cacheKey, user, 5*time.Minute)
    
    return user, nil
}
```

## Security Architecture

### 1. Authentication Middleware

```go
// middleware/auth_middleware.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
            c.Abort()
            return
        }
        
        claims, err := validateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

### 2. Permission Checking

```go
// middleware/permission_middleware.go
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetUint("user_id")
        
        hasPermission, err := checkUserPermission(userID, permission)
        if err != nil || !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

Next: [CLI Guide](cli.md) | [API Documentation](../api/README.md)
