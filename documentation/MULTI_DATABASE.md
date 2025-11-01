# Multi-Database Connection Guide

Complete guide to configuring and using multiple database connections simultaneously in Golang Starter Kit 2025.

## Overview

This project supports concurrent connections to multiple databases using a sophisticated Database Manager pattern. You can connect to MySQL, PostgreSQL, SQLite, and SQL Server simultaneously within a single application.

**Key Features:**
- **Connection Pooling**: Efficient connection management per database
- **Health Monitoring**: Real-time connection status tracking
- **Lazy Loading**: Connections established only when needed
- **Singleton Pattern**: Thread-safe single instance per connection
- **Transaction Support**: Cross-database transaction capabilities
- **CLI Support**: Migrations and seeders work across all connections

## Supported Databases

| Database | Type | Connection Name | Status |
|----------|------|-----------------|--------|
| MySQL (Primary) | mysql | `mysql` | Default |
| PostgreSQL | postgres | `postgres` | Supported |
| MySQL (Secondary) | mysql | `mysql_secondary` | Optional |
| SQLite | sqlite | `sqlite` | Supported |
| SQL Server | sqlserver | `sqlserver` | Supported |

## Environment Configuration

### Basic Setup (.env)

```bash
# Default Database Connection
DB_CONNECTION=mysql  # Which database to use as default

# Connection Pool Settings (applied to all connections)
MAX_IDLE_CONNS=10
MAX_OPEN_CONNS=200
CONN_MAX_LIFETIME=15m
CONN_MAX_IDLE_TIME=5m
```

### MySQL Primary Configuration

```bash
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=root
MYSQL_PASSWORD=your_password
MYSQL_CHARSET=utf8mb4
MYSQL_TIMEZONE=Local  # Local | UTC | Asia/Jakarta

# Connection Pool (optional, uses defaults if not set)
MYSQL_MAX_IDLE_CONNS=10
MYSQL_MAX_OPEN_CONNS=200
MYSQL_CONN_MAX_LIFETIME=15m
MYSQL_CONN_MAX_IDLE_TIME=5m
```

### PostgreSQL Configuration

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=golang_starter_kit_2025_pg
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_SSLMODE=disable  # disable | require | verify-full
POSTGRES_TIMEZONE=UTC

# Connection Pool (optional)
POSTGRES_MAX_IDLE_CONNS=10
POSTGRES_MAX_OPEN_CONNS=200
POSTGRES_CONN_MAX_LIFETIME=15m
POSTGRES_CONN_MAX_IDLE_TIME=5m
```

### MySQL Secondary Configuration (Optional)

Useful for read replicas, separate analytics databases, or multi-tenant setups:

```bash
MYSQL_SECONDARY_HOST=localhost
MYSQL_SECONDARY_PORT=3307
MYSQL_SECONDARY_DB=golang_starter_kit_2025_secondary
MYSQL_SECONDARY_USER=root
MYSQL_SECONDARY_PASSWORD=your_password
MYSQL_SECONDARY_CHARSET=utf8mb4
MYSQL_SECONDARY_TIMEZONE=Local

# Connection Pool (optional)
MYSQL_SECONDARY_MAX_IDLE_CONNS=5
MYSQL_SECONDARY_MAX_OPEN_CONNS=100
```

### SQLite Configuration (Optional)

```bash
SQLITE_DATABASE=storage/database.sqlite
```

### SQL Server Configuration (Optional)

```bash
SQLSERVER_HOST=localhost
SQLSERVER_PORT=1433
SQLSERVER_DB=golang_starter_kit_2025
SQLSERVER_USER=sa
SQLSERVER_PASSWORD=YourStrong@Passw0rd
```

## Database Manager Architecture

### Core Components

```
┌─────────────────────────────────────┐
│    config/database_config.go        │
│  (Load configs from environment)    │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│    database/manager.go              │
│  (Singleton connection manager)     │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│    facades/database.go              │
│  (Simplified access layer)          │
└─────────────────────────────────────┘
```

### Manager Features

**1. Singleton Pattern**
- One manager instance per application
- Thread-safe initialization using `sync.Once`
- Shared across all goroutines

**2. Connection Pool per Database**
- Independent pools for each database
- Configurable via environment variables
- Automatic cleanup on shutdown

**3. Lazy Connection**
- Databases connect only when first accessed
- Reduces startup time
- Saves resources for unused connections

**4. Health Monitoring**
- Ping-based health checks
- Connection statistics (sql.DBStats)
- Exposed via API endpoints

## Accessing Databases

### 1. Default Database (Backward Compatibility)

```go
// Using the default database (defined by DB_CONNECTION)
db := facades.DB

