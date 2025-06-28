# Database Seeders

Panduan lengkap untuk menggunakan database seeders pada Golang Starter Kit.

## Overview

Database Seeders adalah script yang digunakan untuk mengisi database dengan data awal atau data testing. Seeder berguna untuk:

- Mengisi data master (roles, permissions, categories)
- Membuat data testing untuk development
- Setup data awal untuk production
- Populate data untuk demonstration

## Basic Seeder Structure

### 1. Seeder Interface

```go
// interfaces/seeder.go
package interfaces

type SeederInterface interface {
    Run() error
    GetName() string
    GetOrder() int
}
```

### 2. Base Seeder

```go
// app/database/seeders/base_seeder.go
package seeders

import (
    "your-app/facades"
    "gorm.io/gorm"
)

type BaseSeeder struct {
    DB *gorm.DB
}

func NewBaseSeeder() *BaseSeeder {
    return &BaseSeeder{
        DB: facades.DB(),
    }
}

func (s *BaseSeeder) GetDB() *gorm.DB {
    return s.DB
}
```

## Creating Seeders

### 1. Role Seeder

```go
// app/database/seeders/role_seeder.go
package seeders

import (
    "your-app/app/models"
    "your-app/interfaces"
)

type RoleSeeder struct {
    *BaseSeeder
}

func NewRoleSeeder() interfaces.SeederInterface {
    return &RoleSeeder{
        BaseSeeder: NewBaseSeeder(),
    }
}

func (s *RoleSeeder) GetName() string {
    return "RoleSeeder"
}

func (s *RoleSeeder) GetOrder() int {
    return 1
}

func (s *RoleSeeder) Run() error {
    roles := []models.Role{
        {
            Name:        "admin",
            DisplayName: "Administrator",
            Description: "System administrator with full access",
        },
        {
            Name:        "manager",
            DisplayName: "Manager",
            Description: "Manager with limited administrative access",
        },
        {
            Name:        "user",
            DisplayName: "User",
            Description: "Regular user with basic access",
        },
    }

    for _, role := range roles {
        // Check if role already exists
        var existingRole models.Role
        err := s.DB.Where("name = ?", role.Name).First(&existingRole).Error
        
        if err == gorm.ErrRecordNotFound {
            // Role doesn't exist, create it
            if err := s.DB.Create(&role).Error; err != nil {
                return fmt.Errorf("failed to create role %s: %w", role.Name, err)
            }
            fmt.Printf("Created role: %s\n", role.Name)
        } else if err != nil {
            return fmt.Errorf("failed to check role %s: %w", role.Name, err)
        } else {
            fmt.Printf("Role already exists: %s\n", role.Name)
        }
    }

    return nil
}
```

### 2. Permission Seeder

```go
// app/database/seeders/permission_seeder.go
package seeders

import (
    "your-app/app/models"
    "your-app/interfaces"
)

type PermissionSeeder struct {
    *BaseSeeder
}

func NewPermissionSeeder() interfaces.SeederInterface {
    return &PermissionSeeder{
        BaseSeeder: NewBaseSeeder(),
    }
}

func (s *PermissionSeeder) GetName() string {
    return "PermissionSeeder"
}

func (s *PermissionSeeder) GetOrder() int {
    return 2
}

func (s *PermissionSeeder) Run() error {
    permissions := []models.Permission{
        // User permissions
        {Name: "users.create", DisplayName: "Create Users", Description: "Can create new users"},
        {Name: "users.read", DisplayName: "Read Users", Description: "Can view users"},
        {Name: "users.update", DisplayName: "Update Users", Description: "Can update users"},
        {Name: "users.delete", DisplayName: "Delete Users", Description: "Can delete users"},
        
        // Product permissions
        {Name: "products.create", DisplayName: "Create Products", Description: "Can create products"},
        {Name: "products.read", DisplayName: "Read Products", Description: "Can view products"},
        {Name: "products.update", DisplayName: "Update Products", Description: "Can update products"},
        {Name: "products.delete", DisplayName: "Delete Products", Description: "Can delete products"},
        
        // Category permissions
        {Name: "categories.create", DisplayName: "Create Categories", Description: "Can create categories"},
        {Name: "categories.read", DisplayName: "Read Categories", Description: "Can view categories"},
        {Name: "categories.update", DisplayName: "Update Categories", Description: "Can update categories"},
        {Name: "categories.delete", DisplayName: "Delete Categories", Description: "Can delete categories"},
    }

    for _, permission := range permissions {
        var existingPermission models.Permission
        err := s.DB.Where("name = ?", permission.Name).First(&existingPermission).Error
        
        if err == gorm.ErrRecordNotFound {
            if err := s.DB.Create(&permission).Error; err != nil {
                return fmt.Errorf("failed to create permission %s: %w", permission.Name, err)
            }
            fmt.Printf("Created permission: %s\n", permission.Name)
        } else if err != nil {
            return fmt.Errorf("failed to check permission %s: %w", permission.Name, err)
        }
    }

    return nil
}
```

