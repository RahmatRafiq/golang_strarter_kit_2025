# Database Management Guide

Complete guide to database migrations and seeders in Golang Starter Kit 2025. This project implements a Laravel-inspired migration system with batch tracking, **independent multi-database connection tracking**, and comprehensive CLI tools.

## 🎯 Key Features

- **✅ Independent Connection Tracking**: Each database connection has separate migration/seeder tracking
- **✅ Batch System**: Track migrations in batches like Laravel
- **✅ Multi-Database Support**: MySQL, PostgreSQL, SQLite, SQL Server
- **✅ Transaction-Protected Rollbacks**: Atomic rollback operations
- **✅ Bidirectional Migrations**: UP and DOWN migrations
- **✅ Seeder System**: Database seeding with rollback support
- **✅ Status Tracking**: Real-time migration status per connection

## 🔧 Database Schema

### Migrations Table

```sql
CREATE TABLE migrations (
    id INT PRIMARY KEY AUTO_INCREMENT,
    connection_name VARCHAR(50) NOT NULL,    -- NEW: Connection isolation
    filename VARCHAR(255) NOT NULL,
    batch INT NOT NULL,
    migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_migration (connection_name, filename)  -- Prevents duplicates
);
```

**Important**: Each database connection tracks its own migrations independently using `connection_name`.

### Seeds Table

```sql
CREATE TABLE seeds (
    id INT PRIMARY KEY AUTO_INCREMENT,
    connection_name VARCHAR(50) NOT NULL,    -- NEW: Connection isolation
    filename VARCHAR(255) NOT NULL,
    batch BIGINT NOT NULL,
    seeded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_seed (connection_name, filename)
);
```

---

## 📋 Migration Commands

### 1. Create New Migration

```bash
go run main.go make:migration <migration_name>
```

**Creates**: `app/database/migrations/YYYYMMDDHHMMSS_<migration_name>.sql`

**Examples:**
```bash
go run main.go make:migration create_users_table
go run main.go make:migration alter_products_add_status
go run main.go make:migration add_indexes_to_orders
```

**Migration File Structure:**
```sql
-- +++ UP Migration
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- --- DOWN Migration
DROP TABLE IF EXISTS users;
```

---

### 2. Run Specific Migration

```bash
go run main.go migrate --file=<filename> [--connection=mysql]
```

**Parameters:**
- `--file`: Migration filename without `.sql` extension (required)
- `--connection`: Target database (optional, default: `mysql`)

**Available Connections:**
- `mysql` (default)
- `postgres`
- `mysql_secondary`

**Examples:**
```bash
# Run on default MySQL connection
go run main.go migrate --file=20250426184415_create_roles_table

# Run on PostgreSQL
go run main.go migrate --file=20250426184415_create_roles_table --connection=postgres

# Run on secondary MySQL
go run main.go migrate --file=20250426184415_create_roles_table --connection=mysql_secondary
```

**Output:**
```
🚀 Migrate: 20250426184415_create_roles_table on connection postgres
Migrating: 20250426184415_create_roles_table
Migrated: 20250426184415_create_roles_table
```

---

### 3. Run All Pending Migrations

```bash
go run main.go migrate:all [--connection=mysql]
```

**Description**: Runs all pending migrations in chronological order and creates a new batch.

**Process:**
1. Queries `migrations` table filtered by `connection_name`
2. Finds all `.sql` files not yet applied to this connection
3. Creates new batch number (independent per connection)
4. Executes UP section of each migration
5. Records migration with `connection_name` and batch number

**Examples:**
```bash
# Run all pending on MySQL (batch 1, 2, 3...)
go run main.go migrate:all

# Run all pending on PostgreSQL (independent batch 1, 2, 3...)
go run main.go migrate:all --connection=postgres
```

**Output:**
```
🚀 Migrate all on connection postgres
Migrating: 20250426184415_create_roles_table
Migrating: 20250426184424_create_permissions_table
Migrating: 20250426184432_create_users_table
Batch 1 applied.
```

**⚠️ Important**: Each connection has independent batch numbering:
- MySQL may be at batch 5
- PostgreSQL may be at batch 1
- This is correct and prevents cross-contamination

---

### 4. Check Migration Status

```bash
go run main.go migrate:status [--connection=mysql]
```

**Description**: Shows migration status for the specified connection only.