var users []models.User
db.Find(&users)
```

### 2. Specific Connection via Facade Helpers

```go
// MySQL (primary)
mysqlConn, err := facades.MySQL()
if err != nil {
    log.Fatal(err)
}
mysqlConn.DB.Find(&users)

// PostgreSQL
postgresConn, err := facades.PostgreSQL()
if err != nil {
    log.Fatal(err)
}
postgresConn.DB.Find(&products)

// MySQL Secondary
secondaryConn, err := facades.MySQLSecondary()
if err != nil {
    log.Fatal(err)
}
secondaryConn.DB.Find(&analytics)
```

### 3. Generic Connection Access

```go
// Get database manager
manager := facades.GetManager()

// Get specific connection
conn, err := manager.GetConnection("postgres")
if err != nil {
    log.Fatal(err)
}

// Use the connection
conn.DB.Create(&user)
```

### 4. In Repositories (Dependency Injection)

```go
// Repository accepts any *gorm.DB
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepositoryInterface {
    return &userRepository{db: db}
}

// In routes/web.go - choose which database to use
func RegisterRoutes(route *gin.Engine) {
    // Option 1: Use default database
    userRepo := repositories.NewUserRepository(facades.DB)

    // Option 2: Use PostgreSQL
    postgresConn, _ := facades.PostgreSQL()
    userRepo := repositories.NewUserRepository(postgresConn.DB)

    // Option 3: Use secondary MySQL
    secondaryConn, _ := facades.MySQLSecondary()
    analyticsRepo := repositories.NewAnalyticsRepository(secondaryConn.DB)

    // Create services and controllers as usual
    userService := services.NewUserService(userRepo)
    userController := controllers.NewUserController(*userService)
}
```

## CLI Commands with Multi-Database

All migration and seeder commands support the `--connection` flag:

### Migration Commands

```bash
# Run migrations on MySQL (default)
go run main.go migrate:all

# Run migrations on PostgreSQL
go run main.go migrate:all --connection=postgres

# Run migrations on secondary MySQL
go run main.go migrate:all --connection=mysql_secondary

# Check migration status on PostgreSQL
go run main.go migrate:status --connection=postgres

# Rollback on specific database
go run main.go rollback:batch --connection=postgres

# Fresh migration with seed on PostgreSQL
go run main.go migrate:fresh --seed --connection=postgres
```

### Seeder Commands

```bash
# Seed PostgreSQL database
go run main.go db:seed --connection=postgres

# Seed specific seeder on secondary MySQL
go run main.go db:seed --class=UserSeeder --connection=mysql_secondary

# Rollback seeders on PostgreSQL
go run main.go rollback:seeder --connection=postgres
```

### Database Wipe

```bash
# Wipe PostgreSQL database (requires confirmation)
go run main.go db:wipe --connection=postgres

# Force wipe (for CI/CD)
go run main.go db:wipe --connection=postgres --force
```

## Use Cases & Examples

### 1. Read/Write Splitting

**Scenario**: Use primary MySQL for writes, secondary for reads

```go
// In services/user_service.go
type UserService struct {
    writeRepo interfaces.UserRepositoryInterface
    readRepo  interfaces.UserRepositoryInterface
}

func NewUserService(writeDB, readDB *gorm.DB) *UserService {
    return &UserService{
        writeRepo: repositories.NewUserRepository(writeDB),
        readRepo:  repositories.NewUserRepository(readDB),
    }
}

func (s *UserService) CreateUser(user *models.User) error {
    // Write to primary database
    return s.writeRepo.Create(user)
}

func (s *UserService) GetUsers(page, limit int) ([]models.User, error) {
    // Read from secondary (replica)
    return s.readRepo.List(page, limit)
}

// In routes/web.go
mysqlPrimary := facades.DB
mysqlSecondary, _ := facades.MySQLSecondary()

userService := services.NewUserService(
    mysqlPrimary,      // writes
    mysqlSecondary.DB, // reads
)
```

### 2. Data Synchronization Between Databases

**Scenario**: Sync users from MySQL to PostgreSQL

```go
// In services/sync_service.go
type SyncService struct{}

