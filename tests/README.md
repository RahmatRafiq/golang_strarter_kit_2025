# Testing Guide - Golang Starter Kit 2025

Complete testing infrastructure with Laravel-inspired patterns for integration and unit testing.

## Table of Contents
- [Overview](#overview)
- [Quick Start](#quick-start)
- [Test Structure](#test-structure)
- [Running Tests](#running-tests)
- [Writing Tests](#writing-tests)
- [Test Database](#test-database)
- [Factories](#factories)
- [Best Practices](#best-practices)

## Overview

This project implements a comprehensive testing infrastructure inspired by Laravel's testing patterns:

- **Integration Tests**: Test complete workflows with real database (like Laravel Feature Tests)
- **Unit Tests**: Test individual functions with mocks (like Laravel Unit Tests)
- **Test Database**: Automatic setup and teardown of test database (like Laravel RefreshDatabase)
- **Factories**: Convenient test data generation (like Laravel Factories)
- **CLI Command**: Run tests with `go run main.go test` (like `php artisan test`)

### Test Coverage
✅ **7 Integration Test Suites** - 18 Sub-tests (100% passing)
- Authentication flow (Login, Refresh, Logout, Complete Flow)
- User CRUD operations
- Role assignment and management
- Pagination and bulk operations

## Quick Start

### 1. Setup Test Environment

Create `.env.testing` file (already configured):
```bash
# Test database will be auto-created
APP_ENV=testing
MYSQL_DB=golang_starter_kit_2025_test
MYSQL_USER=root
MYSQL_PASSWORD=your_password
```

### 2. Run All Tests
```bash
# Run all tests
go run main.go test

# Run integration tests only
go run main.go test --type integration

# Run with verbose output
go run main.go test -v

# Run with coverage report
go run main.go test --coverage
```

### 3. Run Specific Tests
```bash
# Filter by test name
go run main.go test --filter Auth

# Setup database before running
go run main.go test --setup
```

## Test Structure

```
tests/
├── README.md                    # This file
├── integration/                 # Integration tests (with real DB)
│   ├── auth_test.go            # Auth service tests
│   └── user_test.go            # User service tests
├── unit/                        # Unit tests (with mocks)
│   └── (coming soon)
├── e2e/                         # End-to-end tests
│   └── (coming soon)
├── helpers/                     # Test helpers
│   ├── database.go             # Database setup/teardown
│   └── factories.go            # Test data factories
└── testdata/                    # Test fixtures
```

## Running Tests

### Using CLI Command (Recommended)
```bash
# Show help
go run main.go test --help

# Run all tests
go run main.go test

# Run by type
go run main.go test --type integration
go run main.go test --type unit
go run main.go test --type e2e

# Verbose output
go run main.go test -v

# Coverage report
go run main.go test --coverage

# Filter tests
go run main.go test --filter "Login"
go run main.go test --filter "User.*CRUD"
```

### Using Go Test Directly
```bash
# Run all integration tests
go test ./tests/integration/... -v

# Run specific test
go test ./tests/integration/... -v -run TestAuthService_Login

# Run with coverage
go test ./tests/integration/... -cover
```

## Writing Tests

### Integration Tests

Integration tests use real database connections and test complete workflows.

#### Basic Structure
```go
func TestMyService_Integration(t *testing.T) {
    // Setup: Create test database
    testDB := helpers.SetupTestDB(t)
    defer helpers.CleanupTestDB(t, testDB)

    // Create service with real repository
    repo := repositories.NewMyRepository(testDB.DB)
    service := services.NewMyService(repo)

    t.Run("success - happy path", func(t *testing.T) {
        // Arrange: Create test data
        user := helpers.NewUserFactory(t, testDB).
            WithEmail("test@example.com").
            Create()

        // Act: Execute the operation
        result, err := service.DoSomething(user.ID)

        // Assert: Verify results
        if err != nil {
            t.Errorf("Expected no error, got: %v", err)
        }
        if result == nil {
            t.Error("Expected result, got nil")
        }
    })

    t.Run("error - invalid input", func(t *testing.T) {
        // Test error cases
        result, err := service.DoSomething(0)

        if err == nil {
            t.Error("Expected error, got nil")
        }
    })
}
```

#### Example: Authentication Test
```go
func TestAuthService_Login_Integration(t *testing.T) {
    testDB := helpers.SetupTestDB(t)
    defer helpers.CleanupTestDB(t, testDB)

    authRepo := repositories.NewAuthRepository(testDB.DB)
    authService := services.NewAuthService(authRepo)

    t.Run("success - valid credentials", func(t *testing.T) {
        // Create user with known password
        user, plainPassword := helpers.NewUserFactory(t, testDB).
            WithEmail("john@example.com").
            WithPassword("secure_password_123").
            CreateWithPlainPassword()

        // Attempt login
        token, err := authService.Login(requests.LoginRequest{
            Email:    user.Email,
            Password: plainPassword,
        })

        // Verify success
        if err != nil {
            t.Errorf("Expected successful login, got error: %v", err)
        }
        if token == nil || token.Token == "" {
            t.Error("Expected valid token")
        }
    })
}
```

### Unit Tests

Unit tests use mocks and test individual functions in isolation.

```go
func TestAuthService_Login(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockAuthRepositoryInterface(ctrl)
    authService := services.NewAuthService(mockRepo)

    t.Run("success - valid credentials", func(t *testing.T) {
        expectedUser := &models.User{
            ID:    1,
            Email: "test@example.com",
        }

        mockRepo.EXPECT().
            FindUserByEmail("test@example.com").
            Return(expectedUser, nil)

        // Test implementation
    })
}
```

## Test Database

### Automatic Setup

The test database is automatically managed:

1. **Database Creation**: Fresh database created for each test run
2. **Migrations**: GORM AutoMigrate runs all models
3. **Junction Tables**: Many-to-many tables created automatically
4. **Cleanup**: Database closed after tests (recreated next run)

### Database Helper Functions

```go
// Setup test database
testDB := helpers.SetupTestDB(t)
defer helpers.CleanupTestDB(t, testDB)

// Access database
testDB.DB        // *gorm.DB instance
testDB.Connection // "mysql" or "postgres"
testDB.DBName     // "golang_starter_kit_2025_test"
```

### Configuration

Test database configuration in `.env.testing`:

```env
APP_ENV=testing
DB_CONNECTION=mysql

# MySQL Test Database
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=golang_starter_kit_2025_test
MYSQL_USER=root
MYSQL_PASSWORD=your_password

# Smaller pool for testing
MYSQL_MAX_IDLE_CONNS=5
MYSQL_MAX_OPEN_CONNS=10
```

## Factories

Factories make it easy to create test data with sensible defaults.

### User Factory

```go
// Basic user
user := helpers.NewUserFactory(t, testDB).Create()

// User with custom fields
user := helpers.NewUserFactory(t, testDB).
    WithEmail("john@example.com").
    WithUsername("john_doe").
    WithPassword("secure_password").
    Create()

// User with plain password (for login tests)
user, plainPassword := helpers.NewUserFactory(t, testDB).
    WithEmail("test@example.com").
    WithPassword("password123").
    CreateWithPlainPassword()
```

### Role Factory

```go
// Basic role
role := helpers.NewRoleFactory(t, testDB).Create()

// Role with custom name
role := helpers.NewRoleFactory(t, testDB).
    WithName("admin").
    Create()
```

### Permission Factory

```go
// Basic permission
perm := helpers.NewPermissionFactory(t, testDB).Create()

// Permission with custom name
perm := helpers.NewPermissionFactory(t, testDB).
    WithName("create:users").
    Create()
```

### Helper Functions

```go
// Create user with role assigned
user, role := helpers.CreateUserWithRole(t, testDB, "admin")

// Create role with permissions
role := helpers.CreateRoleWithPermissions(t, testDB, "admin", []string{
    "create:users",
    "delete:users",
})
```

## Best Practices

### 1. Test Organization

```go
func TestFeature_Integration(t *testing.T) {
    // Setup once for all sub-tests
    testDB := helpers.SetupTestDB(t)
    defer helpers.CleanupTestDB(t, testDB)

    t.Run("success - happy path", func(t *testing.T) {
        // Test success cases
    })

    t.Run("error - validation failure", func(t *testing.T) {
        // Test error cases
    })

    t.Run("edge case - boundary condition", func(t *testing.T) {
        // Test edge cases
    })
}
```

### 2. Use Descriptive Test Names

✅ Good:
```go
t.Run("success - user can login with valid credentials", func(t *testing.T) {})
t.Run("error - login fails with invalid password", func(t *testing.T) {})
t.Run("success - refresh token generates new access token", func(t *testing.T) {})
```

❌ Bad:
```go
t.Run("test1", func(t *testing.T) {})
t.Run("login", func(t *testing.T) {})
```

### 3. Follow AAA Pattern

```go
t.Run("success - example", func(t *testing.T) {
    // Arrange: Set up test data
    user := helpers.NewUserFactory(t, testDB).Create()

    // Act: Execute the operation
    result, err := service.DoSomething(user.ID)

    // Assert: Verify results
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
})
```

### 4. Test Error Cases

Always test both success and error scenarios:

```go
t.Run("success - operation succeeds", func(t *testing.T) {
    // Test happy path
})

t.Run("error - invalid input", func(t *testing.T) {
    // Test validation errors
})

t.Run("error - not found", func(t *testing.T) {
    // Test resource not found
})

t.Run("error - unauthorized", func(t *testing.T) {
    // Test permission errors
})
```

### 5. Use t.Helper() in Helper Functions

```go
func createTestUser(t *testing.T, db *TestDB) *models.User {
    t.Helper() // Marks this as helper, errors show caller's line

    user := helpers.NewUserFactory(t, db).Create()
    return user
}
```

### 6. Clean Assertions

✅ Good:
```go
if err != nil {
    t.Errorf("Expected no error, got: %v", err)
}
if user == nil {
    t.Fatal("Expected user, got nil")
}
if len(users) != 3 {
    t.Errorf("Expected 3 users, got %d", len(users))
}
```

❌ Bad:
```go
if err != nil {
    t.Error("error") // Not descriptive
}
if user == nil {
    panic("no user") // Don't use panic
}
```

### 7. Use t.Log for Documentation

```go
t.Run("complete authentication flow", func(t *testing.T) {
    t.Log("Step 1: Login with credentials")
    token, _ := authService.Login(loginReq)

    t.Log("Step 2: Refresh access token")
    newToken, _ := authService.RefreshToken(token.RefreshToken)

    t.Log("Step 3: Logout and clear tokens")
    authService.Logout(newToken.Token)

    t.Log("✓ Complete authentication flow successful")
})
```

### 8. Soft Delete Testing

When testing soft delete, check both:

```go
// Verify DeletedAt is set
var deletedUser models.User
testDB.DB.Unscoped().Where("id = ?", user.ID).First(&deletedUser)
if deletedUser.DeletedAt.Time.IsZero() {
    t.Error("Expected user to be soft deleted")
}

// Verify hidden from normal queries
var count int64
testDB.DB.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
if count != 0 {
    t.Error("Expected soft deleted user to be hidden")
}
```

### 9. Test Pagination

```go
t.Run("success - pagination works correctly", func(t *testing.T) {
    // Create 50 users
    for i := 1; i <= 50; i++ {
        helpers.NewUserFactory(t, testDB).Create()
    }

    // Test first page
    page1, _ := userService.List(1, 20)
    if len(page1) != 20 {
        t.Errorf("Expected 20 users on page 1, got %d", len(page1))
    }

    // Test last page
    page3, _ := userService.List(3, 20)
    if len(page3) != 10 {
        t.Errorf("Expected 10 users on page 3, got %d", len(page3))
    }
})
```

### 10. Database Schema Issues

Always use correct table names:
- ✅ `user_has_roles` (correct)
- ❌ `user_roles` (wrong)
- ✅ `role_has_permissions` (correct)
- ❌ `role_permissions` (wrong)

When querying, use Model() for automatic soft delete handling:
```go
// ✅ Correct - respects soft delete
testDB.DB.Model(&models.User{}).Where("id = ?", id).Count(&count)

// ❌ Wrong - ignores soft delete
testDB.DB.Table("users").Where("id = ?", id).Count(&count)
```

## Troubleshooting

### Test Database Issues

**Problem**: Database connection fails
```bash
# Check MySQL is running
mysql -u root -p -e "SELECT 1"

# Check .env.testing credentials
cat .env.testing | grep MYSQL
```

**Problem**: Tables not found
```bash
# Verify test database exists
mysql -u root -p -e "SHOW DATABASES LIKE '%test%'"

# Manually recreate if needed
bash tests/setup_test_db.sh
```

### Test Failures

**Problem**: User ID conversion errors
```go
// ❌ Wrong - converts to unicode character
userID := string(rune(user.ID))

// ✅ Correct - converts to string representation
userID := strconv.Itoa(int(user.ID))
```

**Problem**: Facades.DB is nil
```go
// Solution: setupFacades() is called automatically in SetupTestDB
// If you get nil pointer, ensure you're using SetupTestDB helper
```

### Common Errors

1. **"Table 'user_roles' doesn't exist"**: Use `user_has_roles` instead
2. **"invalid syntax parsing user ID"**: Use `strconv.Itoa()` not `string(rune())`
3. **"Soft deleted user still visible"**: Use `.Model(&models.User{})` not `.Table("users")`
4. **"Database name error"**: Ensure `.env.testing` is loaded with `godotenv.Overload()`

## Examples

See complete test examples in:
- `tests/integration/auth_test.go` - Authentication flows
- `tests/integration/user_test.go` - User CRUD and role management

## Future Enhancements

Planned improvements:
- [ ] E2E API tests with HTTP requests
- [ ] Performance benchmarks
- [ ] Test coverage reports
- [ ] Parallel test execution
- [ ] Database transaction rollback per test
- [ ] Snapshot testing for responses

## Contributing

When adding new tests:
1. Follow the established patterns
2. Use factories for test data
3. Test both success and error cases
4. Add descriptive test names
5. Update this README if needed

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [GORM Testing](https://gorm.io/docs/testing.html)
- [Laravel Testing](https://laravel.com/docs/testing) (inspiration)
- [Testify Library](https://github.com/stretchr/testify) (optional)