**Output Example:**
```
Migration Status on connection: postgres
================================================================================
Migration                                          Batch      Status
--------------------------------------------------------------------------------
20250426184415_create_roles_table                  1          ✅ Ran
20250426184424_create_permissions_table            1          ✅ Ran
20250426184432_create_users_table                  1          ✅ Ran
20250426184449_create_users_has_roles_table        -          ⏳ Pending
20250426184459_create_role_has_permissions_table   -          ⏳ Pending
================================================================================
Total: 5 migrations (3 ran, 2 pending)
```

**Examples:**
```bash
# Check MySQL status
go run main.go migrate:status

# Check PostgreSQL status
go run main.go migrate:status --connection=postgres

# Compare different connections
go run main.go migrate:status --connection=mysql
go run main.go migrate:status --connection=postgres
# Each shows different migration states
```

---

### 5. Rollback Last Batch

```bash
go run main.go rollback:batch [--connection=mysql]
```

**Description**: Rolls back all migrations in the last batch for the specified connection.

**Process:**
1. Finds highest batch number for the connection
2. Gets all migrations in that batch (filtered by `connection_name`)
3. **Wraps in transaction** for atomic operation
4. Runs DOWN section for each migration (reverse order)
5. Removes records from `migrations` table

**Examples:**
```bash
# Rollback last batch on MySQL
go run main.go rollback:batch

# Rollback last batch on PostgreSQL (doesn't affect MySQL)
go run main.go rollback:batch --connection=postgres
```

**Output:**
```
🔄 Rollback batch 2 on connection postgres
Rolling back: 20250426184459_create_role_has_permissions_table
Rolled back: 20250426184459_create_role_has_permissions_table
Rolling back: 20250426184449_create_users_has_roles_table
Rolled back: 20250426184449_create_users_has_roles_table
Batch 2 rolled back.
```

---

### 6. Rollback Specific Batch

```bash
go run main.go rollback:batch --batch=<number> [--connection=mysql]
```

**Parameters:**
- `--batch`: Specific batch number to rollback (required)
- `--connection`: Target database (optional)

**Examples:**
```bash
# Rollback batch 3 on MySQL
go run main.go rollback:batch --batch=3

# Rollback batch 1 on PostgreSQL
go run main.go rollback:batch --batch=1 --connection=postgres
```

---

### 7. Rollback N Last Batches

```bash
go run main.go rollback:batch --step=<N> [--connection=mysql]
```

**Parameters:**
- `--step`: Number of batches to rollback from latest (required)
- `--connection`: Target database (optional)

**Examples:**
```bash
# Rollback last 3 batches on MySQL
go run main.go rollback:batch --step=3

# Rollback last 2 batches on PostgreSQL
go run main.go rollback:batch --step=2 --connection=postgres
```

**Output:**
```
🔄 Rolling back last 3 batch(es) on connection mysql
Batch 5 rolled back.
Batch 4 rolled back.
Batch 3 rolled back.
✅ Successfully rolled back 3 batch(es).
```

---

### 8. Rollback All Migrations

```bash
go run main.go rollback:all [--connection=mysql]
```

**Description**: Rolls back ALL batches on the specified connection (nuclear option).

**⚠️ Warning**: This removes all migrations but keeps the tables/data. Use `migrate:fresh` for complete reset.

**Examples:**
```bash
# Rollback everything on MySQL
go run main.go rollback:all

# Rollback everything on PostgreSQL (MySQL unaffected)
go run main.go rollback:all --connection=postgres
```

---

### 9. Fresh Migration (Reset + Migrate)

```bash
go run main.go migrate:fresh [--seed] [--connection=mysql]
```

**Description**: Deletes all migration records and re-runs all migrations.

**Process:**
1. Deletes records from `migrations` table WHERE `connection_name = ?`
2. Runs all migrations from scratch (batch 1)
3. Optionally runs seeders on same connection

**Parameters:**
- `--seed`: Run seeders after migrations (optional)
- `--connection`: Target database (optional)

**Examples:**
```bash
# Fresh migration on MySQL
go run main.go migrate:fresh

# Fresh migration with seeders on PostgreSQL
go run main.go migrate:fresh --seed --connection=postgres
```

**Output:**
```
🔄 Fresh: rollback all then migrate all on connection postgres
Batch 3 rolled back.
Batch 2 rolled back.
Batch 1 rolled back.
🚀 Migrate all on connection postgres
Migrating: 20250426184415_create_roles_table
...
Batch 1 applied.
🌱 Running seeders on connection postgres...
✅ Seeders completed successfully!
```

