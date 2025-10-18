.PHONY: help build run test lint lint-fix test-performance clean swagger docker-build docker-run docker-stop docker-clean

# Default target
help:
	@echo "Available targets:"
	@echo "  make build            - Build the application"
	@echo "  make run              - Run the application"
	@echo "  make test             - Run unit tests"
	@echo "  make lint             - Run linter"
	@echo "  make lint-fix         - Run linter with auto-fix"
	@echo "  make swagger          - Generate Swagger documentation"
	@echo "  make test-performance - Run K6 performance tests"
	@echo "  make clean            - Clean build artifacts"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-run       - Run with docker-compose"
	@echo "  make docker-stop      - Stop docker-compose"
	@echo "  make docker-clean     - Remove Docker image and containers"

# Build the application
build:
	@echo "Building..."
	go build -o bin/server cmd/main.go

# Run the application
run:
	@echo "Running server..."
	go run cmd/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "\nCoverage:"
	go tool cover -func=coverage.out

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null || (echo "Error: swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest" && exit 1)
	swag init -g cmd/main.go --output ./docs
	@echo "Swagger docs generated at ./docs"

# Run linter (requires golangci-lint to be installed)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Error: golangci-lint not installed. Run: brew install golangci-lint" && exit 1)
	golangci-lint run

# Run linter with auto-fix (requires golangci-lint to be installed)
lint-fix:
	@echo "Running linter with auto-fix..."
	@which golangci-lint > /dev/null || (echo "Error: golangci-lint not installed. Run: brew install golangci-lint" && exit 1)
	golangci-lint run --fix

# Run K6 performance tests
test-performance:
	@echo "Running K6 performance tests..."
	@which k6 > /dev/null || (echo "Error: k6 not installed. Visit https://k6.io/docs/get-started/installation/" && exit 1)
	@echo "Make sure the server is running on http://localhost:8080"
	@echo "Starting tests in 3 seconds..."
	@sleep 3
	k6 run test-k6.js

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf docs/
	rm -f coverage.out
	go clean

# Docker: Build image
docker-build:
	@echo "Building Docker image..."
	docker build --load -t ip-location-api:latest .
	@echo "✅ Image built successfully!"
	@docker images ip-location-api:latest

# Docker: Run with docker-compose
docker-run:
	@echo "Starting containers with docker-compose..."
	docker compose up -d
	@echo "✅ Server running at http://localhost:8080"
	@echo "   Health check: http://localhost:8080/health"
	@echo "   Swagger UI: http://localhost:8080/swagger/index.html"

# Docker: Stop containers
docker-stop:
	@echo "Stopping containers..."
	docker compose down

# Docker: Clean up
docker-clean:
	@echo "Removing containers and images..."
	docker compose down --rmi all --volumes --remove-orphans
	@echo "Cleanup complete!"