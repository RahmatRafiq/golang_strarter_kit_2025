# Database Management Guide

Complete guide to database migrations and seeders in Golang Starter Kit 2025. This project implements a powerful Laravel-inspired migration system with batch tracking, multi-database support, and comprehensive CLI tools.

## Overview

The migration system provides:
- **Batch Tracking**: Track migrations in batches like Laravel
- **Multi-Database Support**: Run migrations on any configured database
- **Bidirectional Migrations**: UP and DOWN migrations for safe rollbacks
- **Status Tracking**: See which migrations have run and which are pending
- **Seeder System**: Database seeding with rollback support
- **CLI Tools**: Comprehensive command-line interface

## Migration Commands

### 1. Create New Migration

```bash
go run main.go make:migration <migration_name>
```

**Description**: Creates a new migration file with timestamp prefix.

**File Format**: `YYYYMMDDHHMMSS_<migration_name>.sql`

**Location**: `app/database/migrations/`

**Examples:**
```bash
# Create users table
go run main.go make:migration create_users_table

# Alter products table
go run main.go make:migration alter_products_add_status

# Add indexes
go run main.go make:migration add_indexes_to_orders
```

**Best Practices:**
- Use descriptive names: `create_`, `alter_`, `add_`, `drop_`
- One table operation per migration
- Keep migrations atomic and reversible

---

### 2. Run Specific Migration

```bash
go run main.go migrate --file=<filename> [--connection=mysql]
```

**Description**: Run a specific migration file (UP direction only).

**Parameters:**
- `--file`: Migration filename (required)
- `--connection`: Database connection (optional, default: mysql)

**Examples:**
```bash
# Run specific migration on default database
go run main.go migrate --file=20250426184415_create_roles_table.sql

# Run on PostgreSQL
go run main.go migrate --file=20250426184415_create_roles_table.sql --connection=postgres
```

**Use Cases:**
- Testing a new migration
- Re-running a failed migration after fixes
- Applying migrations out of order (not recommended)

---

### 3. Run All Pending Migrations

```bash
go run main.go migrate:all [--connection=mysql]
```

**Description**: Run all pending migrations in chronological order.

**Process:**
1. Creates new batch number
2. Finds all `.sql` files not in `migrations` table
3. Executes UP section of each migration
4. Records each migration with batch number

**Parameters:**
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Run all pending migrations on default database
go run main.go migrate:all

# Run on PostgreSQL
go run main.go migrate:all --connection=postgres

# Run on secondary MySQL
go run main.go migrate:all --connection=mysql_secondary
```

**Output:**
```
Running migrations...
✅ 20250426184415_create_roles_table.sql (Batch: 1)
✅ 20250426184424_create_permissions_table.sql (Batch: 1)
✅ 20250426184432_create_users_table.sql (Batch: 1)

Migrations completed successfully!
```

---

### 4. Check Migration Status

```bash
go run main.go migrate:status [--connection=mysql]
```

**Description**: Display status of all migrations showing which have run and which are pending.

**Output Format:**
```
================================================================================
Migration                                          Batch      Status
--------------------------------------------------------------------------------
20250426184415_create_roles_table                  1          ✅ Ran
20250426184424_create_permissions_table            1          ✅ Ran
20250426184432_create_users_table                  1          ✅ Ran
20250426184440_create_user_has_roles               2          ✅ Ran
20250508221607_create_products_table               -          ⏳ Pending
20250508221615_create_categories_table             -          ⏳ Pending
================================================================================
Total: 6 migrations (4 ran, 2 pending)
```

**Parameters:**
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Check status on default database
go run main.go migrate:status

# Check PostgreSQL status
go run main.go migrate:status --connection=postgres
```

**Use Cases:**
- Pre-deployment checks
- Verify migrations before rollback
- Debug migration issues
- Track deployment progress

---

### 5. Rollback Specific Migration

```bash
go run main.go rollback --file=<filename> [--connection=mysql]
```

**Description**: Rollback a specific migration (run DOWN section) without affecting batch tracking.

**Parameters:**
- `--file`: Migration filename (required)
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Rollback specific migration
go run main.go rollback --file=20250426184415_create_roles_table.sql