func (s *SyncService) SyncUsersToPostgres() error {
    // Get connections
    mysqlConn := facades.DB
    postgresConn, err := facades.PostgreSQL()
    if err != nil {
        return err
    }

    // Fetch from MySQL
    var users []models.User
    if err := mysqlConn.Find(&users).Error; err != nil {
        return err
    }

    // Insert to PostgreSQL
    for _, user := range users {
        // Use FirstOrCreate to avoid duplicates
        var existingUser models.User
        result := postgresConn.DB.Where("email = ?", user.Email).
            FirstOrCreate(&existingUser, user)

        if result.Error != nil {
            log.Printf("Failed to sync user %s: %v", user.Email, result.Error)
        }
    }

    return nil
}
```

### 3. Multi-Tenant Architecture

**Scenario**: Different tenants use different databases

```go
// In middleware/tenant_middleware.go
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetHeader("X-Tenant-ID")

        var db *gorm.DB
        switch tenantID {
        case "tenant1":
            db = facades.DB // MySQL primary
        case "tenant2":
            conn, _ := facades.PostgreSQL()
            db = conn.DB
        case "tenant3":
            conn, _ := facades.MySQLSecondary()
            db = conn.DB
        default:
            c.JSON(400, gin.H{"error": "Invalid tenant"})
            c.Abort()
            return
        }

        c.Set("tenant_db", db)
        c.Next()
    }
}

// In controller
func (ctrl *UserController) List(c *gin.Context) {
    db := c.MustGet("tenant_db").(*gorm.DB)

    var users []models.User
    db.Find(&users)

    helpers.SuccessResponse(c, "Users retrieved", users)
}
```

### 4. Analytics Database

**Scenario**: Store analytics in separate PostgreSQL database

```go
// In services/analytics_service.go
type AnalyticsService struct {
    analyticsDB *gorm.DB
}

func NewAnalyticsService() *AnalyticsService {
    postgresConn, _ := facades.PostgreSQL()
    return &AnalyticsService{
        analyticsDB: postgresConn.DB,
    }
}

func (s *AnalyticsService) LogEvent(event *models.AnalyticsEvent) error {
    return s.analyticsDB.Create(event).Error
}

func (s *AnalyticsService) GetUserStats(userID uint) (*models.UserStats, error) {
    var stats models.UserStats
    err := s.analyticsDB.
        Model(&models.AnalyticsEvent{}).
        Where("user_id = ?", userID).
        Select("COUNT(*) as event_count, MAX(created_at) as last_activity").
        Scan(&stats).Error

    return &stats, err
}
```

### 5. Cross-Database Transactions

**Scenario**: Transaction across multiple databases

```go
func (s *OrderService) CreateOrderWithLog(order *models.Order) error {
    // Transaction on MySQL (orders database)
    err := facades.DB.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(order).Error; err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        return err
    }

    // Log to PostgreSQL (analytics database)
    postgresConn, _ := facades.PostgreSQL()
    return postgresConn.DB.Transaction(func(tx *gorm.DB) error {
        log := &models.OrderLog{
            OrderID:   order.ID,
            Action:    "created",
            Timestamp: time.Now(),
        }
        return tx.Create(log).Error
    })

    // Note: This is not a distributed transaction
    // If PostgreSQL fails, MySQL changes are already committed
    // Consider implementing saga pattern for true distributed transactions
}
```

## Health Check API

### Check All Database Connections

```bash
GET /health/databases
```

**Response:**
```json
{
    "overall_health": true,
    "connections": {
        "mysql": {
            "status": "healthy",
            "stats": {
                "max_open_connections": 200,
                "open_connections": 5,
                "in_use": 2,
                "idle": 3,
                "wait_count": 0,
                "wait_duration": 0,
                "max_idle_closed": 0,
                "max_lifetime_closed": 0
            }
        },
        "postgres": {
            "status": "healthy",
            "stats": {
                "max_open_connections": 200,
                "open_connections": 3,
                "in_use": 1,
                "idle": 2,
                "wait_count": 0,
                "wait_duration": 0,
                "max_idle_closed": 0,
                "max_lifetime_closed": 0
            }
        },
        "mysql_secondary": {
            "status": "disconnected"
        }
    }
}
```

**HTTP Status Codes:**
- `200`: All configured connections are healthy
- `503`: One or more connections unhealthy

### Check Default Database

```bash
GET /health
```

**Response:**
```json
{
    "message": "database is connected",
    "database": "golang_starter_kit_2025"
}
```

## Best Practices

### 1. Connection Pool Configuration

```bash
# For high-traffic production
MAX_IDLE_CONNS=20
MAX_OPEN_CONNS=500
CONN_MAX_LIFETIME=1h
CONN_MAX_IDLE_TIME=10m

# For development
MAX_IDLE_CONNS=5
MAX_OPEN_CONNS=50
CONN_MAX_LIFETIME=15m
CONN_MAX_IDLE_TIME=5m

# For background workers
MAX_IDLE_CONNS=2
MAX_OPEN_CONNS=10
CONN_MAX_LIFETIME=30m
CONN_MAX_IDLE_TIME=15m
```

### 2. Error Handling

```go
// Always check errors when getting connections
mysqlConn, err := facades.MySQL()
if err != nil {
    log.Printf("Failed to get MySQL connection: %v", err)
    // Fallback or return error
    return err
}

