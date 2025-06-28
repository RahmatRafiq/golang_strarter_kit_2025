# Frequently Asked Questions (FAQ)

Kumpulan pertanyaan yang sering diajukan tentang Golang Starter Kit.

## General Questions

### Q: Apa itu Golang Starter Kit?

**A:** Golang Starter Kit adalah boilerplate/template project untuk membangun REST API menggunakan Go dengan arsitektur yang mirip dengan Laravel. Kit ini menyediakan struktur project yang terorganisir, multi-database support, authentication, migration system, dan berbagai fitur lainnya yang diperlukan untuk membangun aplikasi web modern.

### Q: Mengapa memilih arsitektur mirip Laravel?

**A:** Laravel memiliki arsitektur yang proven dan developer-friendly. Dengan mengadopsi pola-pola dari Laravel seperti:
- MVC pattern
- Service layer
- Repository pattern
- Migration system
- Facade pattern

Developer yang familiar dengan Laravel dapat dengan mudah beradaptasi dengan struktur ini.

### Q: Apakah bisa digunakan untuk production?

**A:** Ya, Golang Starter Kit dirancang untuk production-ready. Namun, pastikan untuk:
- Menggunakan environment variables yang aman
- Setup monitoring dan logging
- Konfigurasi database yang optimal
- Implementasi security best practices

## Installation & Setup

### Q: Apa saja requirements untuk menjalankan project ini?

**A:** Requirements minimal:
- Go 1.21 atau lebih baru
- MySQL 8.0+ atau MariaDB 10.5+
- Git

Requirements opsional:
- PostgreSQL 13+ (untuk multi-database)
- Redis 6.0+ (untuk caching)
- Docker & Docker Compose (untuk containerization)

### Q: Bagaimana cara setup multi-database?

**A:** 
1. Configure environment variables di file `.env`:
```bash
# Primary MySQL
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_DATABASE=primary_db

# PostgreSQL
DB_PGSQL_HOST=127.0.0.1
DB_PGSQL_DATABASE=postgres_db

# Secondary MySQL
DB_MYSQL_SECONDARY_HOST=127.0.0.1
DB_MYSQL_SECONDARY_DATABASE=secondary_db
```

2. Jalankan migration untuk setiap database:
```bash
go run cmd/migrate.go up --connection=mysql
go run cmd/migrate.go up --connection=postgresql
go run cmd/migrate.go up --connection=mysql_secondary
```

Lihat: [Multi-Database Guide](../database/multi-database.md)

### Q: Error "table doesn't exist" saat menjalankan aplikasi?

**A:** Pastikan migration sudah dijalankan:
```bash
# Check migration status
go run cmd/migrate.go status

# Run pending migrations
go run cmd/migrate.go up
```

### Q: Bagaimana cara menggunakan Docker?

**A:** 
1. Build dan jalankan dengan docker-compose:
```bash
docker-compose up -d
```

2. Atau build manual:
```bash
docker build -t golang-starter-kit .
docker run -p 8080:8080 golang-starter-kit
```

## Database

### Q: Bagaimana cara membuat migration baru?

**A:** 
```bash
# Create migration
go run cmd/migrate.go create create_products_table

# Edit file di app/database/migrations/
# Kemudian jalankan:
go run cmd/migrate.go up
```

### Q: Bagaimana cara rollback migration?

**A:** 
```bash
# Rollback 1 step
go run cmd/migrate.go down

# Rollback to specific version
go run cmd/migrate.go down --to=20250101000000
```

### Q: Bagaimana cara menggunakan database yang berbeda untuk read/write?

**A:** 
```go
// Write to primary
err := facades.DB("mysql").Create(&user).Error

// Read from replica
var users []models.User
err := facades.DB("mysql_secondary").Find(&users).Error
```

### Q: Error koneksi database?

**A:** Periksa:
1. Database server berjalan
2. Credentials di `.env` benar
3. Database sudah dibuat
4. Firewall tidak memblokir koneksi

```bash
# Test koneksi MySQL
mysql -h 127.0.0.1 -u root -p database_name

# Test koneksi PostgreSQL
psql -h 127.0.0.1 -U postgres -d database_name
```

## Development

### Q: Bagaimana cara menambahkan endpoint API baru?

**A:** 
1. Buat model (jika diperlukan)
2. Buat request validation
3. Buat service untuk business logic
4. Buat controller
5. Register route

Contoh lengkap ada di [Development Guide](development.md)

### Q: Bagaimana cara implementasi authentication custom?

**A:** 
1. Modify `app/services/auth_service.go`
2. Update `app/middleware/auth_middleware.go`
3. Sesuaikan `app/casts/jwt_claims.go`

### Q: Bagaimana cara menambahkan validation rules custom?

**A:** 
```go
// Di request struct
type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email,custom_email_validation"`
}

// Implementasi custom validator
func customEmailValidation(fl validator.FieldLevel) bool {
    email := fl.Field().String()
    // Custom validation logic
    return !strings.Contains(email, "example.com")
}
```

### Q: Bagaimana cara setup hot reload untuk development?

**A:** 
```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run dengan auto-reload
air

# Atau menggunakan Makefile
make dev
```

## Testing

### Q: Bagaimana cara menjalankan tests?

**A:** 
```bash
# Run all tests
go test ./...

