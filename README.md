# Starter Kit Backend Golang 2025

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge)
![GORM](https://img.shields.io/badge/GORM-ORM-00ADD8?style=for-the-badge)
![MySQL](https://img.shields.io/badge/MySQL-Database-4479A1?style=for-the-badge&logo=mysql)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?style=for-the-badge&logo=postgresql)
![Swagger](https://img.shields.io/badge/Swagger-API%20Docs-85EA2D?style=for-the-badge&logo=swagger)

## 🚀 Deskripsi

Starter Kit Backend Golang adalah template lengkap untuk memulai pengembangan aplikasi backend menggunakan Go. Proyek ini menyediakan struktur modular yang siap pakai dengan fitur-fitur modern dan best practices.

## ✨ Fitur Utama

- 🔐 **Autentikasi JWT** - Sistem autentikasi yang aman
- 🗄️ **Multi-Database Support** - MySQL, PostgreSQL, dan lainnya
- 🔄 **Hot Reload** - Development yang efisien dengan Air
- 📚 **Auto-Generated API Docs** - Swagger/OpenAPI integration
- 🏗️ **Modular Architecture** - Clean code structure
- 📊 **Database Management** - Migration & seeding system
- 🛡️ **Middleware Support** - Auth, logging, CORS, dan lainnya
- 🧪 **Testing Ready** - Unit test structure

## 📋 Quick Start

```bash
# Clone repository
git clone https://github.com/RahmatRafiq/golang_starter_kit_2025.git
cd golang_starter_kit_2025

# Install dependencies
go mod tidy
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest

# Setup environment
cp .env.example .env

# Generate API documentation
swag init

# Run application
air
```

🌐 **Akses Aplikasi:**
- API Documentation: http://localhost:8080/swagger/index.html
- Health Check: http://localhost:8080/api/health

## 📖 Dokumentasi

| Topik | Link | Deskripsi |
|-------|------|-----------|
| 🚀 Getting Started | [documentation/GETTING_STARTED.md](documentation/GETTING_STARTED.md) | Panduan instalasi dan setup |
| 🗄️ Database Management | [documentation/DATABASE.md](documentation/DATABASE.md) | Migration, seeder, dan CLI commands |
| 🔗 Multi-Database | [documentation/MULTI_DATABASE.md](documentation/MULTI_DATABASE.md) | Konfigurasi multiple database connections |
| 📚 API Reference | [documentation/API_REFERENCE.md](documentation/API_REFERENCE.md) | Dokumentasi lengkap semua API endpoints |

## 🏗️ Arsitektur Project

```
golang_starter_kit_2025/
├── app/
│   ├── controllers/                  # 🎮 API Controllers
│   │   ├── auth_controllers.go       # Logika autentikasi (login, register)
│   │   ├── category_controller.go    # Manajemen kategori produk
│   │   ├── database_controller.go    # Management database connections
│   │   ├── file_controller.go        # Upload dan management file
│   │   ├── permission_controller.go  # Management permission sistem
│   │   ├── product_controller.go     # CRUD produk dan inventory
│   │   ├── role_controller.go        # Management role pengguna
│   │   └── user_controller.go        # Management pengguna
│   ├── casts/                        # 🔄 Data Transformation
│   │   ├── jwt_claims.go             # JWT claims structure
│   │   └── token.go                  # Token management
│   ├── helpers/                      # 🛠️ Helper Functions
│   │   ├── base64file_helper.go      # Base64 file operations
│   │   ├── env_helper.go             # Environment variable handling
│   │   ├── error_helper.go           # Error handling utilities
│   │   ├── file_helper.go            # File operations
│   │   ├── hash_helper.go            # Password hashing (bcrypt)
│   │   ├── path_helper.go            # Path utilities
│   │   ├── reference_helper.go       # Reference data helpers
│   │   ├── response_helper.go        # API response formatting
│   │   └── url_helper.go             # URL utilities
│   ├── middleware/                   # 🛡️ Middleware Components
│   │   ├── auth_middleware.go        # JWT authentication
│   │   └── logger_middleware.go      # Request/response logging
│   ├── models/                       # 📊 Database Models
│   │   ├── category.go               # Category model
│   │   ├── permission.go             # Permission model
│   │   ├── product.go                # Product model
│   │   ├── role.go                   # Role model
│   │   ├── role_has_permission.go    # Role-Permission pivot
│   │   ├── user.go                   # User model
│   │   ├── user_has_role.go          # User-Role pivot
│   │   └── scopes/                   # Query scopes
│   │       └── pagination.go         # Pagination scope
│   ├── repositories/                 # 🗄️ Data Access Layer
│   │   └── user_repository.go        # User data operations
│   ├── requests/                     # ✅ Request Validation
│   │   ├── category_request.go       # Category validation rules
│   │   ├── filter_request.go         # Filter/search validation
│   │   ├── login_request.go          # Login form validation
│   │   ├── permission_request.go     # Permission validation
│   │   ├── product_request.go        # Product validation rules
│   │   └── role_request.go           # Role validation rules
│   ├── responses/                    # 📤 Response Formatting
│   ├── services/                     # 💼 Business Logic
│   │   ├── auth_service.go           # Authentication business logic
│   │   ├── category_service.go       # Category business logic
│   │   ├── database_service.go       # Multi-database operations
│   │   ├── file_service.go           # File upload/management
│   │   ├── jwt_service.go            # JWT token operations
│   │   ├── permission_service.go     # Permission management
│   │   ├── product_service.go        # Product business logic
│   │   ├── role_service.go           # Role management
│   │   └── user_services.go          # User business logic
│   └── database/                     # 🔧 Database Management
│       ├── migration_manager.go      # Migration management system
│       ├── seeder_manager.go         # Seeder management system
│       ├── migrations/               # SQL Migration Files
│       │   ├── 20250426184415_create_roles_table.sql
│       │   ├── 20250426184424_create_permissions_table.sql
│       │   ├── 20250426184432_create_users_table.sql
│       │   └── ...
│       └── seeds/                    # Database Seeders
│           └── ...
├── bootstrap/                        # 🚀 Application Bootstrap
│   └── main.go                       # Application entry point
├── cmd/                              # 📝 CLI Commands
│   ├── migrate.go                    # Migration commands
│   └── seeder.go                     # Seeder commands
├── config/                           # ⚙️ Configuration
│   └── database.go                   # Multi-database configuration
├── database/                         # 🗃️ Database Core
│   └── manager.go                    # Database connection manager
├── docs/                             # 📋 Swagger Documentation
│   ├── docs.go                       # Generated swagger docs
│   ├── swagger.json                  # Swagger JSON spec
│   └── swagger.yaml                  # Swagger YAML spec
├── documentation/                    # 📚 Project Documentation
│   ├── README.md                     # Documentation index
│   ├── GETTING_STARTED.md            # Setup and installation guide
│   ├── DATABASE.md                   # Migration & seeder guide
│   ├── MULTI_DATABASE.md             # Multi-database configuration
│   └── API_REFERENCE.md              # Complete API documentation
├── examples/                         # � Usage Examples
│   └── multi_database_usage.go       # Multi-database examples
├── facades/                          # 🎭 Facade Pattern
│   ├── database.go                   # Primary database facade
│   └── database2.go                  # Secondary database facade
├── routes/                           # �🛣️ API Routes
│   └── web.go                        # Route definitions
├── storage/                          # 💾 File Storage
├── tmp/                              # 🗂️ Temporary Files
│   ├── build-errors.log              # Build error logs
│   └── main.exe                      # Compiled executable
├── .env                              # 🔐 Environment Configuration
├── .air.toml                         # 🔄 Air (Hot Reload) Configuration
├── docker-compose.yml                # 🐳 Docker Compose Configuration
├── Dockerfile                        # 🐳 Docker Container Configuration
├── go.mod                            # 📦 Go Module Definition
├── go.sum                            # 🔒 Go Module Checksums
├── main.go                           # 🎯 Application Main Entry Point
├── Makefile                          # 🔨 Build Automation
└── README.md                         # 📖 Project Documentation
```

## 🛠️ Tech Stack

- **Framework**: Gin (HTTP Web Framework)
- **ORM**: GORM (Object-Relational Mapping)
- **Database**: MySQL, PostgreSQL, SQLite
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: Swagger/OpenAPI
- **Development**: Air (Hot Reload)
- **Testing**: Go Testing Package

## 🔧 Development Tools

### Migration Commands
```bash
go run main.go make:migration create_users_table  # Buat migration baru
go run main.go migrate:all                        # Jalankan semua migration
go run main.go rollback:batch                     # Rollback batch terakhir
```

### Seeder Commands  
```bash
go run main.go make:seeder --name=users          # Buat seeder baru
go run main.go db:seed                           # Jalankan semua seeder
go run main.go rollback:seeder                   # Rollback seeder
```

## 🧪 Testing

```bash
go test ./...              # Run semua test
go test -cover ./...       # Test dengan coverage
go test ./app/controllers  # Test spesifik package
```

## 🤝 Kontribusi

Aplikasi ini dikembangkan oleh [Dzyfhuba](https://github.com/Dzyfhuba) dan [RahmatRafiq](https://github.com/RahmatRafiq). 

### Contributing Guidelines
1. Fork repository ini
2. Buat feature branch (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buka Pull Request

### Cara Berkontribusi
- 🐛 Laporkan bug melalui [Issues](https://github.com/RahmatRafiq/golang_starter_kit_2025/issues)
- 💡 Ajukan fitur baru via [Discussions](https://github.com/RahmatRafiq/golang_starter_kit_2025/discussions)
- 📖 Perbaiki dokumentasi
- 🧪 Tambahkan test coverage

## 💸 Dukung Proyek Ini

Jika proyek ini membantu Anda, consider untuk memberikan dukungan:

[![Saweria](https://img.shields.io/badge/Saweria-Donate-orange?style=for-the-badge)](https://saweria.co/RahmatRafiq)

## 📝 License

Proyek ini menggunakan MIT License. Lihat file [LICENSE](LICENSE) untuk detail lengkap.

---

### 🚀 Selamat coding dan semoga proyek ini membantu pengembangan aplikasi backend Anda! 

**Made with ❤️ by Indonesian Developers**
