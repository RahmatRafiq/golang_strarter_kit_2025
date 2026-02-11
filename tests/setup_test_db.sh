#!/bin/bash

# Setup Test Database Script
# Similar to Laravel's database setup for testing

set -e  # Exit on error

echo "========================================="
echo "Setting up test database..."
echo "========================================="

# Load .env.testing
if [ ! -f .env.testing ]; then
    echo "❌ Error: .env.testing not found!"
    echo "Please copy .env.example to .env.testing and configure test database settings"
    exit 1
fi

# Parse .env.testing - extract specific variables
MYSQL_HOST=$(grep "^MYSQL_HOST=" .env.testing | cut -d '=' -f2)
MYSQL_PORT=$(grep "^MYSQL_PORT=" .env.testing | cut -d '=' -f2)
MYSQL_DB=$(grep "^MYSQL_DB=" .env.testing | cut -d '=' -f2)
MYSQL_USER=$(grep "^MYSQL_USER=" .env.testing | cut -d '=' -f2)
MYSQL_PASSWORD=$(grep "^MYSQL_PASSWORD=" .env.testing | cut -d '=' -f2)

echo "Database Configuration:"
echo "  Host: $MYSQL_HOST:$MYSQL_PORT"
echo "  Database: $MYSQL_DB"
echo "  User: $MYSQL_USER"
echo ""

# Check MySQL connection
echo "Checking MySQL connection..."
if ! mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "SELECT 1" > /dev/null 2>&1; then
    echo "❌ Error: Cannot connect to MySQL!"
    echo ""
    echo "Please check:"
    echo "  1. MySQL is running: sudo systemctl status mysql"
    echo "  2. Credentials in .env.testing are correct"
    echo "  3. User has permission: mysql -u root -p"
    echo ""
    exit 1
fi

echo "✓ MySQL connection successful"
echo ""

# Create test database
echo "Creating test database: $MYSQL_DB"
mysql -h"$MYSQL_HOST" -P"$MYSQL_PORT" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" <<-EOSQL
    DROP DATABASE IF EXISTS $MYSQL_DB;
    CREATE DATABASE $MYSQL_DB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    USE $MYSQL_DB;
    SHOW TABLES;
EOSQL

echo "✓ Test database created successfully"
echo ""

echo "========================================="
echo "✅ Test database setup complete!"
echo "========================================="
echo ""
echo "You can now run integration tests:"
echo "  go test ./tests/integration/... -v"
echo ""
