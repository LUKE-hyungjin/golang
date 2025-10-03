.PHONY: help run build test clean migrate seed

# Variables
APP_NAME=banking-system
MAIN_PATH=cmd/main.go
BINARY_NAME=banking

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Run the application in development mode
	APP_ENV=development go run $(MAIN_PATH)

run-prod: ## Run the application in production mode
	APP_ENV=production go run $(MAIN_PATH)

build: ## Build the binary
	go build -o $(BINARY_NAME) $(MAIN_PATH)

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -cover ./...

clean: ## Clean build artifacts and logs
	rm -f $(BINARY_NAME)
	rm -f *.db
	rm -rf logs/

deps: ## Install dependencies
	go mod download
	go mod tidy

migrate: ## Run database migrations
	@echo "Migrations will run automatically on startup"

seed: ## Seed the database
	@echo "Seeding will run automatically in development mode"

docker-build: ## Build Docker image
	docker build -t $(APP_NAME):latest .

docker-run: ## Run in Docker
	docker run -p 8080:8080 $(APP_NAME):latest

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .