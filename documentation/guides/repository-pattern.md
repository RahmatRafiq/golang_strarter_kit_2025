# Repository Pattern

Implementasi Repository Pattern pada Golang Starter Kit untuk abstraksi data access layer.

## Overview

Repository Pattern adalah design pattern yang digunakan untuk memisahkan business logic dari data access logic. Pattern ini memberikan interface yang konsisten untuk mengakses data, tidak peduli apakah data tersebut berasal dari database, file, atau sumber lainnya.

## Benefits

- **Separation of Concerns**: Memisahkan business logic dari data access
- **Testability**: Mudah untuk membuat mock repository untuk testing
- **Flexibility**: Mudah mengganti implementasi data source
- **Consistency**: Interface yang konsisten untuk semua data operations
- **Maintainability**: Code lebih mudah dirawat dan dikembangkan

## Basic Implementation

### 1. Repository Interface

```go
// interfaces/user_repository.go
package interfaces

import "your-app/app/models"

type UserRepositoryInterface interface {
    Create(user *models.User) error
    GetByID(id uint) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    GetAll(filters map[string]interface{}) ([]models.User, error)
    Update(id uint, data map[string]interface{}) error
    Delete(id uint) error
    Count(filters map[string]interface{}) (int64, error)
}
```

### 2. Repository Implementation

```go
// app/repositories/user_repository.go
package repositories

import (
    "your-app/app/models"
    "your-app/facades"
    "your-app/interfaces"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository() interfaces.UserRepositoryInterface {
    return &UserRepository{
        db: facades.DB(),
    }
}

func (r *UserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
    var user models.User
    err := r.db.Preload("Roles").First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) GetAll(filters map[string]interface{}) ([]models.User, error) {
    var users []models.User
    query := r.db.Preload("Roles")
    
    // Apply filters
    for key, value := range filters {
        switch key {
        case "name":
            query = query.Where("name LIKE ?", "%"+value.(string)+"%")
        case "email":
            query = query.Where("email LIKE ?", "%"+value.(string)+"%")
        case "active":
            query = query.Where("active = ?", value)
        }
    }
    
    err := query.Find(&users).Error
    return users, err
}

func (r *UserRepository) Update(id uint, data map[string]interface{}) error {
    return r.db.Model(&models.User{}).Where("id = ?", id).Updates(data).Error
}

func (r *UserRepository) Delete(id uint) error {
    return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) Count(filters map[string]interface{}) (int64, error) {
    var count int64
    query := r.db.Model(&models.User{})
    
    // Apply filters
    for key, value := range filters {
        switch key {
        case "active":
            query = query.Where("active = ?", value)
        }
    }
    
    err := query.Count(&count).Error
    return count, err
}
```

### 3. Service Using Repository

```go
// app/services/user_service.go
package services

import (
    "errors"
    "your-app/app/models"
    "your-app/app/repositories"
    "your-app/interfaces"
)

type UserService struct {
    userRepo interfaces.UserRepositoryInterface
}

func NewUserService() *UserService {
    return &UserService{
        userRepo: repositories.NewUserRepository(),
    }
}

// Dependency injection untuk testing
func NewUserServiceWithRepo(repo interfaces.UserRepositoryInterface) *UserService {
    return &UserService{
        userRepo: repo,
    }
}

func (s *UserService) CreateUser(data map[string]interface{}) (*models.User, error) {
    // Business logic validation
    if data["email"] == "" {
        return nil, errors.New("email is required")
    }
    
    // Check if email exists
    existingUser, _ := s.userRepo.GetByEmail(data["email"].(string))
    if existingUser != nil {
        return nil, errors.New("email already exists")
    }
    
    user := &models.User{
        Name:     data["name"].(string),
        Email:    data["email"].(string),
        Password: data["password"].(string), // Hash this in real implementation
        Active:   true,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *UserService) GetUsers(filters map[string]interface{}) ([]models.User, error) {
    return s.userRepo.GetAll(filters)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
    return s.userRepo.GetByID(id)
}
```

## Advanced Patterns

### 1. Generic Repository

```go
// interfaces/base_repository.go
package interfaces

type BaseRepositoryInterface[T any] interface {
    Create(entity *T) error
    GetByID(id uint) (*T, error)
    GetAll(filters map[string]interface{}) ([]T, error)
    Update(id uint, data map[string]interface{}) error
    Delete(id uint) error
    Count(filters map[string]interface{}) (int64, error)
}

// repositories/base_repository.go
package repositories

import (
    "your-app/interfaces"
    "gorm.io/gorm"
)

type BaseRepository[T any] struct {
    db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) interfaces.BaseRepositoryInterface[T] {
    return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(entity *T) error {
    return r.db.Create(entity).Error
}

func (r *BaseRepository[T]) GetByID(id uint) (*T, error) {
    var entity T
    err := r.db.First(&entity, id).Error
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

func (r *BaseRepository[T]) GetAll(filters map[string]interface{}) ([]T, error) {
    var entities []T
    query := r.db
    
    // Apply basic filters
    for key, value := range filters {
        query = query.Where(key+" = ?", value)
    }
    
    err := query.Find(&entities).Error
    return entities, err
}
```

