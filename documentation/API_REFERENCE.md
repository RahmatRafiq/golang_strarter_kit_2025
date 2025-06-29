# API Reference

## Authentication Endpoints

### Login
```http
POST /api/auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "error": false,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "user",
      "email": "user@example.com"
    }
  }
}
```

### Register
```http
POST /api/auth/register
```

**Request Body:**
```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123"
}
```

### Logout
```http
POST /api/auth/logout
```
**Headers:**
```
Authorization: Bearer <token>
```

## User Management

### Get Users
```http
GET /api/users
```

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page
- `search` (optional): Search term

### Get User by ID
```http
GET /api/users/{id}
```

### Create User
```http
POST /api/users
```

**Request Body:**
```json
{
  "username": "newuser",
  "email": "user@example.com",
  "password": "password123"
}
```

### Update User
```http
PUT /api/users/{id}
```

### Delete User
```http
DELETE /api/users/{id}
```

## Product Management

### Get Products
```http
GET /api/products
```

### Get Product by ID
```http
GET /api/products/{id}
```

### Create Product
```http
POST /api/products
```

**Request Body:**
```json
{
  "name": "Product Name",
  "description": "Product description",
  "price": 99.99,
  "category_id": 1
}
```

### Update Product
```http
PUT /api/products/{id}
```

### Delete Product
```http
DELETE /api/products/{id}
```

## Category Management

### Get Categories
```http
GET /api/categories
```

### Create Category
```http
POST /api/categories
```

**Request Body:**
```json
{
  "name": "Category Name",
  "description": "Category description"
}
```

## Role & Permission Management

### Get Roles
```http
GET /api/roles
```

### Create Role
```http
POST /api/roles
```

**Request Body:**
```json
{
  "name": "admin",
  "description": "Administrator role"
}
```

### Assign Role to User
```http
POST /api/users/{user_id}/roles/{role_id}
```

### Get Permissions
```http
GET /api/permissions
```

### Create Permission
```http
POST /api/permissions
```

**Request Body:**
```json
{
  "name": "create_user",
  "description": "Can create users"
}
```

## Database Management API

### Get Database Status
```http
GET /api/database/status
```

**Response:**
```json
{
  "error": false,
  "message": "Database connection status retrieved successfully",
  "data": {
    "mysql": {
      "connected": true,
      "open_connections": 5,
      "in_use": 1,
      "idle": 4
    },
    "postgres": {
      "connected": true,
      "open_connections": 3,
      "in_use": 0,
      "idle": 3
    }
  }
}
```

### Database Health Check
```http
GET /api/database/health
```

### Test Database Connection
```http
GET /api/database/test?connection=mysql
```

**Query Parameters:**
- `connection`: Database connection name (mysql, postgres, mysql_secondary)

### Sync Data Between Databases
```http
POST /api/database/sync?source=mysql&target=postgres
```

**Query Parameters:**
- `source`: Source database connection
- `target`: Target database connection
- `table` (optional): Specific table to sync

## File Upload

### Upload File
```http
POST /api/files/upload
```

**Request:**
- Content-Type: `multipart/form-data`
- Field: `file`

**Response:**
```json
{
  "error": false,
  "message": "File uploaded successfully",
  "data": {
    "filename": "uploaded_file.jpg",
    "url": "/storage/uploads/uploaded_file.jpg",
    "size": 1024000
  }
}
```

## Error Responses

All endpoints may return these error responses:

### 400 Bad Request
```json
{
  "error": true,
  "message": "Invalid request data",
  "details": "Validation error details"
}
```

### 401 Unauthorized
```json
{
  "error": true,
  "message": "Unauthorized access",
  "details": "Invalid or expired token"
}
```

### 403 Forbidden
```json
{
  "error": true,
  "message": "Access forbidden",
  "details": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "error": true,
  "message": "Resource not found",
  "details": "The requested resource does not exist"
}
```

### 500 Internal Server Error
```json
{
  "error": true,
  "message": "Internal server error",
  "details": "Server error details"
}
```

## Rate Limiting

API endpoints are rate limited:
- **Default**: 100 requests per minute
- **Authentication**: 5 requests per minute
- **File Upload**: 10 requests per minute

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## Pagination

List endpoints support pagination:

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `sort`: Sort field (default: id)
- `order`: Sort order (asc/desc, default: asc)

**Response:**
```json
{
  "error": false,
  "message": "Data retrieved successfully",
  "data": [...],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10,
    "has_next": true,
    "has_prev": false
  }
}
```

## WebSocket Endpoints

### Real-time Notifications
```
ws://localhost:8080/ws/notifications
```

### Database Status Updates
```
ws://localhost:8080/ws/database-status
```

## Health Check

### Application Health
```http
GET /api/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-06-29T10:30:00Z",
  "version": "1.0.0",
  "uptime": "2h30m15s",
  "checks": {
    "database": "healthy",
    "redis": "healthy",
    "external_api": "healthy"
  }
}
```
