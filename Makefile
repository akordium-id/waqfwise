.PHONY: help build-community build-enterprise run-community run-enterprise test clean docker-up docker-down migrate-up migrate-down deps

# Variables
COMMUNITY_BINARY=bin/waqfwise-community
ENTERPRISE_BINARY=bin/waqfwise-enterprise
MIGRATION_PATH=migrations/community
DATABASE_URL=postgres://waqfwise:waqfwise_dev_password@localhost:5432/waqfwise?sslmode=disable

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

deps: ## Download Go dependencies
	go mod download
	go mod tidy
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

build-community: ## Build community edition
	@echo "Building WaqfWise Community Edition..."
	@mkdir -p bin
	go build -tags community -o $(COMMUNITY_BINARY) ./cmd/waqfwise-community

build-enterprise: ## Build enterprise edition
	@echo "Building WaqfWise Enterprise Edition..."
	@mkdir -p bin
	go build -tags enterprise -o $(ENTERPRISE_BINARY) ./cmd/waqfwise-enterprise

build: build-community build-enterprise ## Build both editions

run-community: build-community ## Run community edition
	@echo "Starting WaqfWise Community Edition..."
	./$(COMMUNITY_BINARY)

run-enterprise: build-enterprise ## Run enterprise edition
	@echo "Starting WaqfWise Enterprise Edition..."
	./$(ENTERPRISE_BINARY)

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-up: ## Start Docker services
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5

docker-down: ## Stop Docker services
	@echo "Stopping Docker services..."
	docker-compose down

docker-clean: ## Stop and remove Docker volumes
	@echo "Stopping Docker services and removing volumes..."
	docker-compose down -v

migrate-up: ## Run database migrations
	@echo "Running database migrations..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required. Usage: make migrate-force version=1"; \
		exit 1; \
	fi
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" force $(version)

dev-setup: docker-up migrate-up ## Setup development environment
	@echo "Development environment ready!"

dev: docker-up ## Start development environment
	@echo "Starting development environment..."
	$(MAKE) run-community

lint: ## Run linters
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

vet: ## Run go vet
	go vet ./...

check: fmt vet lint test ## Run all checks

.DEFAULT_GOAL := help
