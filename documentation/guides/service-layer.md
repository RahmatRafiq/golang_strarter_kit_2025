# Service Layer

Implementasi Service Layer pattern pada Golang Starter Kit untuk mengorganisir business logic.

## Overview

Service Layer adalah pattern yang digunakan untuk mengorganisir business logic aplikasi. Layer ini bertindak sebagai jembatan antara Controllers dan Repository/Model, serta bertanggung jawab untuk:

- Business logic dan rules
- Data validation dan transformation
- Orchestration of multiple repositories
- Transaction management
- External service integration

## Architecture

```
Controller → Service → Repository → Database
    ↓         ↓           ↓
 HTTP      Business    Data Access
Request    Logic       Layer
```

## Basic Implementation

### 1. Service Interface

```go
// interfaces/user_service.go
package interfaces

import "your-app/app/models"

type UserServiceInterface interface {
    CreateUser(data map[string]interface{}) (*models.User, error)
    GetUser(id uint) (*models.User, error)
    GetUsers(filters map[string]interface{}) ([]models.User, error)
    UpdateUser(id uint, data map[string]interface{}) (*models.User, error)
    DeleteUser(id uint) error
    ChangePassword(id uint, oldPassword, newPassword string) error
    AssignRole(userID, roleID uint) error
}
```

### 2. Service Implementation

```go
// app/services/user_service.go
package services

import (
    "errors"
    "fmt"
    "your-app/app/models"
    "your-app/app/repositories"
    "your-app/app/helpers"
    "your-app/interfaces"
    "gorm.io/gorm"
)

type UserService struct {
    userRepo interfaces.UserRepositoryInterface
    roleRepo interfaces.RoleRepositoryInterface
    db       *gorm.DB
}

func NewUserService() interfaces.UserServiceInterface {
    return &UserService{
        userRepo: repositories.NewUserRepository(),
        roleRepo: repositories.NewRoleRepository(),
        db:       facades.DB(),
    }
}

func (s *UserService) CreateUser(data map[string]interface{}) (*models.User, error) {
    // Validate business rules
    if err := s.validateUserData(data); err != nil {
        return nil, err
    }
    
    // Check if email already exists
    existingUser, _ := s.userRepo.GetByEmail(data["email"].(string))
    if existingUser != nil {
        return nil, errors.New("email already exists")
    }
    
    // Hash password
    hashedPassword, err := helpers.HashPassword(data["password"].(string))
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Create user
    user := &models.User{
        Name:     data["name"].(string),
        Email:    data["email"].(string),
        Password: hashedPassword,
        Active:   true,
    }
    
    // Start transaction
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Create user in transaction
    if err := s.userRepo.CreateWithTx(tx, user); err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // Assign default role if provided
    if roleID, exists := data["role_id"]; exists {
        if err := s.assignRoleInTx(tx, user.ID, uint(roleID.(float64))); err != nil {
            tx.Rollback()
            return nil, fmt.Errorf("failed to assign role: %w", err)
        }
    }
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return user, nil
}

func (s *UserService) GetUser(id uint) (*models.User, error) {
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    return user, nil
}

func (s *UserService) GetUsers(filters map[string]interface{}) ([]models.User, error) {
    // Apply business logic filters
    processedFilters := s.processUserFilters(filters)
    return s.userRepo.GetAll(processedFilters)
}

func (s *UserService) UpdateUser(id uint, data map[string]interface{}) (*models.User, error) {
    // Get existing user
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    
    // Validate update data
    if err := s.validateUpdateData(data, user); err != nil {
        return nil, err
    }
    
    // Update user
    if err := s.userRepo.Update(id, data); err != nil {
        return nil, fmt.Errorf("failed to update user: %w", err)
    }
    
    // Return updated user
    return s.userRepo.GetByID(id)
}

func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
    // Get user
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return fmt.Errorf("user not found: %w", err)
    }
    
    // Verify old password
    if !helpers.CheckPassword(oldPassword, user.Password) {
        return errors.New("old password is incorrect")
    }
    
    // Validate new password
    if err := s.validatePassword(newPassword); err != nil {
        return err
    }
    
    // Hash new password
    hashedPassword, err := helpers.HashPassword(newPassword)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Update password
    updateData := map[string]interface{}{
        "password": hashedPassword,
    }
    
    return s.userRepo.Update(id, updateData)
}

// Private helper methods
func (s *UserService) validateUserData(data map[string]interface{}) error {
    if data["name"] == nil || data["name"].(string) == "" {
        return errors.New("name is required")
    }
    
    if data["email"] == nil || data["email"].(string) == "" {
        return errors.New("email is required")
    }
    
    if data["password"] == nil || data["password"].(string) == "" {
        return errors.New("password is required")
    }
    
    return s.validatePassword(data["password"].(string))
}

func (s *UserService) validatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    // Add more validation rules as needed
    return nil
}
```

