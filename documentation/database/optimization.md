# Database Optimization

Panduan lengkap untuk optimasi database pada Golang Starter Kit.

## Overview

Database optimization adalah proses meningkatkan performa database untuk mengurangi response time, meningkatkan throughput, dan mengoptimalkan resource usage. Panduan ini mencakup berbagai teknik optimasi dari level query hingga infrastructure.

## Query Optimization

### 1. Index Optimization

#### Creating Effective Indexes

```sql
-- Primary key index (automatic)
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) UNIQUE,
    name VARCHAR(255),
    created_at TIMESTAMP
);

-- Single column index
CREATE INDEX idx_users_email ON users(email);

-- Composite index
CREATE INDEX idx_users_name_email ON users(name, email);

-- Partial index (MySQL 8.0+)
CREATE INDEX idx_active_users ON users(email) WHERE active = 1;
```

#### GORM Index Definition

```go
// app/models/user.go
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex;size:255"`
    Name      string    `gorm:"index;size:255"`
    Active    bool      `gorm:"index"`
    CreatedAt time.Time `gorm:"index"`
    UpdatedAt time.Time
}

// Composite index
type Product struct {
    ID         uint    `gorm:"primaryKey"`
    Name       string  `gorm:"size:255"`
    CategoryID uint    `gorm:"index:idx_category_active"`
    Active     bool    `gorm:"index:idx_category_active"`
    Price      float64 `gorm:"index"`
}
```

#### Index Analysis

```go
// Check index usage
func AnalyzeQuery(db *gorm.DB, query string) {
    var result []map[string]interface{}
    
    // MySQL EXPLAIN
    db.Raw("EXPLAIN " + query).Scan(&result)
    
    for _, row := range result {
        fmt.Printf("Table: %v, Type: %v, Key: %v, Rows: %v\n",
            row["table"], row["type"], row["key"], row["rows"])
    }
}

// Usage
AnalyzeQuery(facades.DB(), "SELECT * FROM users WHERE email = 'test@example.com'")
```

### 2. Query Optimization Techniques

#### Efficient WHERE Clauses

```go
// Good: Use indexed columns
var users []models.User
db.Where("email = ?", email).Find(&users)

// Bad: Function in WHERE clause
db.Where("UPPER(email) = ?", strings.ToUpper(email)).Find(&users)

// Good: Use LIKE efficiently
db.Where("name LIKE ?", name+"%").Find(&users) // Can use index

// Bad: Leading wildcard
db.Where("name LIKE ?", "%"+name+"%").Find(&users) // Cannot use index
```

#### Limit and Pagination

```go
// Efficient pagination
func GetUsersPaginated(page, limit int) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    // Count with same conditions
    query := facades.DB().Model(&models.User{}).Where("active = ?", true)
    query.Count(&total)
    
    // Get paginated data
    offset := (page - 1) * limit
    err := query.Offset(offset).Limit(limit).Find(&users).Error
    
    return users, total, err
}

// Cursor-based pagination (better for large datasets)
func GetUsersAfterCursor(cursor uint, limit int) ([]models.User, error) {
    var users []models.User
    err := facades.DB().Where("id > ?", cursor).
        Order("id ASC").
        Limit(limit).
        Find(&users).Error
    
    return users, err
}
```

#### Select Specific Fields

```go
// Good: Select only needed fields
var users []models.User
db.Select("id, name, email").Find(&users)

// Bad: Select all fields
db.Find(&users)

// Select with relationships
db.Select("users.id, users.name, roles.name as role_name").
   Joins("LEFT JOIN user_has_roles ON users.id = user_has_roles.user_id").
   Joins("LEFT JOIN roles ON user_has_roles.role_id = roles.id").
   Find(&users)
```

## Connection Optimization

### 1. Connection Pooling

```go
// config/database.go
func configureMySQLConnection() *gorm.DB {
    dsn := buildMySQLDSN()
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        log.Fatal("Failed to connect to MySQL:", err)
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("Failed to get underlying sql.DB:", err)
    }
    
    // Connection pool settings
    sqlDB.SetMaxIdleConns(10)           // Maximum idle connections
    sqlDB.SetMaxOpenConns(100)          // Maximum open connections
    sqlDB.SetConnMaxLifetime(time.Hour) // Connection lifetime
    sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Idle timeout
    
    return db
}
```

### 2. Connection Monitoring

```go
// database/monitor.go
package database