# Rollback on PostgreSQL
go run main.go rollback --file=20250426184415_create_roles_table.sql --connection=postgres
```

**Warning**: This doesn't update batch tracking. Use `rollback:batch` for proper batch management.

---

### 6. Rollback Last Batch

```bash
go run main.go rollback:batch [--connection=mysql]
```

**Description**: Rollback all migrations in the last batch (like Laravel).

**Process:**
1. Finds highest batch number
2. Gets all migrations in that batch
3. Runs DOWN section for each migration (reverse order)
4. Removes migrations from tracking table

**Parameters:**
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Rollback last batch
go run main.go rollback:batch

# Rollback on PostgreSQL
go run main.go rollback:batch --connection=postgres
```

**Output:**
```
Rolling back batch 3...
✅ Rolled back: 20250508221615_create_categories_table.sql
✅ Rolled back: 20250508221607_create_products_table.sql

Rollback completed successfully!
```

---

### 7. Rollback N Last Batches

```bash
go run main.go rollback:batch --step=<N> [--connection=mysql]
```

**Description**: Rollback the last N batches (Laravel-style step rollback).

**Parameters:**
- `--step`: Number of batches to rollback (required)
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Rollback last 1 batch
go run main.go rollback:batch --step=1

# Rollback last 3 batches
go run main.go rollback:batch --step=3

# Rollback on PostgreSQL
go run main.go rollback:batch --step=2 --connection=postgres
```

**Process:**
- `--step=1`: Rollback batch 5 (if current max is 5)
- `--step=3`: Rollback batches 5, 4, and 3
- Automatically calculates target batches

**Use Cases:**
- Incremental rollbacks
- Safe deployment rollback
- Testing rollback procedures

---

### 8. Rollback Specific Batch

```bash
go run main.go rollback:batch --batch=<number> [--connection=mysql]
```

**Description**: Rollback all migrations in a specific batch number.

**Parameters:**
- `--batch`: Batch number to rollback (required)
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Rollback batch 2
go run main.go rollback:batch --batch=2

# Rollback batch 1 on PostgreSQL
go run main.go rollback:batch --batch=1 --connection=postgres
```

**Use Cases:**
- Rollback specific deployment
- Fix issues in particular batch
- Selective migration management

---

### 9. Rollback All Migrations

```bash
go run main.go rollback:all [--connection=mysql]
# Alias:
go run main.go migrate:reset [--connection=mysql]
```

**Description**: Rollback ALL migrations from highest batch to batch 1.

**Process:**
1. Finds all batches (highest to lowest)
2. Rolls back each batch sequentially
3. Clears entire `migrations` table

**Parameters:**
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Rollback all migrations
go run main.go rollback:all

# Using alias
go run main.go migrate:reset

# On PostgreSQL
go run main.go rollback:all --connection=postgres
```

**Warning**: This removes all database structure created by migrations. Use with caution!

---

### 10. Fresh Migration

```bash
go run main.go migrate:fresh [--connection=mysql]
```

**Description**: Rollback all migrations, then run all migrations again (clean slate).

**Process:**
1. Rollback all migrations (via `rollback:all`)
2. Run all migrations (via `migrate:all`)

**Parameters:**
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Fresh migration on default database
go run main.go migrate:fresh

# Fresh on PostgreSQL
go run main.go migrate:fresh --connection=postgres
```

**Use Cases:**
- Clean development environment
- Reset test database
- Fix migration conflicts

**Warning**: Destroys all data! Don't use in production.

---

### 11. Fresh Migration with Seeding

```bash
go run main.go migrate:fresh --seed [--connection=mysql]
```

**Description**: Fresh migration + automatically run all seeders.

**Process:**
1. Rollback all migrations
2. Run all migrations
3. Run all seeders

**Parameters:**
- `--seed`: Enable seeding after migration
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Fresh migration with seeding
go run main.go migrate:fresh --seed

# On PostgreSQL
go run main.go migrate:fresh --seed --connection=postgres
```

**Use Cases:**
- Setup development environment
- Reset test database with sample data
- Prepare demo environment

**Workflow:**
```
DB State: Empty
    ↓