## Advanced Service Patterns

### 1. Service with Event Handling

```go
// app/services/user_service_with_events.go
package services

import (
    "your-app/app/events"
    "your-app/app/models"
)

type UserServiceWithEvents struct {
    UserService
    eventDispatcher *events.Dispatcher
}

func NewUserServiceWithEvents() *UserServiceWithEvents {
    return &UserServiceWithEvents{
        UserService:     *NewUserService().(*UserService),
        eventDispatcher: events.NewDispatcher(),
    }
}

func (s *UserServiceWithEvents) CreateUser(data map[string]interface{}) (*models.User, error) {
    // Create user using parent method
    user, err := s.UserService.CreateUser(data)
    if err != nil {
        return nil, err
    }
    
    // Dispatch user created event
    s.eventDispatcher.Dispatch(&events.UserCreated{
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.Name,
    })
    
    return user, nil
}
```

### 2. Service with Caching

```go
// app/services/cached_user_service.go
package services

import (
    "fmt"
    "time"
    "your-app/app/models"
    "your-app/cache"
)

type CachedUserService struct {
    UserService
    cache cache.CacheInterface
}

func NewCachedUserService() *CachedUserService {
    return &CachedUserService{
        UserService: *NewUserService().(*UserService),
        cache:       cache.NewRedisCache(),
    }
}

func (s *CachedUserService) GetUser(id uint) (*models.User, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("user:%d", id)
    if cached := s.cache.Get(cacheKey); cached != nil {
        return cached.(*models.User), nil
    }
    
    // Get from database
    user, err := s.UserService.GetUser(id)
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    s.cache.Set(cacheKey, user, 5*time.Minute)
    
    return user, nil
}

func (s *CachedUserService) UpdateUser(id uint, data map[string]interface{}) (*models.User, error) {
    // Update user
    user, err := s.UserService.UpdateUser(id, data)
    if err != nil {
        return nil, err
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("user:%d", id)
    s.cache.Delete(cacheKey)
    
    return user, nil
}
```

### 3. Service with External API Integration

```go
// app/services/user_notification_service.go
package services

import (
    "your-app/app/models"
    "your-app/external"
)

type UserNotificationService struct {
    UserService
    emailService external.EmailServiceInterface
    smsService   external.SMSServiceInterface
}

func NewUserNotificationService() *UserNotificationService {
    return &UserNotificationService{
        UserService:  *NewUserService().(*UserService),
        emailService: external.NewEmailService(),
        smsService:   external.NewSMSService(),
    }
}

func (s *UserNotificationService) CreateUser(data map[string]interface{}) (*models.User, error) {
    // Create user
    user, err := s.UserService.CreateUser(data)
    if err != nil {
        return nil, err
    }
    
    // Send welcome email
    go s.sendWelcomeEmail(user)
    
    return user, nil
}

func (s *UserNotificationService) sendWelcomeEmail(user *models.User) {
    emailData := map[string]interface{}{
        "name":  user.Name,
        "email": user.Email,
    }
    
    s.emailService.SendTemplate("welcome", user.Email, emailData)
}
```

## Service Composition

### 1. Service Container