### 3. User Seeder

```go
// app/database/seeders/user_seeder.go
package seeders

import (
    "your-app/app/models"
    "your-app/app/helpers"
    "your-app/interfaces"
)

type UserSeeder struct {
    *BaseSeeder
}

func NewUserSeeder() interfaces.SeederInterface {
    return &UserSeeder{
        BaseSeeder: NewBaseSeeder(),
    }
}

func (s *UserSeeder) GetName() string {
    return "UserSeeder"
}

func (s *UserSeeder) GetOrder() int {
    return 3
}

func (s *UserSeeder) Run() error {
    // Hash password
    hashedPassword, err := helpers.HashPassword("admin123")
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    // Create admin user
    adminUser := models.User{
        Name:     "Administrator",
        Email:    "admin@example.com",
        Password: hashedPassword,
        Active:   true,
    }

    var existingUser models.User
    err = s.DB.Where("email = ?", adminUser.Email).First(&existingUser).Error
    
    if err == gorm.ErrRecordNotFound {
        if err := s.DB.Create(&adminUser).Error; err != nil {
            return fmt.Errorf("failed to create admin user: %w", err)
        }
        fmt.Printf("Created admin user: %s\n", adminUser.Email)
        
        // Assign admin role
        var adminRole models.Role
        if err := s.DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
            return fmt.Errorf("admin role not found: %w", err)
        }
        
        userRole := models.UserHasRole{
            UserID: adminUser.ID,
            RoleID: adminRole.ID,
        }
        
        if err := s.DB.Create(&userRole).Error; err != nil {
            return fmt.Errorf("failed to assign admin role: %w", err)
        }
        
    } else if err != nil {
        return fmt.Errorf("failed to check admin user: %w", err)
    }

    return nil
}
```

### 4. Category Seeder

```go
// app/database/seeders/category_seeder.go
package seeders

import (
    "your-app/app/models"
    "your-app/interfaces"
)

type CategorySeeder struct {
    *BaseSeeder
}

func NewCategorySeeder() interfaces.SeederInterface {
    return &CategorySeeder{
        BaseSeeder: NewBaseSeeder(),
    }
}

func (s *CategorySeeder) GetName() string {
    return "CategorySeeder"
}

func (s *CategorySeeder) GetOrder() int {
    return 4
}

func (s *CategorySeeder) Run() error {
    categories := []models.Category{
        {
            Name:        "Electronics",
            Description: "Electronic devices and gadgets",
            Active:      true,
        },
        {
            Name:        "Clothing",
            Description: "Fashion and apparel",
            Active:      true,
        },
        {
            Name:        "Books",
            Description: "Books and educational materials",
            Active:      true,
        },
        {
            Name:        "Home & Garden",
            Description: "Home improvement and gardening supplies",
            Active:      true,
        },
    }

    for _, category := range categories {
        var existingCategory models.Category
        err := s.DB.Where("name = ?", category.Name).First(&existingCategory).Error
        
        if err == gorm.ErrRecordNotFound {
            if err := s.DB.Create(&category).Error; err != nil {
                return fmt.Errorf("failed to create category %s: %w", category.Name, err)
            }
            fmt.Printf("Created category: %s\n", category.Name)
        } else if err != nil {
            return fmt.Errorf("failed to check category %s: %w", category.Name, err)
        }
    }

    return nil
}
```

