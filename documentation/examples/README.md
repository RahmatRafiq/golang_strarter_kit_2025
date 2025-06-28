# Code Examples & Usage

## Overview

Koleksi contoh kode dan implementasi untuk berbagai fitur dalam aplikasi. Examples ini dirancang untuk membantu developer memahami cara menggunakan fitur-fitur yang tersedia.

## 📁 Structure

```
examples/
├── basic-usage.md          # Contoh penggunaan dasar
├── advanced-examples.md    # Contoh lanjutan
├── multi-database.md       # Multiple database usage
├── repository-pattern.md   # Repository pattern examples
├── api-integration.md      # API integration examples
└── code-samples/          # Sample code files
    ├── user_service.go
    ├── product_repository.go
    └── database_sync.go
```

## 🚀 Quick Examples

### 1. Basic Database Operations

```go
package main

import (
    "golang_starter_kit_2025/facades"
    "golang_starter_kit_2025/app/models"
)

func main() {
    // Initialize database
    facades.ConnectDB()
    defer facades.CloseDB()
    
    // Get default connection
    db := facades.GetDB()
    
    // Create user
    user := &models.User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    if err := db.Create(user).Error; err != nil {
        panic(err)
    }
    
    // Find user
    var foundUser models.User
    db.First(&foundUser, user.ID)
}
```

### 2. Multiple Database Usage

```go
package main

import (
    "golang_starter_kit_2025/facades"
    "golang_starter_kit_2025/app/models"
)

func main() {
    facades.ConnectDB()
    defer facades.CloseDB()
    
    // Get MySQL connection
    mysqlConn, err := facades.MySQL()
    if err != nil {
        panic(err)
    }
    
    // Get PostgreSQL connection
    postgresConn, err := facades.PostgreSQL()
    if err != nil {
        panic(err)
    }
    
    // Create user in MySQL
    user := &models.User{Name: "Alice", Email: "alice@example.com"}
    mysqlConn.DB.Create(user)
    
    // Create analytics record in PostgreSQL
    analytics := &models.Analytics{
        UserID: user.ID,
        Action: "user_created",
    }
    postgresConn.DB.Create(analytics)
}
```

### 3. Repository Pattern

```go
// repositories/user_repository.go
type UserRepository struct {
    mysqlDB    *database.Connection
    postgresDB *database.Connection
}

func NewUserRepository() *UserRepository {
    mysql, _ := facades.MySQL()
    postgres, _ := facades.PostgreSQL()
    
    return &UserRepository{
        mysqlDB:    mysql,
        postgresDB: postgres,
    }
}

func (r *UserRepository) CreateUser(user *models.User) error {
    // Create in MySQL
    if err := r.mysqlDB.DB.Create(user).Error; err != nil {
        return err
    }
    
    // Log to PostgreSQL
    analytics := &models.UserAnalytics{
        UserID:    user.ID,
        CreatedAt: time.Now(),
    }
    r.postgresDB.DB.Create(analytics)
    
    return nil
}

func (r *UserRepository) SyncUsers() error {
    var users []models.User
    
    // Get users from MySQL
    if err := r.mysqlDB.DB.Find(&users).Error; err != nil {
        return err
    }
    
    // Sync to PostgreSQL
    for _, user := range users {
        var existingUser models.User
        result := r.postgresDB.DB.Where("id = ?", user.ID).First(&existingUser)
        
        if result.Error == gorm.ErrRecordNotFound {
            r.postgresDB.DB.Create(&user)
        } else {
            r.postgresDB.DB.Save(&user)
        }
    }
    
    return nil
}
```

### 4. Service Layer Implementation

```go
// services/user_service.go
type UserService struct {
    userRepo repositories.UserRepositoryInterface
    emailService *EmailService
}

func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
    return &UserService{
        userRepo: userRepo,
        emailService: NewEmailService(),
    }
}

func (s *UserService) RegisterUser(req *requests.RegisterRequest) (*models.User, error) {
    // Validate request
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Check if email exists
    existingUser, _ := s.userRepo.GetByEmail(req.Email)
    if existingUser != nil {
        return nil, errors.New("email already exists")
    }
    
    // Hash password
    hashedPassword, err := helpers.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    // Create user
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    // Send welcome email
    s.emailService.SendWelcomeEmail(user.Email, user.Name)
    
    return user, nil
}
```

