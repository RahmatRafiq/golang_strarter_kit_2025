# API Reference

Complete API endpoint reference for Golang Starter Kit 2025. All endpoints return JSON responses and support standard HTTP status codes.

## Base URL

```
http://localhost:9999
```

For production, replace with your domain.

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

### Getting a Token

Login to get a JWT token:

```bash
curl -X PUT http://localhost:9999/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

## Response Format

### Success Response

```json
{
  "message": "Success message",
  "data": {
    ...
  }
}
```

### Error Response

```json
{
  "reference": "ERROR-001",
  "message": "Error description"
}
```

### Paginated Response

```json
{
  "message": "Data retrieved successfully",
  "data": [...],
  "total": 100,
  "page": 1,
  "limit": 10
}
```

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Missing or invalid token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource doesn't exist |
| 500 | Internal Server Error |
| 503 | Service Unavailable - Database connection issue |

---

## Endpoints

### General

#### Hello World
```http
GET /
```

**Description**: Basic endpoint to verify API is running.

**Response:**
```json
{
  "message": "Hello World"
}
```

---

## Authentication

### Login
```http
PUT /auth/login
```

**Description**: Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Validation:**
- `email`: Required, must be valid email format
- `password`: Required

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expired_at": "2025-01-01T12:00:00Z"
}
```

**Error Responses:**
- `400`: Invalid credentials
- `404`: User not found

**Example:**
```bash
curl -X PUT http://localhost:9999/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

### Logout
```http
GET /auth/logout
```
🔒 **Requires Authentication**

**Description**: Invalidate current JWT token.

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200):**
```json
{
  "message": "Logout successful"
}
```

**Example:**
```bash
curl -X GET http://localhost:9999/auth/logout \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Refresh Token
```http
GET /auth/refresh
```
🔒 **Requires Authentication**

**Description**: Get a new JWT token with fresh expiry.

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expired_at": "2025-01-01T13:00:00Z"
}
```

---

## Users

### List Users
```http
GET /users
```
🔒 **Requires Authentication**

**Description**: Get paginated list of users.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)

**Success Response (200):**
```json
{
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "reference": "USR-20250101-ABC123",
      "username": "admin",
      "email": "admin@example.com",
      "created_at": "2025-01-01T10:00:00Z",
      "roles": [
        {
          "id": 1,
          "name": "admin",
          "group": "system"
        }
      ]
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 10
}
```

**Example:**
```bash
curl -X GET "http://localhost:9999/users?page=1&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get User by ID
```http
GET /users/:id
```
🔒 **Requires Authentication**

**Description**: Get detailed information for a specific user.

**Path Parameters:**
- `id`: User ID (integer)

**Success Response (200):**
```json
{
  "message": "User found",
  "data": {
    "id": 1,
    "reference": "USR-20250101-ABC123",
    "username": "admin",
    "email": "admin@example.com",
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T10:00:00Z",
    "roles": [...]
  }
}
```

**Error Response:**
- `404`: User not found

### Create/Update User
```http
PUT /users
```
🔒 **Requires Authentication**

**Description**: Create new user or update existing user (upsert).

**Request Body:**
```json
{
  "id": 0,
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securepassword123",
  "pin": "123456"
}
```

**Validation:**
- `username`: Required, unique
- `email`: Required, valid email format, unique
- `password`: Required for new users
- `pin`: Optional

**Success Response (200/201):**
```json
{
  "message": "User created/updated successfully",
  "data": {
    "id": 2,
    "reference": "USR-20250102-XYZ789",
    "username": "newuser",
    "email": "newuser@example.com"
  }
}
```

**Note**: Password is automatically hashed using Argon2id.

### Delete User
```http
DELETE /users/:id
```
🔒 **Requires Authentication**

**Description**: Soft delete a user.

**Path Parameters:**
- `id`: User ID (integer)

**Success Response (200):**
```json
{
  "message": "User deleted successfully"
}
```

