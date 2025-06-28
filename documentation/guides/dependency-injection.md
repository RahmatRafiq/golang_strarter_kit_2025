# Dependency Injection

Implementasi Dependency Injection pada Golang Starter Kit untuk loose coupling dan testability.

## Overview

Dependency Injection (DI) adalah design pattern yang memungkinkan kita untuk menginject dependencies ke dalam suatu object, bukan membuat dependencies tersebut di dalam object itu sendiri. Ini meningkatkan testability, maintainability, dan flexibility dari kode.

## Benefits

- **Loose Coupling**: Mengurangi ketergantungan antar komponen
- **Testability**: Mudah untuk membuat mock dependencies
- **Flexibility**: Mudah mengganti implementasi
- **Maintainability**: Kode lebih mudah dirawat
- **Single Responsibility**: Setiap komponen fokus pada tugasnya

## Manual Dependency Injection

### 1. Constructor Injection

```go
// interfaces/user_service.go
package interfaces

type UserRepositoryInterface interface {
    Create(user *models.User) error
    GetByID(id uint) (*models.User, error)
}

type EmailServiceInterface interface {
    SendWelcomeEmail(email, name string) error
}

// app/services/user_service.go
package services

import (
    "your-app/app/models"
    "your-app/interfaces"
)

type UserService struct {
    userRepo     interfaces.UserRepositoryInterface
    emailService interfaces.EmailServiceInterface
}

// Constructor injection
func NewUserService(
    userRepo interfaces.UserRepositoryInterface,
    emailService interfaces.EmailServiceInterface,
) *UserService {
    return &UserService{
        userRepo:     userRepo,
        emailService: emailService,
    }
}

func (s *UserService) CreateUser(data map[string]interface{}) (*models.User, error) {
    user := &models.User{
        Name:  data["name"].(string),
        Email: data["email"].(string),
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    // Send welcome email
    go s.emailService.SendWelcomeEmail(user.Email, user.Name)
    
    return user, nil
}
```

### 2. Setter Injection

```go
type UserService struct {
    userRepo     interfaces.UserRepositoryInterface
    emailService interfaces.EmailServiceInterface
}

func NewUserService() *UserService {
    return &UserService{}
}

// Setter injection
func (s *UserService) SetUserRepository(repo interfaces.UserRepositoryInterface) {
    s.userRepo = repo
}

func (s *UserService) SetEmailService(service interfaces.EmailServiceInterface) {
    s.emailService = service
}
```

### 3. Interface Injection

```go
type UserRepositoryInjector interface {
    InjectUserRepository(repo interfaces.UserRepositoryInterface)
}

type EmailServiceInjector interface {
    InjectEmailService(service interfaces.EmailServiceInterface)
}

type UserService struct {
    userRepo     interfaces.UserRepositoryInterface
    emailService interfaces.EmailServiceInterface
}

func (s *UserService) InjectUserRepository(repo interfaces.UserRepositoryInterface) {
    s.userRepo = repo
}

func (s *UserService) InjectEmailService(service interfaces.EmailServiceInterface) {
    s.emailService = service
}
```

## DI Container Implementation

### 1. Simple Container

```go
// container/container.go
package container

import (
    "fmt"
    "reflect"
    "sync"
)

type Container struct {
    services map[string]interface{}
    mu       sync.RWMutex
}

func NewContainer() *Container {
    return &Container{
        services: make(map[string]interface{}),
    }
}

// Register a service
func (c *Container) Register(name string, service interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.services[name] = service
}

// Get a service
func (c *Container) Get(name string) (interface{}, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    service, exists := c.services[name]
    if !exists {
        return nil, fmt.Errorf("service %s not found", name)
    }
    
    return service, nil
}

// Resolve dependencies using reflection
func (c *Container) Resolve(target interface{}) error {
    targetValue := reflect.ValueOf(target)
    if targetValue.Kind() != reflect.Ptr {
        return fmt.Errorf("target must be a pointer")
    }
    
    targetValue = targetValue.Elem()
    targetType := targetValue.Type()
    
    for i := 0; i < targetValue.NumField(); i++ {
        field := targetValue.Field(i)
        fieldType := targetType.Field(i)
        
        // Check if field has dependency tag
        if tag := fieldType.Tag.Get("inject"); tag != "" {
            service, err := c.Get(tag)
            if err != nil {
                return err
            }
            
            if field.CanSet() {
                field.Set(reflect.ValueOf(service))
            }
        }
    }
    
    return nil
}
```

