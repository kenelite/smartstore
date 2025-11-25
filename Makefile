# SmartStore Makefile
.PHONY: help build run test clean fmt lint vet deps tidy install-tools db-up db-down db-migrate docker-build docker-run

# Variables
APP_NAME := smartstore-gateway
BINARY_DIR := bin
BINARY := $(BINARY_DIR)/$(APP_NAME)
CMD_DIR := ./cmd/gateway
CONFIG_FILE := config.yaml

# Go related variables
GO := go
GOBIN := $(shell $(GO) env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell $(GO) env GOPATH)/bin
endif

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "\
	-X main.Version=$(VERSION) \
	-X main.Commit=$(COMMIT) \
	-X main.BuildTime=$(BUILD_TIME)"

# Database variables (override these via environment or command line)
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= smartstore
DB_PASSWORD ?= smartstore
DB_NAME ?= smartstore
DB_URL := postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Docker variables
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG ?= latest

# Color output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

##@ General

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(YELLOW)<target>$(NC)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(BLUE)%-20s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(GREEN)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

deps: ## Download Go dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@$(GO) mod download
	@$(GO) mod verify

tidy: ## Tidy Go modules
	@echo "$(BLUE)Tidying Go modules...$(NC)"
	@$(GO) mod tidy

fmt: ## Format Go code
	@echo "$(BLUE)Formatting code...$(NC)"
	@$(GO) fmt ./...

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@$(GO) vet ./...

lint: install-tools ## Run golangci-lint
	@echo "$(BLUE)Running golangci-lint...$(NC)"
	@if [ -f $(GOBIN)/golangci-lint ]; then \
		$(GOBIN)/golangci-lint run --timeout=5m ./...; \
	else \
		echo "$(YELLOW)golangci-lint not found, skipping...$(NC)"; \
	fi

##@ Build

build: ## Build the application
	@echo "$(BLUE)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@$(GO) build $(LDFLAGS) -o $(BINARY) $(CMD_DIR)
	@echo "$(GREEN)✓ Binary created: $(BINARY)$(NC)"

build-linux: ## Build for Linux
	@echo "$(BLUE)Building for Linux...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY)-linux-amd64 $(CMD_DIR)
	@echo "$(GREEN)✓ Linux amd64 binary created$(NC)"

build-darwin: ## Build for macOS
	@echo "$(BLUE)Building for macOS...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY)-darwin-amd64 $(CMD_DIR)
	@GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BINARY)-darwin-arm64 $(CMD_DIR)
	@echo "$(GREEN)✓ macOS binaries created$(NC)"

build-windows: ## Build for Windows
	@echo "$(BLUE)Building for Windows...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY)-windows-amd64.exe $(CMD_DIR)
	@echo "$(GREEN)✓ Windows binary created$(NC)"

build-all: ## Build for all platforms (Linux, macOS, Windows)
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@echo "$(YELLOW)Building Linux amd64...$(NC)"
	@GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY)-linux-amd64 $(CMD_DIR)
	@echo "$(YELLOW)Building macOS amd64...$(NC)"
	@GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY)-darwin-amd64 $(CMD_DIR)
	@echo "$(YELLOW)Building macOS arm64...$(NC)"
	@GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BINARY)-darwin-arm64 $(CMD_DIR)
	@echo "$(YELLOW)Building Windows amd64...$(NC)"
	@GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BINARY)-windows-amd64.exe $(CMD_DIR)
	@echo ""
	@echo "$(GREEN)✓ All platform binaries created successfully!$(NC)"
	@echo "$(BLUE)Built binaries:$(NC)"
	@ls -lh $(BINARY_DIR)/$(APP_NAME)-* | awk '{printf "  - %-45s %8s\n", $$9, $$5}'

install: build ## Install the binary to $GOPATH/bin
	@echo "$(BLUE)Installing $(APP_NAME)...$(NC)"
	@cp $(BINARY) $(GOBIN)/
	@echo "$(GREEN)✓ Installed to $(GOBIN)/$(APP_NAME)$(NC)"

##@ Run

run: ## Run the application
	@echo "$(BLUE)Running $(APP_NAME)...$(NC)"
	@$(GO) run $(CMD_DIR)/main.go

