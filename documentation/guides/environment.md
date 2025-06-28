# Environment Configuration

Panduan lengkap untuk konfigurasi environment pada Golang Starter Kit.

## Overview

Golang Starter Kit menggunakan file `.env` untuk manajemen konfigurasi environment. Semua konfigurasi aplikasi diatur melalui environment variables yang mudah dikelola untuk berbagai deployment scenario.

## Setup Environment

### 1. Copy Environment Template

```bash
cp .env.example .env
```

### 2. Basic Configuration

```bash
# Application
APP_NAME="Golang Starter Kit"
APP_ENV=local
APP_DEBUG=true
APP_URL=http://localhost:8080
APP_PORT=8080

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key-here
JWT_EXPIRES_IN=24h
JWT_REFRESH_EXPIRES_IN=168h

# Database Primary (MySQL)
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=golang_starter
DB_USERNAME=root
DB_PASSWORD=

# Database Secondary (PostgreSQL)
DB_PGSQL_HOST=127.0.0.1
DB_PGSQL_PORT=5432
DB_PGSQL_DATABASE=golang_starter_pg
DB_PGSQL_USERNAME=postgres
DB_PGSQL_PASSWORD=

# Database Third (MySQL Secondary)
DB_MYSQL_SECONDARY_HOST=127.0.0.1
DB_MYSQL_SECONDARY_PORT=3307
DB_MYSQL_SECONDARY_DATABASE=golang_starter_secondary
DB_MYSQL_SECONDARY_USERNAME=root
DB_MYSQL_SECONDARY_PASSWORD=

# Redis
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# File Storage
STORAGE_DISK=local
STORAGE_PATH=./storage
```

## Environment Variables

### Application Settings

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `APP_NAME` | string | "Golang Starter Kit" | Nama aplikasi |
| `APP_ENV` | string | local | Environment aplikasi (local, staging, production) |
| `APP_DEBUG` | bool | true | Debug mode |
| `APP_URL` | string | http://localhost:8080 | Base URL aplikasi |
| `APP_PORT` | int | 8080 | Port aplikasi |

### JWT Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `JWT_SECRET_KEY` | string | - | Secret key untuk JWT (wajib diisi) |
| `JWT_EXPIRES_IN` | duration | 24h | Durasi expire access token |
| `JWT_REFRESH_EXPIRES_IN` | duration | 168h | Durasi expire refresh token |

### Database Configuration

#### Primary Database (MySQL)
| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `DB_CONNECTION` | string | mysql | Driver database |
| `DB_HOST` | string | 127.0.0.1 | Host database |
| `DB_PORT` | int | 3306 | Port database |
| `DB_DATABASE` | string | - | Nama database |
| `DB_USERNAME` | string | - | Username database |
| `DB_PASSWORD` | string | - | Password database |

#### PostgreSQL Database
| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `DB_PGSQL_HOST` | string | 127.0.0.1 | Host PostgreSQL |
| `DB_PGSQL_PORT` | int | 5432 | Port PostgreSQL |
| `DB_PGSQL_DATABASE` | string | - | Nama database PostgreSQL |
| `DB_PGSQL_USERNAME` | string | - | Username PostgreSQL |
| `DB_PGSQL_PASSWORD` | string | - | Password PostgreSQL |

#### MySQL Secondary Database
| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `DB_MYSQL_SECONDARY_HOST` | string | 127.0.0.1 | Host MySQL secondary |
| `DB_MYSQL_SECONDARY_PORT` | int | 3307 | Port MySQL secondary |
| `DB_MYSQL_SECONDARY_DATABASE` | string | - | Nama database MySQL secondary |
| `DB_MYSQL_SECONDARY_USERNAME` | string | - | Username MySQL secondary |
| `DB_MYSQL_SECONDARY_PASSWORD` | string | - | Password MySQL secondary |

### Redis Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `REDIS_HOST` | string | 127.0.0.1 | Host Redis |
| `REDIS_PORT` | int | 6379 | Port Redis |
| `REDIS_PASSWORD` | string | - | Password Redis |
| `REDIS_DB` | int | 0 | Database Redis |

### Storage Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `STORAGE_DISK` | string | local | Driver storage (local, s3, etc) |
| `STORAGE_PATH` | string | ./storage | Path untuk local storage |

## Environment Specific Settings

### Development Environment

```bash
APP_ENV=local
APP_DEBUG=true
LOG_LEVEL=debug
```

### Staging Environment

```bash
APP_ENV=staging
APP_DEBUG=false
LOG_LEVEL=info
```

### Production Environment

```bash
APP_ENV=production
APP_DEBUG=false
LOG_LEVEL=error
# Pastikan semua secret keys aman
JWT_SECRET_KEY=production-secure-key
```

## Security Best Practices

### 1. Environment Variables Security

- **Jangan commit** file `.env` ke version control
- Gunakan `.env.example` sebagai template
- Gunakan secret management untuk production
- Rotasi secret keys secara berkala

### 2. Database Security

```bash
# Gunakan password yang kuat
DB_PASSWORD=complex-secure-password-123!

# Batasi akses database
DB_HOST=127.0.0.1  # Jangan gunakan 0.0.0.0 kecuali diperlukan
```

### 3. JWT Security

```bash
# Gunakan secret key yang panjang dan kompleks
JWT_SECRET_KEY=your-very-long-and-complex-secret-key-for-jwt-signing

# Sesuaikan expire time
JWT_EXPIRES_IN=1h     # Untuk production, gunakan waktu yang lebih pendek
JWT_REFRESH_EXPIRES_IN=24h
```

## Docker Environment

Untuk deployment dengan Docker, gunakan environment variables:

```yaml
# docker-compose.yml
environment:
  - APP_ENV=production
  - DB_HOST=mysql
  - DB_PASSWORD=${DB_PASSWORD}
  - JWT_SECRET_KEY=${JWT_SECRET_KEY}
```

## Helper Functions

Aplikasi menyediakan helper functions untuk mengakses environment variables:

```go
import "your-app/app/helpers"

// Get environment variable dengan default value
dbHost := helpers.GetEnv("DB_HOST", "127.0.0.1")

// Get environment variable sebagai integer
dbPort := helpers.GetEnvAsInt("DB_PORT", 3306)

// Get environment variable sebagai boolean
debug := helpers.GetEnvAsBool("APP_DEBUG", false)
```

## Troubleshooting

### Common Issues

1. **Environment file not found**
   ```bash
   # Pastikan file .env ada di root project
   ls -la .env
   ```

2. **Database connection failed**
   ```bash
   # Periksa konfigurasi database
   echo $DB_HOST $DB_PORT $DB_USERNAME
   ```

3. **JWT token invalid**
   ```bash
   # Pastikan JWT_SECRET_KEY sudah diset
   echo $JWT_SECRET_KEY
   ```

## Validation

Aplikasi akan memvalidasi environment variables yang wajib diisi saat startup. Jika ada yang missing, aplikasi akan menampilkan error dan tidak akan start.

### Required Variables

- `JWT_SECRET_KEY`
- `DB_DATABASE` (untuk primary database)
- `DB_USERNAME` (untuk primary database)

### Optional Variables

Semua variable lain memiliki default value yang reasonable untuk development.

---

Untuk informasi lebih lanjut, lihat:
- [Installation Guide](installation.md)
- [Database Configuration](../database/README.md)
- [Multi-Database Setup](../database/multi-database.md)