**⚠️ Important**: `--seed` flag now correctly respects `--connection` parameter (fixed bug #4).

---

### 10. Reset Migrations

```bash
go run main.go migrate:reset [--connection=mysql]
```

**Description**: Alias for `rollback:all`. Rolls back all batches.

---

### 11. Wipe Database (Nuclear Option)

```bash
go run main.go db:wipe [--connection=mysql] [--force]
```

**Description**: Drops ALL tables from the database (data destruction).

**Parameters:**
- `--connection`: Target database (required)
- `--force`: Skip confirmation prompt (optional)

**Examples:**
```bash
# Wipe with confirmation
go run main.go db:wipe --connection=postgres

# Force wipe without prompt
go run main.go db:wipe --connection=postgres --force
```

**Output:**
```
⚠️  WARNING: This will DROP ALL TABLES in the database!
Connection: postgres

Are you sure you want to continue? (type 'yes' to confirm): yes
🗑️  Dropping all tables on connection postgres
Dropping table: users
Dropping table: roles
Dropping table: permissions
✅ Successfully dropped 15 tables.
```

---

### 12. Rollback Specific Migration

```bash
go run main.go rollback --file=<filename> [--connection=mysql]
```

**Description**: Rolls back a specific migration file (runs DOWN section only).

**⚠️ Warning**: This removes the migration record but doesn't affect batch tracking properly. Use `rollback:batch` for proper batch management.

**Examples:**
```bash
go run main.go rollback --file=20250426184415_create_roles_table
go run main.go rollback --file=20250426184415_create_roles_table --connection=postgres
```

---

## 🌱 Seeder Commands

### 1. Run All Seeders

```bash
go run main.go db:seed [--connection=mysql]
```

**Description**: Runs all registered seeders on the specified connection.

**Process:**
1. Checks `seeds` table filtered by `connection_name`
2. Finds seeders not yet applied to this connection
3. Creates batch number (Unix timestamp)
4. Executes seeder functions
5. Records with `connection_name` and batch

**Examples:**
```bash
# Seed MySQL
go run main.go db:seed

# Seed PostgreSQL (independent tracking)
go run main.go db:seed --connection=postgres
```

**Output:**
```
🌱 Running all seeders on connection postgres
🌱 Seeding: UserSeeder
✅ Seed batch 1770731266 applied on connection postgres.
✅ All seeders completed successfully!
```

---

### 2. Run Specific Seeder

```bash
go run main.go db:seed --class=<SeederName> [--connection=mysql]
```

**Parameters:**
- `--class`: Seeder class name (required)
- `--connection`: Target database (optional)

**Examples:**
```bash
# Run specific seeder on MySQL
go run main.go db:seed --class=UserSeeder

# Run on PostgreSQL
go run main.go db:seed --class=ProductSeeder --connection=postgres
```

**Output:**
```
🌱 Running seeder: UserSeeder
✅ Seeder 'UserSeeder' completed successfully on connection postgres
```

---

### 3. Rollback Last Seeder Batch

```bash
go run main.go rollback:seeder [--connection=mysql]
```

**Description**: Rolls back the last seeder batch for the specified connection.

**Process:**
1. Finds highest batch number in `seeds` WHERE `connection_name = ?`
2. Executes rollback function for each seeder
3. Removes records from `seeds` table

**Examples:**
```bash
# Rollback last seeder batch on MySQL
go run main.go rollback:seeder

# Rollback on PostgreSQL
go run main.go rollback:seeder --connection=postgres
```

**Output:**
```
🔄 Rolling back seeder: UserSeeder
✅ Seeder batch 1770731266 rolled back on connection postgres.
```

---

### 4. Rollback Specific Seeder Batch

```bash
go run main.go rollback:seeder --batch=<timestamp> [--connection=mysql]
```

**Parameters:**
- `--batch`: Specific batch number (Unix timestamp) to rollback
- `--connection`: Target database (optional)

**Examples:**
```bash
go run main.go rollback:seeder --batch=1770731266 --connection=postgres
```

---

### 5. Create New Seeder

```bash
go run main.go make:seeder --name=<SeederName>
```

**Creates**: `app/database/seeds/YYYYMMDDHHMMSS_<SeederName>.go`

**Examples:**
```bash
go run main.go make:seeder --name=ProductSeeder
go run main.go make:seeder --name=CategorySeeder
```

**Generated File Structure:**
```go
package seeds

import (
    "golang_starter_kit_2025/app/models"
    "gorm.io/gorm"
)

func SeedProductSeeder(db *gorm.DB) error {
    // Insert seed data here
    products := []models.Product{
        {Name: "Product 1", Price: 100.00},
        {Name: "Product 2", Price: 200.00},
    }
    return db.Create(&products).Error
}

func RollbackProductSeeder(db *gorm.DB) error {
    // Remove seed data here
    return db.Unscoped().Where("name IN ?", []string{"Product 1", "Product 2"}).
        Delete(&models.Product{}).Error
}
```

**Important**: Register seeder in `app/database/seeder_manager.go`:
```go
var SeederList = []Seeder{
    {Name: "ProductSeeder", Run: seeds.SeedProductSeeder, Rollback: seeds.RollbackProductSeeder},
}
```

---

## 🔍 Multi-Connection Isolation Examples

### Example 1: Independent Batch Tracking

```bash
# MySQL: Run migrations (batch 1)
$ go run main.go migrate:all --connection=mysql
Batch 1 applied.

# PostgreSQL: Run migrations (also batch 1, not batch 2!)
$ go run main.go migrate:all --connection=postgres
Batch 1 applied.  # ✅ CORRECT: Independent numbering

# Verify isolation
$ go run main.go migrate:status --connection=mysql
Total: 12 migrations (12 ran, 0 pending)  # Batch 1

$ go run main.go migrate:status --connection=postgres
Total: 12 migrations (12 ran, 0 pending)  # Also Batch 1
```

**Key Point**: Each connection has independent batch numbering starting from 1.

---

### Example 2: Rollback Doesn't Affect Other Connections

```bash
# Rollback MySQL batch
$ go run main.go rollback:batch --connection=mysql
Batch 1 rolled back.

# Check MySQL (now empty)
$ go run main.go migrate:status --connection=mysql
Total: 12 migrations (0 ran, 12 pending)

# Check PostgreSQL (unaffected!)
$ go run main.go migrate:status --connection=postgres
Total: 12 migrations (12 ran, 0 pending)  # ✅ Still intact
```

---

### Example 3: Seeders Per Connection

```bash
# Seed MySQL
$ go run main.go db:seed --connection=mysql
Seed batch 1770731266 applied on connection mysql.

# Seed PostgreSQL (independent tracking)
$ go run main.go db:seed --connection=postgres
Seed batch 1770731280 applied on connection postgres.

# Run again on MySQL (no re-seeding)
$ go run main.go db:seed --connection=mysql
⚠️ Seeder 'UserSeeder' has already been run  # ✅ Tracks per connection
```

---

### Example 4: Fresh with Seed

```bash
# Fresh migration with seeding on PostgreSQL
$ go run main.go migrate:fresh --seed --connection=postgres

# This will:
# 1. Delete postgres migration records
# 2. Re-run all migrations on postgres (batch 1)
# 3. Run seeders on postgres (NOT on mysql!)  # ✅ Fixed bug #4
```

---

## 🗄️ Database Connection Configuration

Available connections are configured in `.env`:

```bash
# Default MySQL Connection
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=root
MYSQL_PASSWORD=secure_password

# PostgreSQL Connection
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=golang_starter_kit_2025_pg
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secure_password

# Secondary MySQL Connection
MYSQL_SECONDARY_HOST=localhost
MYSQL_SECONDARY_PORT=3307
MYSQL_SECONDARY_DB=golang_starter_kit_2025_secondary
MYSQL_SECONDARY_USER=root
MYSQL_SECONDARY_PASSWORD=secure_password
```

**View Available Connections:**
```bash
go run main.go db:connections
```

**Output:**
```
📊 Available Database Connections:
  - mysql
  - postgres
  - mysql_secondary
```

---

## 📊 Database Schema Details

### Current Schema State

Query migrations table:
```sql
SELECT * FROM migrations ORDER BY batch, id;
```

**Sample Output:**
```
| id | connection_name | filename                    | batch | migrated_at         |
|----|----------------|-----------------------------|-------|---------------------|
| 1  | mysql          | create_roles_table          | 1     | 2025-07-22 22:13:34 |
| 2  | mysql          | create_permissions_table    | 1     | 2025-07-22 22:13:34 |
| 3  | postgres       | create_roles_table          | 1     | 2025-07-22 22:15:10 |
| 4  | postgres       | create_permissions_table    | 1     | 2025-07-22 22:15:10 |
```

Notice: Same migration file, different `connection_name`, same batch number per connection.

Query seeds table:
```sql
SELECT * FROM seeds ORDER BY batch, id;
```

**Sample Output:**
```
| id | connection_name | filename    | batch      | seeded_at           |
|----|----------------|-------------|------------|---------------------|
| 1  | mysql          | UserSeeder  | 1770731266 | 2025-07-22 22:13:36 |
| 2  | postgres       | UserSeeder  | 1770731280 | 2025-07-22 22:15:15 |
```

---

## ⚠️ Important Notes & Best Practices

### 1. **Connection Isolation is Mandatory**

Always specify `--connection` when working with non-default databases:

```bash
# ❌ Wrong: Will use default mysql connection
go run main.go migrate:all

# ✅ Correct: Explicitly specify postgres
go run main.go migrate:all --connection=postgres
```

### 2. **Batch Numbers are Independent**

Each connection has its own batch sequence starting from 1. This is **correct behavior**:

- MySQL: batch 1, 2, 3, 4, 5
- PostgreSQL: batch 1, 2, 3
- Secondary MySQL: batch 1

### 3. **Rollbacks are Transaction-Protected**

All rollback operations use database transactions for atomicity:

```go
// Rollback is wrapped in transaction
return conn.DB.Transaction(func(tx *gorm.DB) error {
    // Run DOWN migrations
    // Delete migration records
    return nil  // Commit
})
```

If any step fails, the entire rollback is reverted.

### 4. **Migration File Syntax**

Use database-agnostic SQL or create connection-specific migrations:

**MySQL:**
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ...
);
```

**PostgreSQL:**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    ...
);
```

