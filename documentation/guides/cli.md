# CLI Commands Guide

## Overview

Aplikasi ini dilengkapi dengan CLI tools yang powerful untuk mengelola database, migrations, seeders, dan berbagai operasi development. CLI ini dirancang untuk memudahkan developer dalam mengelola aplikasi.

## Available Commands

### Database Management

| Command | Description | Options |
|---------|-------------|---------|
| `db:connections` | List all available database connections | - |
| `db:status` | Check database connection status | `--connection` |
| `db:stats` | Show database statistics | `--connection` |
| `db:ping` | Test database connectivity | `--connection` |

### Migration Commands

| Command | Description | Options |
|---------|-------------|---------|
| `make:migration` | Create new migration file | `name` |
| `migrate` | Run specific migration | `--file`, `--connection` |
| `migrate:all` | Run all pending migrations | `--connection` |
| `migrate:fresh` | Drop all tables and re-migrate | `--connection` |
| `rollback` | Rollback specific migration | `--file`, `--connection` |
| `rollback:batch` | Rollback specific batch | `--batch`, `--connection` |
| `rollback:all` | Rollback all migrations | `--connection` |

### Seeder Commands

| Command | Description | Options |
|---------|-------------|---------|
| `make:seeder` | Create new seeder file | `--name` |
| `db:seed` | Run all seeders | `--connection` |
| `rollback:seeder` | Rollback seeder batch | `--batch`, `--connection` |

## Database Management Commands

### 1. List Available Connections

```bash
go run main.go db:connections
```

Output:
```
📊 Available Database Connections:
  - mysql
  - postgres  
  - mysql_secondary
```

### 2. Check Connection Status

```bash
# Check default connection
go run main.go db:status

# Check specific connection
go run main.go db:status --connection=postgres
go run main.go db:status --connection=mysql_secondary
```

Output:
```
🔍 Checking connection status for: postgres
✅ Connection 'postgres' is healthy
   Database Type: postgres
   Open Connections: 5
   In Use: 2
   Idle: 3
```

### 3. Database Statistics

```bash
go run main.go db:stats --connection=mysql
```

Output:
```
📊 Database Statistics for 'mysql':
┌─────────────────────┬──────────┐
│ Metric              │ Value    │
├─────────────────────┼──────────┤
│ Open Connections    │ 15       │
│ In Use              │ 7        │
│ Idle                │ 8        │
│ Max Open            │ 200      │
│ Max Idle            │ 10       │
│ Total Queries       │ 1,234    │
│ Average Latency     │ 2.5ms    │
└─────────────────────┴──────────┘
```

### 4. Test Connectivity

```bash
go run main.go db:ping --connection=postgres
```

Output:
```
🏓 Pinging database 'postgres'...
✅ Database is reachable (Response time: 3ms)
```

## Migration Commands

### 1. Create Migration

```bash
# Basic migration
go run main.go make:migration create_users_table

# More examples
go run main.go make:migration alter_products_add_description
go run main.go make:migration add_index_orders_user_id
go run main.go make:migration create_categories_table
```

Generated file: `app/database/migrations/20250629120000_create_users_table.sql`

```sql
-- +++ UP Migration
CREATE TABLE users (
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- --- DOWN Migration
DROP TABLE IF EXISTS users;
```

### 2. Run Single Migration

```bash
# Run on default connection (mysql)
go run main.go migrate --file=20250629120000_create_users_table

# Run on specific connection
go run main.go migrate --file=20250629120000_create_users_table --connection=postgres
go run main.go migrate --file=20250629120000_create_users_table --connection=mysql_secondary
```

Output:
```
🚀 Migrate: 20250629120000_create_users_table on connection postgres
✅ Migration '20250629120000_create_users_table' applied on connection 'postgres'
```

### 3. Run All Pending Migrations

```bash
# Run all on default connection
go run main.go migrate:all

# Run all on specific connection
go run main.go migrate:all --connection=mysql
go run main.go migrate:all --connection=postgres
```

Output:
```
🚀 Migrate all on connection mysql
🚀 Running 20250629120000_create_users_table on connection mysql
🚀 Running 20250629120100_create_products_table on connection mysql
🚀 Running 20250629120200_create_orders_table on connection mysql
✅ Batch 1 applied on connection mysql
```

