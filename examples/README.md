# Examples

This directory contains example code demonstrating how to use various features of the Golang Starter Kit.

## Important Notes

- **Not included in production builds**: All examples use the `//go:build examples` build tag, which excludes them from standard builds
- **For reference only**: These examples are meant for learning and testing purposes
- **Requires proper setup**: Examples need database connections and environment configuration

## Available Examples

### Multi-Database Usage (`multi_database_usage.go`)

Demonstrates how to work with multiple database connections (MySQL, PostgreSQL, MySQL Secondary).

**Features shown:**
- Database service usage
- Executing queries on different databases
- Connection health checks
- Connection statistics

**How to run:**
```bash
# From project root
go run -tags examples examples/multi_database_usage.go
```

**Prerequisites:**
- Valid `.env` file with database configurations
- Running database instances (MySQL, PostgreSQL)
- Proper database credentials

## Build Tag Explanation

The `//go:build examples` tag ensures:
- ✅ Examples are excluded from `go build` and `go install`
- ✅ Reduced production binary size
- ✅ No accidental inclusion in production
- ✅ Examples can still be run explicitly with `-tags examples`

## Adding New Examples

When adding new examples:

1. Add `//go:build examples` as the first line
2. Use package `main` with a `main()` function
3. Use structured logging (zerolog) instead of fmt.Print
4. Document prerequisites and usage in this README
5. Follow existing code style and patterns

## Related Documentation

- [Main README](../README.md)
- [Database Configuration](../config/README.md)
- [Services Documentation](../app/services/README.md)
