# Project variables
BINARY_NAME=project-security
GO=go
GOLINT=golangci-lint
GOFILES=$(shell find . -name "*.go" -type f)

# Colors
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m

.PHONY: all build clean test coverage lint fmt run help install-deps dev

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' Makefile | sed 's/## /  /'

## all: Run lint, test, and build
all: lint test build

## build: Build the application
build:
	@echo "${GREEN}Building...${NC}"
	$(GO) build -v -o $(BINARY_NAME) ./cmd

## run: Run the application
run:
	@echo "${GREEN}Running application...${NC}"
	$(GO) run ./cmd

## dev: Run in development mode with hot reload (requires air)
dev:
	@echo "${GREEN}Running in dev mode...${NC}"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "${RED}Air is not installed. Install with: go install github.com/cosmtrek/air@latest${NC}"; \
		exit 1; \
	fi

## clean: Remove binary and test artifacts
clean:
	@echo "${YELLOW}Cleaning...${NC}"
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

## test: Run tests
test:
	@echo "${GREEN}Running tests...${NC}"
	$(GO) test -v -race ./...

## test-unit: Run unit tests
test-unit:
	@echo "${GREEN}Running unit tests...${NC}"
	$(GO) test -v -race -short ./...

## test-integration: Run integration tests
test-integration:
	@echo "${GREEN}Running integration tests...${NC}"
	$(GO) test -v -race -tags=integration ./tests/integration/...

## coverage: Run tests with coverage
coverage:
	@echo "${GREEN}Running tests with coverage...${NC}"
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}Coverage report: coverage.html${NC}"

## bench: Run benchmarks
bench:
	@echo "${GREEN}Running benchmarks...${NC}"
	$(GO) test -bench=. -benchmem ./...

## fmt: Format code
fmt:
	@echo "${GREEN}Formatting code...${NC}"
	$(GO) fmt ./...
	gofmt -s -w $(GOFILES)

## lint: Run linters
lint:
	@echo "${GREEN}Running linters...${NC}"
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run; \
	else \
		echo "${RED}golangci-lint not installed. Run: make install-tools${NC}"; \
		exit 1; \
	fi

## lint-fix: Fix linting issues
lint-fix:
	@echo "${GREEN}Fixing lint issues...${NC}"
	$(GOLINT) run --fix

## vet: Run go vet
vet:
	@echo "${GREEN}Running go vet...${NC}"
	$(GO) vet ./...

## install-deps: Install dependencies
install-deps:
	@echo "${GREEN}Installing dependencies...${NC}"
	$(GO) mod download
	$(GO) mod tidy

## install-tools: Install development tools
install-tools:
	@echo "${GREEN}Installing tools...${NC}"
	# Install golangci-lint
	@if ! command -v $(GOLINT) > /dev/null; then \
		echo "${YELLOW}Installing golangci-lint...${NC}"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	fi
	# Install air for hot reload
	@if ! command -v air > /dev/null; then \
		echo "${YELLOW}Installing air...${NC}"; \
		go install github.com/cosmtrek/air@latest; \
	fi
	# Install migrate for database migrations
	@if ! command -v migrate > /dev/null; then \
		echo "${YELLOW}Installing migrate...${NC}"; \
		go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@echo "${GREEN}Tools installed!${NC}"

## migrate-up: Run database migrations up
migrate-up:
	@echo "${GREEN}Running migrations...${NC}"
	migrate -path ./migrations -database "sqlite3://./security.db" up

## migrate-down: Run database migrations down
migrate-down:
	@echo "${YELLOW}Rolling back migrations...${NC}"
	migrate -path ./migrations -database "sqlite3://./security.db" down

## migrate-create: Create a new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

## seed: Seed the database
seed:
	@echo "${GREEN}Seeding database...${NC}"
	$(GO) run ./scripts/seed.go

## docker-build: Build Docker image
docker-build:
	@echo "${GREEN}Building Docker image...${NC}"
	docker build -t $(BINARY_NAME):latest .

## docker-run: Run Docker container
docker-run:
	@echo "${GREEN}Running Docker container...${NC}"
	docker run -p 8080:8080 \
		-e ENVIRONMENT=production \
		-v $$(pwd)/config:/app/config \
		$(BINARY_NAME):latest

## docker-compose-up: Start services with docker-compose
docker-compose-up:
	@echo "${GREEN}Starting services...${NC}"
	docker-compose up -d

## docker-compose-down: Stop services
docker-compose-down:
	@echo "${YELLOW}Stopping services...${NC}"
	docker-compose down

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "${GREEN}All checks passed!${NC}"

## ci: Run CI pipeline
ci: install-deps check
	@echo "${GREEN}CI pipeline passed!${NC}"