**Tip**: Check `conn.IsPostgreSQL()` in migration manager for syntax switching.

### 5. **Seeder Rollback Best Practices**

Always use `Unscoped()` for hard deletes in seeder rollbacks:

```go
func RollbackUserSeeder(db *gorm.DB) error {
    // ✅ Hard delete (bypasses soft delete)
    return db.Unscoped().
        Where("username = ?", "admin").
        Delete(&models.User{}).Error
}
```

### 6. **migrate:fresh vs db:wipe**

- `migrate:fresh`: Removes migration records, keeps table structure
- `db:wipe`: Drops ALL tables (data destruction)

### 7. **Check Before Rollback**

Always check migration status before rolling back:

```bash
# Check current state
go run main.go migrate:status --connection=postgres

# Verify batch numbers
# Then rollback
go run main.go rollback:batch --connection=postgres
```

---

## 🐛 Known Issues & Fixes

### Issue #23: Multi-Connection Tracking Bugs (FIXED)

**Status**: ✅ **RESOLVED** (2026-02-10)

**Problems Fixed:**
1. ✅ Migrations table lacked `connection_name` column
2. ✅ Seeds table lacked `connection_name` column
3. ✅ Rollbacks not wrapped in transactions
4. ✅ `migrate:fresh --seed` ignored connection flag

