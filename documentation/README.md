# Multi-Database Golang Starter Kit Documentation

Selamat datang di dokumentasi lengkap untuk Multi-Database Golang Starter Kit 2025. Dokumentasi ini dirancang untuk membantu Anda memahami dan menggunakan semua fitur yang tersedia.

## 📚 Struktur Dokumentasi

```
documentation/
├── guides/           # Panduan penggunaan
├── database/         # Dokumentasi database
├── api/             # Dokumentasi API
└── examples/        # Contoh implementasi
```

## 🚀 Quick Start

1. [Installation Guide](guides/installation.md) - Setup dan instalasi
2. [Architecture Guide](guides/architecture.md) - Memahami arsitektur
3. [Database Configuration](database/README.md) - Setup multiple database
4. [CLI Commands](guides/cli.md) - Command line tools

## 📖 Panduan Utama

### 🛠️ Setup & Installation
- [Installation Guide](guides/installation.md) - Panduan instalasi lengkap
- [Environment Configuration](guides/environment.md) - Konfigurasi environment
- [Development Setup](guides/development.md) - Setup untuk development

### 🏗️ Architecture & Patterns
- [Architecture Overview](guides/architecture.md) - Arsitektur aplikasi
- [Repository Pattern](guides/repository-pattern.md) - Pattern repository
- [Service Layer](guides/service-layer.md) - Business logic layer
- [Dependency Injection](guides/dependency-injection.md) - DI pattern

### 🗄️ Database Management
- [Multi-Database Setup](database/README.md) - Konfigurasi multiple database
- [Migration System](database/migrations.md) - Sistem migrasi
- [Seeder System](database/seeders.md) - Data seeding
- [Query Optimization](database/optimization.md) - Optimasi query

### 🔌 API Development
- [API Guidelines](api/README.md) - Panduan API development
- [Authentication](api/authentication.md) - Sistem autentikasi
- [Error Handling](api/error-handling.md) - Penanganan error
- [Response Format](api/response-format.md) - Format respons

### 💡 Examples & Usage
- [Basic Usage](examples/basic-usage.md) - Contoh penggunaan dasar
- [Advanced Examples](examples/advanced-examples.md) - Contoh lanjutan
- [Best Practices](examples/best-practices.md) - Best practices
- [Code Samples](examples/code-samples.md) - Sample kode

## 🛠️ Tools & CLI

### Command Line Interface
- [CLI Commands](guides/cli.md) - Semua command yang tersedia
- [Migration Commands](database/migrations.md#cli-commands) - Command migrasi
- [Database Commands](guides/cli.md#database-management-commands) - Command database

### Development Tools
- Hot reload dengan Air
- Swagger documentation
- Database migration tools
- Health monitoring

## 🎯 Use Cases

### Scenario Penggunaan
1. **E-commerce Platform**
   - MySQL untuk transactional data
   - PostgreSQL untuk analytics
   - Redis untuk caching

2. **SaaS Application**
   - Multi-tenant database setup
   - Cross-database reporting
   - Data synchronization

3. **Microservices**
   - Database per service
   - Event-driven architecture
   - API gateway integration

## 🔧 Configuration Examples

### Basic Configuration
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
```

### Advanced Configuration
```env
# Connection pooling
MYSQL_MAX_IDLE_CONNS=10
MYSQL_MAX_OPEN_CONNS=200
MYSQL_CONN_MAX_LIFETIME=15m

# PostgreSQL SSL
POSTGRES_SSLMODE=require
POSTGRES_TIMEZONE=UTC
```

## 🚀 Quick Commands

```bash
# Database setup
go run main.go db:connections
go run main.go db:status --connection=mysql

# Migrations
go run main.go migrate:all --connection=mysql
go run main.go migrate:all --connection=postgres

# Development
air                    # Hot reload
swag init             # Update API docs
go test ./...         # Run tests
```

## 📋 Feature Matrix

| Feature | MySQL | PostgreSQL | SQLite | SQL Server |
|---------|-------|------------|--------|------------|
| **CRUD Operations** | ✅ | ✅ | ✅ | ✅ |
| **Transactions** | ✅ | ✅ | ✅ | ✅ |
| **Migrations** | ✅ | ✅ | ✅ | ✅ |
| **Connection Pool** | ✅ | ✅ | ✅ | ✅ |
| **JSON Support** | ✅ | ✅ | ✅ | ✅ |
| **Full-text Search** | ✅ | ✅ | ❌ | ✅ |
| **Geospatial** | ✅ | ✅ | ❌ | ✅ |

## 🤝 Contributing

Ingin berkontribusi? Lihat panduan kontribusi:

1. Fork repository
2. Buat feature branch
3. Commit perubahan
4. Push ke branch
5. Buat Pull Request

## 📞 Support

Butuh bantuan? Hubungi kami:

- 📧 Email: support@example.com
- 💬 Discord: [Join Server](https://discord.gg/example)
- 🐛 Issues: [GitHub Issues](https://github.com/RahmatRafiq/golang_starter_kit_2025/issues)
- 📖 Wiki: [Project Wiki](https://github.com/RahmatRafiq/golang_starter_kit_2025/wiki)

## 📄 License

Proyek ini dilisensikan di bawah [MIT License](../LICENSE).

---

<div align="center">

**⭐ Jangan lupa berikan star jika dokumentasi ini membantu! ⭐**

Dibuat dengan ❤️ oleh [RahmatRafiq](https://github.com/RahmatRafiq) & [Dzyfhuba](https://github.com/Dzyfhuba)

</div>
