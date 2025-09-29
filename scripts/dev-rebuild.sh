#!/bin/bash

# Wonder Development Quick Rebuild Script
# This script rebuilds and redeploys ONLY the application layer
# without affecting infrastructure and monitoring services

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if a service is healthy
check_service_health() {
    local service_name=$1
    local max_attempts=15
    local attempt=1

    print_status "Waiting for $service_name to be ready..."

    while [ $attempt -le $max_attempts ]; do
        if docker inspect --format='{{.State.Health.Status}}' "wonder-$service_name" 2>/dev/null | grep -q "healthy"; then
            print_success "$service_name is healthy!"
            return 0
        fi

        if docker inspect --format='{{.State.Status}}' "wonder-$service_name" 2>/dev/null | grep -q "running"; then
            print_success "$service_name is running!"
            return 0
        fi

        print_status "Attempt $attempt/$max_attempts: $service_name not ready yet..."
        sleep 2
        ((attempt++))
    done

    print_warning "$service_name took longer than expected to be ready"
    return 1
}

print_status "🔄 Quick Rebuild: Application Layer Only..."

# Check if infrastructure is running
if ! docker ps --format "table {{.Names}}" | grep -q "wonder-postgres"; then
    print_warning "Infrastructure layer not running. Starting full environment..."
    ./scripts/dev-setup.sh
    exit 0
fi

# Parse command line arguments
REBUILD_BACKEND=true
REBUILD_FRONTEND=true
NO_CACHE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --backend-only)
            REBUILD_FRONTEND=false
            shift
            ;;
        --frontend-only)
            REBUILD_BACKEND=false
            shift
            ;;
        --no-cache)
            NO_CACHE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --backend-only    Rebuild only the backend service"
            echo "  --frontend-only   Rebuild only the frontend service"
            echo "  --no-cache        Force rebuild without using Docker cache"
            echo "  -h, --help        Show this help message"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Build Docker Compose flags
COMPOSE_FLAGS=""
if [ "$NO_CACHE" = true ]; then
    COMPOSE_FLAGS="--no-cache"
fi

# Stop application services
print_status "🛑 Stopping application services..."
docker-compose -f docker-compose.app.yaml down

# Build and start specific services
if [ "$REBUILD_BACKEND" = true ] && [ "$REBUILD_FRONTEND" = true ]; then
    print_status "🔧 Rebuilding both Backend and Frontend..."
    docker-compose -f docker-compose.app.yaml up -d --build $COMPOSE_FLAGS
elif [ "$REBUILD_BACKEND" = true ]; then
    print_status "🔧 Rebuilding Backend only..."
    docker-compose -f docker-compose.app.yaml up -d --build $COMPOSE_FLAGS wonder-backend
elif [ "$REBUILD_FRONTEND" = true ]; then
    print_status "🔧 Rebuilding Frontend only..."
    docker-compose -f docker-compose.app.yaml up -d --build $COMPOSE_FLAGS wonder-frontend
fi

# Wait for services to be ready
if [ "$REBUILD_BACKEND" = true ]; then
    check_service_health "backend"
fi

if [ "$REBUILD_FRONTEND" = true ]; then
    check_service_health "frontend"
fi

# Show service status
print_status "📋 Current Application Status:"
echo ""

if [ "$REBUILD_BACKEND" = true ]; then
    backend_status=$(docker inspect --format='{{.State.Status}}' "wonder-backend" 2>/dev/null || echo "not found")
    if [ "$backend_status" = "running" ]; then
        echo -e "🔧 Backend:  ${GREEN}$backend_status${NC} - http://localhost:8080"
    else
        echo -e "🔧 Backend:  ${RED}$backend_status${NC}"
    fi
fi

if [ "$REBUILD_FRONTEND" = true ]; then
    frontend_status=$(docker inspect --format='{{.State.Status}}' "wonder-frontend" 2>/dev/null || echo "not found")
    if [ "$frontend_status" = "running" ]; then
        echo -e "🌐 Frontend: ${GREEN}$frontend_status${NC} - http://localhost:3001"
    else
        echo -e "🌐 Frontend: ${RED}$frontend_status${NC}"
    fi
fi

echo ""
print_success "🎉 Application rebuild completed!"
print_status "Infrastructure and monitoring services remain unchanged and preserve all data."