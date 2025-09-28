# Wonder

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](#)
[![Test Coverage](https://img.shields.io/badge/Coverage-93.4%25-brightgreen.svg)](#testing)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](#docker-deployment)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](#)

A modern, scalable Go web service built with Domain-Driven Design (DDD) principles, featuring comprehensive authentication, monitoring, and observability capabilities.

## 🌟 Features

- **🏗️ Domain-Driven Design**: Clean architecture with clear layer separation
- **🔐 Authentication System**: JWT-based authentication with secure password handling
- **📊 Monitoring & Observability**: Integrated Prometheus metrics, Grafana dashboards, and ELK stack
- **🧪 Comprehensive Testing**: 93.4% test coverage with unit, integration, and E2E tests
- **🐳 Containerized**: Full Docker support with multi-stage builds
- **⚡ High Performance**: Built with Gin framework and optimized for scalability
- **📝 Structured Logging**: Configurable logging with file and stdout support
- **🔄 Database Migration**: Automated GORM-based database migrations
- **🎯 ID Generation**: Distributed Snowflake algorithm for unique ID generation

## 🏛️ Architecture

Wonder follows Domain-Driven Design principles with a clean, layered architecture:

```
┌─────────────────────────────────────────────────────────────────┐
│                    Interface Layer                              │
│              HTTP Handlers, Middleware                         │
├─────────────────────────────────────────────────────────────────┤
│                   Application Layer                             │
│              Use Cases, Application Services                    │
├─────────────────────────────────────────────────────────────────┤
│                     Domain Layer                                │
│        Entities, Value Objects, Domain Services                │
├─────────────────────────────────────────────────────────────────┤
│                 Infrastructure Layer                            │
│     Database, External Services, Technical Components          │
└─────────────────────────────────────────────────────────────────┘
```

See [docs/architecture.mermaid](docs/architecture.mermaid) for a detailed architectural diagram.

## 🚀 Quick Start

### Prerequisites

- **Go 1.24+**
- **PostgreSQL 16+**
- **Docker & Docker Compose** (optional)

### Local Development Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/cctw-zed/wonder.git
   cd wonder
   ```

2. **Set up environment:**
   ```bash
   source .envrc  # Load Go environment variables
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Start PostgreSQL:**
   ```bash
   # Using Docker
   docker run --name wonder-postgres \
     -e POSTGRES_DB=wonder_dev \
     -e POSTGRES_USER=dev \
     -e POSTGRES_PASSWORD=dev \
     -p 5432:5432 -d postgres:16-alpine
   ```

5. **Configure environment variables:**
   ```bash
   # Copy and modify configuration
   cp configs/config.yaml configs/config.development.yaml
   # Edit database connection settings
   ```

6. **Run the application:**
   ```bash
   go run cmd/server/main.go
   ```

7. **Verify installation:**
   ```bash
   curl http://localhost:8080/health
   ```

### Docker Deployment

**Quick start with Docker Compose:**

```bash
# Start all services (PostgreSQL + Wonder + Monitoring)
docker-compose up -d

# View logs
docker-compose logs -f wonder

# Stop services
docker-compose down
```

**Manual Docker build:**

```bash
# Build image
docker build -t wonder:latest .

# Run container
docker run -p 8080:8080 \
  -e WONDER_DATABASE_HOST=your-db-host \
  -e WONDER_DATABASE_USERNAME=your-username \
  -e WONDER_DATABASE_PASSWORD=your-password \
  wonder:latest
```

## 📖 API Documentation

Wonder provides a RESTful API with the following core endpoints:

### Authentication
- `POST /api/v1/users/register` - User registration (public)
- `POST /api/v1/auth/login` - User login (public)
- `POST /api/v1/auth/logout` - User logout (authenticated)
- `GET /api/v1/auth/me` - Get current user info (authenticated)

### User Management
- `GET /api/v1/users` - List users (optional auth)
- `GET /api/v1/users/:id` - Get user profile by ID (authenticated)
- `PUT /api/v1/users/:id` - Update user profile (authenticated)
- `DELETE /api/v1/users/:id` - Delete user (authenticated)

### Health & Monitoring
- `GET /health` - Application health check
- `GET /metrics` - Prometheus metrics endpoint

### Interactive API Testing

Use the included `api.http` file with your favorite HTTP client:

```bash
# VS Code with REST Client extension
code api.http

# Or use curl commands from the file
curl -X GET http://localhost:8080/health
```

**Example API calls:**

```bash
# Health check
curl http://localhost:8080/health

# User registration
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","name":"John Doe","password":"password123"}'

# User login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Get current user (authenticated)
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get user profile by ID (authenticated)
curl -X GET http://localhost:8080/api/v1/users/USER_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### API Response Formats

**Registration Response** (`POST /api/v1/users/register`):
```json
{
  "data": {
    "id": "user-id-123",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2025-09-28T10:00:00Z",
    "updated_at": "2025-09-28T10:00:00Z"
  },
  "trace_id": "trace-abc-123"
}
```

**Login Response** (`POST /api/v1/auth/login`):
```json
{
  "data": {
    "user": {
      "id": "user-id-123",
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2025-09-28T10:00:00Z",
      "updated_at": "2025-09-28T10:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  },
  "trace_id": "trace-abc-124"
}
```

**User Profile Response** (`GET /api/v1/auth/me` or `GET /api/v1/users/:id`):
```json
{
  "data": {
    "id": "user-id-123",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2025-09-28T10:00:00Z",
    "updated_at": "2025-09-28T10:00:00Z"
  },
  "trace_id": "trace-abc-125"
}
```

**Error Response Format**:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "details": {
      "field": "email",
      "value": "invalid-email"
    }
  },
  "trace_id": "trace-abc-126"
}
```

## 🧪 Testing

Wonder maintains **93.4% test coverage** with comprehensive testing at all layers:

### Running Tests

```bash
# Set up environment
source .envrc

# Run all tests
./scripts/test.sh all

# Run specific test types
./scripts/test.sh unit        # Unit tests only
./scripts/test.sh integration # Integration tests only
./scripts/test.sh e2e         # End-to-end tests only

# Generate coverage report
./scripts/test.sh coverage
open coverage/coverage.html
```

### Test Coverage by Layer

- **Domain Layer**: 100% - Pure business logic and entities
- **Application Layer**: 94.1% - Use cases and application services
- **Interface Layer**: 100% - HTTP handlers and API contracts
- **Repository Layer**: 84.6% - Database operations and GORM integration
- **E2E Tests**: Complete stack integration testing

### Testing Strategy

- **Unit Tests**: Fast, isolated tests with mocks for external dependencies
- **Integration Tests**: Test component interactions with real dependencies
- **E2E Tests**: Full application flow testing with real HTTP server and database
- **Test Data Builders**: Fluent API for creating test data objects

## 📊 Monitoring & Observability

Wonder includes a comprehensive monitoring stack:

### Metrics Collection
- **Prometheus**: Collects application and system metrics
- **Custom Metrics**: Business metrics and performance indicators
- **Health Checks**: Application and dependency health monitoring

### Visualization
- **Grafana**: Pre-configured dashboards for service monitoring
- **Business Dashboards**: User registration, authentication success rates
- **Infrastructure Dashboards**: System resources, database performance

### Logging
- **Structured Logging**: JSON-formatted logs for easy parsing
- **ELK Stack**: Elasticsearch, Logstash, and Kibana for log aggregation
- **Log Levels**: Configurable logging levels (debug, info, warn, error)

### Accessing Monitoring Tools

```bash
# Start monitoring stack
docker-compose up -d

# Access dashboards
open http://localhost:3000    # Grafana (admin/admin)
open http://localhost:9090    # Prometheus
open http://localhost:5601    # Kibana
```

## ⚙️ Configuration

Wonder uses a flexible, environment-aware configuration system:

### Configuration Files
- `configs/config.yaml` - Base configuration
- `configs/config.development.yaml` - Development overrides
- `configs/config.production.yaml` - Production settings
- `configs/config.testing.yaml` - Test environment settings

### Environment Variables

Key configuration options can be set via environment variables:

```bash
# Application settings
WONDER_APP_ENVIRONMENT=development
WONDER_APP_DEBUG=true
WONDER_SERVER_HOST=0.0.0.0
WONDER_SERVER_PORT=8080

# Database configuration
WONDER_DATABASE_HOST=localhost
WONDER_DATABASE_PORT=5432
WONDER_DATABASE_USERNAME=dev
WONDER_DATABASE_PASSWORD=dev
WONDER_DATABASE_DATABASE=wonder_dev

# JWT configuration
WONDER_JWT_SIGNING_KEY=your-secret-key
WONDER_JWT_EXPIRY=24h

# Logging configuration
WONDER_LOG_LEVEL=info
WONDER_LOG_OUTPUT=stdout
WONDER_LOG_ENABLE_FILE=false
```

See [README_CONFIG.md](docs/README_CONFIG.md) for detailed configuration options.

## 🏗️ Development

### Project Structure

```
wonder/
├── cmd/server/           # Application entry point
├── internal/             # Private application code
│   ├── domain/          # Domain layer (entities, value objects)
│   ├── application/     # Application layer (use cases, services)
│   ├── infrastructure/  # Infrastructure layer (database, external services)
│   ├── interfaces/      # Interface layer (HTTP handlers, DTOs)
│   ├── middleware/      # HTTP middleware components
│   └── server/          # Server configuration and setup
├── pkg/                 # Public/shared packages
│   ├── errors/         # Error handling system
│   ├── logger/         # Logging utilities
│   ├── jwt/           # JWT token management
│   └── snowflake/     # ID generation
├── test/               # Test files
│   ├── e2e/           # End-to-end tests
│   └── integration/   # Integration tests
├── docs/              # Documentation
├── configs/           # Configuration files
└── scripts/           # Build and deployment scripts
```

### Development Commands

```bash
# Set up environment (required before all commands)
source .envrc

# Build application
go build -o bin/server cmd/server/main.go

# Run application
go run cmd/server/main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Static analysis
go vet ./...

# Update dependencies
go mod tidy

# Generate mocks (for testing)
go generate ./...
```

### Code Quality Requirements

- **Test-Driven Development**: Write tests before implementation
- **Domain-Driven Design**: Follow DDD principles and patterns
- **Clean Architecture**: Maintain clear layer separation
- **Code Coverage**: Maintain >= 80% test coverage
- **Error Handling**: Use structured error handling system
- **Security**: Implement proper input validation and authentication

### Git Workflow

```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Make changes and commit
git add .
git commit -m "feat: implement feature description"

# Run tests before pushing
./scripts/test.sh all

# Push changes
git push origin feature/your-feature-name
```

## 🔧 Advanced Topics

### Database Operations

Wonder uses GORM for database operations with automatic migrations:

```bash
# Reset database (development only)
go run scripts/reset_db.go

# Manual migration
# Migrations run automatically on application startup
```

### ID Generation

Wonder uses Snowflake algorithm for distributed ID generation:

```go
// Service-specific ID generation
userID := id.GenerateForService(id.ServiceTypeUser, instanceID)
orderID := id.GenerateForService(id.ServiceTypeOrder, instanceID)
```

### Error Handling

Structured error handling with custom error types:

```go
// Domain errors
if err := user.Validate(); err != nil {
    return errors.NewDomainError("INVALID_USER", err.Error())
}

// Application errors
return errors.NewApplicationError("USER_NOT_FOUND", "User does not exist")

// HTTP errors with proper status codes
return errors.NewHTTPError(http.StatusBadRequest, "INVALID_REQUEST", "Missing required fields")
```

### Custom Middleware

Add custom middleware for cross-cutting concerns:

```go
// Authentication middleware
router.Use(middleware.AuthMiddleware())

// Metrics middleware
router.Use(middleware.MetricsMiddleware())

// Tracing middleware
router.Use(middleware.TracingMiddleware())
```

## 📚 Documentation

- **[Architecture](docs/architecture.mermaid)** - System architecture diagram
- **[Technical Specification](docs/technical.md)** - Detailed technical documentation
- **[Configuration Guide](docs/README_CONFIG.md)** - Configuration options and setup
- **[Docker Guide](docs/docker-guide.md)** - Docker deployment and operations
- **[Monitoring Guide](docs/monitoring.md)** - Monitoring setup and usage
- **[Logging Guide](docs/logging-guide.md)** - Logging configuration and best practices
- **[Testing Guide](docs/others/testing.md)** - Testing strategies and best practices

## 🤝 Contributing

We welcome contributions! Please read our development guidelines:

1. **Follow DDD Principles**: Maintain domain-driven design patterns
2. **Write Tests**: Ensure comprehensive test coverage
3. **Update Documentation**: Keep documentation current with changes
4. **Code Review**: All changes require code review
5. **Commit Messages**: Use conventional commit format

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make changes following coding standards
4. Write/update tests
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Repository**: [https://github.com/cctw-zed/wonder](https://github.com/cctw-zed/wonder)
- **Issues**: [GitHub Issues](https://github.com/cctw-zed/wonder/issues)
- **Documentation**: [Project Documentation](docs/)

## 📞 Support

For support and questions:

- **GitHub Issues**: Report bugs and request features
- **Documentation**: Check the [docs/](docs/) directory for detailed guides
- **API Reference**: Use the [api.http](test/api.http) file for API examples

---

**Wonder** - Building modern, scalable web services with Go and Domain-Driven Design principles.