## Advanced Seeders

### 1. Factory Pattern Seeder

```go
// app/database/seeders/product_seeder.go
package seeders

import (
    "fmt"
    "math/rand"
    "your-app/app/models"
    "your-app/interfaces"
)

type ProductSeeder struct {
    *BaseSeeder
}

func NewProductSeeder() interfaces.SeederInterface {
    return &ProductSeeder{
        BaseSeeder: NewBaseSeeder(),
    }
}

func (s *ProductSeeder) GetName() string {
    return "ProductSeeder"
}

func (s *ProductSeeder) GetOrder() int {
    return 5
}

func (s *ProductSeeder) Run() error {
    // Get all categories
    var categories []models.Category
    if err := s.DB.Find(&categories).Error; err != nil {
        return fmt.Errorf("failed to fetch categories: %w", err)
    }

    if len(categories) == 0 {
        return fmt.Errorf("no categories found, run CategorySeeder first")
    }

    // Generate sample products
    productNames := []string{
        "Smartphone X1", "Laptop Pro", "Wireless Headphones", "Smart Watch",
        "T-Shirt Classic", "Jeans Premium", "Sneakers Sport", "Jacket Winter",
        "Programming Book", "Design Thinking", "Management 101", "History World",
        "Garden Tools Set", "Plant Pot", "Watering Can", "Fertilizer Organic",
    }

    for i, name := range productNames {
        category := categories[i%len(categories)]
        
        product := models.Product{
            Name:        name,
            Description: fmt.Sprintf("High quality %s with excellent features", name),
            Price:       float64(rand.Intn(1000) + 10), // Random price between 10-1010
            CategoryID:  category.ID,
            Stock:       rand.Intn(100) + 1, // Random stock between 1-100
            Active:      true,
        }

        var existingProduct models.Product
        err := s.DB.Where("name = ?", product.Name).First(&existingProduct).Error
        
        if err == gorm.ErrRecordNotFound {
            if err := s.DB.Create(&product).Error; err != nil {
                return fmt.Errorf("failed to create product %s: %w", product.Name, err)
            }
            fmt.Printf("Created product: %s\n", product.Name)
        } else if err != nil {
            return fmt.Errorf("failed to check product %s: %w", product.Name, err)
        }
    }

    return nil
}
```

### 2. CSV-Based Seeder

```go
// app/database/seeders/csv_seeder.go
package seeders

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "your-app/app/models"
    "your-app/interfaces"
)

type CSVProductSeeder struct {
    *BaseSeeder
    FilePath string
}

func NewCSVProductSeeder(filePath string) interfaces.SeederInterface {
    return &CSVProductSeeder{
        BaseSeeder: NewBaseSeeder(),
        FilePath:   filePath,
    }
}

func (s *CSVProductSeeder) GetName() string {
    return "CSVProductSeeder"
}

func (s *CSVProductSeeder) GetOrder() int {
    return 10
}

func (s *CSVProductSeeder) Run() error {
    file, err := os.Open(s.FilePath)
    if err != nil {
        return fmt.Errorf("failed to open CSV file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return fmt.Errorf("failed to read CSV file: %w", err)
    }

    // Skip header
    for i, record := range records[1:] {
        if len(record) < 4 {
            fmt.Printf("Skipping invalid record at line %d\n", i+2)
            continue
        }

        price, err := strconv.ParseFloat(record[2], 64)
        if err != nil {
            fmt.Printf("Invalid price at line %d: %s\n", i+2, record[2])
            continue
        }

        categoryID, err := strconv.ParseUint(record[3], 10, 32)
        if err != nil {
            fmt.Printf("Invalid category ID at line %d: %s\n", i+2, record[3])
            continue
        }

        product := models.Product{
            Name:        record[0],
            Description: record[1],
            Price:       price,
            CategoryID:  uint(categoryID),
            Active:      true,
        }

        if err := s.DB.Create(&product).Error; err != nil {
            fmt.Printf("Failed to create product %s: %v\n", product.Name, err)
            continue
        }

        fmt.Printf("Created product from CSV: %s\n", product.Name)
    }

    return nil
}
```