### 2. Repository with Scopes

```go
// repositories/user_repository.go
func (r *UserRepository) WithScopes(scopes ...func(*gorm.DB) *gorm.DB) interfaces.UserRepositoryInterface {
    newRepo := *r
    newRepo.db = r.db.Scopes(scopes...)
    return &newRepo
}

// Usage
activeUsers := userRepo.WithScopes(scopes.ActiveUsers).GetAll(nil)
```

### 3. Multi-Database Repository

```go
// repositories/multi_db_user_repository.go
package repositories

import (
    "your-app/database"
    "your-app/interfaces"
)

type MultiDBUserRepository struct {
    dbManager *database.Manager
}

func NewMultiDBUserRepository() interfaces.UserRepositoryInterface {
    return &MultiDBUserRepository{
        dbManager: database.GetManager(),
    }
}

func (r *MultiDBUserRepository) Create(user *models.User) error {
    // Use primary database for writes
    return r.dbManager.GetConnection("mysql").Create(user).Error
}

func (r *MultiDBUserRepository) GetByID(id uint) (*models.User, error) {
    // Use read replica for reads
    var user models.User
    err := r.dbManager.GetConnection("mysql_secondary").First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

## Repository Factory

```go
// factories/repository_factory.go
package factories

import (
    "your-app/app/repositories"
    "your-app/interfaces"
)

type RepositoryFactory struct{}

func NewRepositoryFactory() *RepositoryFactory {
    return &RepositoryFactory{}
}

func (f *RepositoryFactory) CreateUserRepository() interfaces.UserRepositoryInterface {
    return repositories.NewUserRepository()
}

func (f *RepositoryFactory) CreateProductRepository() interfaces.ProductRepositoryInterface {
    return repositories.NewProductRepository()
}

func (f *RepositoryFactory) CreateCategoryRepository() interfaces.CategoryRepositoryInterface {
    return repositories.NewCategoryRepository()
}
```

## Unit Testing

### 1. Mock Repository

```go
// mocks/user_repository_mock.go
package mocks

import (
    "your-app/app/models"
    "github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
    mock.Mock
}

func (m *UserRepositoryMock) Create(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *UserRepositoryMock) GetByID(id uint) (*models.User, error) {
    args := m.Called(id)
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByEmail(email string) (*models.User, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}
```

### 2. Service Test

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
    // Setup
    mockRepo := new(mocks.UserRepositoryMock)
    service := NewUserServiceWithRepo(mockRepo)
    
    // Mock expectations
    mockRepo.On("GetByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)
    mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
    
    // Test data
    userData := map[string]interface{}{
        "name":     "Test User",
        "email":    "test@example.com",
        "password": "password123",
    }
    
    // Execute
    user, err := service.CreateUser(userData)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "Test User", user.Name)
    assert.Equal(t, "test@example.com", user.Email)
    
    // Verify mock calls
    mockRepo.AssertExpectations(t)
}
```

## Query Builder Integration

```go
// repositories/advanced_user_repository.go
package repositories

type AdvancedUserRepository struct {
    UserRepository
}

func (r *AdvancedUserRepository) GetUsersWithPagination(page, limit int, filters map[string]interface{}) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    query := r.db.Model(&models.User{})
    
    // Apply filters
    for key, value := range filters {
        switch key {
        case "search":
            query = query.Where("name LIKE ? OR email LIKE ?", "%"+value.(string)+"%", "%"+value.(string)+"%")
        case "role":
            query = query.Joins("JOIN user_has_roles ON users.id = user_has_roles.user_id").
                         Joins("JOIN roles ON user_has_roles.role_id = roles.id").
                         Where("roles.name = ?", value)
        }
    }
    
    // Get total count
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Get paginated data
    offset := (page - 1) * limit
    err := query.Preload("Roles").Offset(offset).Limit(limit).Find(&users).Error
    
    return users, total, err
}
```

## Best Practices

### 1. Interface Segregation

```go
// Split large interfaces into smaller ones
type UserReaderInterface interface {
    GetByID(id uint) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    GetAll(filters map[string]interface{}) ([]models.User, error)
}

type UserWriterInterface interface {
    Create(user *models.User) error
    Update(id uint, data map[string]interface{}) error
    Delete(id uint) error
}

type UserRepositoryInterface interface {
    UserReaderInterface
    UserWriterInterface
}
```

### 2. Error Handling

```go
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
    var user models.User
    err := r.db.First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user with id %d not found", id)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return &user, nil
}
```

### 3. Context Support

```go
import "context"

type UserRepositoryInterface interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uint) (*models.User, error)
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

### 4. Caching Integration

```go
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("user:%d", id)
    if cached := cache.Get(cacheKey); cached != nil {
        return cached.(*models.User), nil
    }
    
    // Get from database
    var user models.User
    err := r.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    cache.Set(cacheKey, &user, 5*time.Minute)
    
    return &user, nil
}
```

---

Untuk informasi lebih lanjut, lihat:
- [Service Layer Guide](service-layer.md)
- [Development Guide](development.md)
- [Architecture Guide](architecture.md)
