#!/bin/bash

# Wonder Docker Cache Fix Script
# This script fixes Docker build cache corruption issues

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

print_status "🔧 Wonder Docker Cache Fix Tool"
echo ""

# Step 1: Stop all Wonder services
print_status "🛑 Stopping all Wonder services..."
./scripts/dev-cleanup.sh --all --force >/dev/null 2>&1 || true

# Step 2: Clean only build cache and containers (preserve base images)
print_status "🧹 Cleaning Docker build cache and stopped containers..."
docker container prune -f >/dev/null 2>&1 || true
docker builder prune -a -f >/dev/null 2>&1 || true

# Step 3: Remove only Wonder-specific images (preserve base images)
print_status "🗑️  Removing only Wonder application images..."
docker rmi -f $(docker images -q --filter "reference=*wonder*") >/dev/null 2>&1 || true
docker rmi -f $(docker images -q --filter "reference=backend-*") >/dev/null 2>&1 || true

print_warning "Preserving base images (PostgreSQL, Elasticsearch, etc.) to avoid re-downloading"

# Step 4: Clear any remaining build cache
print_status "🧽 Clearing build cache..."
docker system df >/dev/null 2>&1 || true

print_success "✅ Docker cache cleanup completed!"
echo ""

print_status "🚀 Ready to rebuild. You can now run:"
echo "  make setup                 # For complete environment setup"
echo "  make rebuild-no-cache      # For no-cache rebuild of applications"
echo ""

print_warning "Note: The first build after cache cleanup will take longer as images need to be downloaded fresh."