### Assign Roles to User
```http
POST /users/:id/roles
```
🔒 **Requires Authentication**

**Description**: Assign multiple roles to a user.

**Path Parameters:**
- `id`: User ID (integer)

**Request Body:**
```json
{
  "role_ids": [1, 2, 3]
}
```

**Success Response (200):**
```json
{
  "message": "Roles assigned successfully"
}
```

### Get User's Roles
```http
GET /users/:id/roles
```
🔒 **Requires Authentication**

**Description**: Get all roles assigned to a user.

**Path Parameters:**
- `id`: User ID (integer)

**Success Response (200):**
```json
{
  "message": "User roles retrieved",
  "data": [
    {
      "id": 1,
      "name": "admin",
      "group": "system"
    },
    {
      "id": 2,
      "name": "editor",
      "group": "content"
    }
  ]
}
```

---

## Roles

### List Roles
```http
GET /roles
```
🔒 **Requires Authentication**

**Description**: Get all roles with their permissions.

**Success Response (200):**
```json
{
  "message": "Roles retrieved",
  "data": [
    {
      "id": 1,
      "name": "admin",
      "group": "system",
      "permissions": [
        {
          "id": 1,
          "name": "users.create",
          "group": "users"
        }
      ]
    }
  ]
}
```

### Create/Update Role
```http
PUT /roles
```
🔒 **Requires Authentication**

**Description**: Create new role or update existing role.

**Request Body:**
```json
{
  "id": 0,
  "name": "moderator",
  "group": "content"
}
```

**Validation:**
- `name`: Required, unique
- `group`: Required

**Success Response (200/201):**
```json
{
  "message": "Role created/updated successfully",
  "data": {
    "id": 3,
    "name": "moderator",
    "group": "content"
  }
}
```

### Delete Role
```http
DELETE /roles/:id
```
🔒 **Requires Authentication**

**Description**: Delete a role.

**Path Parameters:**
- `id`: Role ID (integer)

**Success Response (200):**
```json
{
  "message": "Role deleted successfully"
}
```

**Note**: Deleting a role also removes all user-role and role-permission associations (cascade delete).

### Assign Permissions to Role
```http
POST /roles/:id/permissions
```
🔒 **Requires Authentication**

**Description**: Assign multiple permissions to a role.

**Path Parameters:**
- `id`: Role ID (integer)

**Request Body:**
```json
{
  "permission_ids": [1, 2, 3, 4, 5]
}
```

**Success Response (200):**
```json
{
  "message": "Permissions assigned successfully"
}
```

### Get Role's Permissions
```http
GET /roles/:id/permissions
```
🔒 **Requires Authentication**

**Description**: Get all permissions assigned to a role.

**Path Parameters:**
- `id`: Role ID (integer)

**Success Response (200):**
```json
{
  "message": "Role permissions retrieved",
  "data": [
    {
      "id": 1,
      "name": "users.create",
      "group": "users"
    },
    {
      "id": 2,
      "name": "users.read",
      "group": "users"
    }
  ]
}
```

---

## Permissions

### List Permissions
```http
GET /permissions
```
🔒 **Requires Authentication**

**Description**: Get all available permissions.

**Success Response (200):**
```json
{
  "message": "Permissions retrieved",
  "data": [
    {
      "id": 1,
      "name": "users.create",
      "group": "users"
    },
    {
      "id": 2,
      "name": "users.read",
      "group": "users"
    },
    {
      "id": 3,
      "name": "products.create",
      "group": "products"
    }
  ]
}
```

### Create/Update Permission
```http
PUT /permissions
```
🔒 **Requires Authentication**

**Description**: Create new permission or update existing permission.

**Request Body:**
```json
{
  "id": 0,
  "name": "orders.create",
  "group": "orders"
}
```

**Validation:**
- `name`: Required, unique, format: `resource.action`
- `group`: Required

**Success Response (200/201):**
```json
{
  "message": "Permission created/updated successfully",
  "data": {
    "id": 10,
    "name": "orders.create",
    "group": "orders"
  }
}
```

