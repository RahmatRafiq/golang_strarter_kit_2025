# Multi-Database Support Documentation

Sistem ini mendukung multiple database connections dengan arsitektur yang modular seperti Laravel.

## Supported Databases

- **MySQL/MariaDB** (Primary & Secondary)
- **PostgreSQL**
- **SQLite** (Optional)
- **SQL Server** (Optional)

## Configuration

### Environment Variables

Konfigurasi database dilakukan melalui environment variables dalam file `.env`:

```env
# Default connection
DB_CONNECTION=mysql

# MySQL Primary
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=your_database
MYSQL_USER=your_username
MYSQL_PASSWORD=your_password

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=your_pg_database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_pg_password

# MySQL Secondary (Optional)
MYSQL_SECONDARY_HOST=localhost
MYSQL_SECONDARY_PORT=3307
MYSQL_SECONDARY_DB=secondary_database
```

## Usage Examples

### 1. Basic Database Operations

```go
// Get database service
dbService := services.NewDatabaseService()

// MySQL operations
mysqlDB, err := dbService.GetMySQL()
if err != nil {
    log.Fatal(err)
}

// PostgreSQL operations
postgresDB, err := dbService.GetPostgreSQL()
if err != nil {
    log.Fatal(err)
}

// Execute on specific database
err = dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
    return db.Create(&user).Error
})
```

### 2. Repository Pattern

```go
// Using repository pattern
userRepo := repositories.NewUserRepository()

// Create user in MySQL
err := userRepo.CreateOnMySQL(&user)

// Create user in PostgreSQL
err := userRepo.CreateOnPostgreSQL(&user)

// Create user in both databases
err := userRepo.CreateOnBothDatabases(&user)

// Sync data between databases
err := userRepo.SyncAllUsersFromMySQLToPostgreSQL()
```

### 3. Transactions

```go
// Transaction on MySQL
err := dbService.TransactionOnMySQL(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    return tx.Create(&profile).Error
})

// Transaction on PostgreSQL
err := dbService.TransactionOnPostgreSQL(func(tx *gorm.DB) error {
    // Your transaction logic here
    return nil
})
```

### 4. Data Synchronization

```go
// Sync data between MySQL and PostgreSQL
err := dbService.SyncData(func(mysql, postgres *gorm.DB) error {
    var users []models.User
    mysql.Find(&users)
    
    for _, user := range users {
        postgres.Create(&user)
    }
    return nil
})
```

## CLI Commands

### Migration Commands

```bash
# Run migration on specific connection
go run main.go migrate --file=create_users_table --connection=mysql
go run main.go migrate --file=create_users_table --connection=postgres

# Run all migrations
go run main.go migrate:all --connection=mysql
go run main.go migrate:all --connection=postgres

# Rollback migrations
go run main.go rollback --file=create_users_table --connection=postgres
go run main.go rollback:all --connection=postgres

# Fresh migrations (rollback all then migrate all)
go run main.go migrate:fresh --connection=mysql
```

### Database Management Commands

```bash
# List available connections
go run main.go db:connections

# Check database status
go run main.go db:status --connection=mysql
go run main.go db:status --connection=postgres
```

## API Endpoints

### Health Check

```bash
# Check all database connections
GET /health/databases

# Detailed database status (requires auth)
GET /api/database/status

# Health check for specific connection
GET /api/database/health
```

### Connection Testing

```bash
# Test specific connection
GET /api/database/test?connection=mysql
GET /api/database/test?connection=postgres
```

### Data Synchronization

```bash
# Sync data between databases
POST /api/database/sync?source=mysql&target=postgres
```

## Architecture Benefits

### 1. Modular Design
- Configurable database connections
- Easy to add new database types
- Separation of concerns

### 2. Laravel-like Experience
- Facade pattern for easy access
- Repository pattern for data operations
- Service layer for business logic

### 3. Flexibility
- Multiple connections per database type
- Connection-specific operations
- Cross-database data synchronization

### 4. Production Ready
- Connection pooling
- Health monitoring
- Error handling
- Transaction support

## Database-Specific Features

### MySQL/MariaDB
- Auto-increment primary keys
- UTF8MB4 charset support
- Connection pooling optimized for MySQL

### PostgreSQL
- Serial primary keys
- UUID support
- JSONB data types
- Advanced indexing

## Error Handling

```go
// Connection-specific error handling
conn, err := facades.GetConnection("postgres")
if err != nil {
    // Handle connection error
    log.Printf("PostgreSQL connection failed: %v", err)
    return
}

// Database operation error handling
err = dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
    return db.Create(&user).Error
})
if err != nil {
    log.Printf("PostgreSQL operation failed: %v", err)
}
```

## Best Practices

1. **Use Repository Pattern**: Encapsulate database operations in repositories
2. **Handle Connections Gracefully**: Always check for connection errors
3. **Use Transactions**: For operations that span multiple tables
4. **Monitor Connections**: Use health check endpoints
5. **Connection Pooling**: Configure appropriate pool sizes
6. **Data Synchronization**: Implement proper sync strategies for multi-database scenarios

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Check database server is running
   - Verify connection parameters
   - Check firewall settings

2. **Authentication Failed**
   - Verify username/password
   - Check user permissions
   - Ensure database exists

3. **Migration Errors**
   - Check SQL syntax for specific database
   - Verify table doesn't already exist
   - Check foreign key constraints

### Debugging

```go
// Enable detailed logging
facades.ConnectDB() // Will show connection logs

// Check connection stats
stats, err := facades.GetManager().GetConnectionStats("postgres")
fmt.Printf("PostgreSQL Stats: %+v\n", stats)
```