rollback:all → Empty (cleaned)
    ↓
migrate:all → Structure created
    ↓
db:seed → Data populated
    ↓
DB State: Fresh with data
```

---

### 12. Wipe Database

```bash
go run main.go db:wipe [--connection=mysql] [--force]
```

**Description**: Drop ALL tables from database (nuclear option).

**Parameters:**
- `--connection`: Database connection (optional)
- `--force`: Skip confirmation prompt (for CI/CD)

**Process:**
1. Disables foreign key checks
2. Drops all tables in database
3. Re-enables foreign key checks

**Examples:**
```bash
# Wipe with confirmation
go run main.go db:wipe

# Force wipe (no confirmation)
go run main.go db:wipe --force

# Wipe PostgreSQL
go run main.go db:wipe --connection=postgres --force
```

**Confirmation Prompt:**
```
⚠️  WARNING: This will DROP ALL TABLES in the database!
Type 'yes' to confirm:
```

**Supported Databases:**
- MySQL
- PostgreSQL
- SQLite
- SQL Server

**Use Cases:**
- Clean slate for testing
- Remove all data and structure
- CI/CD pipeline cleanup

**⚠️ WARNING**:
- Extremely destructive operation
- No undo available
- Always backup before wiping
- Never use in production without extreme caution

---

## Seeder Commands

### 1. Create New Seeder

```bash
go run main.go make:seeder --name=<SeederName>
```

**Description**: Create new seeder file with template.

**File Format**: `YYYYMMDDHHMMSS_<SeederName>.go`

**Location**: `app/database/seeds/`

**Generated Template:**
```go
package seeds

import (
    "golang_starter_kit_2025/app/models"
    "gorm.io/gorm"
)

func SeedUserSeeder(db *gorm.DB) error {
    // Seeding logic here
    return nil
}

func RollbackUserSeeder(db *gorm.DB) error {
    // Rollback logic here
    return nil
}
```

**Examples:**
```bash
# Create user seeder
go run main.go make:seeder --name=UserSeeder

# Create product seeder
go run main.go make:seeder --name=ProductSeeder

# Create role seeder
go run main.go make:seeder --name=RoleSeeder
```

**Naming Convention:**
- File: `20250423230248_UserSeeder.go`
- Seed function: `SeedUserSeeder(db *gorm.DB) error`
- Rollback function: `RollbackUserSeeder(db *gorm.DB) error`
- ⚠️ Function names MUST match filename!

---

### 2. Run All Seeders

```bash
go run main.go db:seed [--connection=mysql]
```

**Description**: Run all seeders in `app/database/seeds/` directory.

**Process:**
1. Creates new batch number
2. Finds all seeder files
3. Calls Seed function for each
4. Records in `seeders` table with batch

**Parameters:**
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Run all seeders on default database
go run main.go db:seed

# Run on PostgreSQL
go run main.go db:seed --connection=postgres

# Run on secondary MySQL
go run main.go db:seed --connection=mysql_secondary
```

**Output:**
```
Running seeders...
🌱 UserSeeder seeded successfully (Batch: 1)
🌱 RoleSeeder seeded successfully (Batch: 1)
🌱 ProductSeeder seeded successfully (Batch: 1)

Seeders completed successfully!
```

---

### 3. Run Specific Seeder

```bash
go run main.go db:seed --class=<SeederName> [--connection=mysql]
```

**Description**: Run a single specific seeder (Laravel-style).

**Parameters:**
- `--class`: Seeder class name (required)
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Run user seeder only
go run main.go db:seed --class=UserSeeder

# Run on PostgreSQL
go run main.go db:seed --class=ProductSeeder --connection=postgres
```

**Features:**
- Checks if seeder already ran
- Errors if seeder not found
- Supports multi-database

**Use Cases:**
- Seed specific data only
- Re-run failed seeders
- Test individual seeders

---

### 4. Rollback Last Seeder Batch

```bash
go run main.go rollback:seeder [--connection=mysql]
```

**Description**: Rollback all seeders in the last batch.

**Process:**
1. Finds highest seeder batch number
2. Gets all seeders in that batch
3. Calls Rollback function for each (reverse order)
4. Removes from `seeders` table

**Parameters:**
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Rollback last seeder batch
go run main.go rollback:seeder

# Rollback on PostgreSQL
go run main.go rollback:seeder --connection=postgres
```

