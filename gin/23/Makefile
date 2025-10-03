# Makefile for Go project with linting and formatting

# Variables
BINARY_NAME=app
GO=go
GOLINT=golangci-lint
GOFMT=gofmt
GOFILES=$(shell find . -name "*.go" -type f)

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: all build clean test coverage lint fmt vet help install-tools run

# Default target
all: lint test build

## help: Display this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${GREEN}%-15s${NC} %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## build: Build the binary
build:
	@echo "${GREEN}Building...${NC}"
	$(GO) build -v -o $(BINARY_NAME) .
	@echo "${GREEN}Build complete!${NC}"

## run: Run the application
run:
	@echo "${GREEN}Running application...${NC}"
	$(GO) run .

## clean: Remove binary and cached files
clean:
	@echo "${YELLOW}Cleaning...${NC}"
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "${GREEN}Clean complete!${NC}"

## test: Run tests
test:
	@echo "${GREEN}Running tests...${NC}"
	$(GO) test -v -race ./...

## coverage: Run tests with coverage
coverage:
	@echo "${GREEN}Running tests with coverage...${NC}"
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}Coverage report generated: coverage.html${NC}"

## bench: Run benchmarks
bench:
	@echo "${GREEN}Running benchmarks...${NC}"
	$(GO) test -bench=. -benchmem ./...

## fmt: Format code with gofmt
fmt:
	@echo "${GREEN}Formatting code...${NC}"
	@$(GOFMT) -s -w $(GOFILES)
	@echo "${GREEN}Format complete!${NC}"

## fmt-check: Check if code is formatted
fmt-check:
	@echo "${GREEN}Checking code format...${NC}"
	@diff=$$($(GOFMT) -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "${RED}Please run 'make fmt' to format the following files:${NC}"; \
		echo "$$diff"; \
		exit 1; \
	else \
		echo "${GREEN}Code is properly formatted!${NC}"; \
	fi

## vet: Run go vet
vet:
	@echo "${GREEN}Running go vet...${NC}"
	$(GO) vet ./...
	@echo "${GREEN}Vet complete!${NC}"

## lint: Run golangci-lint
lint:
	@echo "${GREEN}Running golangci-lint...${NC}"
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run --timeout 5m; \
	else \
		echo "${RED}golangci-lint is not installed. Run 'make install-tools' to install it.${NC}"; \
		exit 1; \
	fi

## lint-fix: Run golangci-lint with fix
lint-fix:
	@echo "${GREEN}Running golangci-lint with fix...${NC}"
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run --fix --timeout 5m; \
	else \
		echo "${RED}golangci-lint is not installed. Run 'make install-tools' to install it.${NC}"; \
		exit 1; \
	fi

## install-tools: Install development tools
install-tools:
	@echo "${GREEN}Installing development tools...${NC}"
	# Install golangci-lint
	@if ! command -v $(GOLINT) >/dev/null 2>&1; then \
		echo "${YELLOW}Installing golangci-lint...${NC}"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	else \
		echo "${GREEN}golangci-lint is already installed${NC}"; \
	fi
	# Install other tools
	@echo "${YELLOW}Installing other Go tools...${NC}"
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/kisielk/errcheck@latest
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	$(GO) install github.com/gordonklaus/ineffassign@latest
	@echo "${GREEN}Tools installation complete!${NC}"

## mod: Download and tidy Go modules
mod:
	@echo "${GREEN}Downloading dependencies...${NC}"
	$(GO) mod download
	$(GO) mod tidy
	@echo "${GREEN}Dependencies updated!${NC}"

## check: Run all checks (fmt, vet, lint, test)
check: fmt-check vet lint test
	@echo "${GREEN}All checks passed!${NC}"

## ci: Run CI pipeline (used in CI/CD)
ci: mod fmt-check vet lint test
	@echo "${GREEN}CI pipeline complete!${NC}"

# Docker targets (optional)
## docker-build: Build Docker image
docker-build:
	@echo "${GREEN}Building Docker image...${NC}"
	docker build -t $(BINARY_NAME):latest .

## docker-run: Run Docker container
docker-run:
	@echo "${GREEN}Running Docker container...${NC}"
	docker run -p 8080:8080 $(BINARY_NAME):latest