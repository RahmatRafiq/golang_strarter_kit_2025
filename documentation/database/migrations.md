# Migration Guide

## Overview

Sistem migrasi pada aplikasi ini mendukung multiple database connections dan menyediakan tools untuk mengelola schema database secara version-controlled.

## Migration Structure

```
app/database/migrations/
├── 20250629120000_create_users_table.sql
├── 20250629120100_create_roles_table.sql
├── 20250629120200_create_permissions_table.sql
└── 20250629120300_create_user_roles_table.sql
```

### File Format
- **Timestamp**: `YYYYMMDDHHMMSS` format
- **Description**: Descriptive name dengan underscore
- **Extension**: `.sql`

## Migration File Structure

Setiap file migrasi memiliki dua bagian: UP dan DOWN

```sql
-- +++ UP Migration
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE INDEX idx_users_email ON users(email);

-- --- DOWN Migration
DROP INDEX idx_users_email ON users;
DROP TABLE IF EXISTS users;
```

## CLI Commands

### 1. Create Migration

```bash
# Basic migration
go run main.go make:migration create_users_table

# Table creation migration
go run main.go make:migration create_products_table

# Table alteration migration  
go run main.go make:migration alter_users_add_avatar_column

# Index creation migration
go run main.go make:migration add_index_users_email
```

### 2. Run Migrations

#### Single Migration
```bash
# Run on default connection (MySQL)
go run main.go migrate --file=20250629120000_create_users_table

# Run on specific connection
go run main.go migrate --file=20250629120000_create_users_table --connection=postgres
go run main.go migrate --file=20250629120000_create_users_table --connection=mysql_secondary
```

#### All Pending Migrations
```bash
# Run all pending migrations on default connection
go run main.go migrate:all

# Run all pending migrations on specific connection
go run main.go migrate:all --connection=mysql
go run main.go migrate:all --connection=postgres
go run main.go migrate:all --connection=mysql_secondary
```

### 3. Rollback Migrations

#### Single Migration Rollback
```bash
# Rollback specific migration
go run main.go rollback --file=20250629120000_create_users_table --connection=mysql
```

#### Batch Rollback
```bash
# Rollback last batch on default connection
go run main.go rollback:batch

# Rollback specific batch
go run main.go rollback:batch --batch=3 --connection=postgres

# Rollback all migrations
go run main.go rollback:all --connection=mysql
```

### 4. Fresh Migrations

```bash
# Drop all tables and re-run all migrations
go run main.go migrate:fresh --connection=mysql
go run main.go migrate:fresh --connection=postgres
```

## Database-Specific Migrations

### MySQL/MariaDB Migrations

```sql
-- +++ UP Migration
CREATE TABLE products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    stock INT NOT NULL DEFAULT 0,
    category_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    INDEX idx_products_category (category_id),
    INDEX idx_products_name (name),
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --- DOWN Migration
DROP TABLE IF EXISTS products;
```

### PostgreSQL Migrations

```sql
-- +++ UP Migration
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    stock INTEGER NOT NULL DEFAULT 0,
    category_id BIGINT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_metadata ON products USING GIN(metadata);

ALTER TABLE products ADD CONSTRAINT fk_products_category 
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL;

-- --- DOWN Migration
DROP TABLE IF EXISTS products;
```

## Advanced Migration Patterns

### 1. Data Migration with Schema Changes

```sql
-- +++ UP Migration
-- Add new column
ALTER TABLE users ADD COLUMN full_name VARCHAR(255);

-- Migrate existing data
UPDATE users SET full_name = CONCAT(first_name, ' ', last_name) 
WHERE first_name IS NOT NULL AND last_name IS NOT NULL;

-- Drop old columns
ALTER TABLE users DROP COLUMN first_name;
ALTER TABLE users DROP COLUMN last_name;

-- --- DOWN Migration
-- Add back old columns
ALTER TABLE users ADD COLUMN first_name VARCHAR(255);
ALTER TABLE users ADD COLUMN last_name VARCHAR(255);

-- Migrate data back
UPDATE users SET 
    first_name = SUBSTRING_INDEX(full_name, ' ', 1),
    last_name = SUBSTRING_INDEX(full_name, ' ', -1)
WHERE full_name IS NOT NULL;

-- Drop new column
ALTER TABLE users DROP COLUMN full_name;
```

### 2. Index Management

```sql
-- +++ UP Migration
-- Add composite index for better query performance
CREATE INDEX idx_orders_user_status_date ON orders(user_id, status, created_at);

-- Add unique constraint
ALTER TABLE products ADD CONSTRAINT uk_products_sku UNIQUE (sku);

-- --- DOWN Migration
DROP INDEX idx_orders_user_status_date ON orders;
ALTER TABLE products DROP CONSTRAINT uk_products_sku;
```

### 3. View Creation

```sql
-- +++ UP Migration
CREATE VIEW user_stats AS
SELECT 
    u.id,
    u.name,
    u.email,
    COUNT(o.id) as total_orders,
    SUM(o.total_amount) as total_spent,
    MAX(o.created_at) as last_order_date
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.name, u.email;

-- --- DOWN Migration
DROP VIEW IF EXISTS user_stats;
```

