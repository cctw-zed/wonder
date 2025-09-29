#!/bin/bash

# Wonder Development Environment Setup Script
# This script sets up the entire development environment in the correct order

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
    local max_attempts=30
    local attempt=1

    print_status "Waiting for $service_name to be healthy..."

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

    print_error "$service_name failed to become healthy within expected time"
    return 1
}

print_status "🚀 Setting up Wonder Development Environment..."

# Step 1: Start Infrastructure Layer
print_status "📦 Starting Infrastructure Layer (Databases & Storage)..."
docker-compose -f docker-compose.infrastructure.yaml up -d

# Wait for infrastructure services to be ready
check_service_health "postgres"
check_service_health "elasticsearch"
check_service_health "grafana"

print_success "✅ Infrastructure layer is ready!"

# Step 2: Start Monitoring Layer
print_status "📊 Starting Monitoring Layer (Metrics & Logging)..."
docker-compose -f docker-compose.monitoring.yaml up -d

# Wait for monitoring services
sleep 10
print_success "✅ Monitoring layer is ready!"

# Step 3: Start Application Layer
print_status "🔧 Starting Application Layer (Backend & Frontend)..."
docker-compose -f docker-compose.app.yaml up -d --build

# Wait for application services
check_service_health "backend"
check_service_health "frontend"

print_success "✅ Application layer is ready!"

# Final status check
print_status "🔍 Final System Status Check..."

echo ""
echo "=== Wonder Development Environment Status ==="
echo ""

# Check all services
services=("postgres" "elasticsearch" "grafana" "prometheus" "logstash" "kibana" "cadvisor" "backend" "frontend")

for service in "${services[@]}"; do
    if docker ps --format "table {{.Names}}" | grep -q "wonder-$service"; then
        status=$(docker inspect --format='{{.State.Status}}' "wonder-$service" 2>/dev/null || echo "unknown")
        if [ "$status" = "running" ]; then
            echo -e "✅ wonder-$service: ${GREEN}$status${NC}"
        else
            echo -e "❌ wonder-$service: ${RED}$status${NC}"
        fi
    else
        echo -e "❌ wonder-$service: ${RED}not found${NC}"
    fi
done

echo ""
echo "=== Service URLs ==="
echo "🌐 Frontend:     http://localhost:3001"
echo "🔧 Backend API:  http://localhost:8080"
echo "📊 Grafana:      http://localhost:3000 (admin/admin)"
echo "📈 Prometheus:   http://localhost:9090"
echo "📋 Kibana:       http://localhost:5601"
echo "🗄️  PostgreSQL:   localhost:5432 (dev/dev/wonder_dev)"
echo "🔍 Elasticsearch: http://localhost:9200"
echo "📦 cAdvisor:     http://localhost:8081"
echo ""

print_success "🎉 Wonder Development Environment is fully ready!"
print_status "To rebuild and redeploy applications only, use: ./scripts/dev-rebuild.sh"