**Changes Made:**
- Added `connection_name VARCHAR(50)` to both tables
- Added UNIQUE constraints `(connection_name, filename)`
- All queries now filter by `connection_name`
- Rollbacks use transactions
- `migrate:fresh --seed` respects `--connection` parameter

**Migration**: Schema changes applied via `20260210_fix_multi_connection_tracking.sql`

**GitHub Issue**: https://github.com/RahmatRafiq/golang_strarter_kit_2025/issues/23

---

## 🧪 Testing Migration System

### Test Suite

```bash
# Test 1: Independent tracking
go run main.go migrate:all --connection=mysql
go run main.go migrate:all --connection=postgres
# Both should start at batch 1

# Test 2: Status isolation
go run main.go migrate:status --connection=mysql
go run main.go migrate:status --connection=postgres
# Different ran/pending counts

# Test 3: Rollback isolation
go run main.go rollback:batch --connection=mysql
go run main.go migrate:status --connection=postgres
# Postgres migrations unaffected

# Test 4: Seeder tracking
go run main.go db:seed --connection=mysql
go run main.go db:seed --connection=postgres
# Independent seeder batches

# Test 5: Fresh with seed
go run main.go migrate:fresh --seed --connection=postgres
# Seeders run on postgres only
```

---

## 📚 Related Documentation

- [Multi-Database Configuration](./MULTI_DATABASE.md)
- [API Reference](./API_REFERENCE.md)
- [Getting Started Guide](./GETTING_STARTED.md)

---

**Last Updated**: 2026-02-10
**Version**: 2.0 (Multi-Connection Isolation)