### 2. Advanced Container with Factories

```go
// container/advanced_container.go
package container

import (
    "fmt"
    "reflect"
    "sync"
)

type ServiceFactory func(*Container) (interface{}, error)

type AdvancedContainer struct {
    services   map[string]interface{}
    factories  map[string]ServiceFactory
    singletons map[string]interface{}
    mu         sync.RWMutex
}

func NewAdvancedContainer() *AdvancedContainer {
    return &AdvancedContainer{
        services:   make(map[string]interface{}),
        factories:  make(map[string]ServiceFactory),
        singletons: make(map[string]interface{}),
    }
}

// Register singleton service
func (c *AdvancedContainer) RegisterSingleton(name string, factory ServiceFactory) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.factories[name] = factory
}

// Register transient service
func (c *AdvancedContainer) RegisterTransient(name string, factory ServiceFactory) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.factories[name] = factory
}

// Register instance
func (c *AdvancedContainer) RegisterInstance(name string, instance interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.services[name] = instance
}

// Get service
func (c *AdvancedContainer) Get(name string) (interface{}, error) {
    c.mu.RLock()
    
    // Check for registered instance
    if service, exists := c.services[name]; exists {
        c.mu.RUnlock()
        return service, nil
    }
    
    // Check for singleton
    if singleton, exists := c.singletons[name]; exists {
        c.mu.RUnlock()
        return singleton, nil
    }
    
    // Check for factory
    factory, exists := c.factories[name]
    c.mu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("service %s not found", name)
    }
    
    // Create service using factory
    service, err := factory(c)
    if err != nil {
        return nil, err
    }
    
    // Store as singleton if needed
    c.mu.Lock()
    c.singletons[name] = service
    c.mu.Unlock()
    
    return service, nil
}

// Get with type assertion
func Get[T any](c *AdvancedContainer, name string) (T, error) {
    var zero T
    service, err := c.Get(name)
    if err != nil {
        return zero, err
    }
    
    if typed, ok := service.(T); ok {
        return typed, nil
    }
    
    return zero, fmt.Errorf("service %s is not of type %T", name, zero)
}
```

## Service Registration

### 1. Bootstrap Container

```go
// bootstrap/container.go
package bootstrap

import (
    "your-app/app/repositories"
    "your-app/app/services"
    "your-app/container"
    "your-app/external"
    "your-app/facades"
    "your-app/interfaces"
)

func RegisterServices(c *container.AdvancedContainer) {
    // Register repositories
    c.RegisterSingleton("user_repository", func(c *container.AdvancedContainer) (interface{}, error) {
        return repositories.NewUserRepository(), nil
    })
    
    c.RegisterSingleton("product_repository", func(c *container.AdvancedContainer) (interface{}, error) {
        return repositories.NewProductRepository(), nil
    })
    
    // Register external services
    c.RegisterSingleton("email_service", func(c *container.AdvancedContainer) (interface{}, error) {
        return external.NewEmailService(), nil
    })
    
    // Register business services
    c.RegisterTransient("user_service", func(c *container.AdvancedContainer) (interface{}, error) {
        userRepo, err := container.Get[interfaces.UserRepositoryInterface](c, "user_repository")
        if err != nil {
            return nil, err
        }
        
        emailService, err := container.Get[interfaces.EmailServiceInterface](c, "email_service")
        if err != nil {
            return nil, err
        }
        
        return services.NewUserService(userRepo, emailService), nil
    })
    
    c.RegisterTransient("product_service", func(c *container.AdvancedContainer) (interface{}, error) {
        productRepo, err := container.Get[interfaces.ProductRepositoryInterface](c, "product_repository")
        if err != nil {
            return nil, err
        }
        
        return services.NewProductService(productRepo), nil
    })
}
```

