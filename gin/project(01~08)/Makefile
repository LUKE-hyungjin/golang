.PHONY: run build clean test help install

# Default target
help:
	@echo "Available commands:"
	@echo "  make install   - Install dependencies"
	@echo "  make run       - Run the application"
	@echo "  make build     - Build the application"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run the application
run:
	@echo "Starting Blog Platform..."
	go run cmd/main.go

# Build the application
build:
	@echo "Building application..."
	go build -o bin/blog-platform cmd/main.go
	@echo "Build complete! Binary: bin/blog-platform"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf uploads/*
	@echo "Clean complete!"

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development mode..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air
