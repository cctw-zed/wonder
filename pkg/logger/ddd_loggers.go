package logger

import (
	"context"
	"time"
)

// DomainLogger provides logging functionality for the domain layer
type DomainLogger struct {
	Logger
}

// NewDomainLogger creates a domain-specific logger
func NewDomainLogger(baseLogger Logger) *DomainLogger {
	return &DomainLogger{
		Logger: baseLogger.WithLayer(DomainLayer),
	}
}

// LogBusinessRule logs domain business rule validation
func (d *DomainLogger) LogBusinessRule(ctx context.Context, rule string, entityType string, entityID string, passed bool, fields ...Field) {
	allFields := append(fields,
		String("rule", rule),
		String("entity_type", entityType),
		String("entity_id", entityID),
		Bool("passed", passed),
	)

	if passed {
		d.Debug(ctx, "Business rule validation passed", allFields...)
	} else {
		d.Warn(ctx, "Business rule validation failed", allFields...)
	}
}

// LogDomainEvent logs domain events
func (d *DomainLogger) LogDomainEvent(ctx context.Context, eventType string, aggregateID string, fields ...Field) {
	allFields := append(fields,
		String("event_type", eventType),
		String("aggregate_id", aggregateID),
		Time("event_time", time.Now()),
	)

	d.Info(ctx, "Domain event occurred", allFields...)
}

// LogAggregateChange logs changes to aggregates
func (d *DomainLogger) LogAggregateChange(ctx context.Context, aggregateType string, aggregateID string, operation string, fields ...Field) {
	allFields := append(fields,
		String("aggregate_type", aggregateType),
		String("aggregate_id", aggregateID),
		String("operation", operation),
	)

	d.Debug(ctx, "Aggregate state changed", allFields...)
}

// ApplicationLogger provides logging functionality for the application layer
type ApplicationLogger struct {
	Logger
}

// NewApplicationLogger creates an application-specific logger
func NewApplicationLogger(baseLogger Logger) *ApplicationLogger {
	return &ApplicationLogger{
		Logger: baseLogger.WithLayer(ApplicationLayer),
	}
}

// LogUseCase logs use case execution
func (a *ApplicationLogger) LogUseCase(ctx context.Context, useCaseName string, startTime time.Time, success bool, fields ...Field) {
	duration := time.Since(startTime)
	allFields := append(fields,
		String("use_case", useCaseName),
		Duration("duration", duration),
		Bool("success", success),
	)

	if success {
		a.Info(ctx, "Use case completed successfully", allFields...)
	} else {
		a.Error(ctx, "Use case failed", allFields...)
	}
}

// LogServiceCall logs application service method calls
func (a *ApplicationLogger) LogServiceCall(ctx context.Context, serviceName string, method string, startTime time.Time, fields ...Field) {
	duration := time.Since(startTime)
	allFields := append(fields,
		String("service", serviceName),
		String("method", method),
		Duration("duration", duration),
	)

	a.Debug(ctx, "Service method called", allFields...)
}

// LogValidation logs application-level validation
func (a *ApplicationLogger) LogValidation(ctx context.Context, validationType string, passed bool, errors []string, fields ...Field) {
	allFields := append(fields,
		String("validation_type", validationType),
		Bool("passed", passed),
	)

	if len(errors) > 0 {
		allFields = append(allFields, Any("validation_errors", errors))
	}

	if passed {
		a.Debug(ctx, "Application validation passed", allFields...)
	} else {
		a.Warn(ctx, "Application validation failed", allFields...)
	}
}

// InfrastructureLogger provides logging functionality for the infrastructure layer
type InfrastructureLogger struct {
	Logger
}

// NewInfrastructureLogger creates an infrastructure-specific logger
func NewInfrastructureLogger(baseLogger Logger) *InfrastructureLogger {
	return &InfrastructureLogger{
		Logger: baseLogger.WithLayer(InfrastructureLayer),
	}
}