## Seeder Manager

### 1. Seeder Manager Implementation

```go
// app/database/seeder_manager.go
package database

import (
    "fmt"
    "sort"
    "your-app/app/database/seeders"
    "your-app/interfaces"
)

type SeederManager struct {
    seeders []interfaces.SeederInterface
}

func NewSeederManager() *SeederManager {
    return &SeederManager{
        seeders: make([]interfaces.SeederInterface, 0),
    }
}

func (sm *SeederManager) RegisterSeeder(seeder interfaces.SeederInterface) {
    sm.seeders = append(sm.seeders, seeder)
}

func (sm *SeederManager) RegisterAllSeeders() {
    sm.RegisterSeeder(seeders.NewRoleSeeder())
    sm.RegisterSeeder(seeders.NewPermissionSeeder())
    sm.RegisterSeeder(seeders.NewUserSeeder())
    sm.RegisterSeeder(seeders.NewCategorySeeder())
    sm.RegisterSeeder(seeders.NewProductSeeder())
}

func (sm *SeederManager) RunAll() error {
    // Sort seeders by order
    sort.Slice(sm.seeders, func(i, j int) bool {
        return sm.seeders[i].GetOrder() < sm.seeders[j].GetOrder()
    })

    fmt.Println("Starting database seeding...")
    
    for _, seeder := range sm.seeders {
        fmt.Printf("Running seeder: %s\n", seeder.GetName())
        if err := seeder.Run(); err != nil {
            return fmt.Errorf("seeder %s failed: %w", seeder.GetName(), err)
        }
        fmt.Printf("Seeder %s completed successfully\n", seeder.GetName())
    }

    fmt.Println("Database seeding completed successfully!")
    return nil
}

func (sm *SeederManager) RunSeeder(name string) error {
    for _, seeder := range sm.seeders {
        if seeder.GetName() == name {
            fmt.Printf("Running seeder: %s\n", name)
            if err := seeder.Run(); err != nil {
                return fmt.Errorf("seeder %s failed: %w", name, err)
            }
            fmt.Printf("Seeder %s completed successfully\n", name)
            return nil
        }
    }
    
    return fmt.Errorf("seeder %s not found", name)
}

func (sm *SeederManager) ListSeeders() {
    fmt.Println("Available seeders:")
    
    // Sort by order
    sort.Slice(sm.seeders, func(i, j int) bool {
        return sm.seeders[i].GetOrder() < sm.seeders[j].GetOrder()
    })
    
    for _, seeder := range sm.seeders {
        fmt.Printf("  %d. %s\n", seeder.GetOrder(), seeder.GetName())
    }
}
```

## Using Seeders

### 1. Command Line Interface

```go
// cmd/seeder.go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "your-app/app/database"
    "your-app/bootstrap"
)

func main() {
    // Initialize application
    bootstrap.InitializeApp()
    
    // Command line flags
    var (
        action = flag.String("action", "run", "Action: run, list")
        seeder = flag.String("seeder", "", "Specific seeder to run")
        conn   = flag.String("connection", "mysql", "Database connection")
    )
    flag.Parse()

    // Initialize seeder manager
    seederManager := database.NewSeederManager()
    seederManager.RegisterAllSeeders()

    switch *action {
    case "run":
        if *seeder != "" {
            // Run specific seeder
            if err := seederManager.RunSeeder(*seeder); err != nil {
                log.Fatalf("Failed to run seeder: %v", err)
            }
        } else {
            // Run all seeders
            if err := seederManager.RunAll(); err != nil {
                log.Fatalf("Failed to run seeders: %v", err)
            }
        }
    case "list":
        seederManager.ListSeeders()
    default:
        fmt.Println("Usage: go run cmd/seeder.go [options]")
        fmt.Println("Options:")
        fmt.Println("  -action=run|list    Action to perform")
        fmt.Println("  -seeder=name        Run specific seeder")
        fmt.Println("  -connection=name    Database connection")
        os.Exit(1)
    }
}
```

