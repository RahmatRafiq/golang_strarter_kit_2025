# Golang Starter Kit 2025

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)
![Gin](https://img.shields.io/badge/Gin-v1.10-00ADD8?style=for-the-badge)
![GORM](https://img.shields.io/badge/GORM-v1.26-00ADD8?style=for-the-badge)
![MySQL](https://img.shields.io/badge/MySQL-Database-4479A1?style=for-the-badge&logo=mysql)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?style=for-the-badge&logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

A production-ready Go backend starter template with clean architecture, JWT authentication, multi-database support, and comprehensive API documentation. Built with modern best practices and designed for scalability.

## ✨ Key Features

- 🔐 **JWT Authentication** - Secure authentication with Argon2id password hashing
- 🗄️ **Multi-Database Support** - MySQL, PostgreSQL, SQLite with connection pooling
- 🏗️ **Clean Architecture** - Layered pattern (Controller → Service → Repository → Database)
- 🔄 **Hot Reload** - Fast development with Air
- 📚 **Auto-Generated API Docs** - Swagger/OpenAPI integration
- 🛡️ **RBAC System** - Role-Based Access Control with Users, Roles, Permissions
- 📊 **Database Management** - Powerful migration & seeding system with batch tracking
- 🧪 **Testing Ready** - BDD testing with Ginkgo & Gomega
- 🚀 **Production Ready** - Docker support, health checks, monitoring

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or higher
- MySQL 8.0+ or PostgreSQL 13+
- Git

### Installation

```bash
# Clone repository
git clone https://github.com/RahmatRafiq/golang_starter_kit_2025.git
cd golang_starter_kit_2025

# Install dependencies
go mod tidy

# Install development tools
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest

# Setup environment
cp .env.example .env
# Edit .env and configure your database & JWT secret

# Generate JWT secret (recommended)
openssl rand -base64 48

# Run database migrations
go run main.go migrate:all

# Generate API documentation
swag init

# Run with hot reload (development)
air

# OR run directly
go run main.go
```

### Access the Application

- **API Server**: http://localhost:9999
- **Swagger UI**: http://localhost:9999/swagger/index.html
- **Health Check**: http://localhost:9999/health
- **Multi-DB Health**: http://localhost:9999/health/databases

### First API Request

```bash
# Login (default credentials from seeder)
curl -X PUT http://localhost:9999/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "admin123"}'

# Use the returned token for authenticated requests
curl -X GET http://localhost:9999/users \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 🏗️ Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│         Middleware Layer            │
│  (CORS, Auth, Logger)               │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│      Controller Layer               │
│  (HTTP Request/Response Handling)   │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│       Service Layer                 │
│  (Business Logic)                   │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│      Repository Layer               │
│  (Database Operations)              │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Database                    │
│  (MySQL / PostgreSQL / SQLite)      │
└─────────────────────────────────────┘
```

**Key Architectural Patterns:**
- **Dependency Injection** - Constructor-based injection for testability
- **Repository Pattern** - Database abstraction with interfaces
- **Facade Pattern** - Simplified database access
- **Singleton Pattern** - Database connection manager

## 📁 Project Structure

```
golang_starter_kit_2025/
├── app/
│   ├── controllers/         # HTTP request handlers
│   ├── services/            # Business logic layer
│   ├── repositories/        # Data access layer
│   │   └── interfaces/      # Repository contracts
│   ├── models/              # GORM database models
│   ├── middleware/          # HTTP middleware (auth, logging)
│   ├── requests/            # Request validation structs
│   ├── responses/           # Response formatting
│   ├── helpers/             # Utility functions
│   ├── casts/               # Data transformation objects
│   └── database/            # Migration & seeder management
│       ├── migrations/      # SQL migration files
│       └── seeds/           # Go seeder files
├── bootstrap/               # Application initialization
├── cmd/                     # CLI command definitions
├── config/                  # Configuration loaders
├── database/                # Database connection manager
├── facades/                 # Database access facade
├── routes/                  # API route definitions
├── storage/                 # File upload storage
└── documentation/           # Detailed documentation
```

See [CLAUDE.md](CLAUDE.md) for comprehensive architecture documentation.

## 🛠️ Tech Stack

| Category | Technology | Version |
|----------|------------|---------|
| **Language** | Go | 1.24.3 |
| **Framework** | Gin | 1.10.0 |
| **ORM** | GORM | 1.26.1 |
| **Databases** | MySQL, PostgreSQL, SQLite | - |
| **Auth** | JWT (golang-jwt/jwt) | 5.2.2 |
| **Password Hashing** | Argon2id | - |
| **Validation** | go-playground/validator | 10.26.0 |
| **API Docs** | Swaggo (Swagger) | - |
| **Testing** | Ginkgo & Gomega | 2.20.2 / 1.34.2 |
| **CLI** | urfave/cli | 2.27.6 |
| **Hot Reload** | Air | latest |

## 📊 Database Management

### Migration Commands

```bash
# Create new migration
go run main.go make:migration create_products_table

# Run all migrations
go run main.go migrate:all

# Run on specific database connection
go run main.go migrate:all --connection=postgres

# View migration status
go run main.go migrate:status

# Rollback last batch
go run main.go rollback:batch

# Rollback all migrations
go run main.go rollback:all

# Fresh migration (drop all + migrate)
go run main.go migrate:fresh
```

### Seeder Commands

```bash
# Create new seeder
go run main.go make:seeder --name=ProductsSeeder

# Run all seeders
go run main.go db:seed

# Run specific seeder
go run main.go db:seed --class=UserSeeder

# Rollback seeders
go run main.go rollback:seeder
```

**Multi-Database Support**: Add `--connection=postgres` or `--connection=mysql_secondary` to any command.

## 🔐 Authentication

The project uses JWT-based authentication with Argon2id password hashing:

```go
// Login
PUT /auth/login
Body: {"email": "user@example.com", "password": "password123"}

// Logout (requires authentication)
GET /auth/logout
Header: Authorization: Bearer {token}

// Refresh token
GET /auth/refresh
Header: Authorization: Bearer {token}
```

**Security Features:**
- Argon2id password hashing (winner of Password Hashing Competition 2015)
- JWT tokens with configurable expiry
- Token stored in database (enables logout)
- Development mode: Set `SKIP_AUTH=true` to bypass authentication (⚠️ never in production!)

## 🎯 API Endpoints

### Authentication
- `PUT /auth/login` - User login
- `GET /auth/logout` - User logout (protected)
- `GET /auth/refresh` - Refresh JWT token (protected)

### Users
- `GET /users` - List all users (paginated)
- `GET /users/:id` - Get user by ID
- `PUT /users` - Create/Update user
- `DELETE /users/:id` - Delete user
- `POST /users/:id/roles` - Assign roles to user
- `GET /users/:id/roles` - Get user's roles

### Roles & Permissions
- `GET /roles` - List all roles
- `PUT /roles` - Create/Update role
- `POST /roles/:id/permissions` - Assign permissions to role
- `GET /permissions` - List all permissions

### Products & Categories
- `GET /products` - List products (paginated)
- `GET /categories` - List categories
- Full CRUD operations available

### Health Checks
- `GET /health` - Basic health check
- `GET /health/databases` - Multi-database health with statistics

See full API documentation at `/swagger/index.html` when running the application.

## 🧪 Testing

```bash
# Run all tests with race detector
ginkgo -r --race

# Run tests with coverage
go test -cover ./...

# Test specific package
go test ./app/controllers

# Watch mode (re-run on changes)
ginkgo watch -r
```

## 🐳 Docker Support

```bash
# Build and run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop containers
docker-compose down
```

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [CLAUDE.md](CLAUDE.md) | Comprehensive technical documentation (2476 lines) |
| [documentation/GETTING_STARTED.md](documentation/GETTING_STARTED.md) | Detailed setup guide |
| [documentation/DATABASE.md](documentation/DATABASE.md) | Migration & seeder documentation |
| [documentation/MULTI_DATABASE.md](documentation/MULTI_DATABASE.md) | Multi-database configuration |
| [documentation/API_REFERENCE.md](documentation/API_REFERENCE.md) | Complete API endpoint reference |

## ⚙️ Configuration

Key environment variables in `.env`:

```bash
# Application
APP_PORT=9999
APP_ENV=development                    # development | production
JWT_SECRET_KEY=<generate-strong-key>   # openssl rand -base64 48
JWT_EXPIRE_MINUTES=60

# Default Database Connection
DB_CONNECTION=mysql                    # mysql | postgres | sqlite

# MySQL Configuration
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025
MYSQL_USER=root
MYSQL_PASSWORD=your_password

# PostgreSQL Configuration (optional)
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=golang_starter_kit_2025_pg
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password

# Development (⚠️ Never in production!)
SKIP_AUTH=false                        # Set true to bypass authentication
```

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Development Guidelines:**
- Follow clean architecture principles
- Write tests for new features
- Update documentation
- Run `swag init` after API changes
- Ensure all tests pass before submitting PR

## 🐛 Issues & Support

- Report bugs via [GitHub Issues](https://github.com/RahmatRafiq/golang_starter_kit_2025/issues)
- Feature requests via [GitHub Discussions](https://github.com/RahmatRafiq/golang_starter_kit_2025/discussions)
- Read [CLAUDE.md](CLAUDE.md) for detailed technical guidance

## 🌟 Features Roadmap

- [x] JWT Authentication
- [x] Multi-Database Support
- [x] RBAC System
- [x] Migration & Seeder System
- [x] Swagger Documentation
- [x] Docker Support
- [ ] GraphQL Support
- [ ] WebSocket Support
- [ ] Rate Limiting
- [ ] Caching Layer (Redis)
- [ ] Message Queue Integration
- [ ] Microservices Architecture Support

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

Developed with ❤️ by:
- [Dzyfhuba](https://github.com/Dzyfhuba)
- [RahmatRafiq](https://github.com/RahmatRafiq)

## 💝 Support

If this project helps you, consider supporting us:

[![Saweria](https://img.shields.io/badge/Saweria-Support%20Us-orange?style=for-the-badge&logo=ko-fi)](https://saweria.co/RahmatRafiq)

---

**Made with ❤️ by Indonesian Developers**

⭐ Star this repository if you find it helpful!