**Permission Naming Convention:**
- Format: `resource.action`
- Examples: `users.create`, `products.update`, `orders.delete`

### Delete Permission
```http
DELETE /permissions/:id
```
🔒 **Requires Authentication**

**Description**: Delete a permission.

**Path Parameters:**
- `id`: Permission ID (integer)

**Success Response (200):**
```json
{
  "message": "Permission deleted successfully"
}
```

---

## Products

### List Products
```http
GET /products
```
🔒 **Requires Authentication**

**Description**: Get paginated list of products with categories.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)

**Success Response (200):**
```json
{
  "message": "Products retrieved",
  "data": [
    {
      "id": 1,
      "reference": "PRD-20250101-ABC123",
      "name": "Laptop",
      "description": "High-performance laptop",
      "price": 1200.50,
      "stock": 50,
      "category_id": 1,
      "category": {
        "id": 1,
        "category": "Electronics"
      },
      "images": [
        "http://localhost:9999/file/products/image1.jpg",
        "http://localhost:9999/file/products/image2.jpg"
      ],
      "created_at": "2025-01-01T10:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "limit": 10
}
```

**Note**: Image URLs are automatically generated from stored paths.

### Get Product by ID
```http
GET /products/:id
```
🔒 **Requires Authentication**

**Description**: Get detailed product information.

**Path Parameters:**
- `id`: Product ID (integer)

**Success Response (200):**
```json
{
  "message": "Product found",
  "data": {
    "id": 1,
    "reference": "PRD-20250101-ABC123",
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 1200.50,
    "margin": 20.0,
    "stock": 50,
    "sold": 10,
    "category_id": 1,
    "category": {...},
    "images": [...]
  }
}
```

### Create/Update Product
```http
PUT /products
```
🔒 **Requires Authentication**

**Description**: Create new product or update existing product.

**Request Body:**
```json
{
  "id": 0,
  "name": "Smartphone",
  "description": "Latest model smartphone",
  "price": 899.99,
  "margin": 15.0,
  "stock": 100,
  "category_id": 1,
  "images": [
    "base64-encoded-image-data",
    "base64-encoded-image-data"
  ]
}
```

**Validation:**
- `name`: Required
- `price`: Required, must be >= 0
- `stock`: Required, must be >= 0
- `category_id`: Required, must exist

**Success Response (200/201):**
```json
{
  "message": "Product created/updated successfully",
  "data": {
    "id": 2,
    "reference": "PRD-20250102-XYZ789",
    "name": "Smartphone",
    ...
  }
}
```

**Note**: Images can be sent as base64-encoded strings and will be automatically processed.

### Delete Product
```http
DELETE /products/:id
```
🔒 **Requires Authentication**

**Description**: Soft delete a product.

**Path Parameters:**
- `id`: Product ID (integer)

**Success Response (200):**
```json
{
  "message": "Product deleted successfully"
}
```

---

## Categories

### List Categories
```http
GET /categories
```
🔒 **Requires Authentication**

**Description**: Get all categories with product count.

**Success Response (200):**
```json
{
  "message": "Categories retrieved",
  "data": [
    {
      "id": 1,
      "category": "Electronics",
      "created_at": "2025-01-01T10:00:00Z",
      "products": [
        {
          "id": 1,
          "name": "Laptop",
          "price": 1200.50
        }
      ]
    }
  ]
}
```

### Get Category by ID
```http
GET /categories/:id
```
🔒 **Requires Authentication**

**Description**: Get category with all its products.

**Path Parameters:**
- `id`: Category ID (integer)

**Success Response (200):**
```json
{
  "message": "Category found",
  "data": {
    "id": 1,
    "category": "Electronics",
    "products": [...]
  }
}
```

### Create/Update Category
```http
PUT /categories
```
🔒 **Requires Authentication**

**Description**: Create new category or update existing category.

