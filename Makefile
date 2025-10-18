.PHONY: help build run test lint lint-fix clean

# Default target
help:
	@echo "Available targets:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make test      - Run tests"
	@echo "  make lint      - Run linter"
	@echo "  make lint-fix  - Run linter with auto-fix"
	@echo "  make clean     - Clean build artifacts"

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

# Run linter (requires golangci-lint to be installed)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Error: golangci-lint not installed. Run: brew install golangci-lint" && exit 1)
	golangci-lint run

# Run linter with auto-fix
lint-fix:
	@echo "Running linter with auto-fix..."
	golangci-lint run --fix ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out
	go clean