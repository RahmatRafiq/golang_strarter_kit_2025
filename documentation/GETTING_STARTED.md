# Quick Start Guide

## Instalasi dan Setup

### 1. Clone Repository
```bash
git clone https://github.com/RahmatRafiq/golang_starter_kit_2025.git
```

### 2. Masuk ke Direktori Project
```bash
cd golang_starter_kit_2025
```

### 3. Install Dependencies
```bash
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

### 4. Setup Environment
Salin file `.env.example` menjadi `.env` dan sesuaikan konfigurasi:
```bash
cp .env.example .env
```

### 5. Generate Dokumentasi Swagger
```bash
swag init
```

### 6. Jalankan Aplikasi
```bash
air
```

### 7. Akses Aplikasi
- **API Documentation**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/api/health
- **Database Status**: http://localhost:8080/api/database/status

## Konfigurasi Environment

Proyek ini menggunakan file `.env` untuk mengatur konfigurasi. Berikut adalah contoh konfigurasi:

```env
# Application Configuration
APP_NAME="Golang Starter Kit 2025"
APP_ENV=development
APP_SCHEME=http
APP_HOST=localhost
APP_PORT=8080
APP_VERSION=1.0.0

# Database Configuration
DB_CONNECTION=mysql
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=root
MYSQL_PASSWORD=root

# JWT Configuration
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE_MINUTES=60

# File Upload Configuration
IMAGE_EXPIRE_MINUTES=2
```

## Struktur Project

```
golang_starter_kit_2025/
├── app/
│   ├── controllers/     # Controllers untuk handle request
│   ├── middleware/      # Middleware aplikasi
│   ├── models/         # Model database
│   ├── services/       # Business logic
│   ├── repositories/   # Data access layer
│   ├── requests/       # Request validation
│   ├── responses/      # Response format
│   └── database/       # Migration & seeder
├── config/             # Konfigurasi aplikasi
├── docs/               # Swagger documentation
├── documentation/      # Project documentation
├── routes/             # Route definitions
└── storage/            # File storage
```

## Testing

### Run Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Run Specific Test
```bash
go test ./app/controllers -v
```

## Development

### Hot Reload
Aplikasi menggunakan Air untuk hot reload. Setiap perubahan code akan otomatis merestart server.

### API Documentation
Swagger documentation akan otomatis ter-generate dari comment di controller. Format comment:
```go
// @Summary User login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body requests.LoginRequest true "Login request"
// @Success 200 {object} responses.LoginResponse
// @Router /api/auth/login [post]
```

### Database Migration
Lihat dokumentasi lengkap di [DATABASE.md](./DATABASE.md)

### Multi-Database
Lihat dokumentasi lengkap di [MULTI_DATABASE.md](./MULTI_DATABASE.md)

## Production Deployment

### 1. Build Application
```bash
go build -o main .
```

### 2. Run with Environment
```bash
APP_ENV=production ./main
```

### 3. Using Docker
```bash
docker build -t golang-starter-kit .
docker run -p 8080:8080 golang-starter-kit
```

### 4. Using Docker Compose
```bash
docker-compose up -d
```

## Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Check what's using port 8080
netstat -tulpn | grep 8080

# Kill process
kill -9 <PID>
```

#### 2. Database Connection Error
- Check database server is running
- Verify credentials in `.env`
- Test connection manually

#### 3. Swagger Not Loading
```bash
# Regenerate swagger docs
swag init

# Check if swagger docs folder exists
ls docs/
```

## Getting Help

- 📖 [Database Documentation](./DATABASE.md)
- 🔗 [Multi-Database Guide](./MULTI_DATABASE.md)
- 🐛 [Issue Tracker](https://github.com/RahmatRafiq/golang_starter_kit_2025/issues)
- 💬 [Discussions](https://github.com/RahmatRafiq/golang_starter_kit_2025/discussions)