### 2. Controller with DI

```go
// app/controllers/user_controller.go
package controllers

import (
    "net/http"
    "strconv"
    "your-app/app/handlers"
    "your-app/container"
    "your-app/interfaces"
    "github.com/gin-gonic/gin"
)

type UserController struct {
    userService interfaces.UserServiceInterface
}

func NewUserController(c *container.AdvancedContainer) (*UserController, error) {
    userService, err := container.Get[interfaces.UserServiceInterface](c, "user_service")
    if err != nil {
        return nil, err
    }
    
    return &UserController{
        userService: userService,
    }, nil
}

func (uc *UserController) Index(c *gin.Context) {
    users, err := uc.userService.GetUsers(nil)
    if err != nil {
        handlers.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users", err)
        return
    }
    
    handlers.SuccessResponse(c, "Users retrieved successfully", users)
}
```

## Wire Integration

[Wire](https://github.com/google/wire) adalah dependency injection framework dari Google yang menggunakan code generation.

### 1. Install Wire

```bash
go install github.com/google/wire/cmd/wire@latest
```

### 2. Define Providers

```go
// wire/providers.go
//go:build wireinject

package wire

import (
    "your-app/app/repositories"
    "your-app/app/services"
    "your-app/external"
    "your-app/interfaces"
    "github.com/google/wire"
)

// Provider sets
var RepositorySet = wire.NewSet(
    repositories.NewUserRepository,
    wire.Bind(new(interfaces.UserRepositoryInterface), new(*repositories.UserRepository)),
    
    repositories.NewProductRepository,
    wire.Bind(new(interfaces.ProductRepositoryInterface), new(*repositories.ProductRepository)),
)

var ServiceSet = wire.NewSet(
    services.NewUserService,
    wire.Bind(new(interfaces.UserServiceInterface), new(*services.UserService)),
    
    services.NewProductService,
    wire.Bind(new(interfaces.ProductServiceInterface), new(*services.ProductService)),
    
    external.NewEmailService,
    wire.Bind(new(interfaces.EmailServiceInterface), new(*external.EmailService)),
)

var ControllerSet = wire.NewSet(
    controllers.NewUserController,
    controllers.NewProductController,
)
```

### 3. Wire Injectors

```go
// wire/injector.go
//go:build wireinject

package wire

import (
    "your-app/app/controllers"
    "github.com/google/wire"
)

func InitializeUserController() (*controllers.UserController, error) {
    wire.Build(RepositorySet, ServiceSet, controllers.NewUserController)
    return &controllers.UserController{}, nil
}

func InitializeProductController() (*controllers.ProductController, error) {
    wire.Build(RepositorySet, ServiceSet, controllers.NewProductController)
    return &controllers.ProductController{}, nil
}

func InitializeApplication() (*Application, error) {
    wire.Build(
        RepositorySet,
        ServiceSet,
        ControllerSet,
        NewApplication,
    )
    return &Application{}, nil
}
```

### 4. Generate Wire Code

```bash
# Generate wire code
cd wire
wire
```

## Testing with DI

### 1. Mock Dependencies

```go
// mocks/user_service_mock.go
package mocks

import (
    "your-app/app/models"
    "github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
    mock.Mock
}

func (m *UserServiceMock) CreateUser(data map[string]interface{}) (*models.User, error) {
    args := m.Called(data)
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserServiceMock) GetUser(id uint) (*models.User, error) {
    args := m.Called(id)
    return args.Get(0).(*models.User), args.Error(1)
}
```

### 2. Test with DI

```go
// app/controllers/user_controller_test.go
package controllers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "your-app/app/models"
    "your-app/mocks"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestUserController_Store(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    mockUserService := new(mocks.UserServiceMock)
    
    controller := &UserController{
        userService: mockUserService,
    }
    
    // Test data
    userData := map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "password123",
    }
    
    expectedUser := &models.User{
        ID:    1,
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Mock expectations
    mockUserService.On("CreateUser", mock.MatchedBy(func(data map[string]interface{}) bool {
        return data["name"] == "John Doe" && data["email"] == "john@example.com"
    })).Return(expectedUser, nil)
    
    // Setup HTTP request
    jsonData, _ := json.Marshal(userData)
    req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req
    
    // Execute
    controller.Store(c)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    mockUserService.AssertExpectations(t)
}
```

### 3. Test Container

```go
// container/test_container.go
package container

import (
    "your-app/mocks"
    "your-app/interfaces"
)

func NewTestContainer() *AdvancedContainer {
    c := NewAdvancedContainer()
    
    // Register mocks
    c.RegisterInstance("user_repository", new(mocks.UserRepositoryMock))
    c.RegisterInstance("email_service", new(mocks.EmailServiceMock))
    
    // Register services with mocks
    c.RegisterTransient("user_service", func(c *AdvancedContainer) (interface{}, error) {
        userRepo, _ := Get[interfaces.UserRepositoryInterface](c, "user_repository")
        emailService, _ := Get[interfaces.EmailServiceInterface](c, "email_service")
        return services.NewUserService(userRepo, emailService), nil
    })
    
    return c
}
```

## Middleware for DI

### 1. Container Middleware

```go
// middleware/container_middleware.go
package middleware

import (
    "your-app/container"
    "github.com/gin-gonic/gin"
)

func ContainerMiddleware(c *container.AdvancedContainer) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        ctx.Set("container", c)
        ctx.Next()
    }
}

// Helper function to get container from context
func GetContainer(c *gin.Context) *container.AdvancedContainer {
    return c.MustGet("container").(*container.AdvancedContainer)
}
```

### 2. Usage in Routes

```go
// routes/web.go
func RegisterRoutes(router *gin.Engine, c *container.AdvancedContainer) {
    // Add container middleware
    router.Use(middleware.ContainerMiddleware(c))
    
    api := router.Group("/api")
    {
        users := api.Group("/users")
        {
            users.GET("", func(ctx *gin.Context) {
                container := middleware.GetContainer(ctx)
                controller, _ := NewUserController(container)
                controller.Index(ctx)
            })
        }
    }
}
```

## Best Practices

### 1. Interface Segregation

```go
// Good: Small, focused interfaces
type UserReader interface {
    GetByID(id uint) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
}

type UserWriter interface {
    Create(user *models.User) error
    Update(id uint, data map[string]interface{}) error
}

// Bad: Large interface
type UserRepository interface {
    Create(user *models.User) error
    GetByID(id uint) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    Update(id uint, data map[string]interface{}) error
    Delete(id uint) error
    GetWithRoles(id uint) (*models.User, error)
    GetActiveUsers() ([]models.User, error)
    // ... many more methods
}
```

### 2. Avoid God Objects

```go
// Good: Focused services
type UserService struct {
    userRepo interfaces.UserRepositoryInterface
}

type EmailService struct {
    client EmailClient
}

// Bad: God object
type ApplicationService struct {
    userRepo    UserRepository
    productRepo ProductRepository
    emailClient EmailClient
    fileStorage FileStorage
    cache       Cache
    logger      Logger
    // ... many more dependencies
}
```

### 3. Lazy Initialization

```go
type LazyService struct {
    factory func() SomeService
    service SomeService
    once    sync.Once
}

func (ls *LazyService) GetService() SomeService {
    ls.once.Do(func() {
        ls.service = ls.factory()
    })
    return ls.service
}
```

### 4. Lifecycle Management

```go
type ServiceWithLifecycle interface {
    Start() error
    Stop() error
}

type Container struct {
    services []ServiceWithLifecycle
}

func (c *Container) StartAll() error {
    for _, service := range c.services {
        if err := service.Start(); err != nil {
            return err
        }
    }
    return nil
}

func (c *Container) StopAll() error {
    for i := len(c.services) - 1; i >= 0; i-- {
        if err := c.services[i].Stop(); err != nil {
            return err
        }
    }
    return nil
}
```

---

Untuk informasi lebih lanjut, lihat:
- [Service Layer Guide](service-layer.md)
- [Repository Pattern](repository-pattern.md)
- [Development Guide](development.md)
