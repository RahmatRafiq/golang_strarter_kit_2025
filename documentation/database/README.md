# Database Configuration Guide

## Overview

Aplikasi ini mendukung multiple database connections dengan arsitektur yang fleksibel dan mudah dikonfigurasi. Anda dapat menggunakan berbagai jenis database secara bersamaan untuk kebutuhan yang berbeda.

## Supported Databases

| Database | Status | Primary | Secondary | Use Case |
|----------|--------|---------|-----------|----------|
| **MySQL/MariaDB** | ✅ Production | ✅ | ✅ | General purpose, high performance |
| **PostgreSQL** | ✅ Production | ✅ | ❌ | Advanced features, JSON support |
| **SQLite** | ✅ Development | ✅ | ❌ | Testing, prototyping |
| **SQL Server** | 🧪 Beta | ✅ | ❌ | Enterprise integration |

## Database Manager Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Database Manager                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐   │
│  │   MySQL     │  │ PostgreSQL  │  │  MySQL Secondary    │   │
│  │ (Primary)   │  │             │  │   (Optional)        │   │
│  │ Port: 3306  │  │ Port: 5432  │  │   Port: 3307        │   │
│  └─────────────┘  └─────────────┘  └─────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│ • Connection Pooling  • Health Monitoring  • Auto Failover  │
│ • Load Balancing     • Transaction Support • Query Logging  │
└─────────────────────────────────────────────────────────────┘
```

## Environment Configuration

### Complete .env Configuration

```env
# Default Database Connection
DB_CONNECTION=mysql

# MySQL Primary Configuration
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=app_user
MYSQL_PASSWORD=secure_password
MYSQL_CHARSET=utf8mb4
MYSQL_TIMEZONE=Local
MYSQL_MAX_IDLE_CONNS=10
MYSQL_MAX_OPEN_CONNS=200
MYSQL_CONN_MAX_LIFETIME=15m
MYSQL_CONN_MAX_IDLE_TIME=5m

# PostgreSQL Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=golang_starter_kit_2025_pg
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secure_password
POSTGRES_SSLMODE=disable
POSTGRES_TIMEZONE=UTC
POSTGRES_MAX_IDLE_CONNS=10
POSTGRES_MAX_OPEN_CONNS=200
POSTGRES_CONN_MAX_LIFETIME=15m
POSTGRES_CONN_MAX_IDLE_TIME=5m

# MySQL Secondary Configuration (Optional)
MYSQL_SECONDARY_HOST=localhost
MYSQL_SECONDARY_PORT=3307
MYSQL_SECONDARY_DB=golang_starter_kit_2025_secondary
MYSQL_SECONDARY_USER=app_user
MYSQL_SECONDARY_PASSWORD=secure_password
MYSQL_SECONDARY_CHARSET=utf8mb4
MYSQL_SECONDARY_TIMEZONE=Local
MYSQL_SECONDARY_MAX_IDLE_CONNS=10
MYSQL_SECONDARY_MAX_OPEN_CONNS=200
MYSQL_SECONDARY_CONN_MAX_LIFETIME=15m
MYSQL_SECONDARY_CONN_MAX_IDLE_TIME=5m
```

### Configuration Parameters Explained

#### Connection Pool Settings
- **MAX_IDLE_CONNS**: Jumlah maksimum koneksi idle dalam pool
- **MAX_OPEN_CONNS**: Jumlah maksimum koneksi yang dapat dibuka
- **CONN_MAX_LIFETIME**: Berapa lama koneksi dapat digunakan sebelum ditutup
- **CONN_MAX_IDLE_TIME**: Berapa lama koneksi idle sebelum ditutup

#### Database-Specific Settings
- **CHARSET** (MySQL): Character set yang digunakan
- **TIMEZONE**: Timezone untuk timestamp
- **SSLMODE** (PostgreSQL): Mode SSL connection (disable, require, verify-ca, verify-full)

## Using Multiple Databases

### 1. Basic Usage

```go
package main

import (
    "golang_starter_kit_2025/facades"
)

func main() {
    // Default connection (MySQL)
    db := facades.GetDB()
    
    // Specific connection
    mysqlConn, err := facades.MySQL()
    if err != nil {
        panic(err)
    }
    
    postgresConn, err := facades.PostgreSQL()
    if err != nil {
        panic(err)
    }
    
    // Use connections...
}
```

### 2. Service Implementation

```go
type UserService struct {
    primaryDB   *database.Connection
    analyticsDB *database.Connection
}

func NewUserService() *UserService {
    primary, _ := facades.MySQL()
    analytics, _ := facades.PostgreSQL()
    
    return &UserService{
        primaryDB:   primary,
        analyticsDB: analytics,
    }
}

func (s *UserService) CreateUser(user *models.User) error {
    // Create in primary database
    if err := s.primaryDB.DB.Create(user).Error; err != nil {
        return err
    }
    
    // Log to analytics database
    analytics := &models.UserAnalytics{
        UserID:    user.ID,
        Action:    "created",
        Timestamp: time.Now(),
    }
    s.analyticsDB.DB.Create(analytics)
    
    return nil
}
```

### 3. Migration on Multiple Databases

```bash
# Run migrations on all databases
go run main.go migrate:all --connection=mysql
go run main.go migrate:all --connection=postgres
go run main.go migrate:all --connection=mysql_secondary

