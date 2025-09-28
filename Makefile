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
.PHONY: build build-all test run run-test clean kill help

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

# Show help
help:
	@echo "Wonder Project Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build      Build server binary to bin/ directory"
	@echo "  build-all  Build server for all platforms"
	@echo "  test       Run tests"
	@echo "  run        Run server in development mode"
	@echo "  run-test   Run server in testing mode"
	@echo "  clean      Clean build artifacts"
	@echo "  kill       Kill all wonder server processes"
	@echo "  help       Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build              # Build server"
	@echo "  make run-test           # Start testing environment"
	@echo "  make VERSION=2.0.0 build # Build with specific version"