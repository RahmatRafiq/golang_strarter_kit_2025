# Development Guide

Panduan lengkap untuk development dengan Golang Starter Kit.

## Getting Started

### Prerequisites

- Go 1.21 atau lebih baru
- MySQL 8.0+ atau MariaDB 10.5+
- PostgreSQL 13+ (opsional untuk multi-database)
- Redis 6.0+ (opsional)
- Git

### Development Setup

1. **Clone Repository**
   ```bash
   git clone <repository-url>
   cd golang_strarter_kit_2025
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Setup Environment**
   ```bash
   cp .env.example .env
   # Edit .env sesuai konfigurasi local
   ```

4. **Setup Database**
   ```bash
   # Migrate database
   go run cmd/migrate.go up
   
   # Seed data (opsional)
   go run cmd/seeder.go run
   ```

5. **Run Application**
   ```bash
   # Development dengan auto-reload
   make dev
   
   # Atau manual
   go run main.go
   ```

## Project Structure

```
golang_strarter_kit_2025/
├── app/                    # Application logic
│   ├── controllers/        # HTTP Controllers
│   ├── services/          # Business logic
│   ├── models/            # Database models
│   ├── requests/          # Request validation
│   ├── responses/         # Response structures
│   ├── middleware/        # HTTP middleware
│   ├── handlers/          # Response handlers
│   ├── helpers/           # Utility functions
│   ├── casts/            # Type casting
│   └── database/         # Database utilities
├── bootstrap/             # Application bootstrap
├── cmd/                  # CLI commands
├── config/               # Configuration
├── database/             # Database manager
├── documentation/        # Project documentation
├── examples/             # Usage examples
├── facades/              # Service facades
├── routes/               # Route definitions
├── storage/              # File storage
└── tmp/                  # Temporary files
```

## Development Workflow

### 1. Feature Development

#### a. Create Migration
```bash
# Membuat migration baru
go run cmd/migrate.go create create_products_table

# Edit file migration di app/database/migrations/
# Kemudian jalankan:
go run cmd/migrate.go up
```

#### b. Create Model
```go
// app/models/product.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"not null"`
    Description string         `json:"description"`
    Price       float64        `json:"price" gorm:"not null"`
    CategoryID  uint           `json:"category_id"`
    Category    Category       `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
```

#### c. Create Request Validation
```go
// app/requests/product_request.go
package requests

type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Description string  `json:"description" validate:"max=500"`
    Price       float64 `json:"price" validate:"required,min=0"`
    CategoryID  uint    `json:"category_id" validate:"required,exists=categories,id"`
}

type UpdateProductRequest struct {
    Name        string  `json:"name" validate:"min=3,max=100"`
    Description string  `json:"description" validate:"max=500"`
    Price       float64 `json:"price" validate:"min=0"`
    CategoryID  uint    `json:"category_id" validate:"exists=categories,id"`
}
```

#### d. Create Service
```go
// app/services/product_service.go
package services

import (
    "your-app/app/models"
    "your-app/facades"
    "gorm.io/gorm"
)

type ProductService struct{}

func NewProductService() *ProductService {
    return &ProductService{}
}

func (s *ProductService) Create(data map[string]interface{}) (*models.Product, error) {
    product := &models.Product{
        Name:        data["name"].(string),
        Description: data["description"].(string),
        Price:       data["price"].(float64),
        CategoryID:  uint(data["category_id"].(float64)),
    }
    
    if err := facades.DB().Create(product).Error; err != nil {
        return nil, err
    }
    
    return product, nil
}

func (s *ProductService) GetAll() ([]models.Product, error) {
    var products []models.Product
    err := facades.DB().Preload("Category").Find(&products).Error
    return products, err
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
    var product models.Product
    err := facades.DB().Preload("Category").First(&product, id).Error
    if err != nil {
        return nil, err
    }
    return &product, nil
}
```

#### e. Create Controller
```go
// app/controllers/product_controller.go
package controllers

import (
    "net/http"
    "strconv"
    
    "your-app/app/requests"
    "your-app/app/services"
    "your-app/app/handlers"
    
    "github.com/gin-gonic/gin"
)

type ProductController struct {
    productService *services.ProductService
}

func NewProductController() *ProductController {
    return &ProductController{
        productService: services.NewProductService(),
    }
}

func (pc *ProductController) Index(c *gin.Context) {
    products, err := pc.productService.GetAll()
    if err != nil {
        handlers.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch products", err)
        return
    }
    
    handlers.SuccessResponse(c, "Products retrieved successfully", products)
}

func (pc *ProductController) Store(c *gin.Context) {
    var req requests.CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        handlers.ValidationErrorResponse(c, err)
        return
    }
    
    data := map[string]interface{}{
        "name":        req.Name,
        "description": req.Description,
        "price":       req.Price,
        "category_id": req.CategoryID,
    }
    
    product, err := pc.productService.Create(data)
    if err != nil {
        handlers.ErrorResponse(c, http.StatusInternalServerError, "Failed to create product", err)
        return
    }
    
    handlers.SuccessResponse(c, "Product created successfully", product)
}
```