# Check migration status
go run main.go db:status --connection=mysql
go run main.go db:status --connection=postgres
```

## Database-Specific Features

### MySQL/MariaDB Features

#### Advantages
- **High Performance**: Optimized for speed
- **Wide Compatibility**: Supported everywhere
- **Mature Ecosystem**: Extensive tooling
- **Replication**: Master-slave setup

#### Use Cases
- Web applications
- E-commerce platforms
- Content management systems
- High-traffic applications

#### Configuration Tips
```env
# Performance tuning
MYSQL_MAX_OPEN_CONNS=100
MYSQL_MAX_IDLE_CONNS=20
MYSQL_CONN_MAX_LIFETIME=30m

# For high traffic
MYSQL_MAX_OPEN_CONNS=500
MYSQL_MAX_IDLE_CONNS=50
```

### PostgreSQL Features

#### Advantages
- **Advanced Features**: JSON, arrays, custom types
- **ACID Compliance**: Full transaction support
- **Extensibility**: Custom functions, extensions
- **Analytics**: Window functions, CTEs

#### Use Cases
- Analytics applications
- Complex queries
- JSON document storage
- Geospatial applications

#### Configuration Tips
```env
# Analytics workload
POSTGRES_MAX_OPEN_CONNS=50
POSTGRES_MAX_IDLE_CONNS=10
POSTGRES_CONN_MAX_LIFETIME=1h

# JSON-heavy applications
POSTGRES_TIMEZONE=UTC
POSTGRES_SSLMODE=require
```

## Database Management Commands

### Connection Management

```bash
# List all available connections
go run main.go db:connections

# Check connection health
go run main.go db:status --connection=mysql
go run main.go db:status --connection=postgres

# Connection statistics
curl http://localhost:8080/api/admin/db/stats
```

### Migration Commands

```bash
# Create migration for specific database
go run main.go make:migration create_users_table

# Run migration on specific connection
go run main.go migrate --file=20240629_create_users_table --connection=mysql
go run main.go migrate --file=20240629_create_users_table --connection=postgres

# Rollback on specific connection
go run main.go rollback --file=20240629_create_users_table --connection=mysql
```

### Data Synchronization

```bash
# Sync data between databases
go run main.go db:sync --from=mysql --to=postgres --table=users

# Backup and restore
go run main.go db:backup --connection=mysql --file=backup.sql
go run main.go db:restore --connection=postgres --file=backup.sql
```

## Performance Optimization

### Connection Pool Tuning

#### For High Traffic Applications
```env
# MySQL High Traffic
MYSQL_MAX_OPEN_CONNS=1000
MYSQL_MAX_IDLE_CONNS=100
MYSQL_CONN_MAX_LIFETIME=10m
MYSQL_CONN_MAX_IDLE_TIME=2m

# PostgreSQL Analytics
POSTGRES_MAX_OPEN_CONNS=200
POSTGRES_MAX_IDLE_CONNS=50
POSTGRES_CONN_MAX_LIFETIME=30m
```

#### For Low Resource Applications
```env
# Minimal configuration
MYSQL_MAX_OPEN_CONNS=25
MYSQL_MAX_IDLE_CONNS=5
POSTGRES_MAX_OPEN_CONNS=10
POSTGRES_MAX_IDLE_CONNS=2
```

### Monitoring and Debugging

#### Enable Query Logging
```go
// In database manager configuration
gormConfig := &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
    PrepareStmt: true,
}
```

#### Health Check Endpoint
```bash
curl http://localhost:8080/health/database
```

Response:
```json
{
  "databases": {
    "mysql": {
      "status": "healthy",
      "connections": {
        "open": 15,
        "idle": 8,
        "in_use": 7
      },
      "response_time": "2ms"
    },
    "postgres": {
      "status": "healthy",
      "connections": {
        "open": 5,
        "idle": 3,
        "in_use": 2
      },
      "response_time": "3ms"
    }
  }
}
```

## Best Practices

### 1. Connection Selection Strategy

```go
// Use appropriate database for the task
func (s *AnalyticsService) RecordEvent(event *Event) error {
    // Use PostgreSQL for analytics (better for complex queries)
    return s.postgresDB.DB.Create(event).Error
}

func (s *UserService) GetUser(id uint) (*User, error) {
    // Use MySQL for transactional data (faster for simple queries)
    var user User
    err := s.mysqlDB.DB.First(&user, id).Error
    return &user, err
}
```

### 2. Transaction Management

```go
func (s *OrderService) CreateOrder(order *Order) error {
    // Use transaction on primary database
    return s.primaryDB.Transaction(func(tx *gorm.DB) error {
        // Create order
        if err := tx.Create(order).Error; err != nil {
            return err
        }
        
        // Update inventory
        if err := tx.Model(&Product{}).Where("id = ?", order.ProductID).
            Update("stock", gorm.Expr("stock - ?", order.Quantity)).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

### 3. Error Handling

```go
func (s *DatabaseService) GetConnection(name string) (*database.Connection, error) {
    conn, err := facades.GetConnection(name)
    if err != nil {
        // Log error and fallback to default
        log.Printf("Failed to get connection '%s': %v", name, err)
        return facades.GetDefaultConnection()
    }
    return conn, nil
}
```

## Troubleshooting

### Common Issues

#### 1. Connection Pool Exhaustion
```
Error: database connection pool exhausted
```
**Solution**: Increase MAX_OPEN_CONNS or optimize query performance

#### 2. Connection Timeout
```
Error: dial tcp: connection timeout
```
**Solution**: Check network connectivity and database server status

#### 3. Authentication Failed
```
Error: Access denied for user
```
**Solution**: Verify username, password, and database permissions

### Debugging Tools

```bash
# Check active connections
go run main.go db:stats --connection=mysql

# Test connectivity
go run main.go db:ping --connection=postgres

# Monitor performance
go run main.go db:monitor --connection=mysql --duration=5m
```

---

Next: [Migration Guide](migrations.md) | [Seeder Guide](seeders.md)
