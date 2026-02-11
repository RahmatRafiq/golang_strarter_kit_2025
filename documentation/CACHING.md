# Caching Strategy Documentation

## Overview

This application implements a comprehensive multi-layer caching strategy using Redis as the caching backend. The caching system significantly improves performance by reducing database queries and response times.

## Table of Contents

- [Architecture](#architecture)
- [Cache Layers](#cache-layers)
- [Configuration](#configuration)
- [Cache Service](#cache-service)
- [Service Layer Caching](#service-layer-caching)
- [HTTP Response Caching](#http-response-caching)
- [Cache Invalidation](#cache-invalidation)
- [Best Practices](#best-practices)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                   CLIENT REQUEST                     │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│              HTTP CACHE MIDDLEWARE                   │
│  (Caches entire HTTP responses for GET requests)    │
└──────────────────────┬──────────────────────────────┘
                       │
                    Miss│
                       ▼
┌─────────────────────────────────────────────────────┐
│                  CONTROLLER                         │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│              SERVICE LAYER CACHE                    │
│  (Caches database query results: users, roles)      │
└──────────────────────┬──────────────────────────────┘
                       │
                    Miss│
                       ▼
┌─────────────────────────────────────────────────────┐
│                   DATABASE                          │
│              (PostgreSQL/MySQL)                     │
└─────────────────────────────────────────────────────┘
```

---

## Cache Layers

### 1. HTTP Response Cache (Layer 1)

**Purpose:** Cache entire HTTP responses to avoid controller and service execution

**TTL:** Configurable per endpoint (default: 5 minutes)

**Use Cases:**
- Public API endpoints
- Read-heavy endpoints
- List endpoints with pagination

**Implementation:**
```go
// Cache all GET requests for 5 minutes
router.Use(middleware.HTTPCacheMiddleware(5 * time.Minute))

// Or cache specific endpoints
router.Use(middleware.CachableEndpoints(
    5 * time.Minute,
    "/api/v1/users",
    "/api/v1/roles",
))
```

### 2. Service Layer Cache (Layer 2)

**Purpose:** Cache database query results at the service level

**TTL:**
- Short: 5 minutes (frequently changing data)
- Medium: 30 minutes (moderate data)
- Long: 2 hours (stable data)
- Day: 24 hours (very stable data)

**Use Cases:**
- User data: `UserService.FindByID()`
- Role data: `RoleService.FindByID()`
- Permission data: `RoleService.GetPermissions()`

**Implementation:**
```go
// Automatic caching in service layer
user, err := userService.FindByID(1)
// First call: Database query
// Subsequent calls: Cache hit (30 min TTL)
```

---

## Configuration

### Environment Variables

Add to `.env` file:

```bash
# Redis Configuration (Required for caching)
REDIS_ADDR=localhost:6379      # Redis server address
REDIS_PASSWORD=                # Redis password (optional)
REDIS_DB=0                     # Redis database number

# Cache Configuration (Optional)
CACHE_ENABLED=true             # Enable/disable caching
CACHE_DEFAULT_TTL=1800         # Default TTL in seconds (30 minutes)
```

### Redis Connection Settings

```go
// app/helpers/redis.go
DialTimeout:  5 * time.Second  // Connection timeout
ReadTimeout:  3 * time.Second  # Read operation timeout
WriteTimeout: 3 * time.Second  # Write operation timeout
PoolSize:     10               # Connection pool size
MinIdleConns: 5                # Minimum idle connections
```

---

## Cache Service

### Core Methods

#### Get

```go
var user models.User
err := cacheService.Get("user:1", &user)
if err != nil {
    // Cache miss or error
}
```

#### Set

```go
user := models.User{ID: 1, Email: "user@example.com"}
err := cacheService.Set("user:1", user, services.CacheTTLMedium)
```

#### Delete

```go
err := cacheService.Delete("user:1")
```

#### Delete Pattern

```go
// Delete all user caches
err := cacheService.DeletePattern("user:*")
```

#### Remember (Cache-Aside Pattern)

```go
result, err := cacheService.Remember(
    "expensive_query",
    30 * time.Minute,
    func() (interface{}, error) {
        // Expensive database query here
        return db.Query("SELECT * FROM ..."), nil
    },
)
```

#### Flush (Clear All Cache)

```go
err := cacheService.Flush()
```

### Helper Methods

```go
// Check if key exists
exists := cacheService.Exists("user:1")

// Get TTL
ttl, err := cacheService.GetTTL("user:1")

// Increment counter
count, err := cacheService.Increment("page_views", 1)

// Set only if not exists (distributed locking)
success, err := cacheService.SetNX("lock:user:1", "locked", 10*time.Second)
```

---

## Service Layer Caching

### User Service

**Cached Operations:**

```go
// FindByID - TTL: 30 minutes
user, err := userService.FindByID(1)
// Cache Key: "user:1"

// GetRoles - TTL: 30 minutes
roles, err := userService.GetRoles(1)
// Cache Key: "user:1:roles"
```

**Cache Invalidation:**

```go
// Automatic invalidation on:
userService.Update(user)      // Invalidates: user:1, user:1:roles
userService.DeleteByID(1)     // Invalidates: user:1, user:1:roles
userService.AssignRoles(1, roleIDs) // Invalidates: user:1:roles
```

### Role Service

**Cached Operations:**

```go
// FindByID - TTL: 2 hours (stable data)
role, err := roleService.FindByID(1)
// Cache Key: "role:1"

// GetPermissions - TTL: 2 hours
permissions, err := roleService.GetPermissions(1)
// Cache Key: "role:1:permissions"
```

**Cache Invalidation:**

```go
// Automatic invalidation on:
roleService.Update(role)      // Invalidates: role:1, role:1:permissions
roleService.DeleteByID(1)     // Invalidates: role:1, role:1:permissions
roleService.AssignPermissions(1, permIDs) // Invalidates: role:1:permissions, user:*:roles
```

---

## HTTP Response Caching

### Basic Usage

```go
// In routes/api_v1.go
router.GET("/users", middleware.HTTPCacheMiddleware(5*time.Minute), userController.List)
```

### Cache-Control Headers

Cached responses include:

```http
X-Cache: HIT          # Response served from cache
X-Cache: MISS         # Response not in cache
```

### Cache Key Generation

Cache keys are generated from:
- HTTP Method (GET only)
- Request Path
- Query Parameters
- User ID (if authenticated)

Example:
```
GET /api/v1/users?page=1&limit=10 (user_id: 5)
Cache Key: http_cache:a3f5d8e... (SHA256 hash)
```

### Invalidation

```go
// Invalidate specific pattern
middleware.InvalidateHTTPCache("/users")
// Clears all caches matching: http_cache:*users*
```

---

## Cache Invalidation

### Automatic Invalidation

The system automatically invalidates cache on data mutations:

| Operation | Invalidated Keys |
|-----------|------------------|
| `User Update` | `user:{id}`, `user:{id}:roles` |
| `User Delete` | `user:{id}`, `user:{id}:roles` |
| `Role Update` | `role:{id}`, `role:{id}:permissions`, `user:*:roles` |
| `Role Delete` | `role:{id}`, `role:{id}:permissions`, `user:*:roles` |
| `Assign Roles` | `user:{id}:roles` |
| `Assign Permissions` | `role:{id}:permissions`, `user:*:roles` |

### Manual Invalidation

```go
// Invalidate specific user
cacheService.InvalidateUserCache(userID)

// Invalidate specific role
cacheService.InvalidateRoleCache(roleID)

// Invalidate pattern
cacheService.DeletePattern("user:*")

// Clear all cache
cacheService.Flush()
```

### Cache Stampede Prevention

Use `SetNX` for distributed locking:

```go
// Prevent multiple simultaneous cache rebuilds
locked, err := cacheService.SetNX("rebuild:users", "locked", 10*time.Second)
if locked {
    defer cacheService.Delete("rebuild:users")
    // Rebuild cache here
}
```

---

## Best Practices

### 1. Choose Appropriate TTL

```go
// Frequently changing data
CacheTTLShort = 5 * time.Minute

// Moderate change frequency
CacheTTLMedium = 30 * time.Minute

// Stable data
CacheTTLLong = 2 * time.Hour

// Very stable configuration data
CacheTTLDay = 24 * time.Hour
```

### 2. Cache Key Naming Convention

Follow consistent naming:
```
{entity}:{id}                    # user:1
{entity}:{id}:{relation}         # user:1:roles
{entity}:{id}:{relation}:{id}    # user:1:permission:5
```

### 3. Graceful Degradation

Cache failures should not break the application:

```go
var user models.User
if err := cache.Get(key, &user); err != nil {
    // Cache miss/error, query database
    user, err = db.FindByID(id)
}
```

### 4. Avoid Caching Sensitive Data

**Don't cache:**
- Passwords (even hashed)
- Payment information
- Personal identification numbers
- Temporary tokens

**Safe to cache:**
- User profiles (without sensitive fields)
- Public roles and permissions
- Product catalogs
- Static content

### 5. Monitor Cache Hit Rates

```go
// Log cache hits/misses
log.Debug().
    Str("cache_key", key).
    Bool("hit", hit).
    Msg("Cache access")
```

---

## Monitoring

### Metrics to Track

1. **Hit Rate**
   ```
   Hit Rate = Cache Hits / (Cache Hits + Cache Misses) * 100%
   Target: > 80%
   ```

2. **Response Time Improvement**
   ```
   Improvement = (DB Query Time - Cache Hit Time) / DB Query Time * 100%
   Expected: 80-95%
   ```

3. **Memory Usage**
   ```bash
   redis-cli INFO memory
   ```

### Redis CLI Commands

```bash
# Monitor cache operations in real-time
redis-cli MONITOR

# Check memory usage
redis-cli INFO memory

# Count keys
redis-cli DBSIZE

# View specific key
redis-cli GET user:1

# View all keys matching pattern
redis-cli KEYS "user:*"

# Check TTL
redis-cli TTL user:1

# Clear all cache (DANGER!)
redis-cli FLUSHDB
```

### Performance Benchmarks

Expected performance improvements:

| Operation | Without Cache | With Cache | Improvement |
|-----------|--------------|------------|-------------|
| User FindByID | ~15-20ms | ~1-2ms | **90%** |
| User GetRoles | ~25-30ms | ~1-2ms | **93%** |
| Role GetPermissions | ~20-25ms | ~1-2ms | **92%** |
| HTTP Response | ~50-100ms | ~2-5ms | **95%** |

---

## Troubleshooting

### Cache Not Working

**Symptom:** Always cache miss

**Solutions:**
1. Check Redis connection:
   ```bash
   redis-cli ping
   # Expected: PONG
   ```

2. Check environment variables:
   ```bash
   echo $REDIS_ADDR
   echo $CACHE_ENABLED
   ```

3. Check logs:
   ```bash
   grep "cache" application.log
   ```

### High Memory Usage

**Symptom:** Redis using too much memory

**Solutions:**
1. Reduce TTL values
2. Implement cache eviction policy:
   ```bash
   redis-cli CONFIG SET maxmemory-policy allkeys-lru
   redis-cli CONFIG SET maxmemory 256mb
   ```

3. Clear old cache:
   ```bash
   redis-cli FLUSHDB
   ```

### Stale Cache

**Symptom:** Seeing old data after updates

**Solutions:**
1. Check cache invalidation logic
2. Reduce TTL for frequently changing data
3. Manual flush:
   ```go
   cacheService.InvalidateUserCache(userID)
   ```

### Redis Connection Failures

**Symptom:** "Failed to connect to Redis"

**Solutions:**
1. Ensure Redis is running:
   ```bash
   redis-cli ping
   # or
   docker ps | grep redis
   ```

2. Check firewall:
   ```bash
   telnet localhost 6379
   ```

3. Verify connection settings in `.env`

4. The application gracefully degrades:
   ```
   [WARN] Redis not available, caching disabled
   ```

---

## Advanced Topics

### Cache Warming

Pre-populate cache with frequently accessed data:

```go
func WarmCache() {
    // Cache top 100 users
    users, _ := userRepo.List(1, 100)
    for _, user := range users {
        _ = cache.Set(UserCacheKey(user.ID), user, CacheTTLLong)
    }
}
```

### Distributed Caching

For multi-instance deployments, Redis acts as a shared cache:

```
Instance 1 ──┐
Instance 2 ──┼──▶ Redis Cluster
Instance 3 ──┘
```

### Cache Partitioning

Split cache by data type:

```go
// Use different Redis databases
REDIS_DB=0  // User data
REDIS_DB=1  // Session data
REDIS_DB=2  // Queue data
```

---

## Migration Guide

### Enabling Cache on Existing Project

1. **Install Redis**
   ```bash
   # Docker
   docker run -d -p 6379:6379 redis:7-alpine

   # Or docker-compose (already included)
   docker-compose up -d redis
   ```

2. **Update Environment**
   ```bash
   # Add to .env
   REDIS_ADDR=localhost:6379
   REDIS_PASSWORD=
   REDIS_DB=0
   ```

3. **Restart Application**
   ```bash
   # Caching is automatically enabled
   go run main.go
   ```

4. **Verify**
   ```bash
   # Check logs
   grep "Redis connected" application.log

   # Test cache
   redis-cli KEYS "*"
   ```

### Disabling Cache

```bash
# Option 1: Stop Redis
docker-compose stop redis

# Option 2: Set environment
CACHE_ENABLED=false

# Application will gracefully degrade
```

---

## Performance Tips

1. **Use appropriate TTL** - Balance freshness vs. performance
2. **Cache high-traffic endpoints** - Focus on 80/20 rule
3. **Implement cache warming** - Pre-load frequently accessed data
4. **Monitor hit rates** - Adjust strategy based on metrics
5. **Use Redis clustering** - For high-traffic production
6. **Implement circuit breaker** - Handle Redis failures gracefully
7. **Regular maintenance** - Monitor memory and evict old keys

---

## References

- [Redis Documentation](https://redis.io/docs/)
- [go-redis Library](https://github.com/redis/go-redis)
- [Caching Best Practices](https://aws.amazon.com/caching/best-practices/)
- [Cache Stampede Prevention](https://en.wikipedia.org/wiki/Cache_stampede)

---

**Last Updated:** 2026-02-11
**Version:** 1.0.0
