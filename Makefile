# Wonder Project Makefile - Root Level
# Manages the entire Wonder development environment from project root

SHELL := /bin/bash

.PHONY: help setup rebuild rebuild-backend rebuild-frontend cleanup logs status urls
.PHONY: stop-apps stop-all clean-all test test-backend test-frontend

# Default target
help:
	@echo "🎯 Wonder Project - Development Environment Management"
	@echo ""
	@echo "🚀 Environment Management:"
	@echo "  make setup              - Set up complete development environment"
	@echo "  make status             - Show status of all services"
	@echo "  make urls               - Show all service URLs"
	@echo "  make logs               - Show logs from all services"
	@echo ""
	@echo "🔄 Quick Development:"
	@echo "  make rebuild            - Rebuild and redeploy applications (preserves data)"
	@echo "  make rebuild-backend    - Rebuild only backend service"
	@echo "  make rebuild-frontend   - Rebuild only frontend service"
	@echo "  make rebuild-no-cache   - Rebuild without Docker cache"
	@echo ""
	@echo "🧹 Cleanup:"
	@echo "  make stop-apps          - Stop application services only"
	@echo "  make stop-all           - Stop all services (preserves data)"
	@echo "  make clean-all          - Stop all services and remove volumes (⚠️  DESTRUCTIVE)"
	@echo "  make fix-docker-cache   - Fix Docker cache corruption errors"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  make test               - Run all tests"
	@echo "  make test-backend       - Run backend tests only"
	@echo "  make test-frontend      - Run frontend tests only"
	@echo ""
	@echo "📋 Examples:"
	@echo "  make setup                      # First time setup"
	@echo "  make rebuild                    # Quick redeploy after code changes"
	@echo "  make rebuild-backend            # Redeploy backend only"
	@echo "  make status                     # Check what's running"

# Environment setup
setup:
	@echo "🚀 Setting up Wonder development environment..."
	@./scripts/dev-setup.sh

# Quick rebuild commands
rebuild:
	@echo "🔄 Quick rebuild: Application layer..."
	@./scripts/dev-rebuild.sh

rebuild-backend:
	@echo "🔧 Rebuilding backend only..."
	@./scripts/dev-rebuild.sh --backend-only

rebuild-frontend:
	@echo "🌐 Rebuilding frontend only..."
	@./scripts/dev-rebuild.sh --frontend-only

rebuild-no-cache:
	@echo "🔄 Rebuilding without cache..."
	@./scripts/dev-rebuild.sh --no-cache

# Cleanup commands
stop-apps:
	@echo "🛑 Stopping application services..."
	@./scripts/dev-cleanup.sh --app

stop-all:
	@echo "🛑 Stopping all services..."
	@./scripts/dev-cleanup.sh --all --force

clean-all:
	@echo "🗑️  Complete cleanup (removes all data)..."
	@./scripts/dev-cleanup.sh --all --volumes

fix-docker-cache:
	@echo "🔧 Fixing Docker cache corruption..."
	@./scripts/fix-docker-cache.sh

# Status and monitoring
status:
	@echo "📋 Wonder Environment Status:"
	@echo ""
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep wonder- || echo "No Wonder services running"

logs:
	@echo "📜 Wonder Service Logs:"
	@echo ""
	@echo "=== Backend Logs ==="
	@docker logs wonder-backend --tail 20 2>/dev/null || echo "Backend not running"
	@echo ""
	@echo "=== Frontend Logs ==="
	@docker logs wonder-frontend --tail 20 2>/dev/null || echo "Frontend not running"
	@echo ""
	@echo "=== PostgreSQL Logs ==="
	@docker logs wonder-postgres --tail 10 2>/dev/null || echo "PostgreSQL not running"

urls:
	@echo "🌐 Wonder Development Environment URLs:"
	@echo ""
	@echo "📱 Application Services:"
	@echo "  🌐 Frontend:     http://localhost:3001"
	@echo "  🔧 Backend API:  http://localhost:8080"
	@echo "  📋 API Docs:     http://localhost:8080/health"
	@echo ""
	@echo "📊 Monitoring & Analytics:"
	@echo "  📊 Grafana:      http://localhost:3000 (admin/admin)"
	@echo "  📈 Prometheus:   http://localhost:9090"
	@echo "  📋 Kibana:       http://localhost:5601"
	@echo "  📦 cAdvisor:     http://localhost:8081"
	@echo ""
	@echo "🗄️  Infrastructure:"
	@echo "  🗄️  PostgreSQL:   localhost:5432 (dev/dev/wonder_dev)"
	@echo "  🔍 Elasticsearch: http://localhost:9200"

# Testing commands
test:
	@echo "🧪 Running all tests..."
	@$(MAKE) test-backend
	@$(MAKE) test-frontend

test-backend:
	@echo "🧪 Running backend tests..."
	@cd backend && source .envrc && go test ./...

test-frontend:
	@echo "🧪 Running frontend tests..."
	@cd frontend && npm test --passWithNoTests

# Development utilities
dev-shell-backend:
	@echo "🐚 Opening shell in backend container..."
	@docker exec -it wonder-backend /bin/sh || echo "Backend container not running. Use 'make setup' first."

dev-shell-postgres:
	@echo "🐚 Opening PostgreSQL shell..."
	@docker exec -it wonder-postgres psql -U dev -d wonder_dev || echo "PostgreSQL not running. Use 'make setup' first."

dev-reset-db:
	@echo "🗑️  Resetting database..."
	@cd backend && source .envrc && go run scripts/reset_db.go

# Layer-specific commands for advanced users
start-infrastructure:
	@echo "📦 Starting infrastructure layer only..."
	@docker-compose -f docker-compose.infrastructure.yaml up -d

start-monitoring:
	@echo "📊 Starting monitoring layer only..."
	@docker-compose -f docker-compose.monitoring.yaml up -d

start-apps:
	@echo "🔧 Starting application layer only..."
	@docker-compose -f docker-compose.app.yaml up -d

# Health checks
health-check:
	@echo "🏥 Health Check Summary:"
	@echo ""
	@echo "Backend Health:"
	@curl -s http://localhost:8080/health 2>/dev/null | jq . || echo "❌ Backend not responding"
	@echo ""
	@echo "Frontend Health:"
	@curl -s http://localhost:3001/api/health 2>/dev/null | jq . || echo "❌ Frontend not responding"
	@echo ""

# Quick development commands
dev:
	@echo "🔧 Starting development mode..."
	@echo "This will:"
	@echo "  1. Set up the complete environment if not running"
	@echo "  2. Show service URLs"
	@echo "  3. Start watching for changes"
	@$(MAKE) setup
	@$(MAKE) urls
	@echo ""
	@echo "✨ Development environment ready!"
	@echo "💡 Use 'make rebuild' after making code changes"
	@echo "💡 Use 'make logs' to monitor service logs"