import (
    "context"
    "time"
    "database/sql"
)

type ConnectionMonitor struct {
    db *sql.DB
}

func NewConnectionMonitor(db *sql.DB) *ConnectionMonitor {
    return &ConnectionMonitor{db: db}
}

func (cm *ConnectionMonitor) GetStats() sql.DBStats {
    return cm.db.Stats()
}

func (cm *ConnectionMonitor) Monitor(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := cm.GetStats()
        log.Printf("DB Stats - Open: %d, InUse: %d, Idle: %d",
            stats.OpenConnections,
            stats.InUse,
            stats.Idle)
    }
}

func (cm *ConnectionMonitor) HealthCheck(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    return cm.db.PingContext(ctx)
}
```

## Caching Strategies

### 1. Query Result Caching

```go
// services/cached_user_service.go
package services

import (
    "encoding/json"
    "fmt"
    "time"
    "your-app/app/models"
    "your-app/cache"
)

type CachedUserService struct {
    userRepo interfaces.UserRepositoryInterface
    cache    cache.CacheInterface
}

func NewCachedUserService() *CachedUserService {
    return &CachedUserService{
        userRepo: repositories.NewUserRepository(),
        cache:    cache.NewRedisCache(),
    }
}

func (s *CachedUserService) GetUser(id uint) (*models.User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    
    // Try cache first
    if cached := s.cache.Get(cacheKey); cached != nil {
        var user models.User
        if err := json.Unmarshal(cached.([]byte), &user); err == nil {
            return &user, nil
        }
    }
    
    // Get from database
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    if data, err := json.Marshal(user); err == nil {
        s.cache.Set(cacheKey, data, 5*time.Minute)
    }
    
    return user, nil
}

func (s *CachedUserService) UpdateUser(id uint, data map[string]interface{}) error {
    err := s.userRepo.Update(id, data)
    if err != nil {
        return err
    }
    
    // Invalidate cache
    cacheKey := fmt.Sprintf("user:%d", id)
    s.cache.Delete(cacheKey)
    
    return nil
}
```

### 2. Database Query Caching

```go
// database/query_cache.go
package database

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "time"
    "your-app/cache"
)

type QueryCache struct {
    cache cache.CacheInterface
    ttl   time.Duration
}

func NewQueryCache(ttl time.Duration) *QueryCache {
    return &QueryCache{
        cache: cache.NewRedisCache(),
        ttl:   ttl,
    }
}

func (qc *QueryCache) getCacheKey(query string, args ...interface{}) string {
    data := fmt.Sprintf("%s:%v", query, args)
    hash := md5.Sum([]byte(data))
    return "query:" + hex.EncodeToString(hash[:])
}

func (qc *QueryCache) Get(query string, args []interface{}, dest interface{}) bool {
    key := qc.getCacheKey(query, args...)
    
    cached := qc.cache.Get(key)
    if cached == nil {
        return false
    }
    
    return json.Unmarshal(cached.([]byte), dest) == nil
}

func (qc *QueryCache) Set(query string, args []interface{}, data interface{}) {
    key := qc.getCacheKey(query, args...)
    
    if jsonData, err := json.Marshal(data); err == nil {
        qc.cache.Set(key, jsonData, qc.ttl)
    }
}

// Usage in repository
func (r *UserRepository) GetUsersCached(filters map[string]interface{}) ([]models.User, error) {
    query := "SELECT * FROM users WHERE active = ?"
    args := []interface{}{true}
    
    var users []models.User
    if r.queryCache.Get(query, args, &users) {
        return users, nil
    }
    
    err := r.db.Where("active = ?", true).Find(&users).Error
    if err != nil {
        return nil, err
    }
    
    r.queryCache.Set(query, args, users)
    return users, nil
}
```

## Read/Write Splitting

### 1. Master-Slave Configuration

```go
// database/read_write_manager.go
package database

import (
    "gorm.io/gorm"
    "gorm.io/plugin/dbresolver"
)