**Request Body:**
```json
{
  "id": 0,
  "category": "Books"
}
```

**Validation:**
- `category`: Required, unique

**Success Response (200/201):**
```json
{
  "message": "Category created/updated successfully",
  "data": {
    "id": 2,
    "category": "Books"
  }
}
```

### Delete Category
```http
DELETE /categories/:id
```
🔒 **Requires Authentication**

**Description**: Delete a category.

**Path Parameters:**
- `id`: Category ID (integer)

**Success Response (200):**
```json
{
  "message": "Category deleted successfully"
}
```

**Note**: Cannot delete category if it has associated products.

---

## Files

### Serve File
```http
GET /file/:key/:filename
```

**Description**: Serve uploaded files (images, documents, etc.).

**Path Parameters:**
- `key`: Storage folder (e.g., "products", "users")
- `filename`: File name

**Success Response (200):**
- Returns file binary with appropriate Content-Type header

**Error Response:**
- `404`: File not found

**Example:**
```
http://localhost:9999/file/products/image123.jpg
http://localhost:9999/file/users/avatar456.png
```

**Note**: This endpoint is public (no authentication required) for serving static files.

---

## Database Management

### Get Connection Status
```http
GET /api/database/status
```

**Description**: Get status and statistics for all database connections.

**Success Response (200):**
```json
{
  "connections": {
    "mysql": {
      "status": "connected",
      "stats": {
        "max_open_connections": 200,
        "open_connections": 5,
        "in_use": 2,
        "idle": 3,
        "wait_count": 0,
        "wait_duration": 0
      }
    },
    "postgres": {
      "status": "connected",
      "stats": {...}
    },
    "mysql_secondary": {
      "status": "disconnected"
    }
  }
}
```

### Health Check (All Databases)
```http
GET /api/database/health
```

**Description**: Check health of all configured database connections.

**Success Response (200):**
```json
{
  "overall_health": true,
  "databases": {
    "mysql": "healthy",
    "postgres": "healthy",
    "mysql_secondary": "disconnected"
  }
}
```

**Error Response (503):**
```json
{
  "overall_health": false,
  "databases": {
    "mysql": "unhealthy",
    "error": "connection refused"
  }
}
```

### Test Specific Connection
```http
GET /api/database/test?connection=mysql
```

**Description**: Test a specific database connection.

**Query Parameters:**
- `connection`: Connection name (mysql, postgres, mysql_secondary)

**Success Response (200):**
```json
{
  "connection": "postgres",
  "status": "healthy",
  "response_time": "5ms"
}
```

---

## Health Checks

### Basic Health Check
```http
GET /health
```

**Description**: Check if API server and default database are healthy.

**Success Response (200):**
```json
{
  "message": "database is connected",
  "database": "golang_starter_kit_2025"
}
```

**Error Response (500):**
```json
{
  "message": "database connection failed",
  "error": "connection refused"
}
```

### Multi-Database Health Check
```http
GET /health/databases
```

**Description**: Comprehensive health check for all database connections.

**Success Response (200):**
```json
{
  "overall_health": true,
  "connections": {
    "mysql": {
      "status": "healthy",
      "stats": {
        "max_open_connections": 200,
        "open_connections": 5,
        "in_use": 2,
        "idle": 3,
        "wait_count": 0,
        "wait_duration": 0,
        "max_idle_closed": 0,
        "max_lifetime_closed": 0
      }
    },
    "postgres": {
      "status": "healthy",
      "stats": {...}
    },
    "mysql_secondary": {
      "status": "disconnected"
    }
  }
}
```

**Error Response (503):**
```json
{
  "overall_health": false,
  "connections": {
    "mysql": {
      "status": "unhealthy",
      "error": "connection timeout"
    }
  }
}
```

---

## Test Endpoints (Development)

### List Test Data
```http
GET /tests
```

**Description**: Get test data from PostgreSQL database.

### Get Test by ID
```http
GET /tests/:id
```

