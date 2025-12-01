# Variables
GOOSE_CMD=goose
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=internal/db/migrations
GOOSE_DBSTRING=${DB_SOURCE_LOCAL}
SQLC_CMD=sqlc

# Load .env file
include .env
export $(shell sed 's/=.*//' .env)

# Default target
.DEFAULT_GOAL := help

.PHONY: help build run clean test up down bash db-migrate db-rollback create-migration migration-status gen-sqlc lint-sqlc


help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build             Build the program"
	@echo "  run               Run the program"
	@echo "  clean             Clean build artifacts"
	@echo "  test              Run tests"
	@echo "  up                Start the application with Docker Compose"
	@echo "  down              Stop the application with Docker Compose"
	@echo "  bash              Open bash shell in Docker container"
	@echo "  gen-sqlc          Generate SQLC code"
	@echo "  lint-sqlc         Lint SQLC generated code"
	@echo "  db-migrate        Run database migrations"
	@echo "  db-rollback       Rollback last database migration"
	@echo "  create-migration  Create a migration"
	@echo "  migrate-status    Show migration status"
	@echo "	 gen-docs		 Generate API documentation"

# Targets
build:
	@echo "Building the program..."
	@go build -o bin/main  cmd/api/main.go

run:
	@echo "Running the program..."
	./bin/main.go

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin

test:
	@echo "Running tests..."
	go test -v ./...

# Docker targets
up:
	@echo "Starting the application with Docker Compose..."
	docker compose up

down:
	@echo "Stopping the application with Docker Compose..."
	docker compose down --remove-orphans

bash:
	@echo "Opening a bash shell in the running Docker container..."
	docker compose exec muslim_tech sh

# Database migration targets
db-migrate:
	@echo "Running database migrations..."
	$(GOOSE_CMD) -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

db-rollback:
	@echo "Rolling back the last migration..."
	$(GOOSE_CMD) -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down

create-migration:
	@if [ -z "$(name)" ]; then \
		echo "Provide a migration name: make create-migration name=add_users_table"; \
		exit 1; \
	fi
	@echo "Creating migration..."
	$(GOOSE_CMD) -dir $(GOOSE_MIGRATION_DIR) create $(name) sql

migration-status:
	@echo "Checking migration status..."
	$(GOOSE_CMD) -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) status

# SQLC targets
gen-sqlc:
	@echo "Generating SQLC database code..."
	$(SQLC_CMD) generate

lint-sqlc:
	@echo "Linting SQLC code..."
	$(SQLC_CMD) vet
