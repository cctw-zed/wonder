# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

This is a Go project using standard Go toolchain:

- **Build**: `go build` or `go build ./cmd/server` for the server binary
- **Run**: `go run main.go` or `go run ./cmd/server/main.go` 
- **Test**: `go test ./...` (run all tests)
- **Test single package**: `go test ./internal/application/service`
- **Format**: `go fmt ./...`
- **Lint**: `go vet ./...`
- **Tidy dependencies**: `go mod tidy`

## Architecture

This is a Go web service following **Hexagonal Architecture** (Clean Architecture) principles:

### Directory Structure
- `cmd/server/` - Application entry point and main binary
- `internal/domain/` - Domain entities and business logic (User aggregate)
- `internal/application/service/` - Application services implementing domain interfaces
- `internal/infrastructure/` - External concerns (cache, mq, repository implementations)
- `internal/interfaces/` - Interface adapters (HTTP handlers, gRPC)
- `internal/container/` - Dependency injection container
- `pkg/` - Shared packages (snowflake ID generator)
- `configs/` - Configuration files
- `deployments/` - Deployment configurations

### Key Components

**Domain Layer** (`internal/domain/user/`):
- `User` struct - Domain entity with business rules
- `UserRepository` interface - Repository contract
- `UserService` interface - Domain service contract

**Application Layer** (`internal/application/service/`):
- `userService` - Implements domain service interface
- Contains business logic and orchestration

**Infrastructure Layer** (`internal/infrastructure/`):
- Repository implementations (not yet implemented)
- Cache and message queue adapters
- External service integrations

**Interface Layer** (`internal/interfaces/http/`):
- `UserHandler` - HTTP REST API handlers using Gin framework
- Request/response DTOs

**Dependencies**:
- **Gin** (`github.com/gin-gonic/gin`) - HTTP web framework
- **Snowflake** (`github.com/bwmarrin/snowflake`) - Distributed ID generation

### Distributed ID Generation
The project uses an optimized Snowflake algorithm with service-based node ID segmentation for distributed deployment:

**Node ID Allocation by Service Type:**
- User Service: Node IDs 0-1023 (ServiceTypeUser)
- Order Service: Node IDs 1024-2047 (ServiceTypeOrder) 
- Payment Service: Node IDs 2048-3071 (ServiceTypePayment)
- Auth Service: Node IDs 3072-4095 (ServiceTypeAuth)
- Gateway Service: Node IDs 4096-5119 (ServiceTypeGateway)

**Environment Configuration:**
1. **Preferred**: Use `SERVICE_TYPE` + `INSTANCE_ID`
   - `SERVICE_TYPE`: user|order|payment|auth|gateway
   - `INSTANCE_ID`: 0-1023 (instance within service)
2. **Legacy**: Use `NODE_ID` directly (0-10239)

**Usage:**
- Initialize: `id.InitDefaultFromEnv()` (auto-detects from env vars)
- Generate: `id.Generate()` (string) or `id.GenerateInt64()`
- Service-specific: `id.InitDefaultForService(serviceType, instanceID)`

**Examples:**
```bash
# User service, instance 5
SERVICE_TYPE=user INSTANCE_ID=5  # Results in node ID 5

# Order service, instance 100  
SERVICE_TYPE=order INSTANCE_ID=100  # Results in node ID 1124

# Legacy mode
NODE_ID=42  # Direct node ID assignment
```

### Container Pattern
`internal/container/container.go` implements dependency injection:
- Initializes all services and handlers
- Manages component lifecycle
- Currently missing repository layer implementation

## Current State
This appears to be a newly scaffolded project with:
- Basic user domain model and service interfaces defined
- HTTP handler for user registration
- Snowflake ID generator implemented
- Infrastructure layer structure created but not implemented
- No database connections or repository implementations yet