### Create Test
```http
POST /tests
```

### Update Test
```http
PUT /tests/:id
```

### Delete Test
```http
DELETE /tests/:id
```

**Note**: These endpoints are for testing multi-database functionality.

---

## Error Reference

### Common Error Codes

| Reference | Description | HTTP Status |
|-----------|-------------|-------------|
| ERROR-1 | Token required | 401 |
| ERROR-2 | Invalid token format | 401 |
| ERROR-3 | Invalid token signature | 401 |
| ERROR-4 | Token expired | 401 |
| ERROR-USER-001 | User not found | 404 |
| ERROR-PRODUCT-001 | Product not found | 404 |
| ERROR-ROLE-001 | Role not found | 404 |

### Example Error Response

```json
{
  "reference": "ERROR-4",
  "message": "Token sudah kadaluarsa"
}
```

---

## Pagination

All list endpoints support pagination with these query parameters:

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `page` | integer | 1 | - | Page number |
| `limit` | integer | 10 | 100 | Items per page |
| `offset` | integer | - | - | Manual offset (overrides page) |

**Example:**
```bash
# Get page 2 with 20 items
GET /users?page=2&limit=20

# Manual offset
GET /users?offset=40&limit=20
```

**Response includes pagination metadata:**
```json
{
  "data": [...],
  "total": 150,
  "page": 2,
  "limit": 20
}
```

---

## Rate Limiting

Currently not implemented. Consider adding rate limiting for production deployments.

**Recommended limits:**
- Authentication endpoints: 5 requests/minute
- List endpoints: 100 requests/minute
- File uploads: 10 requests/minute

---

## CORS

CORS is configured via `ALLOWED_ORIGINS` environment variable:

```bash
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,https://yourdomain.com
```

**Allowed methods:** GET, POST, PUT, PATCH, DELETE, OPTIONS

**Allowed headers:** Origin, Content-Type, Authorization, X-Api-Key, X-Request-ID

---

## Development Mode

### Skip Authentication

Set `SKIP_AUTH=true` in `.env` to bypass authentication for all endpoints:

```bash
SKIP_AUTH=true  # ⚠️ DEVELOPMENT ONLY!
```

When enabled:
- No Authorization header required
- Automatically sets `user_id=1` in context
- ⚠️ **NEVER use in production!**

---

## Swagger Documentation

Interactive API documentation is available at:

```
http://localhost:9999/swagger/index.html
```

**Features:**
- Try out API requests directly
- View request/response schemas
- See validation rules
- Copy curl commands

**Update Swagger docs:**
```bash
swag init
```

---

## Best Practices

### 1. Always Use HTTPS in Production

```bash
APP_SCHEME=https
```

### 2. Validate Input

All endpoints use `go-playground/validator` for input validation.

### 3. Handle Errors Gracefully

Check HTTP status codes and error messages in responses.

### 4. Use Pagination

Always use pagination for list endpoints to avoid performance issues.

### 5. Secure Your JWT Secret

```bash
# Generate strong secret
openssl rand -base64 48

# Set in .env
JWT_SECRET_KEY=<generated-key>
```

### 6. Monitor Rate Limits

Implement rate limiting for production deployments.

### 7. Use Appropriate HTTP Methods

- `GET`: Retrieve resources
- `POST`: Create specific resources
- `PUT`: Create or update (upsert)
- `DELETE`: Remove resources

---

## Further Reading

- [GETTING_STARTED.md](./GETTING_STARTED.md) - Setup guide
- [DATABASE.md](./DATABASE.md) - Migration & seeder commands
- [MULTI_DATABASE.md](./MULTI_DATABASE.md) - Multi-database configuration
- [CLAUDE.md](../CLAUDE.md) - Complete technical documentation
- [Swagger UI](http://localhost:9999/swagger/index.html) - Interactive API docs

---

**Need help?** Check our [GitHub Discussions](https://github.com/RahmatRafiq/golang_starter_kit_2025/discussions)