func SetupReadWriteSplitting(primary, replica *gorm.DB) *gorm.DB {
    primary.Use(dbresolver.Register(dbresolver.Config{
        // Replica databases for read operations
        Replicas: []gorm.Dialector{replica.Dialector},
        
        // Load balancing policy
        Policy: dbresolver.RandomPolicy{},
        
        // Tables that should use replicas for read operations
        Sources: []string{"users", "products", "categories"},
    }))
    
    return primary
}

// Usage
func setupDatabases() {
    primary := connectToMaster()
    replica := connectToReplica()
    
    db := SetupReadWriteSplitting(primary, replica)
    
    // Writes go to master
    db.Create(&user)
    
    // Reads can go to replica
    db.Find(&users)
    
    // Force master for reads
    db.Clauses(dbresolver.Write).Find(&users)
}
```

### 2. Manual Read/Write Splitting

```go
// database/connection_manager.go
package database

type ConnectionManager struct {
    master   *gorm.DB
    replicas []*gorm.DB
    current  int
}

func NewConnectionManager(master *gorm.DB, replicas []*gorm.DB) *ConnectionManager {
    return &ConnectionManager{
        master:   master,
        replicas: replicas,
    }
}

func (cm *ConnectionManager) Master() *gorm.DB {
    return cm.master
}

func (cm *ConnectionManager) Replica() *gorm.DB {
    if len(cm.replicas) == 0 {
        return cm.master
    }
    
    // Round-robin load balancing
    replica := cm.replicas[cm.current]
    cm.current = (cm.current + 1) % len(cm.replicas)
    
    return replica
}

// Usage in repository
func (r *UserRepository) Create(user *models.User) error {
    return r.connManager.Master().Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
    var user models.User
    err := r.connManager.Replica().First(&user, id).Error
    return &user, err
}
```

## Database Sharding

### 1. Horizontal Sharding

```go
// database/shard_manager.go
package database

import (
    "fmt"
    "hash/crc32"
)

type ShardManager struct {
    shards []*gorm.DB
}

func NewShardManager(shards []*gorm.DB) *ShardManager {
    return &ShardManager{shards: shards}
}

func (sm *ShardManager) GetShard(key string) *gorm.DB {
    hash := crc32.ChecksumIEEE([]byte(key))
    index := hash % uint32(len(sm.shards))
    return sm.shards[index]
}

// Usage for user sharding by email
func (r *UserRepository) CreateWithSharding(user *models.User) error {
    shard := r.shardManager.GetShard(user.Email)
    return shard.Create(user).Error
}

func (r *UserRepository) GetByEmailWithSharding(email string) (*models.User, error) {
    var user models.User
    shard := r.shardManager.GetShard(email)
    err := shard.Where("email = ?", email).First(&user).Error
    return &user, err
}
```

### 2. Vertical Sharding

```go
// Split tables across different databases
type DatabaseShards struct {
    UserDB    *gorm.DB // Users and authentication
    ProductDB *gorm.DB // Products and catalog
    OrderDB   *gorm.DB // Orders and transactions
}

func (ds *DatabaseShards) GetUserDB() *gorm.DB {
    return ds.UserDB
}

func (ds *DatabaseShards) GetProductDB() *gorm.DB {
    return ds.ProductDB
}

func (ds *DatabaseShards) GetOrderDB() *gorm.DB {
    return ds.OrderDB
}
```

## Performance Monitoring

### 1. Query Performance Monitoring

```go
// middleware/db_monitor.go
package middleware

import (
    "context"
    "log"
    "time"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type QueryMonitor struct {
    slowQueryThreshold time.Duration
}

func NewQueryMonitor(threshold time.Duration) *QueryMonitor {
    return &QueryMonitor{
        slowQueryThreshold: threshold,
    }
}

func (qm *QueryMonitor) Monitor() logger.Interface {
    return logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             qm.slowQueryThreshold,
            LogLevel:                  logger.Warn,
            IgnoreRecordNotFoundError: true,
            Colorful:                  true,
        },
    )
}

// Custom callback for detailed monitoring
func (qm *QueryMonitor) RegisterCallbacks(db *gorm.DB) {
    db.Callback().Query().Before("gorm:query").Register("monitor:before", qm.beforeQuery)
    db.Callback().Query().After("gorm:query").Register("monitor:after", qm.afterQuery)
}