// Check database health before critical operations
if err := mysqlConn.DB.Exec("SELECT 1").Error; err != nil {
    log.Printf("Database health check failed: %v", err)
    return err
}
```

### 3. Resource Management

```go
// Defer cleanup for temporary connections
func ProcessData() error {
    manager := facades.GetManager()
    conn, err := manager.GetConnection("postgres")
    if err != nil {
        return err
    }

    // Process data...

    // Connection automatically returned to pool
    return nil
}
```

### 4. Monitoring

```go
// Regularly check connection stats
func MonitorDatabaseHealth() {
    manager := facades.GetManager()

    for _, connName := range []string{"mysql", "postgres", "mysql_secondary"} {
        conn, err := manager.GetConnection(connName)
        if err != nil {
            log.Printf("%s: disconnected", connName)
            continue
        }

        sqlDB, _ := conn.DB.DB()
        stats := sqlDB.Stats()

        log.Printf("%s stats: open=%d, in_use=%d, idle=%d, wait_count=%d",
            connName,
            stats.OpenConnections,
            stats.InUse,
            stats.Idle,
            stats.WaitCount,
        )

        // Alert if wait count is high
        if stats.WaitCount > 100 {
            log.Printf("WARNING: High wait count for %s: %d", connName, stats.WaitCount)
        }
    }
}
```

### 5. Data Consistency

```go
// Use transactions for related operations
func (s *UserService) CreateUserWithProfile(user *models.User, profile *models.Profile) error {
    return facades.DB.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(user).Error; err != nil {
            return err
        }

        profile.UserID = user.ID
        if err := tx.Create(profile).Error; err != nil {
            return err
        }

        return nil
    })
}
```

### 6. Performance Optimization

```go
// Use batch operations for large datasets
func (s *SyncService) BatchSyncToPostgres(users []models.User) error {
    postgresConn, _ := facades.PostgreSQL()

    // Batch size
    batchSize := 1000

    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }

        batch := users[i:end]
        if err := postgresConn.DB.CreateInBatches(batch, batchSize).Error; err != nil {
            return err
        }
    }

    return nil
}
```

## Troubleshooting

### Connection Refused

```bash
# Check database is running
systemctl status mysql
systemctl status postgresql

# Test connection manually
mysql -h localhost -P 3306 -u root -p
psql -h localhost -p 5432 -U postgres

# Check firewall
sudo ufw status
```

### Too Many Connections

```bash
# Reduce MAX_OPEN_CONNS in .env
MAX_OPEN_CONNS=100

# Check current connections on database
# MySQL:
SHOW PROCESSLIST;

# PostgreSQL:
SELECT * FROM pg_stat_activity;
```

### Slow Queries

```bash
# Enable query logging
# MySQL: Add to my.cnf
slow_query_log = 1
long_query_time = 2

# PostgreSQL: Add to postgresql.conf
log_statement = 'all'
log_duration = on
```

### Connection Pool Exhausted

```bash
# Increase pool size
MAX_OPEN_CONNS=500

# Or optimize queries to reduce connection time
# Use indexes, optimize joins, etc.
```

## Migration Considerations

### Database-Specific SQL

**MySQL:**
```sql
-- AUTO_INCREMENT for primary keys
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    ...
);
```

**PostgreSQL:**
```sql
-- SERIAL for primary keys
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    ...
);
```

### Handling Differences

Create separate migration files for each database type or use conditional SQL:

```sql
-- +++ UP Migration
-- This migration works for both MySQL and PostgreSQL
CREATE TABLE users (
    id BIGINT PRIMARY KEY,  -- Use GORM to handle auto-increment
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- --- DOWN Migration
DROP TABLE IF EXISTS users;
```

## Security Considerations

1. **Environment Variables**: Never commit credentials
2. **SSL/TLS**: Use encrypted connections in production
3. **Least Privilege**: Database users should have minimal permissions
4. **Connection Limits**: Set reasonable pool sizes
5. **Monitoring**: Log and alert on connection issues
6. **Backups**: Regular backups for all databases
7. **Firewalls**: Restrict database access by IP

## Further Reading

- [CLAUDE.md](../CLAUDE.md) - Complete technical documentation
- [DATABASE.md](./DATABASE.md) - Migration and seeder guide
- [GETTING_STARTED.md](./GETTING_STARTED.md) - Setup guide
- [GORM Documentation](https://gorm.io/docs/) - GORM features and best practices

---

**Need help?** Check our [GitHub Discussions](https://github.com/RahmatRafiq/golang_starter_kit_2025/discussions)