### 4. Fresh Migration

```bash
# Drop all tables and re-run all migrations
go run main.go migrate:fresh --connection=mysql
go run main.go migrate:fresh --connection=postgres
```

Output:
```
🔄 Fresh: rollback all then migrate all on connection mysql
🔄 Rollback all on connection mysql
✅ All batches rolled back on connection mysql
🚀 Migrate all on connection mysql
✅ Batch 1 applied on connection mysql
```

### 5. Rollback Migrations

#### Single Migration Rollback
```bash
go run main.go rollback --file=20250629120000_create_users_table --connection=mysql
```

#### Batch Rollback
```bash
# Rollback last batch
go run main.go rollback:batch --connection=mysql

# Rollback specific batch
go run main.go rollback:batch --batch=2 --connection=postgres
```

#### Rollback All
```bash
go run main.go rollback:all --connection=mysql
```

Output:
```
🔄 Rollback all on connection mysql
🔄 Rollback 20250629120200_create_orders_table on connection mysql
🔄 Rollback 20250629120100_create_products_table on connection mysql
🔄 Rollback 20250629120000_create_users_table on connection mysql
✅ All batches rolled back on connection mysql
```

## Seeder Commands

### 1. Create Seeder

```bash
go run main.go make:seeder --name=users_seeder
go run main.go make:seeder --name=products_seeder
go run main.go make:seeder --name=categories_seeder
```

Generated file: `app/database/seeds/20250629120000_users_seeder.go`

### 2. Run Seeders

```bash
# Run all seeders on default connection
go run main.go db:seed

# Run seeders on specific connection
go run main.go db:seed --connection=postgres
```

### 3. Rollback Seeders

```bash
# Rollback last batch
go run main.go rollback:seeder

# Rollback specific batch
go run main.go rollback:seeder --batch=2 --connection=mysql
```

## Advanced Usage

### 1. Multi-Database Migration Workflow

```bash
#!/bin/bash
# migrate-all-databases.sh

echo "🚀 Starting multi-database migration..."

# Check all connections first
echo "📊 Checking database connections..."
go run main.go db:status --connection=mysql
go run main.go db:status --connection=postgres

# Run migrations on all databases
echo "🔄 Running migrations on MySQL..."
go run main.go migrate:all --connection=mysql

echo "🔄 Running migrations on PostgreSQL..."
go run main.go migrate:all --connection=postgres

echo "🔄 Running migrations on MySQL Secondary..."
go run main.go migrate:all --connection=mysql_secondary

echo "✅ All migrations completed!"

# Verify migrations
echo "🔍 Verifying migrations..."
go run main.go migrate:status --connection=mysql
go run main.go migrate:status --connection=postgres
```

### 2. Database Backup Before Migration

```bash
#!/bin/bash
# safe-migrate.sh

CONNECTION=$1
if [ -z "$CONNECTION" ]; then
    echo "Usage: $0 <connection>"
    exit 1
fi

echo "🔒 Creating backup before migration..."
BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S)_${CONNECTION}.sql"
go run main.go db:backup --connection=$CONNECTION --file=$BACKUP_FILE

echo "🚀 Running migrations..."
go run main.go migrate:all --connection=$CONNECTION

echo "✅ Migration completed. Backup saved as: $BACKUP_FILE"
```

### 3. Development Workflow

```bash
#!/bin/bash
# dev-setup.sh

echo "🛠️  Setting up development environment..."

# Fresh migration on development databases
echo "🔄 Fresh migration on development MySQL..."
go run main.go migrate:fresh --connection=mysql

echo "🔄 Fresh migration on development PostgreSQL..."
go run main.go migrate:fresh --connection=postgres

# Seed development data
echo "🌱 Seeding development data..."
go run main.go db:seed --connection=mysql
go run main.go db:seed --connection=postgres

echo "✅ Development environment ready!"
```

## Configuration Options

### 1. Connection-Specific Commands

Semua database commands mendukung flag `--connection` untuk menentukan database mana yang akan digunakan:

```bash
# Available connections
--connection=mysql           # MySQL primary
--connection=postgres        # PostgreSQL
--connection=mysql_secondary # MySQL secondary
```

### 2. Environment-Based Execution