// LogDatabaseOperation logs database operations
func (i *InfrastructureLogger) LogDatabaseOperation(ctx context.Context, operation string, table string, duration time.Duration, rowsAffected int64, fields ...Field) {
	allFields := append(fields,
		String("operation", operation),
		String("table", table),
		Duration("duration", duration),
		Int64("rows_affected", rowsAffected),
	)

	if duration > 100*time.Millisecond {
		i.Warn(ctx, "Slow database operation", allFields...)
	} else {
		i.Debug(ctx, "Database operation completed", allFields...)
	}
}

// LogExternalServiceCall logs calls to external services
func (i *InfrastructureLogger) LogExternalServiceCall(ctx context.Context, serviceName string, endpoint string, method string, statusCode int, duration time.Duration, fields ...Field) {
	allFields := append(fields,
		String("external_service", serviceName),
		String("endpoint", endpoint),
		String("method", method),
		Int("status_code", statusCode),
		Duration("duration", duration),
	)

	if statusCode >= 400 {
		i.Error(ctx, "External service call failed", allFields...)
	} else if duration > 5*time.Second {
		i.Warn(ctx, "Slow external service call", allFields...)
	} else {
		i.Info(ctx, "External service call completed", allFields...)
	}
}

// LogCacheOperation logs cache operations
func (i *InfrastructureLogger) LogCacheOperation(ctx context.Context, operation string, key string, hit bool, duration time.Duration, fields ...Field) {
	allFields := append(fields,
		String("operation", operation),
		String("cache_key", key),
		Bool("cache_hit", hit),
		Duration("duration", duration),
	)

	i.Debug(ctx, "Cache operation performed", allFields...)
}

// InterfaceLogger provides logging functionality for the interface layer
type InterfaceLogger struct {
	Logger
}

// NewInterfaceLogger creates an interface-specific logger
func NewInterfaceLogger(baseLogger Logger) *InterfaceLogger {
	return &InterfaceLogger{
		Logger: baseLogger.WithLayer(InterfaceLayer),
	}
}

// LogHTTPRequest logs HTTP requests
func (i *InterfaceLogger) LogHTTPRequest(ctx context.Context, method string, path string, statusCode int, duration time.Duration, userAgent string, fields ...Field) {
	allFields := append(fields,
		String("method", method),
		String("path", path),
		Int("status_code", statusCode),
		Duration("duration", duration),
		String("user_agent", userAgent),
	)

	if statusCode >= 500 {
		i.Error(ctx, "HTTP request resulted in server error", allFields...)
	} else if statusCode >= 400 {
		i.Warn(ctx, "HTTP request resulted in client error", allFields...)
	} else if duration > 2*time.Second {
		i.Warn(ctx, "Slow HTTP request", allFields...)
	} else {
		i.Info(ctx, "HTTP request completed", allFields...)
	}
}

// LogAuthentication logs authentication attempts
func (i *InterfaceLogger) LogAuthentication(ctx context.Context, userID string, method string, success bool, reason string, fields ...Field) {
	allFields := append(fields,
		String("user_id", userID),
		String("auth_method", method),
		Bool("success", success),
		String("reason", reason),
	)

	if success {
		i.Info(ctx, "Authentication successful", allFields...)
	} else {
		i.Warn(ctx, "Authentication failed", allFields...)
	}
}

// LogAuthorization logs authorization checks
func (i *InterfaceLogger) LogAuthorization(ctx context.Context, userID string, resource string, action string, allowed bool, fields ...Field) {
	allFields := append(fields,
		String("user_id", userID),
		String("resource", resource),
		String("action", action),
		Bool("allowed", allowed),
	)

	if allowed {
		i.Debug(ctx, "Authorization granted", allFields...)
	} else {
		i.Warn(ctx, "Authorization denied", allFields...)
	}
}

// Convenience functions for creating layer-specific loggers using global factory
func GetDomainLogger() *DomainLogger {
	return NewDomainLogger(NewLogger())
}

func GetApplicationLogger() *ApplicationLogger {
	return NewApplicationLogger(NewLogger())
}

func GetInfrastructureLogger() *InfrastructureLogger {
	return NewInfrastructureLogger(NewLogger())
}

func GetInterfaceLogger() *InterfaceLogger {
	return NewInterfaceLogger(NewLogger())
}