run-bin: build ## Build and run the binary
	@echo "$(BLUE)Running binary...$(NC)"
	@$(BINARY)

dev: ## Run with air for hot reload
	@if [ -f $(GOBIN)/air ]; then \
		$(GOBIN)/air; \
	else \
		echo "$(YELLOW)air not found. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(BLUE)Running normally...$(NC)"; \
		$(GO) run $(CMD_DIR)/main.go; \
	fi

##@ Test

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@$(GO) test -v -race -coverprofile=coverage.out ./...

test-short: ## Run short tests
	@echo "$(BLUE)Running short tests...$(NC)"
	@$(GO) test -short -v ./...

test-coverage: test ## Run tests with coverage report
	@echo "$(BLUE)Generating coverage report...$(NC)"
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

test-bench: ## Run benchmark tests
	@echo "$(BLUE)Running benchmarks...$(NC)"
	@$(GO) test -bench=. -benchmem ./...

##@ Database

db-create: ## Create database
	@echo "$(BLUE)Creating database...$(NC)"
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);" 2>/dev/null || echo "$(YELLOW)Database may already exist$(NC)"

db-drop: ## Drop database
	@echo "$(RED)Dropping database...$(NC)"
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"

db-migrate: ## Run database migrations
	@echo "$(BLUE)Running database migrations...$(NC)"
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f db/schema.sql
	@echo "$(GREEN)✓ Migrations completed$(NC)"

db-reset: db-drop db-create db-migrate ## Reset database (drop, create, migrate)
	@echo "$(GREEN)✓ Database reset complete$(NC)"

db-shell: ## Connect to database shell
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME)

##@ Docker

docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)✓ Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

docker-run: ## Run Docker container
	@echo "$(BLUE)Running Docker container...$(NC)"
	@docker run --rm -p 8080:8080 -v $(PWD)/config.yaml:/app/config.yaml $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up: ## Start services with docker-compose
	@echo "$(BLUE)Starting services with docker-compose...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✓ Services started$(NC)"

docker-compose-down: ## Stop services with docker-compose
	@echo "$(BLUE)Stopping services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

docker-compose-logs: ## Show docker-compose logs
	@docker-compose logs -f

##@ Tools

install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install github.com/cosmtrek/air@latest
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

##@ Scripts

setup: ## Run initial setup script
	@bash scripts/setup.sh

verify: ## Verify project setup
	@bash scripts/verify-setup.sh

test-all: ## Run comprehensive test script
	@bash scripts/test.sh

test-verbose: ## Run tests with verbose output
	@bash scripts/test.sh --verbose

test-with-coverage: ## Run tests with coverage report
	@bash scripts/test.sh --coverage

clean-deep: ## Run deep clean script
	@bash scripts/clean.sh --deep

clean-docker: ## Clean and stop Docker services
	@bash scripts/clean.sh --docker

deploy-prepare: ## Prepare deployment
	@bash scripts/deploy.sh

##@ Clean

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Clean complete$(NC)"

clean-all: clean ## Clean all generated files
	@echo "$(BLUE)Deep cleaning...$(NC)"
	@$(GO) clean -cache -testcache -modcache
	@echo "$(GREEN)✓ Deep clean complete$(NC)"

##@ Release

release-check: fmt vet test ## Run pre-release checks
	@echo "$(GREEN)✓ All checks passed$(NC)"

release-build: release-check build-all ## Build release binaries
	@echo "$(GREEN)✓ Release build complete$(NC)"

##@ Info

info: ## Display project information
	@echo "$(BLUE)Project Information:$(NC)"
	@echo "  App Name:      $(APP_NAME)"
	@echo "  Version:       $(VERSION)"
	@echo "  Commit:        $(COMMIT)"
	@echo "  Build Time:    $(BUILD_TIME)"
	@echo "  Go Version:    $(shell $(GO) version)"
	@echo "  Database URL:  $(DB_URL)"

check-deps: ## Check for outdated dependencies
	@echo "$(BLUE)Checking for outdated dependencies...$(NC)"
	@$(GO) list -u -m all

.DEFAULT_GOAL := help