### 2. Usage Examples

```bash
# Run all seeders
go run cmd/seeder.go

# Run specific seeder
go run cmd/seeder.go -seeder=RoleSeeder

# List available seeders
go run cmd/seeder.go -action=list

# Run with specific database connection
go run cmd/seeder.go -connection=postgresql
```

## Environment-Specific Seeders

### 1. Development Seeders

```go
// app/database/seeders/development_seeder.go
package seeders

type DevelopmentSeeder struct {
    *BaseSeeder
}

func (s *DevelopmentSeeder) Run() error {
    if helpers.GetEnv("APP_ENV", "local") != "local" {
        fmt.Println("Skipping development seeder (not in local environment)")
        return nil
    }

    // Create test users, dummy data, etc.
    return s.createTestData()
}
```

### 2. Production Seeders

```go
// app/database/seeders/production_seeder.go
package seeders

type ProductionSeeder struct {
    *BaseSeeder
}

func (s *ProductionSeeder) Run() error {
    if helpers.GetEnv("APP_ENV", "local") != "production" {
        fmt.Println("Skipping production seeder (not in production environment)")
        return nil
    }

    // Only essential data for production
    return s.createEssentialData()
}
```

## Best Practices

### 1. Idempotent Seeders

```go
func (s *RoleSeeder) Run() error {
    // Always check if data exists before creating
    var existingRole models.Role
    err := s.DB.Where("name = ?", "admin").First(&existingRole).Error
    
    if err == gorm.ErrRecordNotFound {
        // Create only if doesn't exist
        return s.createRole()
    }
    
    return nil
}
```

### 2. Transaction Support

```go
func (s *UserSeeder) Run() error {
    return s.DB.Transaction(func(tx *gorm.DB) error {
        // Create user
        if err := tx.Create(&user).Error; err != nil {
            return err
        }
        
        // Assign roles
        if err := tx.Create(&userRoles).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

### 3. Error Handling

```go
func (s *ProductSeeder) Run() error {
    products := s.getProductData()
    
    for _, product := range products {
        if err := s.createProduct(product); err != nil {
            // Log error but continue with next product
            fmt.Printf("Warning: Failed to create product %s: %v\n", product.Name, err)
            continue
        }
    }
    
    return nil
}
```

### 4. Performance Optimization

```go
func (s *BulkSeeder) Run() error {
    // Use batch insert for better performance
    batchSize := 1000
    
    for i := 0; i < len(data); i += batchSize {
        end := i + batchSize
        if end > len(data) {
            end = len(data)
        }
        
        batch := data[i:end]
        if err := s.DB.CreateInBatches(batch, batchSize).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

## Testing Seeders

### 1. Seeder Tests

```go
// app/database/seeders/role_seeder_test.go
package seeders

import (
    "testing"
    "your-app/test"
    "github.com/stretchr/testify/assert"
)

func TestRoleSeeder_Run(t *testing.T) {
    // Setup test database
    testDB := test.SetupTestDB()
    defer test.CleanupTestDB(testDB)
    
    // Create seeder
    seeder := &RoleSeeder{
        BaseSeeder: &BaseSeeder{DB: testDB},
    }
    
    // Run seeder
    err := seeder.Run()
    assert.NoError(t, err)
    
    // Verify data was created
    var roles []models.Role
    err = testDB.Find(&roles).Error
    assert.NoError(t, err)
    assert.Len(t, roles, 3) // admin, manager, user
    
    // Verify specific role
    var adminRole models.Role
    err = testDB.Where("name = ?", "admin").First(&adminRole).Error
    assert.NoError(t, err)
    assert.Equal(t, "Administrator", adminRole.DisplayName)
}
```

---

Untuk informasi lebih lanjut, lihat:
- [Database Guide](README.md)
- [Migrations Guide](migrations.md)
- [Multi-Database Guide](multi-database.md)