```bash
# Development
export APP_ENV=development
go run main.go migrate:all --connection=mysql

# Staging
export APP_ENV=staging
go run main.go migrate:all --connection=mysql

# Production
export APP_ENV=production
go run main.go migrate:all --connection=mysql
```

### 3. Verbose Output

```bash
# Enable verbose logging
export LOG_LEVEL=debug
go run main.go migrate:all --connection=mysql

# Enable SQL query logging
export DB_LOG_QUERIES=true
go run main.go migrate:all --connection=postgres
```

## Batch Processing

### 1. Migration Batches

Migration system menggunakan batch untuk mengelompokkan migrations:

```bash
# Check current batch
go run main.go migrate:status --connection=mysql
```

Output:
```
📊 Migration Status for 'mysql':
┌─────────┬────────────────────────────────────────┬─────────────────────┐
│ Batch   │ Migration                              │ Applied At          │
├─────────┼────────────────────────────────────────┼─────────────────────┤
│ 1       │ 20250629120000_create_users_table      │ 2025-06-29 12:00:01 │
│ 1       │ 20250629120100_create_products_table   │ 2025-06-29 12:00:02 │
│ 2       │ 20250629120200_create_orders_table     │ 2025-06-29 12:05:01 │
│ 2       │ 20250629120300_create_categories_table │ 2025-06-29 12:05:02 │
└─────────┴────────────────────────────────────────┴─────────────────────┘
```

### 2. Rollback by Batch

```bash
# Rollback last batch (batch 2)
go run main.go rollback:batch --connection=mysql

# Rollback specific batch
go run main.go rollback:batch --batch=1 --connection=mysql
```

## Error Handling

### 1. Common CLI Errors

#### Migration File Not Found
```
Error: failed to read migration file: open app/database/migrations/invalid_migration.sql: no such file or directory
```

**Solution**: Check migration filename and ensure it exists

#### Database Connection Failed
```
Error: failed to get connection 'mysql': dial tcp: connection refused
```

**Solution**: Check database server status and connection configuration

#### Migration Already Applied
```
Error: migration already applied in batch 1
```

**Solution**: Check migration status and skip if already applied

### 2. Debug Mode

```bash
# Enable debug mode for detailed error information
export DEBUG=true
go run main.go migrate:all --connection=mysql
```

### 3. Dry Run Mode

```bash
# Preview what will be executed without actually running
go run main.go migrate:dry-run --connection=mysql
```

Output:
```
🔍 Dry run mode - showing what would be executed:

Migration: 20250629120000_create_users_table
SQL:
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

Migration: 20250629120100_create_products_table
SQL: ...
```

## Best Practices

### 1. Pre-Production Checklist

```bash
#!/bin/bash
# production-deploy.sh

echo "🔍 Pre-production migration checklist..."

# 1. Backup production database
echo "📦 Creating production backup..."
go run main.go db:backup --connection=mysql --file=prod_backup_$(date +%Y%m%d).sql

# 2. Test migrations on staging
echo "🧪 Testing migrations on staging..."
go run main.go migrate:all --connection=staging

# 3. Verify staging
echo "🔍 Verifying staging data..."
go run main.go db:verify --connection=staging

# 4. Run on production
echo "🚀 Running production migration..."
go run main.go migrate:all --connection=mysql

echo "✅ Production deployment completed!"
```

### 2. Development Best Practices

```bash
# Daily development workflow
go run main.go db:status --connection=mysql  # Check status
go run main.go migrate:all --connection=mysql # Apply new migrations
go run main.go db:seed --connection=mysql     # Refresh test data
```

### 3. Team Collaboration

```bash
# Before starting work
git pull origin main
go run main.go migrate:all --connection=mysql

# After creating new migration
go run main.go migrate --file=your_new_migration --connection=mysql
git add app/database/migrations/
git commit -m "Add: new migration for feature X"
```

## Performance Tips

### 1. Large Dataset Migrations

```bash
# For large tables, use batched operations
export MIGRATION_BATCH_SIZE=1000
go run main.go migrate:all --connection=mysql
```

### 2. Parallel Database Setup

```bash
# Run migrations on multiple databases in parallel
(go run main.go migrate:all --connection=mysql &)
(go run main.go migrate:all --connection=postgres &)
wait
echo "All databases migrated!"
```

---

Next: [API Documentation](../api/README.md) | [Examples](../examples/README.md)
