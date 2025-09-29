#!/bin/bash

# Wonder Development Environment Cleanup Script
# This script provides options to clean up different layers of the environment

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

# Function to confirm action
confirm_action() {
    local message=$1

    if [ "$FORCE_CONFIRM" = true ]; then
        print_warning "--force supplied; auto-confirming: $message"
        return 0
    fi

    read -p "$(echo -e "${YELLOW}[CONFIRM]${NC} $message (y/N): ")" -n 1 -r
    echo
    [[ $REPLY =~ ^[Yy]$ ]]
}

# Parse command line arguments
CLEANUP_APP=false
CLEANUP_MONITORING=false
CLEANUP_INFRASTRUCTURE=false
CLEANUP_ALL=false
CLEANUP_VOLUMES=false
FORCE_CONFIRM=false

if [ $# -eq 0 ]; then
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --app               Stop application layer only (preserves data)"
    echo "  --monitoring        Stop monitoring layer only"
    echo "  --infrastructure    Stop infrastructure layer (databases, storage)"
    echo "  --all               Stop all services (preserves volumes)"
    echo "  --volumes           Remove all volumes (⚠️  DESTRUCTIVE - loses all data)"
    echo "  --force             Auto-confirm prompts (use with caution)"
    echo "  -h, --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --app                    # Stop app services for quick rebuild"
    echo "  $0 --monitoring             # Restart monitoring if needed"
    echo "  $0 --all                    # Stop everything but keep data"
    echo "  $0 --all --volumes          # Complete reset (⚠️  loses all data)"
    exit 0
fi

while [[ $# -gt 0 ]]; do
    case $1 in
        --app)
            CLEANUP_APP=true
            shift
            ;;
        --monitoring)
            CLEANUP_MONITORING=true
            shift
            ;;
        --infrastructure)
            CLEANUP_INFRASTRUCTURE=true
            shift
            ;;
        --all)
            CLEANUP_ALL=true
            shift
            ;;
        --volumes)
            CLEANUP_VOLUMES=true
            shift
            ;;
        --force|-f)
            FORCE_CONFIRM=true
            shift
            ;;
        -h|--help)
            echo "Help already shown above"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            exit 1
            ;;
    esac
done

print_status "🧹 Wonder Development Environment Cleanup..."

# Application Layer Cleanup
if [ "$CLEANUP_APP" = true ] || [ "$CLEANUP_ALL" = true ]; then
    print_status "🔧 Stopping Application Layer..."
    docker-compose -f docker-compose.app.yaml down
    print_success "✅ Application services stopped"
fi

# Monitoring Layer Cleanup
if [ "$CLEANUP_MONITORING" = true ] || [ "$CLEANUP_ALL" = true ]; then
    print_status "📊 Stopping Monitoring Layer..."
    docker-compose -f docker-compose.monitoring.yaml down
    print_success "✅ Monitoring services stopped"
fi

# Infrastructure Layer Cleanup
if [ "$CLEANUP_INFRASTRUCTURE" = true ] || [ "$CLEANUP_ALL" = true ]; then
    if [ "$CLEANUP_ALL" = true ]; then
        print_warning "Stopping infrastructure will shut down databases and storage!"
        if ! confirm_action "Are you sure you want to stop infrastructure services?"; then
            print_status "Skipping infrastructure cleanup"
        else
            print_status "🗄️  Stopping Infrastructure Layer..."
            docker-compose -f docker-compose.infrastructure.yaml down
            print_success "✅ Infrastructure services stopped"
        fi
    else
        print_warning "Stopping infrastructure will shut down databases and storage!"
        if confirm_action "Are you sure you want to stop infrastructure services?"; then
            print_status "🗄️  Stopping Infrastructure Layer..."
            docker-compose -f docker-compose.infrastructure.yaml down
            print_success "✅ Infrastructure services stopped"
        fi
    fi
fi

# Volume Cleanup (DESTRUCTIVE)
if [ "$CLEANUP_VOLUMES" = true ]; then
    print_error "⚠️  DESTRUCTIVE ACTION: This will permanently delete all data!"
    print_error "This includes:"
    print_error "  - All PostgreSQL databases and user data"
    print_error "  - All Elasticsearch indices and logs"
    print_error "  - All Grafana dashboards and configurations"
    echo ""

    if confirm_action "Are you ABSOLUTELY SURE you want to delete all data?"; then
        print_status "🗑️  Removing all volumes..."

        # Stop all services first to release volume locks
        docker-compose -f docker-compose.app.yaml down 2>/dev/null || true
        docker-compose -f docker-compose.monitoring.yaml down 2>/dev/null || true
        docker-compose -f docker-compose.infrastructure.yaml down 2>/dev/null || true

        # Remove volumes
        docker volume rm wonder-postgres-data 2>/dev/null || print_warning "PostgreSQL volume not found"
        docker volume rm wonder-grafana-data 2>/dev/null || print_warning "Grafana volume not found"
        docker volume rm wonder-elasticsearch-data 2>/dev/null || print_warning "Elasticsearch volume not found"

        # Remove network
        docker network rm wonder-network 2>/dev/null || print_warning "Network not found"

        print_success "✅ All volumes and data have been deleted"
        print_status "Use './scripts/dev-setup.sh' to recreate the environment"
    else
        print_status "Volume cleanup cancelled"
    fi
fi

# Clean up unused Docker resources
print_status "🧹 Cleaning up unused Docker resources..."
docker system prune -f > /dev/null 2>&1 || true

print_success "🎉 Cleanup completed!"

# Show remaining services
remaining_services=$(docker ps --format "table {{.Names}}" | grep "wonder-" || true)
if [ -n "$remaining_services" ]; then
    echo ""
    print_status "📋 Remaining Wonder services:"
    echo "$remaining_services"
else
    echo ""
    print_status "🏁 All Wonder services have been stopped"
fi