#### f. Register Routes
```go
// routes/web.go
func RegisterRoutes(router *gin.Engine) {
    // ... existing routes
    
    productController := controllers.NewProductController()
    
    api := router.Group("/api")
    {
        products := api.Group("/products")
        {
            products.GET("", productController.Index)
            products.POST("", productController.Store)
            products.GET("/:id", productController.Show)
            products.PUT("/:id", productController.Update)
            products.DELETE("/:id", productController.Destroy)
        }
    }
}
```

### 2. Testing

#### Unit Testing
```bash
# Run all tests
go test ./...

# Run tests dengan coverage
go test -cover ./...

# Run specific test
go test ./app/services -v
```

#### Integration Testing
```bash
# Setup test database
export APP_ENV=testing
go run cmd/migrate.go up --connection=testing

# Run integration tests
go test ./app/controllers -v
```

### 3. Database Operations

#### Migration Commands
```bash
# Create new migration
go run cmd/migrate.go create migration_name

# Run migrations
go run cmd/migrate.go up

# Rollback migrations
go run cmd/migrate.go down

# Check migration status
go run cmd/migrate.go status

# Multi-database migrations
go run cmd/migrate.go up --connection=postgresql
go run cmd/migrate.go up --connection=mysql_secondary
```

#### Seeding Commands
```bash
# Run all seeders
go run cmd/seeder.go run

# Run specific seeder
go run cmd/seeder.go run --seeder=UserSeeder

# Multi-database seeding
go run cmd/seeder.go run --connection=postgresql
```

## Best Practices

### 1. Code Organization

- **Controllers**: Hanya handle HTTP requests dan responses
- **Services**: Business logic dan data processing
- **Models**: Data structure dan database relationships
- **Repositories**: Data access layer (jika diperlukan)

### 2. Error Handling

```go
// Gunakan error handling yang konsisten
if err != nil {
    log.Printf("Error occurred: %v", err)
    handlers.ErrorResponse(c, http.StatusInternalServerError, "Internal server error", err)
    return
}
```

### 3. Database Queries

```go
// Gunakan preload untuk eager loading
var products []models.Product
facades.DB().Preload("Category").Find(&products)

// Gunakan transactions untuk multiple operations
tx := facades.DB().Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&product).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

### 4. Validation

```go
// Gunakan struct tags untuk validation
type UserRequest struct {
    Name     string `json:"name" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

### 5. Configuration

```go
// Gunakan environment variables
dbHost := helpers.GetEnv("DB_HOST", "127.0.0.1")
dbPort := helpers.GetEnvAsInt("DB_PORT", 3306)
```

## Development Tools

### 1. Hot Reload

```bash
# Install Air untuk auto-reload
go install github.com/cosmtrek/air@latest

# Run dengan auto-reload
air
# atau
make dev
```

### 2. Code Formatting

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run
```

### 3. Dependency Management

```bash
# Add dependency
go get github.com/pkg/errors

# Update dependencies
go get -u ./...

# Tidy up modules
go mod tidy
```

## Debugging

### 1. Logging

```go
import "log"

// Basic logging
log.Printf("Debug: %v", data)

// Structured logging (implement jika diperlukan)
logger.Info("User created", map[string]interface{}{
    "user_id": user.ID,
    "email": user.Email,
})
```

### 2. Database Debugging

```go
// Enable SQL logging
db := facades.DB().Debug()
db.Find(&users)
```

### 3. Profiling

```bash
# CPU profiling
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap
```

## Performance Optimization

### 1. Database

- Gunakan indexes yang tepat
- Optimize query dengan EXPLAIN
- Gunakan connection pooling
- Implement caching untuk data yang sering diakses

### 2. HTTP

- Gunakan middleware untuk caching
- Implement rate limiting
- Compress responses
- Use HTTP/2

### 3. Memory

- Avoid memory leaks
- Use sync.Pool untuk object reuse
- Profile memory usage regularly

## Deployment

### 1. Build

```bash
# Build for production
go build -o bin/app main.go

# Build with optimizations
go build -ldflags="-w -s" -o bin/app main.go
```

### 2. Docker

```bash
# Build docker image
docker build -t app-name .

# Run with docker-compose
docker-compose up -d
```

---

Untuk informasi lebih lanjut, lihat:
- [Architecture Guide](architecture.md)
- [Environment Configuration](environment.md)
- [Database Guide](../database/README.md)
