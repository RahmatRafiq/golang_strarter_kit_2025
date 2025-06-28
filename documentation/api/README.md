# API Documentation

## Overview

API ini dibangun menggunakan Gin framework dengan dukungan multiple database connections. Semua endpoint menggunakan JSON untuk request dan response, dengan sistem autentikasi JWT yang aman.

## Base URL

```
Development: http://localhost:8080
Production: https://your-domain.com
```

## Authentication

API menggunakan JWT (JSON Web Token) untuk autentikasi. Token harus disertakan dalam header:

```http
Authorization: Bearer <your-jwt-token>
```

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data here
  },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": {
    "code": "ERROR_CODE",
    "details": "Detailed error information"
  }
}
```

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
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "jwt-token-here",
    "user": {
      "id": 1,
      "name": "John Doe",
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
  "name": "John Doe",
  "email": "user@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}
```

### Logout
```http
POST /api/auth/logout
```
*Requires Authentication*

## User Management

### Get Current User
```http
GET /api/user/profile
```
*Requires Authentication*

### Update Profile
```http
PUT /api/user/profile
```
*Requires Authentication*

**Request Body:**
```json
{
  "name": "Updated Name",
  "email": "updated@example.com"
}
```

### Get All Users
```http
GET /api/users?page=1&limit=10&search=john
```
*Requires Authentication & Admin Role*

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `search` (optional): Search term

## Database Management Endpoints

### Check Database Connections
```http
GET /api/database/connections
```
*Requires Authentication & Admin Role*

**Response:**
```json
{
  "success": true,
  "data": {
    "connections": [
      {
        "name": "mysql",
        "type": "mysql",
        "status": "connected",
        "stats": {
          "open_connections": 5,
          "in_use": 2,
          "idle": 3
        }
      },
      {
        "name": "postgres",
        "type": "postgres", 
        "status": "connected",
        "stats": {
          "open_connections": 3,
          "in_use": 1,
          "idle": 2
        }
      }
    ]
  }
}
```

### Test Database Connection
```http
GET /api/database/test?connection=mysql
```
*Requires Authentication & Admin Role*

**Query Parameters:**
- `connection`: Connection name (mysql, postgres, mysql_secondary)

### Database Health Check
```http
GET /health/database
```

**Response:**
```json
{
  "status": "healthy",
  "databases": {
    "mysql": {
      "status": "connected",
      "response_time": "2ms"
    },
    "postgres": {
      "status": "connected", 
      "response_time": "3ms"
    }
  }
}
```

### Sync Data Between Databases
```http
POST /api/database/sync
```
*Requires Authentication & Admin Role*

**Request Body:**
```json
{
  "source": "mysql",
  "target": "postgres",
  "tables": ["users", "products"],
  "options": {
    "truncate_target": false,
    "batch_size": 1000
  }
}
```

## Product Management

### Get All Products
```http
GET /api/products?page=1&limit=10&category=electronics
```

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page
- `category` (optional): Filter by category
- `search` (optional): Search term

### Get Product by ID
```http
GET /api/products/{id}
```

### Create Product
```http
POST /api/products
```
*Requires Authentication*

**Request Body:**
```json
{
  "name": "Product Name",
  "description": "Product description",
  "price": 99.99,
  "category_id": 1,
  "stock": 100
}
```

### Update Product
```http
PUT /api/products/{id}
```
*Requires Authentication*

### Delete Product
```http
DELETE /api/products/{id}
```
*Requires Authentication*

## Category Management

### Get All Categories
```http
GET /api/categories
```

### Create Category
```http
POST /api/categories
```
*Requires Authentication*

**Request Body:**
```json
{
  "name": "Category Name",
  "description": "Category description"
}
```

## File Upload

### Upload File
```http
POST /api/upload
```
*Requires Authentication*

**Request:** Multipart form data
- `file`: File to upload
- `type` (optional): File type category

**Response:**
```json
{
  "success": true,
  "data": {
    "filename": "uploaded_file.jpg",
    "url": "/storage/uploads/uploaded_file.jpg",
    "size": 1024000,
    "type": "image/jpeg"
  }
}
```

## Error Codes

| Code | Description |
|------|-------------|
| `VALIDATION_ERROR` | Request validation failed |
| `AUTHENTICATION_REQUIRED` | Authentication token required |
| `INVALID_TOKEN` | JWT token is invalid or expired |
| `INSUFFICIENT_PERMISSIONS` | User lacks required permissions |
| `RESOURCE_NOT_FOUND` | Requested resource not found |
| `DATABASE_ERROR` | Database operation failed |
| `INTERNAL_ERROR` | Internal server error |

## Rate Limiting

API endpoints are rate limited:
- **Public endpoints**: 100 requests per minute
- **Authenticated endpoints**: 1000 requests per minute
- **Admin endpoints**: 5000 requests per minute

Rate limit headers are included in responses:
```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

## Swagger Documentation

Interactive API documentation is available at:
```
http://localhost:8080/swagger/index.html
```

## SDK & Client Libraries

### JavaScript/TypeScript
```javascript
import { ApiClient } from 'golang-starter-kit-client';

const client = new ApiClient({
  baseURL: 'http://localhost:8080',
  token: 'your-jwt-token'
});

// Login
const loginResponse = await client.auth.login({
  email: 'user@example.com',
  password: 'password123'
});

// Get products
const products = await client.products.getAll({
  page: 1,
  limit: 10
});
```

### cURL Examples

#### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### Get Products with Authentication
```bash
curl -X GET http://localhost:8080/api/products \
  -H "Authorization: Bearer your-jwt-token"
```

#### Create Product
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Product",
    "description": "Product description",
    "price": 99.99,
    "category_id": 1,
    "stock": 100
  }'
```

## Testing

### Postman Collection
Download the Postman collection: [API Collection](../examples/postman_collection.json)

### Automated Tests
```bash
# Run API tests
go test ./tests/api/...

# Run integration tests
go test ./tests/integration/...
```

---

Next: [Error Handling](error-handling.md) | [Response Format](response-format.md)