## Multi-Database Migration Strategy

### 1. Shared Schema Migrations

Untuk tabel yang perlu ada di multiple databases:

```bash
# Run same migration on multiple connections
go run main.go migrate --file=20250629_create_users_table --connection=mysql
go run main.go migrate --file=20250629_create_users_table --connection=postgres
go run main.go migrate --file=20250629_create_users_table --connection=mysql_secondary
```

### 2. Database-Specific Migrations

Buat migration files terpisah untuk database-specific features:

```bash
# MySQL-specific features
go run main.go make:migration create_mysql_specific_indexes
go run main.go migrate --file=20250629_create_mysql_specific_indexes --connection=mysql

# PostgreSQL-specific features  
go run main.go make:migration create_postgres_jsonb_columns
go run main.go migrate --file=20250629_create_postgres_jsonb_columns --connection=postgres
```

### 3. Migration Scripts for Multiple Databases

```bash
#!/bin/bash
# migrate-all-dbs.sh

echo "Running migrations on all databases..."

# MySQL Primary
echo "Migrating MySQL Primary..."
go run main.go migrate:all --connection=mysql

# PostgreSQL
echo "Migrating PostgreSQL..."
go run main.go migrate:all --connection=postgres

# MySQL Secondary (if configured)
echo "Migrating MySQL Secondary..."
go run main.go migrate:all --connection=mysql_secondary

echo "All migrations completed!"
```

## Migration Management

### 1. Check Migration Status

```bash
# Check which migrations have been applied
go run main.go migrate:status --connection=mysql
go run main.go migrate:status --connection=postgres
```

### 2. Migration Validation

```bash
# Validate migration files before running
go run main.go migrate:validate

# Dry run (show what would be executed)
go run main.go migrate:dry-run --connection=mysql
```

### 3. Migration Rollback Planning

```bash
# Show rollback plan for specific batch
go run main.go rollback:plan --batch=5 --connection=mysql

# Show rollback plan for all batches
go run main.go rollback:plan --all --connection=postgres
```

## Best Practices

### 1. Migration Naming Convention

```bash
# Good naming
create_users_table
alter_products_add_description
add_index_orders_user_id
drop_unused_categories_table

# Avoid
migration1
fix_users
update_stuff
```

### 2. Backward Compatibility

```sql
-- Good: Always provide DOWN migration
-- +++ UP Migration
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);

-- --- DOWN Migration  
ALTER TABLE users DROP COLUMN avatar_url;

-- Bad: Empty DOWN migration
-- --- DOWN Migration
-- TODO: Add rollback logic
```

### 3. Data Safety

```sql
-- Good: Preserve data during schema changes
-- +++ UP Migration
-- Rename column instead of drop/create
ALTER TABLE products CHANGE COLUMN old_name new_name VARCHAR(255);

-- Bad: Destructive changes without data preservation
-- +++ UP Migration
ALTER TABLE products DROP COLUMN important_data;
ALTER TABLE products ADD COLUMN important_data VARCHAR(255);
```

### 4. Performance Considerations

```sql
-- Good: Add indexes for large tables
-- +++ UP Migration
CREATE TABLE large_table (...);
CREATE INDEX idx_large_table_lookup ON large_table(lookup_column);

-- Good: Use appropriate column types
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,  -- BIGINT for scalability
    email VARCHAR(255) NOT NULL,          -- Reasonable length
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Testing Migrations

### 1. Local Testing

```bash
# Test on local development database
go run main.go migrate:fresh --connection=mysql
go run main.go db:seed

# Verify data integrity
go run main.go migrate:rollback --batch=1
go run main.go migrate:all
```

### 2. Staging Environment

```bash
# Test migration on staging
export DB_CONNECTION=mysql_staging
go run main.go migrate:all

# Verify application functionality
curl http://staging.example.com/health
```

### 3. Production Deployment

```bash
# Backup before migration
go run main.go db:backup --connection=mysql --file=pre_migration_backup.sql

# Run migration
go run main.go migrate:all --connection=mysql

# Verify migration success
go run main.go db:status --connection=mysql
```

## Troubleshooting

### Common Migration Issues

#### 1. Migration Already Applied
```
Error: migration already applied
```
**Solution**: Check migration status and skip if already applied

#### 2. Foreign Key Constraint Error
```
Error: Cannot add foreign key constraint
```
**Solution**: Ensure referenced table exists and has correct indexes

#### 3. Column Already Exists
```
Error: Duplicate column name
```
**Solution**: Add IF NOT EXISTS clause or check column existence

#### 4. Rollback Failure
```
Error: Cannot rollback migration
```
**Solution**: Review DOWN migration logic and ensure it's reversible

### Migration Recovery

```bash
# Mark migration as applied without running
go run main.go migrate:mark-applied --file=20250629_problematic_migration

# Force rollback specific migration
go run main.go rollback:force --file=20250629_problematic_migration

# Reset migration state
go run main.go migrate:reset --connection=mysql
```

---

Next: [Seeder Guide](seeders.md) | [Architecture Guide](../guides/architecture.md)