func (qm *QueryMonitor) beforeQuery(db *gorm.DB) {
    db.InstanceSet("start_time", time.Now())
}

func (qm *QueryMonitor) afterQuery(db *gorm.DB) {
    startTime, exists := db.InstanceGet("start_time")
    if !exists {
        return
    }
    
    duration := time.Since(startTime.(time.Time))
    
    if duration > qm.slowQueryThreshold {
        log.Printf("SLOW QUERY [%v]: %s", duration, db.Statement.SQL.String())
    }
}
```

### 2. Database Metrics Collection

```go
// monitoring/db_metrics.go
package monitoring

import (
    "database/sql"
    "time"
)

type DBMetrics struct {
    db *sql.DB
}

type Metrics struct {
    OpenConnections     int
    InUseConnections    int
    IdleConnections     int
    WaitCount           int64
    WaitDuration        time.Duration
    MaxIdleClosed       int64
    MaxLifetimeClosed   int64
    MaxOpenConnections  int
    MaxIdleConnections  int
    ConnMaxLifetime     time.Duration
}

func NewDBMetrics(db *sql.DB) *DBMetrics {
    return &DBMetrics{db: db}
}

func (dm *DBMetrics) CollectMetrics() Metrics {
    stats := dm.db.Stats()
    
    return Metrics{
        OpenConnections:    stats.OpenConnections,
        InUseConnections:   stats.InUse,
        IdleConnections:    stats.Idle,
        WaitCount:          stats.WaitCount,
        WaitDuration:       stats.WaitDuration,
        MaxIdleClosed:      stats.MaxIdleClosed,
        MaxLifetimeClosed:  stats.MaxLifetimeClosed,
        MaxOpenConnections: stats.MaxOpenConnections,
        MaxIdleConnections: stats.MaxIdleConnections,
        ConnMaxLifetime:    stats.ConnMaxLifetime,
    }
}

func (dm *DBMetrics) StartMonitoring(interval time.Duration) {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            metrics := dm.CollectMetrics()
            // Send to monitoring system (Prometheus, etc.)
            dm.sendToMonitoring(metrics)
        }
    }()
}
```

## Batch Operations

### 1. Bulk Insert

```go
// Efficient bulk insert
func BulkCreateUsers(users []models.User) error {
    batchSize := 1000
    
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        if err := facades.DB().CreateInBatches(batch, batchSize).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

### 2. Bulk Update

```go
// Bulk update with case statement
func BulkUpdateUserStatus(updates map[uint]bool) error {
    if len(updates) == 0 {
        return nil
    }
    
    var ids []uint
    var cases []string
    
    for id, status := range updates {
        ids = append(ids, id)
        cases = append(cases, fmt.Sprintf("WHEN %d THEN %t", id, status))
    }
    
    query := fmt.Sprintf(`
        UPDATE users 
        SET active = CASE id %s END 
        WHERE id IN (%s)`,
        strings.Join(cases, " "),
        strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]"))
    
    return facades.DB().Exec(query).Error
}
```

## Best Practices

### 1. Database Design

- Use appropriate data types
- Normalize when appropriate, denormalize for performance
- Create proper foreign key constraints
- Use appropriate storage engines (InnoDB for MySQL)

### 2. Index Strategy

- Index frequently queried columns
- Create composite indexes for multi-column queries
- Avoid over-indexing (impacts write performance)
- Monitor index usage and remove unused indexes

### 3. Query Optimization

- Use EXPLAIN to analyze query execution plans
- Avoid N+1 queries (use eager loading)
- Use appropriate LIMIT clauses
- Cache frequently accessed data

### 4. Connection Management

- Configure appropriate connection pool sizes
- Monitor connection usage
- Set proper connection timeouts
- Use connection health checks

### 5. Monitoring and Alerting

- Monitor slow queries
- Track database performance metrics
- Set up alerts for performance degradation
- Regular database maintenance

---

Untuk informasi lebih lanjut, lihat:
- [Database Guide](README.md)
- [Multi-Database Guide](multi-database.md)
- [Development Guide](../guides/development.md)