**Output:**
```
Rolling back seeder batch 2...
🗑️  ProductSeeder rolled back successfully
🗑️  RoleSeeder rolled back successfully

Seeder rollback completed!
```

---

### 5. Rollback Specific Seeder Batch

```bash
go run main.go rollback:seeder --batch=<number> [--connection=mysql]
```

**Description**: Rollback all seeders in a specific batch number.

**Parameters:**
- `--batch`: Batch number to rollback (required)
- `--connection`: Database connection (optional)

**Examples:**
```bash
# Rollback batch 1
go run main.go rollback:seeder --batch=1

# Rollback on PostgreSQL
go run main.go rollback:seeder --batch=2 --connection=postgres
```

**Use Cases:**
- Remove specific data set
- Rollback problematic seeders
- Clean test data

---

## Multi-Database Support

All migration and seeder commands support the `--connection` flag to target specific databases.

### Available Connections

| Connection Name | Database Type | Default |
|----------------|---------------|---------|
| `mysql` | MySQL (Primary) | ✅ Yes |
| `postgres` | PostgreSQL | No |
| `mysql_secondary` | MySQL (Secondary) | No |

### Usage Examples

```bash
# Migrations on different databases
go run main.go migrate:all --connection=mysql
go run main.go migrate:all --connection=postgres
go run main.go migrate:all --connection=mysql_secondary

# Seeders on different databases
go run main.go db:seed --connection=postgres
go run main.go db:seed --class=UserSeeder --connection=mysql_secondary

# Check status on different databases
go run main.go migrate:status --connection=postgres
go run main.go migrate:status --connection=mysql_secondary

# Wipe different databases
go run main.go db:wipe --connection=postgres --force
```

### Multi-Database Workflow

**Scenario**: Migrate both MySQL and PostgreSQL

```bash
# 1. Migrate primary MySQL
go run main.go migrate:all

# 2. Migrate PostgreSQL
go run main.go migrate:all --connection=postgres

# 3. Check both statuses
go run main.go migrate:status
go run main.go migrate:status --connection=postgres

# 4. Seed both databases
go run main.go db:seed
go run main.go db:seed --connection=postgres
```

---

## Migration File Format

### Basic Structure

All migration files must have UP and DOWN sections:

```sql
-- +++ UP Migration
-- SQL statements to apply the migration

-- --- DOWN Migration
-- SQL statements to reverse the migration
```

### MySQL Example

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

### PostgreSQL Example

```sql
-- +++ UP Migration
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    reference VARCHAR(100) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- --- DOWN Migration
DROP TABLE IF EXISTS users;
```

### Key Differences

| Feature | MySQL | PostgreSQL |
|---------|-------|------------|
| Auto Increment | `AUTO_INCREMENT` | `SERIAL` |
| Update Timestamp | `ON UPDATE CURRENT_TIMESTAMP` | Requires trigger |
| Boolean Type | `TINYINT(1)` | `BOOLEAN` |
| String Type | `VARCHAR(n)` | `VARCHAR(n)` or `TEXT` |

### Complex Migration Example

```sql
-- +++ UP Migration
-- 1. Create table
CREATE TABLE products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    reference VARCHAR(100) UNIQUE NOT NULL,
    category_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- 2. Add foreign key
ALTER TABLE products
ADD CONSTRAINT fk_products_category
FOREIGN KEY (category_id) REFERENCES categories(id)
ON DELETE RESTRICT;

-- 3. Add indexes
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_price ON products(price);

-- --- DOWN Migration
-- Drop table (automatically drops constraints and indexes)
DROP TABLE IF EXISTS products;
```

---

## Seeder File Format

### Template Structure