# Run tests dengan coverage
go test -cover ./...

# Run specific package tests
go test ./app/services -v
```

### Q: Bagaimana cara setup test database?

**A:** 
1. Buat database terpisah untuk testing
2. Set environment variable:
```bash
export APP_ENV=testing
export DB_DATABASE=test_database
```

3. Run migration untuk test database:
```bash
go run cmd/migrate.go up
```

### Q: Bagaimana cara membuat mock untuk testing?

**A:** Gunakan testify/mock atau mockgen:
```go
// Manual mock
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}
```

## Deployment

### Q: Bagaimana cara build untuk production?

**A:** 
```bash
# Build binary
go build -ldflags="-w -s" -o bin/app main.go

# Set production environment
export APP_ENV=production
export APP_DEBUG=false

# Run
./bin/app
```

### Q: Bagaimana cara deploy ke cloud provider?

**A:** 
1. **Docker deployment:**
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

2. **Environment variables:** Set melalui cloud provider interface

3. **Database:** Gunakan managed database service

### Q: Rekomendasi untuk production setup?

**A:** 
- **Reverse Proxy:** Nginx atau Traefik
- **Database:** Managed service (AWS RDS, GCP Cloud SQL)
- **Caching:** Redis cluster
- **Monitoring:** Prometheus + Grafana
- **Logging:** ELK stack atau cloud logging
- **Security:** SSL/TLS, rate limiting, CORS

## Performance

### Q: Aplikasi terasa lambat, apa yang harus dilakukan?

**A:** 
1. **Database optimization:**
   - Add indexes
   - Optimize queries
   - Use connection pooling
   - Consider read replicas

2. **Application optimization:**
   - Profile dengan `go tool pprof`
   - Implement caching
   - Optimize goroutine usage

3. **Infrastructure:**
   - Scale horizontally
   - Use CDN
   - Optimize network

### Q: Bagaimana cara implementasi caching?

**A:** 
```go
// Redis caching
func (s *UserService) GetUser(id uint) (*models.User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    
    // Try cache first
    if cached := cache.Get(cacheKey); cached != nil {
        return cached.(*models.User), nil
    }
    
    // Get from database
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    cache.Set(cacheKey, user, 5*time.Minute)
    
    return user, nil
}
```

### Q: Memory usage tinggi, bagaimana cara debug?

**A:** 
```bash
# Memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap

# Check for memory leaks
go tool pprof -alloc_space http://localhost:8080/debug/pprof/allocs
```

## Security

### Q: Bagaimana cara mengamankan API endpoints?

**A:** 
1. **Authentication:** JWT tokens
2. **Authorization:** Role-based access control
3. **Rate limiting:** Implement middleware
4. **Input validation:** Validate semua input
5. **CORS:** Configure proper CORS headers

### Q: Bagaimana cara mengubah JWT secret secara aman?

**A:** 
1. Generate secret baru
2. Deploy dengan secret baru
3. Gradual rotation untuk backward compatibility
4. Invalidate old tokens setelah grace period

### Q: Bagaimana cara handle SQL injection?

**A:** GORM secara default melindungi dari SQL injection jika menggunakan:
```go
// Safe
db.Where("name = ?", userInput).Find(&users)

// Unsafe - hindari
db.Where(fmt.Sprintf("name = '%s'", userInput)).Find(&users)
```

## Troubleshooting

### Q: Error "bind: address already in use"?

**A:** 
```bash
# Check port usage
lsof -i :8080

# Kill process
kill -9 <PID>

# Atau ubah port di .env
APP_PORT=8081
```

### Q: Error "panic: runtime error: invalid memory address"?

**A:** Biasanya karena nil pointer. Debug dengan:
1. Check logs untuk stack trace
2. Validate semua pointer sebelum digunakan
3. Gunakan proper error handling

### Q: Application tidak start?

**A:** Check:
1. Environment variables loaded
2. Database connection
3. File permissions
4. Port availability
5. Dependencies installed

### Q: Migration failed?

**A:** 
1. Check database connection
2. Verify migration syntax
3. Check database permissions
4. Look at migration logs

```bash
# Check migration status
go run cmd/migrate.go status

# Rollback and retry
go run cmd/migrate.go down
go run cmd/migrate.go up
```

## Contributing

### Q: Bagaimana cara contribute to this project?

**A:** 
1. Fork repository
2. Create feature branch
3. Make changes dengan tests
4. Submit pull request
5. Follow coding standards

### Q: Coding standards yang digunakan?

**A:** 
- Gunakan `gofmt` untuk formatting
- Follow Go naming conventions
- Write tests untuk new features
- Document public functions
- Use meaningful commit messages

### Q: Bagaimana cara report bugs?

**A:** 
1. Check existing issues
2. Create detailed bug report
3. Include steps to reproduce
4. Provide environment information
5. Add relevant logs/screenshots

---

## Need More Help?

Jika pertanyaan Anda tidak terjawab di sini, silakan:

1. **Check Documentation:**
   - [Installation Guide](installation.md)
   - [Development Guide](development.md)
   - [Database Guide](../database/README.md)

2. **Search Issues:** Cek GitHub issues yang sudah ada

3. **Create Issue:** Buat issue baru dengan detail yang lengkap

4. **Community:** Join discussion forums atau chat groups

---

**Last Updated:** January 2025
