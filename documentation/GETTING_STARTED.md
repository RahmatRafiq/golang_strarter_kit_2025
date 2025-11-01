# Getting Started Guide

Complete guide to set up and run the Golang Starter Kit 2025 on your local machine.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24+** - [Download Go](https://golang.org/dl/)
- **MySQL 8.0+** or **PostgreSQL 13+** - Database server
- **Git** - Version control
- **Make** (optional) - Build automation

**Verify installations:**
```bash
go version        # Should show go1.24 or higher
mysql --version   # or: psql --version
git --version
```

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/RahmatRafiq/golang_starter_kit_2025.git
cd golang_starter_kit_2025
```

### 2. Install Dependencies

```bash
# Install Go module dependencies
go mod tidy

# Install development tools
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest

# For testing (optional)
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

### 3. Environment Configuration

```bash
# Copy example environment file
cp .env.example .env

# Edit .env file with your configuration
nano .env  # or use your preferred editor
```

**Required environment variables:**

```bash
# Application
APP_PORT=9999
APP_ENV=development
JWT_SECRET_KEY=<generate-strong-key>  # See below

# Database
DB_CONNECTION=mysql
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=root
MYSQL_PASSWORD=your_password
```

**Generate a strong JWT secret key:**
```bash
openssl rand -base64 48
```
Copy the output and paste it to `JWT_SECRET_KEY` in `.env`.

### 4. Database Setup

**Create database:**

```bash
# MySQL
mysql -u root -p
CREATE DATABASE golang_starter_kit_2025;
EXIT;

# PostgreSQL
psql -U postgres
CREATE DATABASE golang_starter_kit_2025;
\q
```

**Run migrations:**

```bash
# Run all migrations
go run main.go migrate:all

# Check migration status
go run main.go migrate:status
```

**Seed database (optional):**

```bash
# Run all seeders
go run main.go db:seed
```

### 5. Generate API Documentation

```bash
swag init
```

This generates Swagger documentation in the `docs/` folder.

### 6. Run the Application

**Development mode (with hot reload):**
```bash
air
```

**Standard run:**
```bash
go run main.go
```

**Production build:**
```bash
# Build executable
go build -o main .

# Run executable
./main
```

## Access the Application

Once running, you can access:

- **API Server**: http://localhost:9999
- **Swagger UI**: http://localhost:9999/swagger/index.html
- **Health Check**: http://localhost:9999/health
- **Multi-DB Health**: http://localhost:9999/health/databases

## First Steps

### 1. Test the API

**Health check:**
```bash
curl http://localhost:9999/health
```

**Login (using seeded admin account):**
```bash
curl -X PUT http://localhost:9999/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

**Use the token for authenticated requests:**
```bash
# Save token from login response
TOKEN="your_token_here"

# Get users list
curl http://localhost:9999/users \
  -H "Authorization: Bearer $TOKEN"
```

### 2. Explore Swagger Documentation

Open http://localhost:9999/swagger/index.html in your browser to:
- View all available endpoints
- Try out API requests interactively
- See request/response schemas

### 3. Check Database Health

```bash
curl http://localhost:9999/health/databases
```

This shows connection status for all configured databases.

## Project Structure

```
golang_starter_kit_2025/
├── app/
│   ├── controllers/         # HTTP request handlers
│   ├── services/            # Business logic layer
│   ├── repositories/        # Data access layer
│   │   └── interfaces/      # Repository contracts
│   ├── models/              # GORM database models
│   ├── middleware/          # HTTP middleware
│   ├── requests/            # Request validation
│   ├── responses/           # Response formatting
│   ├── helpers/             # Utility functions
│   ├── casts/               # Data transformations
│   └── database/
│       ├── migrations/      # SQL migration files
│       └── seeds/           # Database seeders
├── bootstrap/               # Application initialization
├── cmd/                     # CLI commands
├── config/                  # Configuration loaders
├── database/                # Database manager
├── facades/                 # Database facades
├── routes/                  # Route definitions
├── docs/                    # Swagger documentation
├── storage/                 # File uploads
├── .env                     # Environment variables
├── main.go                  # Application entry point
└── go.mod                   # Go module definition
```

## Development Workflow

### Hot Reload with Air

Air automatically restarts the server when you make code changes:

```bash
# Start with hot reload
air

# Air configuration is in .air.toml
# Customize as needed
```

### Running Tests

```bash
# Run all tests
ginkgo -r --race

# Run with coverage
go test -cover ./...

# Run specific package
go test ./app/controllers

# Watch mode
ginkgo watch -r
```

### Database Management

**Migrations:**
```bash
# Create new migration
go run main.go make:migration create_products_table

# Run migrations
go run main.go migrate:all

# Check status
go run main.go migrate:status

# Rollback last batch
go run main.go rollback:batch

# Fresh migration (reset)
go run main.go migrate:fresh
```

**Seeders:**
```bash
# Create new seeder
go run main.go make:seeder --name=ProductsSeeder

# Run seeders
go run main.go db:seed

# Rollback seeders
go run main.go rollback:seeder
```

See [DATABASE.md](./DATABASE.md) for complete documentation.

### API Documentation

After modifying API endpoints:

```bash
# Update Swagger annotations in controllers
# Then regenerate docs
swag init

# Swagger will be available at /swagger/index.html
```

**Swagger annotation example:**
```go
// @Summary Get user by ID
// @Description Get detailed user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} responses.UserResponse
// @Router /users/{id} [get]
// @Security Bearer
```

## Docker Deployment

### Using Docker

```bash
# Build image
docker build -t golang-starter-kit .

# Run container
docker run -p 9999:9999 \
  -e DB_CONNECTION=mysql \
  -e MYSQL_HOST=host.docker.internal \
  golang-starter-kit
```

### Using Docker Compose

```bash
# Start all services (app + database)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild and start
docker-compose up --build
```

## Environment Configuration

### Development (.env)

```bash
APP_ENV=development
APP_PORT=9999
JWT_SECRET_KEY=<your-secret-key>
JWT_EXPIRE_MINUTES=60

# Development helper (⚠️ NEVER in production!)
SKIP_AUTH=false  # Set true to bypass authentication

DB_CONNECTION=mysql
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=root
MYSQL_PASSWORD=root

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### Production (.env)

```bash
APP_ENV=production
APP_PORT=9999
JWT_SECRET_KEY=<strong-random-key>
JWT_EXPIRE_MINUTES=30

SKIP_AUTH=false  # MUST be false

DB_CONNECTION=mysql
MYSQL_HOST=production-db-host
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=production_user
MYSQL_PASSWORD=<secure-password>

ALLOWED_ORIGINS=https://yourdomain.com
```

## Common Development Tasks

### 1. Adding a New Feature

```bash
# 1. Create migration
go run main.go make:migration create_orders_table

# 2. Edit migration file
nano app/database/migrations/YYYYMMDDHHMMSS_create_orders_table.sql

# 3. Run migration
go run main.go migrate:all

# 4. Create model (app/models/order.go)
# 5. Create repository interface (app/repositories/interfaces/)
# 6. Implement repository (app/repositories/)
# 7. Create service (app/services/)
# 8. Create controller (app/controllers/)
# 9. Register routes (routes/web.go)
# 10. Update Swagger
swag init
```

See [CLAUDE.md](../CLAUDE.md) for detailed guide on adding features.

### 2. Multi-Database Setup

```bash
# Configure PostgreSQL in .env
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=golang_starter_kit_2025_pg
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres

# Run migrations on PostgreSQL
go run main.go migrate:all --connection=postgres
```

See [MULTI_DATABASE.md](./MULTI_DATABASE.md) for complete guide.

### 3. Testing Workflow

```bash
# Run tests before committing
ginkgo -r --race

# Check coverage
go test -cover ./...

# Run specific test
go test ./app/services -v -run TestUserService
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using port 9999
lsof -i :9999  # macOS/Linux
netstat -ano | findstr :9999  # Windows

# Kill the process
kill -9 <PID>  # macOS/Linux
taskkill /PID <PID> /F  # Windows
```

### Database Connection Error

**Check database is running:**
```bash
# MySQL
systemctl status mysql  # Linux
brew services list      # macOS

# PostgreSQL
systemctl status postgresql
```

**Test connection manually:**
```bash
# MySQL
mysql -h localhost -P 3306 -u root -p

# PostgreSQL
psql -h localhost -p 5432 -U postgres
```

**Verify credentials in `.env`**

### Swagger Not Loading

```bash
# Regenerate Swagger docs
swag init

# Check docs folder exists
ls docs/

# Clear browser cache and reload
```

### Migration Errors

```bash
# Check migration status
go run main.go migrate:status

# Rollback failed migration
go run main.go rollback:batch

# Fix migration file
# Run again
go run main.go migrate:all
```

### Hot Reload Not Working

```bash
# Check Air is installed
air -v

# Reinstall if needed
go install github.com/air-verse/air@latest

# Check .air.toml configuration
# Restart Air
air
```

## Security Considerations

### Development

- ✅ Use strong JWT secret key
- ✅ Set `SKIP_AUTH=false` (unless debugging)
- ✅ Use `.env` for credentials (never commit)
- ✅ Enable CORS only for trusted origins

### Production

- ✅ Use environment-specific configuration
- ✅ Enable HTTPS/TLS
- ✅ Use strong database passwords
- ✅ Set `APP_ENV=production`
- ✅ **NEVER set `SKIP_AUTH=true`**
- ✅ Monitor application logs
- ✅ Regular security updates
- ✅ Database backups

## Next Steps

1. **Read Documentation:**
   - [DATABASE.md](./DATABASE.md) - Migration & seeder guide
   - [MULTI_DATABASE.md](./MULTI_DATABASE.md) - Multi-database setup
   - [API_REFERENCE.md](./API_REFERENCE.md) - API endpoints
   - [CLAUDE.md](../CLAUDE.md) - Complete technical documentation

2. **Explore Examples:**
   - Check `examples/` folder for usage examples
   - Try multi-database operations

3. **Customize:**
   - Modify models for your use case
   - Add custom middleware
   - Implement your business logic

4. **Deploy:**
   - Set up CI/CD pipeline
   - Configure production environment
   - Monitor and maintain

## Getting Help

- 📖 **Documentation**: Check [CLAUDE.md](../CLAUDE.md) for detailed technical docs
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/RahmatRafiq/golang_starter_kit_2025/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/RahmatRafiq/golang_starter_kit_2025/discussions)
- 📧 **Email**: [Contact maintainers](mailto:rahmatrafiq@example.com)

## Quick Reference

| Command | Description |
|---------|-------------|
| `air` | Start with hot reload |
| `go run main.go` | Start server |
| `go run main.go migrate:all` | Run migrations |
| `go run main.go db:seed` | Run seeders |
| `swag init` | Generate API docs |
| `ginkgo -r --race` | Run tests |
| `go build -o main .` | Build for production |

---

**Happy coding!** 🚀