### 5. API Controller Example

```go
// controllers/user_controller.go
type UserController struct {
    userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}

func (c *UserController) Register(ctx *gin.Context) {
    var req requests.RegisterRequest
    
    if err := ctx.ShouldBindJSON(&req); err != nil {
        responses.ValidationErrorResponse(ctx, err)
        return
    }
    
    user, err := c.userService.RegisterUser(&req)
    if err != nil {
        responses.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
        return
    }
    
    responses.SuccessResponse(ctx, "User registered successfully", user)
}

func (c *UserController) GetUsers(ctx *gin.Context) {
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
    
    users, total, err := c.userService.GetUsersWithPagination(page, limit)
    if err != nil {
        responses.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    
    responses.PaginatedResponse(ctx, "Users retrieved successfully", users, page, limit, total)
}
```

### 6. Database Migration Example

```sql
-- app/database/migrations/20250629120000_create_users_table.sql

-- +++ UP Migration
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email_verified_at TIMESTAMP NULL,
    remember_token VARCHAR(100) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- --- DOWN Migration
DROP INDEX idx_users_deleted_at ON users;
DROP INDEX idx_users_email ON users;
DROP TABLE IF EXISTS users;
```

### 7. CLI Command Usage

```bash
# Database management
go run main.go db:connections
go run main.go db:status --connection=mysql
go run main.go db:status --connection=postgres

# Migration commands
go run main.go make:migration create_products_table
go run main.go migrate:all --connection=mysql
go run main.go migrate:all --connection=postgres
go run main.go rollback:batch --batch=1 --connection=mysql

# Seeder commands
go run main.go make:seeder --name=users_seeder
go run main.go db:seed --connection=mysql
```

### 8. Testing Examples

```go
// user_service_test.go
func TestUserService_RegisterUser(t *testing.T) {
    // Setup
    mockRepo := &mocks.UserRepositoryInterface{}
    service := services.NewUserService(mockRepo)
    
    // Mock expectations
    mockRepo.On("GetByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)
    mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
    
    // Test data
    req := &requests.RegisterRequest{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    // Execute
    user, err := service.RegisterUser(req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "Test User", user.Name)
    assert.Equal(t, "test@example.com", user.Email)
    
    mockRepo.AssertExpectations(t)
}
```

### 9. Environment Configuration Example

```env
# .env file example

# Application
APP_NAME="My Golang App"
APP_ENV=development
APP_PORT=8080
APP_DEBUG=true

# Database connections
DB_CONNECTION=mysql

# MySQL Primary
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=app_user
MYSQL_PASSWORD=secure_password

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=golang_starter_kit_2025_pg
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secure_password

# JWT
JWT_SECRET_KEY=your-super-secret-jwt-key
JWT_EXPIRES_IN=24h

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

### 10. Docker Compose Example

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_CONNECTION=mysql
      - MYSQL_HOST=mysql
      - POSTGRES_HOST=postgres
    depends_on:
      - mysql
      - postgres
    
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: golang_starter_kit_2025
      MYSQL_USER: app_user
      MYSQL_PASSWORD: secure_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      
  postgres:
    image: postgres:13
    environment:
      POSTGRES_DB: golang_starter_kit_2025_pg
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secure_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  mysql_data:
  postgres_data:
```

## 🔗 Related Documentation

- [Installation Guide](../guides/installation.md)
- [Architecture Guide](../guides/architecture.md)
- [Database Guide](../database/README.md)
- [API Documentation](../api/README.md)
- [CLI Commands](../guides/cli.md)

## 📞 Need Help?

Jika Anda membutuhkan bantuan atau memiliki pertanyaan tentang contoh-contoh ini:

1. Buka [GitHub Issues](https://github.com/RahmatRafiq/golang_starter_kit_2025/issues)
2. Lihat [FAQ](../guides/faq.md)
3. Join [Discord Community](https://discord.gg/example)

---

**💡 Tip**: Mulai dengan contoh basic-usage terlebih dahulu, kemudian lanjut ke advanced-examples setelah familiar dengan konsep dasarnya.