```go
// services/container.go
package services

import "your-app/interfaces"

type Container struct {
    userService     interfaces.UserServiceInterface
    productService  interfaces.ProductServiceInterface
    categoryService interfaces.CategoryServiceInterface
}

func NewContainer() *Container {
    return &Container{
        userService:     NewUserService(),
        productService:  NewProductService(),
        categoryService: NewCategoryService(),
    }
}

func (c *Container) UserService() interfaces.UserServiceInterface {
    return c.userService
}

func (c *Container) ProductService() interfaces.ProductServiceInterface {
    return c.productService
}

func (c *Container) CategoryService() interfaces.CategoryServiceInterface {
    return c.categoryService
}
```

### 2. Service Factory

```go
// factories/service_factory.go
package factories

import (
    "your-app/app/services"
    "your-app/interfaces"
)

type ServiceFactory struct{}

func NewServiceFactory() *ServiceFactory {
    return &ServiceFactory{}
}

func (f *ServiceFactory) CreateUserService(serviceType string) interfaces.UserServiceInterface {
    switch serviceType {
    case "cached":
        return services.NewCachedUserService()
    case "with_events":
        return services.NewUserServiceWithEvents()
    case "with_notifications":
        return services.NewUserNotificationService()
    default:
        return services.NewUserService()
    }
}
```

## Service Testing

### 1. Unit Testing

```go
// app/services/user_service_test.go
package services

import (
    "testing"
    "your-app/app/models"
    "your-app/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestUserService_CreateUser(t *testing.T) {
    // Setup mocks
    mockUserRepo := new(mocks.UserRepositoryMock)
    mockRoleRepo := new(mocks.RoleRepositoryMock)
    mockDB := new(mocks.DBMock)
    
    service := &UserService{
        userRepo: mockUserRepo,
        roleRepo: mockRoleRepo,
        db:       mockDB,
    }
    
    // Test data
    userData := map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "password123",
    }
    
    // Mock expectations
    mockUserRepo.On("GetByEmail", "john@example.com").Return(nil, gorm.ErrRecordNotFound)
    mockDB.On("Begin").Return(mockDB)
    mockUserRepo.On("CreateWithTx", mockDB, mock.AnythingOfType("*models.User")).Return(nil)
    mockDB.On("Commit").Return(mockDB)
    mockDB.On("Error").Return(nil)
    
    // Execute
    user, err := service.CreateUser(userData)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "John Doe", user.Name)
    assert.Equal(t, "john@example.com", user.Email)
    
    // Verify mocks
    mockUserRepo.AssertExpectations(t)
    mockDB.AssertExpectations(t)
}

func TestUserService_CreateUser_EmailExists(t *testing.T) {
    // Setup mocks
    mockUserRepo := new(mocks.UserRepositoryMock)
    
    service := &UserService{
        userRepo: mockUserRepo,
    }
    
    // Test data
    userData := map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "password123",
    }
    
    existingUser := &models.User{
        ID:    1,
        Email: "john@example.com",
    }
    
    // Mock expectations
    mockUserRepo.On("GetByEmail", "john@example.com").Return(existingUser, nil)
    
    // Execute
    user, err := service.CreateUser(userData)
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "email already exists")
    
    // Verify mocks
    mockUserRepo.AssertExpectations(t)
}
```

### 2. Integration Testing

```go
// app/services/user_service_integration_test.go
package services

import (
    "testing"
    "your-app/database"
    "your-app/test"
    "github.com/stretchr/testify/assert"
)

func TestUserService_Integration(t *testing.T) {
    // Setup test database
    testDB := test.SetupTestDB()
    defer test.CleanupTestDB(testDB)
    
    // Create service with real dependencies
    service := NewUserService()
    
    t.Run("CreateUser", func(t *testing.T) {
        userData := map[string]interface{}{
            "name":     "Integration Test User",
            "email":    "integration@example.com",
            "password": "password123",
        }
        
        user, err := service.CreateUser(userData)
        
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, "Integration Test User", user.Name)
        assert.NotEmpty(t, user.Password) // Should be hashed
    })
    
    t.Run("GetUser", func(t *testing.T) {
        // Create user first
        userData := map[string]interface{}{
            "name":     "Get Test User",
            "email":    "get@example.com",
            "password": "password123",
        }
        
        createdUser, _ := service.CreateUser(userData)
        
        // Get user
        user, err := service.GetUser(createdUser.ID)
        
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, createdUser.ID, user.ID)
    })
}
```

