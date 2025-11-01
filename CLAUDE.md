kiki# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based backend starter kit built with Gin framework, GORM ORM, and JWT authentication. It implements a modular architecture with support for multiple database connections (MySQL, PostgreSQL, SQLite, SQL Server) and includes a custom migration/seeding system.

## Development Commands

### Running the Application

```bash
# Development with hot reload
air

# Production build and run
go build -o main .
./main

# Direct run without build
go run main.go
```

The application runs on port 8080 by default (configurable via `APP_PORT` in `.env`).

### Testing

```bash
# Run all tests
ginkgo -r --race

# Alternative: use Go test
go test ./...

# Run tests with coverage
go test -cover ./...

# Test specific package
go test ./app/controllers
```

### Database Migrations

```bash
# Create new migration file
go run main.go make:migration <prefix_name>
# Example: go run main.go make:migration create_users_table

# Run all pending migrations
go run main.go migrate:all

# Run specific migration file
go run main.go migrate --file <filename>

# Rollback last batch
go run main.go rollback:batch

# Rollback specific batch
go run main.go rollback:batch --batch=<number>

# Rollback specific file
go run main.go rollback --file=<filename>

# Rollback all migrations
go run main.go rollback:all

# Fresh migration (rollback all and re-migrate)
go run main.go migrate:fresh
```

Migration files are located in `app/database/migrations/` with format `YYYYMMDDHHMMSS_<name>.sql`. Each file must have `-- UP` and `-- DOWN` sections.

### Database Seeders

```bash
# Create new seeder
go run main.go make:seeder --name=<seeder_name>
# Example: go run main.go make:seeder --name=users

# Run all seeders
go run main.go db:seed

# Rollback last seeder batch
go run main.go rollback:seeder

# Rollback specific seeder batch
go run main.go rollback:seeder --batch=<number>
```

Seeder files are located in `app/database/seeds/` with format `YYYYMMDDHHMMSS_<name>.go`.

### API Documentation

```bash
# Generate/update Swagger docs
swag init

# Install swag if not available
go install github.com/swaggo/swag/cmd/swag@latest
```

Access Swagger UI at: `http://localhost:8080/swagger/index.html`

## Architecture

### Layered Architecture

The codebase follows a clean layered architecture:

**Controllers → Services → Repositories → Models → Database**

1. **Controllers** (`app/controllers/`): Handle HTTP requests/responses, call services
2. **Services** (`app/services/`): Business logic layer
3. **Repositories** (`app/repositories/`): Data access layer with interfaces in `app/repositories/interfaces/`
4. **Models** (`app/models/`): GORM database models with relationships
5. **Requests** (`app/requests/`): Request validation using `go-playground/validator`
6. **Middleware** (`app/middleware/`): Auth, logging, etc.

### Dependency Injection

Dependencies are injected through constructors in `routes/web.go`:

```go
// Repository instantiation
userRepo := repositories.NewUserRepository(facades.DB)

// Service instantiation with repository injection
userService := services.NewUserService(userRepo)

// Controller instantiation with service injection
userController := controllers.NewUserController(*userService)
```

When creating new features, follow this pattern:
1. Create repository interface in `app/repositories/interfaces/`
2. Implement repository in `app/repositories/`
3. Create service that uses the repository
4. Create controller that uses the service
5. Register routes in `routes/web.go`

### Multi-Database Support

The application uses a database manager pattern (`database/manager.go`) to handle multiple database connections:

- **Primary database**: Accessed via `facades.DB` (GORM instance)
- **Named connections**: Accessed via `facades.GetConnection(name)`
- **Available connections**: "mysql", "postgres", "mysql_secondary"

Configuration is in `config/database_config.go` and reads from environment variables.

To use a different database in a repository:
```go
conn, err := facades.GetConnection("postgres")
if err != nil {
    return err
}
db := conn.DB // Use this GORM instance
```

### Authentication

JWT-based authentication with middleware in `app/middleware/auth_middleware.go`:

- Token passed via `Authorization: Bearer <token>` header
- Middleware sets `user_id` in Gin context
- Can skip auth in development by setting `SKIP_AUTH=true` in `.env`
- JWT service in `app/services/jwt_service.go`

Protected routes use `middleware.AuthMiddleware()` in route groups.

### Model Relationships

Models use GORM associations:
- **User ↔ Role**: Many-to-Many via `user_has_role` pivot table
- **Role ↔ Permission**: Many-to-Many via `role_has_permission` pivot table
- **Category ↔ Product**: One-to-Many relationship
- Use `Preload()` to eager load relationships

### Pagination

Pagination helper is in `app/models/scopes/pagination.go`:

```go
db.Scopes(scopes.Paginate(page, limit)).Find(&results)
```

### Response Format

Standard response helper in `app/helpers/response_helper.go`:

```go
helpers.ResponseSuccess(c, &helpers.ResponseParams[DataType]{
    Message: "Success message",
    Data: data,
})

helpers.ResponseError(c, &helpers.ResponseParams[any]{
    Reference: "ERROR-1",
    Message: "Error message",
}, http.StatusBadRequest)
```

## Important Files

- `main.go`: Entry point, delegates to `bootstrap/main.go`
- `bootstrap/main.go`: Initializes app, handles CLI commands vs HTTP server
- `routes/web.go`: Route registration and dependency injection
- `facades/database.go`: Database facade for easy access to connections
- `config/database_config.go`: Multi-database configuration
- `database/manager.go`: Database connection manager (singleton pattern)

## Environment Configuration

Copy `.env.example` to `.env` and configure:

- Database credentials: `MYSQL_*`, `POSTGRES_*`, etc.
- JWT secret: Configure in code or env
- Application settings: `APP_NAME`, `APP_PORT`, `APP_HOST`, etc.
- Development: `SKIP_AUTH=true` to bypass authentication

## Key Dependencies

- **Gin**: HTTP web framework
- **GORM**: ORM with support for MySQL, PostgreSQL, SQLite, SQL Server
- **JWT**: `github.com/golang-jwt/jwt/v5`
- **Validator**: `github.com/go-playground/validator/v10`
- **Swagger**: `github.com/swaggo/gin-swagger`
- **CLI**: `github.com/urfave/cli/v2` for migration/seeder commands
- **Testing**: Ginkgo & Gomega

## Code Patterns

### Creating New Endpoints

1. Define repository interface in `app/repositories/interfaces/`
2. Implement repository in `app/repositories/`
3. Create service using the repository
4. Create controller using the service
5. Add request validation structs in `app/requests/`
6. Register routes in `routes/web.go`
7. Add Swagger annotations to controller methods
8. Run `swag init` to update API docs

### Migration File Format

```sql
-- UP
CREATE TABLE example (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- DOWN
DROP TABLE example;
```

### Working with Transactions

```go
err := facades.DB.Transaction(func(tx *gorm.DB) error {
    // Database operations
    if err := tx.Create(&record).Error; err != nil {
        return err // Rollback
    }
    return nil // Commit
})
```

### Security Considerations

- JWT secret key is hardcoded in `app/middleware/auth_middleware.go:17` - should be moved to environment variable
- CORS is configured to allow all origins (`*`) in `bootstrap/main.go:76` - restrict in production
- Security warning helper in `app/helpers/security_warning.go` checks environment on startup
