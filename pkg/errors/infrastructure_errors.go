package errors

import "fmt"

// DatabaseError represents database operation failures
type DatabaseError struct {
	ErrorCode   ErrorCode
	Operation   string
	Table       string
	Query       string
	Cause       error
	IsRetryable bool
	Context     map[string]interface{}
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("database error in operation '%s' on table '%s': %v", e.Operation, e.Table, e.Cause)
}

func (e DatabaseError) Code() ErrorCode {
	return e.ErrorCode
}

func (e DatabaseError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"operation": e.Operation,
		"table":     e.Table,
		"retryable": e.IsRetryable,
	}
	if e.Query != "" {
		details["query"] = e.Query
	}
	if e.Cause != nil {
		details["cause"] = e.Cause.Error()
	}
	for k, v := range e.Context {
		details[k] = v
	}
	return details
}

func (e DatabaseError) Type() ErrorType {
	return ErrorTypeInfrastructure
}

func (e DatabaseError) WithContext(key string, value interface{}) BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

func (e DatabaseError) Retryable() bool {
	return e.IsRetryable
}

func (e DatabaseError) Unwrap() error {
	return e.Cause
}

// NewDatabaseError creates a new database error
func NewDatabaseError(operation, table string, cause error, retryable bool, context ...map[string]interface{}) *DatabaseError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	return &DatabaseError{
		ErrorCode:   CodeDatabaseError,
		Operation:   operation,
		Table:       table,
		Cause:       cause,
		IsRetryable: retryable,
		Context:     ctx,
	}
}

// NetworkError represents network operation failures
type NetworkError struct {
	ErrorCode   ErrorCode
	Service     string
	Endpoint    string
	Method      string
	Cause       error
	IsRetryable bool
	Context     map[string]interface{}
}

func (e NetworkError) Error() string {
	return fmt.Sprintf("network error calling service '%s' at endpoint '%s': %v", e.Service, e.Endpoint, e.Cause)
}

func (e NetworkError) Code() ErrorCode {
	return e.ErrorCode
}

func (e NetworkError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"service":   e.Service,
		"endpoint":  e.Endpoint,
		"method":    e.Method,
		"retryable": e.IsRetryable,
	}
	if e.Cause != nil {
		details["cause"] = e.Cause.Error()
	}
	for k, v := range e.Context {
		details[k] = v
	}
	return details
}

func (e NetworkError) Type() ErrorType {
	return ErrorTypeInfrastructure
}

func (e NetworkError) WithContext(key string, value interface{}) BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

func (e NetworkError) Retryable() bool {
	return e.IsRetryable
}

func (e NetworkError) Unwrap() error {
	return e.Cause
}

// NewNetworkError creates a new network error
func NewNetworkError(service, endpoint, method string, cause error, retryable bool, context ...map[string]interface{}) *NetworkError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	return &NetworkError{
		ErrorCode:   CodeNetworkError,
		Service:     service,
		Endpoint:    endpoint,
		Method:      method,
		Cause:       cause,
		IsRetryable: retryable,
		Context:     ctx,
	}
}

// ExternalServiceError represents external service integration failures
type ExternalServiceError struct {
	ErrorCode   ErrorCode
	Service     string
	Operation   string
	StatusCode  int
	Response    string
	Cause       error
	IsRetryable bool
	Context     map[string]interface{}
}

func (e ExternalServiceError) Error() string {
	return fmt.Sprintf("external service error in service '%s' operation '%s': %v", e.Service, e.Operation, e.Cause)
}

func (e ExternalServiceError) Code() ErrorCode {
	return e.ErrorCode
}

func (e ExternalServiceError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"service":     e.Service,
		"operation":   e.Operation,
		"status_code": e.StatusCode,
		"retryable":   e.IsRetryable,
	}
	if e.Response != "" {
		details["response"] = e.Response
	}
	if e.Cause != nil {
		details["cause"] = e.Cause.Error()
	}
	for k, v := range e.Context {
		details[k] = v
	}
	return details
}

func (e ExternalServiceError) Type() ErrorType {
	return ErrorTypeInfrastructure
}

func (e ExternalServiceError) WithContext(key string, value interface{}) BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

func (e ExternalServiceError) Retryable() bool {
	return e.IsRetryable
}

func (e ExternalServiceError) Unwrap() error {
	return e.Cause
}

// NewExternalServiceError creates a new external service error
func NewExternalServiceError(service, operation string, statusCode int, response string, cause error, retryable bool, context ...map[string]interface{}) *ExternalServiceError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	return &ExternalServiceError{
		ErrorCode:   CodeExternalServiceError,
		Service:     service,
		Operation:   operation,
		StatusCode:  statusCode,
		Response:    response,
		Cause:       cause,
		IsRetryable: retryable,
		Context:     ctx,
	}
}

// ConfigurationError represents configuration related errors
type ConfigurationError struct {
	ErrorCode ErrorCode
	Component string
	Parameter string
	Value     interface{}
	Reason    string
	Context   map[string]interface{}
}

func (e ConfigurationError) Error() string {
	return fmt.Sprintf("configuration error in component '%s' parameter '%s': %s", e.Component, e.Parameter, e.Reason)
}

func (e ConfigurationError) Code() ErrorCode {
	return e.ErrorCode
}

func (e ConfigurationError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"component": e.Component,
		"parameter": e.Parameter,
		"value":     e.Value,
		"reason":    e.Reason,
		"retryable": e.Retryable(),
	}
	for k, v := range e.Context {
		details[k] = v
	}
	return details
}

func (e ConfigurationError) Type() ErrorType {
	return ErrorTypeInfrastructure
}

func (e ConfigurationError) WithContext(key string, value interface{}) BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

func (e ConfigurationError) Retryable() bool {
	return false // Configuration errors are typically not retryable
}

// NewConfigurationError creates a new configuration error
func NewConfigurationError(component, parameter string, value interface{}, reason string, context ...map[string]interface{}) *ConfigurationError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	return &ConfigurationError{
		ErrorCode: CodeConfigurationError,
		Component: component,
		Parameter: parameter,
		Value:     value,
		Reason:    reason,
		Context:   ctx,
	}
}