## Service Decorators

### 1. Logging Decorator

```go
// app/services/decorators/logging_decorator.go
package decorators

import (
    "log"
    "time"
    "your-app/app/models"
    "your-app/interfaces"
)

type LoggingUserService struct {
    service interfaces.UserServiceInterface
}

func NewLoggingUserService(service interfaces.UserServiceInterface) interfaces.UserServiceInterface {
    return &LoggingUserService{service: service}
}

func (s *LoggingUserService) CreateUser(data map[string]interface{}) (*models.User, error) {
    start := time.Now()
    log.Printf("Creating user with email: %s", data["email"])
    
    user, err := s.service.CreateUser(data)
    
    duration := time.Since(start)
    if err != nil {
        log.Printf("Failed to create user: %v (took %v)", err, duration)
    } else {
        log.Printf("User created successfully: ID=%d (took %v)", user.ID, duration)
    }
    
    return user, err
}
```

### 2. Validation Decorator

```go
// app/services/decorators/validation_decorator.go
package decorators

import (
    "errors"
    "regexp"
    "your-app/app/models"
    "your-app/interfaces"
)

type ValidatingUserService struct {
    service interfaces.UserServiceInterface
}

func NewValidatingUserService(service interfaces.UserServiceInterface) interfaces.UserServiceInterface {
    return &ValidatingUserService{service: service}
}

func (s *ValidatingUserService) CreateUser(data map[string]interface{}) (*models.User, error) {
    // Enhanced validation
    if err := s.validateEmail(data["email"].(string)); err != nil {
        return nil, err
    }
    
    if err := s.validateName(data["name"].(string)); err != nil {
        return nil, err
    }
    
    return s.service.CreateUser(data)
}

func (s *ValidatingUserService) validateEmail(email string) error {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(email) {
        return errors.New("invalid email format")
    }
    return nil
}

func (s *ValidatingUserService) validateName(name string) error {
    if len(name) < 2 {
        return errors.New("name must be at least 2 characters")
    }
    if len(name) > 50 {
        return errors.New("name must be less than 50 characters")
    }
    return nil
}
```

## Best Practices

### 1. Single Responsibility

```go
// Good: Service focused on user operations
type UserService struct {
    userRepo interfaces.UserRepositoryInterface
}

// Bad: Service doing too many things
type UserEverythingService struct {
    userRepo    interfaces.UserRepositoryInterface
    emailSender EmailSender
    fileSaver   FileSaver
    logger      Logger
    cache       Cache
}
```

### 2. Dependency Injection

```go
// Good: Dependencies injected
func NewUserService(userRepo interfaces.UserRepositoryInterface) *UserService {
    return &UserService{userRepo: userRepo}
}

// Bad: Dependencies hardcoded
func NewUserService() *UserService {
    return &UserService{
        userRepo: repositories.NewUserRepository(), // hardcoded
    }
}
```

### 3. Error Handling

```go
func (s *UserService) CreateUser(data map[string]interface{}) (*models.User, error) {
    // Validate input
    if err := s.validateUserData(data); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // Business logic
    user, err := s.userRepo.Create(user)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}
```

### 4. Transaction Management

```go
func (s *UserService) CreateUserWithProfile(userData, profileData map[string]interface{}) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // Create user
        user := &models.User{...}
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        
        // Create profile
        profile := &models.Profile{
            UserID: user.ID,
            ...
        }
        if err := tx.Create(profile).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

---

Untuk informasi lebih lanjut, lihat:
- [Repository Pattern](repository-pattern.md)
- [Dependency Injection](dependency-injection.md)
- [Development Guide](development.md)
