# Installation Guide

## Prasyarat Sistem

Sebelum memulai, pastikan sistem Anda memiliki komponen berikut:

### Wajib
- **Go 1.24.3+** - [Download Go](https://golang.org/dl/)
- **Git** - Version control system
- **Database Server** - Minimal satu dari:
  - MySQL 8.0+ atau MariaDB 10.5+
  - PostgreSQL 13+

### Opsional
- **Air** - Hot reload untuk development
- **Postman/Insomnia** - API testing
- **Docker** - Containerization (opsional)

## Instalasi Step-by-Step

### 1. Clone Repository

```bash
git clone https://github.com/RahmatRafiq/golang_starter_kit_2025.git
cd golang_starter_kit_2025
```

### 2. Install Dependencies

```bash
# Install Go modules
go mod tidy

# Install development tools
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

### 3. Setup Database

#### MySQL/MariaDB
```sql
-- Buat database utama
CREATE DATABASE golang_starter_kit_2025;

-- Buat database secondary (opsional)
CREATE DATABASE golang_starter_kit_2025_secondary;

-- Buat user dan berikan permissions
CREATE USER 'app_user'@'localhost' IDENTIFIED BY 'secure_password';
GRANT ALL PRIVILEGES ON golang_starter_kit_2025.* TO 'app_user'@'localhost';
GRANT ALL PRIVILEGES ON golang_starter_kit_2025_secondary.* TO 'app_user'@'localhost';
FLUSH PRIVILEGES;
```

#### PostgreSQL
```sql
-- Buat database
CREATE DATABASE golang_starter_kit_2025_pg;

-- Buat user dan berikan permissions
CREATE USER app_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE golang_starter_kit_2025_pg TO app_user;
```

### 4. Konfigurasi Environment

```bash
# Copy file environment example
cp .env.example .env

# Edit file .env dengan konfigurasi Anda
nano .env  # atau editor lainnya
```

#### Konfigurasi .env Minimum
```env
# Application
APP_NAME="My Golang App"
APP_PORT=8080

# Default Database Connection
DB_CONNECTION=mysql

# MySQL Configuration
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=app_user
MYSQL_PASSWORD=secure_password

# PostgreSQL Configuration (jika digunakan)
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=golang_starter_kit_2025_pg
POSTGRES_USER=app_user
POSTGRES_PASSWORD=secure_password

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key-here
```

### 5. Generate API Documentation

```bash
swag init
```

### 6. Run Migrations

```bash
# Check available connections
go run main.go db:connections

# Run migrations on MySQL
go run main.go migrate:all --connection=mysql

# Run migrations on PostgreSQL (jika dikonfigurasi)
go run main.go migrate:all --connection=postgres
```

### 7. Seed Data (Opsional)

```bash
# Run seeders untuk data testing
go run main.go db:seed
```

### 8. Start Application

#### Development Mode (Hot Reload)
```bash
air
```

#### Production Mode
```bash
go build -o app main.go
./app
```

## Verifikasi Instalasi

### 1. Check Application Health
```bash
curl http://localhost:8080/health
```

Response yang diharapkan:
```json
{
  "status": "ok",
  "timestamp": "2025-06-29T10:30:00Z",
  "database": {
    "mysql": "connected",
    "postgres": "connected"
  }
}
```

### 2. Check Database Connections
```bash
go run main.go db:status --connection=mysql
go run main.go db:status --connection=postgres
```

### 3. Access API Documentation
Buka browser dan akses: http://localhost:8080/swagger/index.html

## Troubleshooting

### Common Issues

#### 1. Database Connection Error
```
Error: failed to connect to database 'mysql': dial tcp: connection refused
```

**Solusi:**
- Pastikan database server berjalan
- Check konfigurasi host, port, username, password di .env
- Pastikan firewall tidak memblok koneksi

#### 2. Migration Error
```
Error: failed to get connection 'mysql': invalid configuration
```

**Solusi:**
- Verifikasi konfigurasi database di .env
- Pastikan database dan user sudah dibuat
- Check permissions user database

#### 3. Module Import Error
```
go: module not found
```

**Solusi:**
```bash
go mod tidy
go clean -modcache
go mod download
```

#### 4. Air Not Found
```
air: command not found
```

**Solusi:**
```bash
# Install Air globally
go install github.com/air-verse/air@latest

# Atau jalankan langsung dengan go run
go run main.go
```

#### 5. Swagger Documentation Empty
**Solusi:**
```bash
# Regenerate swagger docs
swag init --parseInternal

# Restart application
air
```

### Database-Specific Issues

#### MySQL Issues
```bash
# Check MySQL service
sudo systemctl status mysql

# Start MySQL service
sudo systemctl start mysql

# Check MySQL logs
sudo tail -f /var/log/mysql/error.log
```

#### PostgreSQL Issues
```bash
# Check PostgreSQL service
sudo systemctl status postgresql

# Start PostgreSQL service
sudo systemctl start postgresql

# Check PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-13-main.log
```

## Next Steps

Setelah instalasi berhasil:

1. 📚 Baca [Database Guide](../database/README.md) untuk memahami multi-database setup
2. 🏗️ Pelajari [Architecture Guide](architecture.md) untuk memahami struktur aplikasi
3. 💡 Lihat [Examples](../examples/README.md) untuk contoh implementasi
4. 🔌 Eksplorasi [API Documentation](../api/README.md) untuk endpoint yang tersedia

## Development Tools

### Recommended VS Code Extensions
- **Go** - Rich Go language support
- **Thunder Client** - API testing
- **Database Client** - Database management
- **GitLens** - Git integration
- **Better Comments** - Comment highlighting

### Useful Commands
```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Build for production
go build -ldflags="-s -w" -o app main.go

# Generate test coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

**🎉 Selamat! Aplikasi Anda sudah siap untuk development.**

Jika mengalami masalah, silakan buka [issue di GitHub](https://github.com/RahmatRafiq/golang_starter_kit_2025/issues) atau lihat dokumentasi lainnya.
