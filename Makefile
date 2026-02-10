.PHONY: help build test test-coverage lint fmt clean docker-build docker-run docker-stop swagger install-tools

# Default target
help:
	@echo "Available commands:"
	@echo "  make build           - Build the application"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage"
	@echo "  make lint            - Run linter"
	@echo "  make fmt             - Format code"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make docker-run      - Run Docker container"
	@echo "  make docker-stop     - Stop Docker container"
	@echo "  make swagger         - Generate Swagger docs"
	@echo "  make install-tools   - Install development tools"

# Build binary
build:
	@echo "Building application..."
	@go build -o main .
	@echo "Build complete: ./main"

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run --timeout=5m

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f main coverage.out coverage.html test-report.json
	@rm -rf tmp/
	@echo "Clean complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t golang-starter-kit-2025:latest .
	@echo "Docker image built: golang-starter-kit-2025:latest"

# Run Docker container
docker-run:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Containers started. Check with: docker-compose ps"

# Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped"

# Generate Swagger docs
swagger:
	@echo "Generating Swagger documentation..."
	@swag init --pd
	@echo "Swagger docs generated"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/onsi/ginkgo/v2/ginkgo@latest
	@echo "Tools installed successfully"

# Run application
run:
	@echo "Running application..."
	@go run main.go

# Database migrations
migrate-up:
	@echo "Running migrations..."
	@go run main.go migrate:all

migrate-down:
	@echo "Rolling back migrations..."
	@go run main.go migrate:rollback --all

# CI/CD helpers
ci-test: test-coverage lint
	@echo "CI tests complete"

# Development setup
dev-setup: install-tools
	@echo "Setting up development environment..."
	@cp -n .env.example .env || true
	@echo "Development setup complete. Please edit .env file."