```go
package seeds

import (
    "golang_starter_kit_2025/app/helpers"
    "golang_starter_kit_2025/app/models"
    "log"
    "time"
    "gorm.io/gorm"
)

// Seed function: Insert data
func SeedUserSeeder(db *gorm.DB) error {
    log.Println("🌱 Seeding UserSeeder...")

    // Create data
    users := []models.User{
        {
            Reference: helpers.GenerateReference("USR"),
            Username:  "admin",
            Email:     "admin@example.com",
            Password:  "admin123", // Will be hashed by BeforeCreate hook
            Pin:       "123456",   // Will be hashed by BeforeCreate hook
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        {
            Reference: helpers.GenerateReference("USR"),
            Username:  "user",
            Email:     "user@example.com",
            Password:  "user123",
            Pin:       "654321",
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
    }

    // Insert data
    return db.Create(&users).Error
}

// Rollback function: Remove seeded data
func RollbackUserSeeder(db *gorm.DB) error {
    log.Println("🗑️ Rolling back UserSeeder...")

    // Use Unscoped() for hard delete (ignore soft delete)
    return db.Unscoped().
        Where("username IN ?", []string{"admin", "user"}).
        Delete(&models.User{}).
        Error
}
```

### Important Notes

**1. Function Naming:**
- Seed: `Seed<ClassName>(db *gorm.DB) error`
- Rollback: `Rollback<ClassName>(db *gorm.DB) error`
- Must match filename exactly!

**2. Using Helpers:**
```go
// Generate unique references
reference := helpers.GenerateReference("USR") // USR-20250101-ABC123

// Hash passwords (automatic in BeforeCreate hook)
// No need to manually hash if model has BeforeCreate hook
```

**3. Rollback Best Practices:**
```go
// Use Unscoped() to hard delete
db.Unscoped().Where(...).Delete(...)

// Alternative: Use unique identifiers
db.Unscoped().Where("email IN ?", seededEmails).Delete(...)

// Or use references
db.Unscoped().Where("reference LIKE ?", "USR-20250101%").Delete(...)
```

**4. Idempotent Seeders:**
```go
func SeedUserSeeder(db *gorm.DB) error {
    // Check if already seeded
    var count int64
    db.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&count)

    if count > 0 {
        log.Println("UserSeeder already seeded, skipping...")
        return nil
    }

    // Proceed with seeding
    // ...
}
```

---

## Best Practices

### 1. Check Status Before Deploy

```bash
# Always check migration status before deployment
go run main.go migrate:status

# Expected output: All migrations should show ✅ Ran
# If pending migrations exist, run migrate:all
```

### 2. Test with Fresh Migration

```bash
# Development: Test with fresh database + seed data
go run main.go migrate:fresh --seed

# Verify everything works from scratch
```

### 3. Incremental Rollback

```bash
# Don't rollback all at once, use steps
go run main.go rollback:batch --step=1

# Test after each rollback
go run main.go migrate:status
```

### 4. CI/CD Pipeline

```bash
# Automated testing pipeline
go run main.go db:wipe --force --connection=test_db
go run main.go migrate:fresh --seed --connection=test_db
# Run tests
go test ./...
```

### 5. Development Workflow

```bash
# 1. Create migration
go run main.go make:migration create_orders_table

# 2. Edit migration file
nano app/database/migrations/YYYYMMDDHHMMSS_create_orders_table.sql

# 3. Run migration
go run main.go migrate:all

# 4. Check status
go run main.go migrate:status

# 5. If error, rollback and fix
go run main.go rollback:batch
# Fix the SQL
go run main.go migrate:all
```

### 6. Database-Specific Migrations

```bash
# If using different databases, test on all:
go run main.go migrate:fresh --connection=mysql
go run main.go migrate:fresh --connection=postgres

# Ensure SQL syntax works on both
```

### 7. Version Control

```bash
# Commit migrations with descriptive messages
git add app/database/migrations/20250426184415_create_roles_table.sql
git commit -m "feat: add roles table migration"

# Never modify migrations that have been deployed
# Create new migration for changes
```

### 8. Backup Before Rollback

```bash
# Always backup production database before rollback
mysqldump -u root -p golang_starter_kit_2025 > backup_$(date +%Y%m%d).sql

# Then rollback
go run main.go rollback:batch
```

