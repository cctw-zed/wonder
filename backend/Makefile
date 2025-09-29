# Wonder Project Makefile

# Project information
PROJECT_NAME := wonder
VERSION ?= 1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Directories
BIN_DIR := ./bin
CMD_DIR := ./cmd

# Go build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# Default target
.DEFAULT_GOAL := build

# Phony targets
.PHONY: build build-all test run run-test clean kill help codex-context docker-build docker-up docker-down docker-logs
.PHONY: dev-setup dev-rebuild dev-rebuild-backend dev-rebuild-frontend dev-status dev-urls dev-stop-apps dev-clean-all

# Create bin directory
$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

# Build server binary
build: $(BIN_DIR)
	@echo "🚀 Building $(PROJECT_NAME) server..."
	@source .envrc && go build $(LDFLAGS) -o $(BIN_DIR)/server $(CMD_DIR)/server
	@echo "✅ Build completed: $(BIN_DIR)/server"

# Build for all platforms
build-all: $(BIN_DIR)
	@echo "🚀 Building $(PROJECT_NAME) for all platforms..."
	@source .envrc && GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/server-linux-amd64 $(CMD_DIR)/server
	@source .envrc && GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/server-darwin-amd64 $(CMD_DIR)/server
	@source .envrc && GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/server-darwin-arm64 $(CMD_DIR)/server
	@source .envrc && GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/server-windows-amd64.exe $(CMD_DIR)/server
	@echo "✅ All builds completed!"
	@ls -la $(BIN_DIR)/

# Run tests
test:
	@echo "🧪 Running tests..."
	@source .envrc && go test ./...

# Run server in development mode
run:
	@echo "🏃 Starting server in development mode..."
	@source .envrc && go run $(CMD_DIR)/server/main.go

# Run server in testing mode
run-test:
	@echo "🏃 Starting server in testing mode..."
	@source .envrc && go run $(CMD_DIR)/server/main.go -env=testing

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)/*
	@echo "✅ Clean completed!"

# Kill wonder server processes
kill:
	@echo "🔫 Killing wonder server processes..."
	@pkill -f "wonder" 2>/dev/null || true
	@pkill -f "go run.*cmd/server" 2>/dev/null || true
	@pkill -f "make run" 2>/dev/null || true
	@pkill -f "bin/server" 2>/dev/null || true
	@echo "✅ Wonder processes terminated!"

codex-context:
	@echo "Loading core documentation bundle..."
	@for file in docs/architecture.mermaid docs/technical.md docs/tasks/tasks.md docs/status.md; do \
		printf '\n===== %s =====\n' "$$file"; \
		cat "$$file"; \
		printf '\n'; \
	done

docker-build:
	@echo "🐳 Building Docker images..."
	@docker compose build

docker-up:
	@echo "🚀 Starting Docker stack..."
	@docker compose up -d

docker-down:
	@echo "🛑 Stopping Docker stack..."
	@docker compose down

docker-logs:
	@echo "📜 Streaming Docker logs..."
	@docker compose logs -f

# === NEW DEVELOPMENT WORKFLOW COMMANDS ===

# Set up complete development environment
dev-setup:
	@echo "🚀 Setting up Wonder development environment..."
	@./scripts/dev-setup.sh

# Quick rebuild commands
dev-rebuild:
	@echo "🔄 Quick rebuild: Application layer..."
	@./scripts/dev-rebuild.sh

dev-rebuild-backend:
	@echo "🔧 Rebuilding backend only..."
	@./scripts/dev-rebuild.sh --backend-only

dev-rebuild-frontend:
	@echo "🌐 Rebuilding frontend only..."
	@./scripts/dev-rebuild.sh --frontend-only

dev-rebuild-no-cache:
	@echo "🔄 Rebuilding without cache..."
	@./scripts/dev-rebuild.sh --no-cache

# Status and monitoring
dev-status:
	@echo "📋 Wonder Environment Status:"
	@echo ""
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep wonder- || echo "No Wonder services running"

dev-urls:
	@echo "🌐 Wonder Development Environment URLs:"
	@echo ""
	@echo "Application Services:"
	@echo "  🌐 Frontend:     http://localhost:3001"
	@echo "  🔧 Backend API:  http://localhost:8080"
	@echo ""
	@echo "Monitoring & Analytics:"
	@echo "  📊 Grafana:      http://localhost:3000 (admin/admin)"
	@echo "  📈 Prometheus:   http://localhost:9090"
	@echo "  📋 Kibana:       http://localhost:5601"
	@echo "  📦 cAdvisor:     http://localhost:8081"
	@echo ""
	@echo "Infrastructure:"
	@echo "  🗄️  PostgreSQL:   localhost:5432 (dev/dev/wonder_dev)"
	@echo "  🔍 Elasticsearch: http://localhost:9200"

# Cleanup commands
dev-stop-apps:
	@echo "🛑 Stopping application services..."
	@./scripts/dev-cleanup.sh --app

dev-stop-all:
	@echo "🛑 Stopping all services..."
	@./scripts/dev-cleanup.sh --all

dev-clean-all:
	@echo "🗑️  Complete cleanup (removes all data)..."
	@./scripts/dev-cleanup.sh --all --volumes

# Show help
help:
	@echo "Wonder Project Makefile"
	@echo ""
	@echo "🏗️  Build Commands:"
	@echo "  build                  Build server binary to bin/ directory"
	@echo "  build-all              Build server for all platforms"
	@echo "  test                   Run tests"
	@echo "  run                    Run server in development mode"
	@echo "  run-test               Run server in testing mode"
	@echo "  clean                  Clean build artifacts"
	@echo "  kill                   Kill all wonder server processes"
	@echo ""
	@echo "🚀 Development Environment (NEW):"
	@echo "  dev-setup              Set up complete development environment"
	@echo "  dev-rebuild            Rebuild applications (preserves data)"
	@echo "  dev-rebuild-backend    Rebuild only backend service"
	@echo "  dev-rebuild-frontend   Rebuild only frontend service"
	@echo "  dev-rebuild-no-cache   Rebuild without Docker cache"
	@echo ""
	@echo "📊 Monitoring & Status:"
	@echo "  dev-status             Show status of all services"
	@echo "  dev-urls               Show all service URLs"
	@echo ""
	@echo "🧹 Cleanup:"
	@echo "  dev-stop-apps          Stop application services only"
	@echo "  dev-stop-all           Stop all services (preserves data)"
	@echo "  dev-clean-all          Complete reset (⚠️  removes all data)"
	@echo ""
	@echo "🐳 Legacy Docker Commands:"
	@echo "  docker-build           Build Docker images"
	@echo "  docker-up              Start Docker stack"
	@echo "  docker-down            Stop Docker stack"
	@echo "  docker-logs            Stream Docker logs"
	@echo ""
	@echo "📋 Examples:"
	@echo "  make dev-setup                    # First time setup"
	@echo "  make dev-rebuild                  # Quick redeploy after code changes"
	@echo "  make dev-rebuild-backend          # Redeploy backend only"
	@echo "  make dev-status                   # Check what's running"
	@echo "  make VERSION=2.0.0 build          # Build with specific version"
