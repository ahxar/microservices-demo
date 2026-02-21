.PHONY: help proto-gen build up down logs clean test migrate seed

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

proto-gen: ## Generate protobuf code for all services
	@echo "Generating protobuf code..."
	./scripts/proto-gen.sh

build: ## Build all Docker images
	@echo "Building all services..."
	docker compose -f deployments/docker/docker-compose.yml build

up: ## Start all services
	@echo "Starting all services..."
	docker compose -f deployments/docker/docker-compose.yml up -d

down: ## Stop all services
	@echo "Stopping all services..."
	docker compose -f deployments/docker/docker-compose.yml down

logs: ## Tail logs from all services
	docker compose -f deployments/docker/docker-compose.yml logs -f

clean: ## Remove all containers, volumes, and generated files
	@echo "Cleaning up..."
	docker compose -f deployments/docker/docker-compose.yml down -v
	find proto -name "*.pb.go" -delete

migrate: ## Run database migrations
	@echo "Running database migrations..."
	./scripts/apply-migrations.sh

seed: ## Seed database with sample data
	@echo "Seeding database..."
	./scripts/seed-data.sh

test: ## Run all tests
	@echo "Running tests..."
	cd gateway && go test ./...
	cd services/user && go test ./...
	cd services/catalog && go test ./...
	cd services/cart && go test ./...
	cd services/order && go test ./...
	cd services/shipping && go test ./...
	cd services/notification && go test ./...
	cd services/payment && cargo test

dev-user: ## Run user service locally
	cd services/user && go run cmd/user/main.go

dev-catalog: ## Run catalog service locally
	cd services/catalog && go run cmd/catalog/main.go

dev-cart: ## Run cart service locally
	cd services/cart && go run cmd/cart/main.go

dev-gateway: ## Run API gateway locally
	cd gateway && go run cmd/gateway/main.go