---

## Troubleshooting

### Migration Failed Mid-Execution

**Problem**: Migration partially executed, database in inconsistent state.

**Solution:**
```bash
# 1. Check status
go run main.go migrate:status

# 2. Rollback problematic batch
go run main.go rollback:batch

# 3. Fix SQL in migration file

# 4. Re-run migration
go run main.go migrate:all

# 5. Verify
go run main.go migrate:status
```

### Seeder Already Exists

**Problem**: Trying to seed data that already exists.

**Solution:**
```bash
# 1. Check if data exists in database manually

# 2. Rollback seeder
go run main.go rollback:seeder

# 3. Re-run seeder
go run main.go db:seed

# Or use --class for specific seeder
go run main.go db:seed --class=UserSeeder
```

### Foreign Key Constraint Error

**Problem**: Cannot drop table due to foreign key constraints.

**Solution:**
```bash
# Use db:wipe which handles foreign keys
go run main.go db:wipe --force

# Or manually in MySQL:
# SET FOREIGN_KEY_CHECKS = 0;
# DROP TABLE ...;
# SET FOREIGN_KEY_CHECKS = 1;
```

### Migration Out of Sync

**Problem**: Migration file exists but not recorded in database.

**Solution:**
```bash
# 1. Check what's in database vs files
go run main.go migrate:status

# 2. Run specific missing migration
go run main.go migrate --file=<filename>

# Or run all pending
go run main.go migrate:all
```

### Duplicate Seeder Execution

**Problem**: Seeder runs twice creating duplicate data.

**Solution:**

Make seeders idempotent:
```go
func SeedUserSeeder(db *gorm.DB) error {
    // Use FirstOrCreate to avoid duplicates
    admin := models.User{Email: "admin@example.com"}
    db.FirstOrCreate(&admin, models.User{
        Email:    "admin@example.com",
        Username: "admin",
        Password: "admin123",
    })

    return nil
}
```

### Database Connection Error

**Problem**: Cannot connect to database.

**Solution:**
```bash
# 1. Check database is running
systemctl status mysql
systemctl status postgresql

# 2. Test connection manually
mysql -h localhost -u root -p
psql -h localhost -U postgres

# 3. Verify .env credentials
cat .env | grep MYSQL
cat .env | grep POSTGRES

# 4. Check multi-database health
curl http://localhost:9999/health/databases
```

---

## Advanced Usage

### Running Migrations Programmatically

```go
package main

import (
    "golang_starter_kit_2025/app/database"
    "golang_starter_kit_2025/facades"
)

func main() {
    // Initialize database
    facades.ConnectDB()
    defer facades.CloseDB()

    // Run migrations
    manager := database.NewMigrationManager(facades.DB, "app/database/migrations")

    // Run all migrations
    err := manager.MigrateAll()
    if err != nil {
        log.Fatal(err)
    }

    // Check status
    status, _ := manager.GetStatus()
    for _, migration := range status {
        fmt.Printf("%s: %s\n", migration.Name, migration.Status)
    }
}
```

### Custom Seeder Logic

```go
func SeedProductSeeder(db *gorm.DB) error {
    // Get category first
    var category models.Category
    if err := db.Where("category = ?", "Electronics").First(&category).Error; err != nil {
        return err
    }

    // Create products with relationship
    products := []models.Product{
        {
            Reference:   helpers.GenerateReference("PRD"),
            CategoryID:  category.ID,
            Name:        "Laptop",
            Price:       1200.50,
            Stock:       100,
        },
    }

    return db.Create(&products).Error
}
```

---

## Further Reading

- [MULTI_DATABASE.md](./MULTI_DATABASE.md) - Multi-database configuration guide
- [GETTING_STARTED.md](./GETTING_STARTED.md) - Initial setup guide
- [API_REFERENCE.md](./API_REFERENCE.md) - API endpoint documentation
- [CLAUDE.md](../CLAUDE.md) - Complete technical documentation

---

**Need help?** Check our [GitHub Discussions](https://github.com/RahmatRafiq/golang_starter_kit_2025/discussions)
