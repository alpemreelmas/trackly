package errors

import "net/http"

// Predefined error definitions for common scenarios

// Validation Errors
var (
	ErrInvalidInput = New(
		ErrorTypeValidation,
		"INVALID_INPUT",
		"Invalid input provided",
		http.StatusBadRequest,
	)

	ErrMissingRequiredField = New(
		ErrorTypeValidation,
		"MISSING_REQUIRED_FIELD",
		"Required field is missing",
		http.StatusBadRequest,
	)

	ErrInvalidFormat = New(
		ErrorTypeValidation,
		"INVALID_FORMAT",
		"Invalid format provided",
		http.StatusBadRequest,
	)

	ErrInvalidID = New(
		ErrorTypeValidation,
		"INVALID_ID",
		"Invalid ID format",
		http.StatusBadRequest,
	)
)

// Not Found Errors
var (
	ErrResourceNotFound = New(
		ErrorTypeNotFound,
		"RESOURCE_NOT_FOUND",
		"Requested resource not found",
		http.StatusNotFound,
	)

	ErrProductNotFound = New(
		ErrorTypeNotFound,
		"PRODUCT_NOT_FOUND",
		"Product not found",
		http.StatusNotFound,
	)

	ErrUserNotFound = New(
		ErrorTypeNotFound,
		"USER_NOT_FOUND",
		"User not found",
		http.StatusNotFound,
	)
)

// Authorization Errors
var (
	ErrUnauthorized = New(
		ErrorTypeUnauthorized,
		"UNAUTHORIZED",
		"Authentication required",
		http.StatusUnauthorized,
	)

	ErrInvalidToken = New(
		ErrorTypeUnauthorized,
		"INVALID_TOKEN",
		"Invalid or expired token",
		http.StatusUnauthorized,
	)

	ErrForbidden = New(
		ErrorTypeForbidden,
		"FORBIDDEN",
		"Access denied",
		http.StatusForbidden,
	)

	ErrInsufficientPermissions = New(
		ErrorTypeForbidden,
		"INSUFFICIENT_PERMISSIONS",
		"Insufficient permissions to perform this action",
		http.StatusForbidden,
	)
)

// Conflict Errors
var (
	ErrResourceExists = New(
		ErrorTypeConflict,
		"RESOURCE_EXISTS",
		"Resource already exists",
		http.StatusConflict,
	)

	ErrProductExists = New(
		ErrorTypeConflict,
		"PRODUCT_EXISTS",
		"Product already exists",
		http.StatusConflict,
	)

	ErrConcurrentModification = New(
		ErrorTypeConflict,
		"CONCURRENT_MODIFICATION",
		"Resource was modified by another request",
		http.StatusConflict,
	)
)

// Internal Errors
var (
	ErrInternalServer = New(
		ErrorTypeInternal,
		"INTERNAL_SERVER_ERROR",
		"Internal server error occurred",
		http.StatusInternalServerError,
	)

	ErrDatabaseConnection = New(
		ErrorTypeInternal,
		"DATABASE_CONNECTION_ERROR",
		"Database connection failed",
		http.StatusInternalServerError,
	)

	ErrDatabaseQuery = New(
		ErrorTypeInternal,
		"DATABASE_QUERY_ERROR",
		"Database query failed",
		http.StatusInternalServerError,
	)

	ErrConfigurationError = New(
		ErrorTypeInternal,
		"CONFIGURATION_ERROR",
		"Configuration error",
		http.StatusInternalServerError,
	)
)

// External Service Errors
var (
	ErrExternalService = New(
		ErrorTypeExternal,
		"EXTERNAL_SERVICE_ERROR",
		"External service error",
		http.StatusBadGateway,
	)

	ErrExternalServiceTimeout = New(
		ErrorTypeTimeout,
		"EXTERNAL_SERVICE_TIMEOUT",
		"External service timeout",
		http.StatusGatewayTimeout,
	)

	ErrExternalServiceUnavailable = New(
		ErrorTypeUnavailable,
		"EXTERNAL_SERVICE_UNAVAILABLE",
		"External service unavailable",
		http.StatusServiceUnavailable,
	)
)

// Rate Limiting Errors
var (
	ErrRateLimitExceeded = New(
		ErrorTypeRateLimit,
		"RATE_LIMIT_EXCEEDED",
		"Rate limit exceeded",
		http.StatusTooManyRequests,
	)
)

// Timeout Errors
var (
	ErrRequestTimeout = New(
		ErrorTypeTimeout,
		"REQUEST_TIMEOUT",
		"Request timeout",
		http.StatusRequestTimeout,
	)

	ErrOperationTimeout = New(
		ErrorTypeTimeout,
		"OPERATION_TIMEOUT",
		"Operation timeout",
		http.StatusRequestTimeout,
	)
)

// Service Unavailable Errors
var (
	ErrServiceUnavailable = New(
		ErrorTypeUnavailable,
		"SERVICE_UNAVAILABLE",
		"Service temporarily unavailable",
		http.StatusServiceUnavailable,
	)

	ErrMaintenanceMode = New(
		ErrorTypeUnavailable,
		"MAINTENANCE_MODE",
		"Service is in maintenance mode",
		http.StatusServiceUnavailable,
	)
)

// Helper functions to create specific errors with context

// NewValidationError creates a validation error with custom message
func NewValidationError(field, message string) *AppError {
	return ErrInvalidInput.WithDetails(map[string]string{
		"field":   field,
		"message": message,
	})
}

// NewNotFoundError creates a not found error for a specific resource
func NewNotFoundError(resource, id string) *AppError {
	return ErrResourceNotFound.WithDetails(map[string]string{
		"resource": resource,
		"id":       id,
	})
}

// NewConflictError creates a conflict error with custom message
func NewConflictError(resource, message string) *AppError {
	return ErrResourceExists.WithDetails(map[string]string{
		"resource": resource,
		"message":  message,
	})
}

// NewExternalServiceError creates an external service error
func NewExternalServiceError(service string, err error) *AppError {
	return ErrExternalService.WithCause(err).WithDetails(map[string]string{
		"service": service,
	})
}

// NewDatabaseError creates a database error
func NewDatabaseError(operation string, err error) *AppError {
	return ErrDatabaseQuery.WithCause(err).WithDetails(map[string]string{
		"operation": operation,
	})
}