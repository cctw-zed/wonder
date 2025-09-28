# Configuration Usage Guide

This project uses a comprehensive configuration system that supports multiple sources and environments.

## Configuration Sources (Priority Order)

1. **Environment Variables** (highest priority)
2. **Configuration Files** (YAML format)
3. **Default Values** (lowest priority)

## Configuration Files

The configuration files are located in the `configs/` directory:

### Default Configuration
- `configs/config.yaml` - Default configuration for development

### Environment-Specific Configuration
- `configs/config.development.yaml` - Development environment (same as default)
- `configs/config.testing.yaml` - Testing environment
- `configs/config.production.yaml` - Production environment (uses environment variable placeholders)

## Running the Application

### Using Default Configuration
```bash
# Uses config.yaml
go run cmd/server/main.go
```

### Using Custom Config Path
```bash
# Use a specific config file
go run cmd/server/main.go -config=/path/to/config.yaml

# Use a different configs directory
go run cmd/server/main.go -config=/path/to/configs/directory
```

### Using Environment-Specific Config
```bash
# Uses config.production.yaml
go run cmd/server/main.go -env=production

# Uses config.testing.yaml
go run cmd/server/main.go -env=testing
```

## Environment Variables Override

Any configuration can be overridden using environment variables:

```bash
# App configuration (with WONDER_ prefix)
export WONDER_APP_NAME="my-custom-app"

# Database settings (standard prefixes)
export DB_HOST="prod-database.example.com"
export DB_PORT="5432"
export DB_USERNAME="prod_user"
export DB_PASSWORD="secure_password"

# Server settings (standard prefixes)
export SERVER_HOST="0.0.0.0"
export SERVER_PORT="8080"

# ID generator settings (for production config placeholders)
export ID_SERVICE_TYPE="order"
export ID_INSTANCE_ID="42"
export ID_NODE_ID="1"

# External services (for production config placeholders)
export REDIS_HOST="redis.example.com"
export REDIS_PASSWORD="redis_password"
export EMAIL_HOST="smtp.example.com"
export EMAIL_USERNAME="noreply@example.com"
export EMAIL_PASSWORD="email_password"
```

**Note**: The production config uses environment variable placeholders like `${DB_HOST}` that must be set when running in production.

## Configuration Structure

```yaml
app:
  name: "wonder"                 # Application name
  version: "1.0.0"              # Application version
  environment: "development"     # Environment (development/testing/production)
  debug: true                   # Debug mode

server:
  host: "localhost"             # Server bind address
  port: 8080                    # Server port
  read_timeout: "30s"           # HTTP read timeout
  write_timeout: "30s"          # HTTP write timeout
  idle_timeout: "60s"           # HTTP idle timeout
  enable_cors: true             # Enable CORS middleware

database:
  host: "localhost"             # Database host
  port: 5432                    # Database port
  username: "dev"               # Database username
  password: "dev"               # Database password
  database: "wonder_dev"        # Database name
  ssl_mode: "disable"           # SSL mode (disable/require/verify-ca/verify-full)
  timezone: "UTC"               # Database timezone
  max_open_conns: 25            # Maximum open connections
  max_idle_conns: 10            # Maximum idle connections
  conn_max_lifetime: "1h"       # Connection maximum lifetime
  conn_max_idle_time: "30m"     # Connection maximum idle time
  log_level: "info"             # Database log level

log:
  level: "info"                 # Log level (debug/info/warn/error/fatal)
  format: "json"                # Log format (json/text/console)
  output: "stdout"              # Log output (stdout/stderr/file)
  enable_file: false            # Enable file logging
  file_path: "logs/app.log"     # Log file path

id:
  service_type: "user"          # Service type for ID generation
  instance_id: 0                # Instance ID for distributed ID generation
  node_id: 1                    # Node ID for snowflake algorithm

external:
  redis:                        # Redis configuration (future use)
    host: "localhost"
    port: 6379
    password: ""
    database: 0
  email:                        # Email service configuration (future use)
    provider: "smtp"
    host: "smtp.gmail.com"
    port: 587
    username: ""
    password: ""
```

## Health Check Endpoint

The application provides a health check endpoint that returns configuration information:

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "app": "wonder",
  "version": "1.0.0",
  "environment": "development"
}
```

## Advanced Usage

### Loading Config in Code

```go
package main

import (
    "github.com/cctw-zed/wonder/internal/infrastructure/config"
)

func main() {
    // Load default config
    cfg, err := config.Load()

    // Load from specific path
    cfg, err := config.Load("/path/to/config.yaml")

    // Load environment-specific config
    cfg, err := config.LoadForEnvironment("production")

    // Write config to file
    err := config.WriteConfig(cfg, "/path/to/output.yaml")
}
```

### Container with Config

```go
package main

import (
    "context"
    "github.com/cctw-zed/wonder/internal/container"
)

func main() {
    ctx := context.Background()

    // Default container (uses default config)
    container, err := container.NewContainer()

    // Container with custom config path
    container, err := container.NewContainerWithConfig(ctx, "/path/to/config.yaml")

    // Access loaded config
    cfg := container.Config
    fmt.Printf("App: %s, Env: %s\n", cfg.App.Name, cfg.App.Environment)
}
```

## Environment-Specific Deployment

### Development
```bash
go run cmd/server/main.go
# or
go run cmd/server/main.go -env=development
```

### Testing
```bash
go run cmd/server/main.go -env=testing
```

### Production
```bash
# Build binary
go build -o bin/server cmd/server/main.go

# Run with production config
./bin/server -env=production

# Or with custom config
./bin/server -config=/etc/wonder/production.yaml
```

## Testing

### Running Tests

The project includes a comprehensive test runner script that supports different testing scenarios:

```bash
# Make script executable (first time only)
chmod +x scripts/test.sh

# Run different types of tests
./scripts/test.sh unit           # Unit tests only
./scripts/test.sh integration    # Integration tests only
./scripts/test.sh e2e           # End-to-end tests only
./scripts/test.sh all           # All tests (unit + integration + e2e)
./scripts/test.sh coverage     # Generate coverage report
```

### End-to-End Testing

E2E tests require a test database. Set up environment variables:

```bash
# Set up test database
export DB_USERNAME="test"
export DB_PASSWORD="test"
export DB_DATABASE="wonder_test"

# Run E2E tests
./scripts/test.sh e2e

# Or run directly with Go
DB_USERNAME="test" DB_PASSWORD="test" DB_DATABASE="wonder_test" go test ./test/e2e/... -v
```

### Test Options

```bash
# Verbose output
./scripts/test.sh unit -v

# Short mode (skips long-running tests)
./scripts/test.sh all -s

# Race condition detection
./scripts/test.sh all --race

# Clean test cache
./scripts/test.sh e2e --clean

# Show help
./scripts/test.sh --help
```

For detailed testing documentation, see [docs/testing.md](others/testing.md).