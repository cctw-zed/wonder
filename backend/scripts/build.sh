#!/bin/bash

# Wonder Project Build Script
# Usage: ./scripts/build.sh [target] [platform]
# Examples:
#   ./scripts/build.sh server          # Build server for current platform
#   ./scripts/build.sh server linux    # Build server for Linux
#   ./scripts/build.sh all             # Build all targets

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project information
PROJECT_NAME="wonder"
VERSION=${VERSION:-"1.0.0"}
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build directories
BIN_DIR="./bin"
BUILD_DIR="./build"

# Source environment
source .envrc

echo -e "${BLUE}🚀 Building Wonder Project${NC}"
echo -e "${YELLOW}Version: ${VERSION}${NC}"
echo -e "${YELLOW}Build Time: ${BUILD_TIME}${NC}"
echo -e "${YELLOW}Git Commit: ${GIT_COMMIT}${NC}"
echo ""

# Create directories
mkdir -p $BIN_DIR
mkdir -p $BUILD_DIR

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

# Function to build for specific platform
build_server() {
    local goos=${1:-$(go env GOOS)}
    local goarch=${2:-$(go env GOARCH)}
    local suffix=""

    if [ "$goos" = "windows" ]; then
        suffix=".exe"
    fi

    local output_name="server"
    if [ "$goos" != "$(go env GOOS)" ] || [ "$goarch" != "$(go env GOARCH)" ]; then
        output_name="server-${goos}-${goarch}"
    fi

    echo -e "${GREEN}Building server for ${goos}/${goarch}...${NC}"
    GOOS=$goos GOARCH=$goarch go build \
        -ldflags "$LDFLAGS" \
        -o "${BIN_DIR}/${output_name}${suffix}" \
        ./cmd/server

    echo -e "${GREEN}✅ Built: ${BIN_DIR}/${output_name}${suffix}${NC}"
}

# Function to build all platforms
build_all() {
    echo -e "${BLUE}Building for all platforms...${NC}"

    # Common platforms
    build_server "linux" "amd64"
    build_server "darwin" "amd64"
    build_server "darwin" "arm64"  # Apple Silicon
    build_server "windows" "amd64"

    echo -e "${GREEN}✅ All builds completed!${NC}"
    ls -la $BIN_DIR/
}

# Function to clean build artifacts
clean() {
    echo -e "${YELLOW}🧹 Cleaning build artifacts...${NC}"
    rm -rf $BIN_DIR/*
    rm -rf $BUILD_DIR/*
    echo -e "${GREEN}✅ Clean completed!${NC}"
}

# Function to show usage
usage() {
    echo "Usage: $0 [command] [platform] [architecture]"
    echo ""
    echo "Commands:"
    echo "  server              Build server for current platform"
    echo "  all                 Build server for all platforms"
    echo "  clean               Clean build artifacts"
    echo ""
    echo "Platform examples:"
    echo "  $0 server linux amd64"
    echo "  $0 server darwin arm64"
    echo "  $0 server windows amd64"
    echo ""
    echo "Available platforms:"
    echo "  linux/amd64, darwin/amd64, darwin/arm64, windows/amd64"
}

# Main logic
case ${1:-server} in
    "server")
        build_server $2 $3
        ;;
    "all")
        build_all
        ;;
    "clean")
        clean
        ;;
    "help"|"-h"|"--help")
        usage
        ;;
    *)
        echo -e "${RED}❌ Unknown command: $1${NC}"
        usage
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}🎉 Build process completed successfully!${NC}"