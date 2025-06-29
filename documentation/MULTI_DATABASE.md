# Multi-Database Connection

Proyek ini mendukung koneksi ke beberapa database secara bersamaan. Anda dapat menghubungkan aplikasi ke MySQL, PostgreSQL, dan database lainnya dalam satu aplikasi.

## Konfigurasi Environment untuk Multi-Database

Dalam file `.env`, Anda dapat mengkonfigurasi beberapa koneksi database:

```env
# Default Database Connection
DB_CONNECTION=mysql

# MySQL Primary Connection
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=root
MYSQL_PASSWORD=root

# MySQL Secondary Connection (optional)
MYSQL2_HOST=localhost
MYSQL2_PORT=3307
MYSQL2_DB=golang_starter_kit_2025_secondary
MYSQL2_USER=root
MYSQL2_PASSWORD=root

# PostgreSQL Connection
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=golang_starter_kit_2025
POSTGRES_USER=postgres
POSTGRES_PASSWORD=root
POSTGRES_SSLMODE=disable

# MongoDB Connection (optional)
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DB=golang_starter_kit_2025
MONGO_USER=root
MONGO_PASS=secure_password_here
MONGO_COLL="${APP_ENV}"
```

## Database Manager

Proyek ini menggunakan Database Manager untuk mengelola multiple koneksi database. Manager ini menyediakan:

- **Connection Pooling**: Pengelolaan pool koneksi yang efisien
- **Health Monitoring**: Pemantauan status kesehatan koneksi
- **Automatic Reconnection**: Reconnect otomatis jika koneksi terputus
- **Transaction Support**: Dukungan transaksi di berbagai database

### Fitur-Fitur Database Manager

#### 1. Koneksi Otomatis
```go
// Mendapatkan koneksi default
db := facades.DB()

// Mendapatkan koneksi spesifik
mysqlDB := facades.MySQL()
postgresDB := facades.PostgreSQL()
```

#### 2. Penggunaan Repository Pattern
```go
userRepo := repositories.NewUserRepository()

// Operasi di MySQL
err := userRepo.CreateOnMySQL(user)

// Operasi di PostgreSQL
err := userRepo.CreateOnPostgreSQL(user)
```

#### 3. Database Service untuk Operasi Cross-Database
```go
dbService := services.NewDatabaseService()

// Eksekusi di MySQL
err := dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
    return db.Create(&user).Error
})

// Eksekusi di PostgreSQL
err := dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
    return db.Create(&user).Error
})
```

## API Endpoints untuk Database Management

### 1. Status Koneksi Database
```http
GET /api/database/status
```
Mengembalikan status dan statistik semua koneksi database.

### 2. Health Check Database
```http
GET /api/database/health
```
Melakukan health check pada semua koneksi database yang dikonfigurasi.

### 3. Test Koneksi Spesifik
```http
GET /api/database/test?connection=mysql
```
Menguji koneksi database tertentu.

**Parameter:**
- `connection`: Nama koneksi (mysql, postgres, mysql_secondary)

### 4. Sinkronisasi Data
```http
POST /api/database/sync?source=mysql&target=postgres
```
Melakukan sinkronisasi data antara dua database.

**Parameter:**
- `source`: Database sumber
- `target`: Database tujuan
- `table`: Tabel spesifik (opsional)

## Contoh Penggunaan Multi-Database

### 1. Membuat User di Multiple Database
```go
func CreateUserMultiDB(user *models.User) error {
    dbService := services.NewDatabaseService()
    
    // Create di MySQL
    err := dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
        return db.Create(user).Error
    })
    if err != nil {
        return err
    }
    
    // Create di PostgreSQL
    err = dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
        return db.Create(user).Error
    })
    
    return err
}
```

### 2. Sinkronisasi Data Antar Database
```go
func SyncUsers() error {
    dbService := services.NewDatabaseService()
    
    return dbService.SyncData(func(mysql, postgres *gorm.DB) error {
        var users []models.User
        
        // Ambil data dari MySQL
        if err := mysql.Find(&users).Error; err != nil {
            return err
        }
        
        // Insert ke PostgreSQL
        for _, user := range users {
            postgres.FirstOrCreate(&user, models.User{Email: user.Email})
        }
        
        return nil
    })
}
```

### 3. Menggunakan Transaction di Multiple Database
```go
func TransferData(userID uint) error {
    dbService := services.NewDatabaseService()
    
    // Transaction di MySQL
    err := dbService.ExecuteOnMySQL(func(db *gorm.DB) error {
        return db.Transaction(func(tx *gorm.DB) error {
            // Operasi di MySQL
            return tx.Model(&models.User{}).Where("id = ?", userID).Update("status", "transferred").Error
        })
    })
    
    if err != nil {
        return err
    }
    
    // Transaction di PostgreSQL
    return dbService.ExecuteOnPostgreSQL(func(db *gorm.DB) error {
        return db.Transaction(func(tx *gorm.DB) error {
            // Operasi di PostgreSQL
            return tx.Create(&models.TransferLog{UserID: userID}).Error
        })
    })
}
```

## Best Practices

### 1. Connection Management
- Gunakan connection pooling untuk performa optimal
- Tutup koneksi yang tidak digunakan
- Monitor health database secara berkala

### 2. Data Consistency
- Gunakan transaction untuk operasi critical
- Implementasikan retry mechanism untuk operasi gagal
- Backup data secara berkala

### 3. Performance Optimization
- Gunakan index yang tepat di setiap database
- Optimalkan query berdasarkan karakteristik database
- Monitor query performance

### 4. Security
- Gunakan environment variables untuk credentials
- Implementasikan proper authentication
- Encrypt sensitive data

## Troubleshooting

### Connection Issues
```bash
# Test koneksi manual
go run examples/multi_database_usage.go

# Check health via API
curl http://localhost:8080/api/database/health
```

### Performance Issues
```bash
# Monitor connection stats
curl http://localhost:8080/api/database/status
```

## File Contoh

Lihat file `examples/multi_database_usage.go` untuk contoh lengkap penggunaan multi-database connection.
