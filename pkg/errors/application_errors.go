package errors

import "fmt"

// EntityNotFoundError represents entity not found errors
type EntityNotFoundError struct {
	ErrorCode  ErrorCode
	EntityType string
	EntityID   string
	Context    map[string]interface{}
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("%s with ID '%s' not found", e.EntityType, e.EntityID)
}

func (e *EntityNotFoundError) Code() ErrorCode {
	return e.ErrorCode
}

func (e *EntityNotFoundError) Type() ErrorType {
	return ErrorTypeApplication
}

func (e *EntityNotFoundError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"entity_type": e.EntityType,
		"entity_id":   e.EntityID,
		"type":        e.Type(),
	}
	for k, v := range e.Context {
		details[k] = v
	}
	return details
}

func (e *EntityNotFoundError) WithContext(key string, value interface{}) BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewEntityNotFoundError creates a new entity not found error
func NewEntityNotFoundError(entityType, entityID string, context ...map[string]interface{}) *EntityNotFoundError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	return &EntityNotFoundError{
		ErrorCode:  CodeEntityNotFound,
		EntityType: entityType,
		EntityID:   entityID,
		Context:    ctx,
	}
}

// ConflictError represents resource conflict errors
type ConflictError struct {
	ErrorCode  ErrorCode
	Resource   string
	Reason     string
	ExistingID string
	Context    map[string]interface{}
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("conflict with existing %s: %s", e.Resource, e.Reason)
}

func (e ConflictError) Code() ErrorCode {
	return e.ErrorCode
}

func (e ConflictError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"resource":    e.Resource,
		"reason":      e.Reason,
		"existing_id": e.ExistingID,
	}
	for k, v := range e.Context {
		details[k] = v
	}
	return details
}

func (e ConflictError) Type() ErrorType {
	return ErrorTypeApplication
}

func (e ConflictError) WithContext(key string, value interface{}) BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewConflictError creates a new conflict error
func NewConflictError(resource, reason, existingID string, context ...map[string]interface{}) *ConflictError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	return &ConflictError{
		ErrorCode:  CodeResourceConflict,
		Resource:   resource,
		Reason:     reason,
		ExistingID: existingID,
		Context:    ctx,
	}
}

// UnauthorizedError represents unauthorized access errors
type UnauthorizedError struct {
	ErrorCode ErrorCode
	Operation string
	UserID    string
	Reason    string
	Context   map[string]interface{}
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized to perform operation '%s': %s", e.Operation, e.Reason)
}

func (e UnauthorizedError) Code() ErrorCode {
	return e.ErrorCode
}

func (e UnauthorizedError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"operation": e.Operation,
		"user_id":   e.UserID,
		"reason":    e.Reason,
	}
	for k, v := range e.Context {
		details[k] = v
	}
	return details
}

func (e UnauthorizedError) Type() ErrorType {
	return ErrorTypeApplication
}

func (e UnauthorizedError) WithContext(key string, value interface{}) BaseError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(operation, userID, reason string, context ...map[string]interface{}) *UnauthorizedError {
	var ctx map[string]interface{}
	if len(context) > 0 {
		ctx = context[0]
	}
	return &UnauthorizedError{
		ErrorCode: CodeUnauthorized,
		Operation: operation,
		UserID:    userID,
		Reason:    reason,
		Context:   ctx,
	}
}

// BusinessLogicError represents business logic violations at application layer
type BusinessLogicError struct {
	ErrorCode         ErrorCode
	Operation         string
	Reason            string
	AdditionalDetails map[string]interface{}
}

func (e BusinessLogicError) Error() string {
	return fmt.Sprintf("business logic error in operation '%s': %s", e.Operation, e.Reason)
}

func (e BusinessLogicError) Code() ErrorCode {
	return e.ErrorCode
}

func (e BusinessLogicError) Details() map[string]interface{} {
	details := map[string]interface{}{
		"operation": e.Operation,
		"reason":    e.Reason,
	}
	for k, v := range e.AdditionalDetails {
		details[k] = v
	}
	return details
}

func (e BusinessLogicError) Type() ErrorType {
	return ErrorTypeApplication
}

func (e BusinessLogicError) WithContext(key string, value interface{}) BaseError {
	if e.AdditionalDetails == nil {
		e.AdditionalDetails = make(map[string]interface{})
	}
	e.AdditionalDetails[key] = value
	return e
}

// NewBusinessLogicError creates a new business logic error
func NewBusinessLogicError(operation, reason string, details ...map[string]interface{}) *BusinessLogicError {
	var additionalDetails map[string]interface{}
	if len(details) > 0 {
		additionalDetails = details[0]
	}
	return &BusinessLogicError{
		ErrorCode:         CodeBusinessLogicError,
		Operation:         operation,
		Reason:            reason,
		AdditionalDetails: additionalDetails,
	}
}

// Convenience constructors for common application errors
func NewDuplicateEntryError(entityType, field string, value interface{}, existingID string) *ConflictError {
	return &ConflictError{
		ErrorCode:  CodeDuplicateEntry,
		Resource:   entityType,
		Reason:     fmt.Sprintf("%s '%v' already exists", field, value),
		ExistingID: existingID,
		Context: map[string]interface{}{
			field: value,
		},
	}
}

func NewResourceLockedError(entityType, entityID, reason string) *ConflictError {
	return &ConflictError{
		ErrorCode:  CodeResourceLocked,
		Resource:   entityType,
		Reason:     reason,
		ExistingID: entityID,
	}
}
