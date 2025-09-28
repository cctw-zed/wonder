package errors

import "fmt"

// ValidationError represents validation failures in domain entities
type ValidationError struct {
	ErrorCode ErrorCode
	Field     string
	Value     interface{}
	Message   string
	Context   map[string]interface{}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

func (e *ValidationError) Code() ErrorCode {
	return e.ErrorCode
}

func (e *ValidationError) Type() ErrorType {
	return ErrorTypeDomain
}

func (e *ValidationError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"field":   e.Field,
		"value":   e.Value,
		"message": e.Message,
		"type":    e.Type(),
	}
	for k, v := range e.Context {
		details[k] = v
	}
	return details
}

func (e *ValidationError) WithContext(key string, value interface{}) BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewValidationError creates a new validation error with specific code
func NewValidationError(code ErrorCode, field string, value interface{}, message string, context ...map[string]interface{}) *ValidationError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	return &ValidationError{
		ErrorCode: code,
		Field:     field,
		Value:     value,
		Message:   message,
		Context:   ctx,
	}
}

// Convenience constructors for common validation errors
func NewRequiredFieldError(field string, value interface{}) *ValidationError {
	return NewValidationError(CodeRequiredField, field, value, fmt.Sprintf("%s is required", field))
}

func NewInvalidFormatError(field string, value interface{}, expectedFormat string) *ValidationError {
	return NewValidationError(CodeInvalidFormat, field, value,
		fmt.Sprintf("invalid format for %s, expected: %s", field, expectedFormat))
}

func NewInvalidValueError(field string, value interface{}, reason string) *ValidationError {
	return NewValidationError(CodeInvalidValue, field, value,
		fmt.Sprintf("invalid value for %s: %s", field, reason))
}

func NewOutOfRangeError(field string, value interface{}, min, max interface{}) *ValidationError {
	return NewValidationError(CodeOutOfRange, field, value,
		fmt.Sprintf("%s must be between %v and %v", field, min, max))
}

// DomainRuleError represents business rule violations
type DomainRuleError struct {
	ErrorCode ErrorCode
	Rule      string
	Context   map[string]interface{}
	Message   string
}

func (e *DomainRuleError) Error() string {
	return fmt.Sprintf("domain rule violation '%s': %s", e.Rule, e.Message)
}

func (e *DomainRuleError) Code() ErrorCode {
	return e.ErrorCode
}

func (e *DomainRuleError) Type() ErrorType {
	return ErrorTypeDomain
}

func (e *DomainRuleError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"rule":    e.Rule,
		"message": e.Message,
		"type":    e.Type(),
	}
	for k, v := range e.Context {
		details[k] = v
	}
	return details
}

func (e *DomainRuleError) WithContext(key string, value interface{}) BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewDomainRuleError creates a new domain rule error
func NewDomainRuleError(code ErrorCode, rule, message string, context ...map[string]interface{}) *DomainRuleError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	return &DomainRuleError{
		ErrorCode: code,
		Rule:      rule,
		Message:   message,
		Context:   ctx,
	}
}

// Convenience constructors for domain rule errors
func NewBusinessRuleError(rule, message string, context ...map[string]interface{}) *DomainRuleError {
	return NewDomainRuleError(CodeBusinessRuleError, rule, message, context...)
}

func NewInvariantViolationError(rule, message string, context ...map[string]interface{}) *DomainRuleError {
	return NewDomainRuleError(CodeInvariantViolation, rule, message, context...)
}

// InvalidStateError represents invalid domain object state
type InvalidStateError struct {
	ErrorCode ErrorCode
	Entity    string
	State     string
	Message   string
	Context   map[string]interface{}
}

func (e *InvalidStateError) Error() string {
	return fmt.Sprintf("invalid state for %s: %s", e.Entity, e.Message)
}

func (e *InvalidStateError) Code() ErrorCode {
	return e.ErrorCode
}

func (e *InvalidStateError) Type() ErrorType {
	return ErrorTypeDomain
}

func (e *InvalidStateError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"entity":  e.Entity,
		"state":   e.State,
		"message": e.Message,
		"type":    e.Type(),
	}
	for k, v := range e.Context {
		details[k] = v
	}
	return details
}

func (e *InvalidStateError) WithContext(key string, value interface{}) BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewInvalidStateError creates a new invalid state error
func NewInvalidStateError(code ErrorCode, entity, state, message string, context ...map[string]interface{}) *InvalidStateError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	return &InvalidStateError{
		ErrorCode: code,
		Entity:    entity,
		State:     state,
		Message:   message,
		Context:   ctx,
	}
}

// Convenience constructors for state errors
func NewStateTransitionError(entity, fromState, toState, reason string) *InvalidStateError {
	return NewInvalidStateError(CodeStateTransition, entity, fromState,
		fmt.Sprintf("cannot transition from %s to %s: %s", fromState, toState, reason))
}

func NewPreconditionError(entity, state, precondition string) *InvalidStateError {
	return NewInvalidStateError(CodePreconditionError, entity, state,
		fmt.Sprintf("precondition not met: %s", precondition))
}
