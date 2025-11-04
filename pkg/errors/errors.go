package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents a custom application error with additional context
type AppError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	Code       string    `json:"code"`
	HTTPStatus int       `json:"http_status"`
	Details    any       `json:"details,omitempty"`
	Cause      error     `json:"-"`
}

// ErrorType represents the category of error
type ErrorType string

const (
	ErrorTypeValidation    ErrorType = "validation"
	ErrorTypeNotFound      ErrorType = "not_found"
	ErrorTypeUnauthorized  ErrorType = "unauthorized"
	ErrorTypeForbidden     ErrorType = "forbidden"
	ErrorTypeConflict      ErrorType = "conflict"
	ErrorTypeInternal      ErrorType = "internal"
	ErrorTypeExternal      ErrorType = "external"
	ErrorTypeTimeout       ErrorType = "timeout"
	ErrorTypeRateLimit     ErrorType = "rate_limit"
	ErrorTypeBadRequest    ErrorType = "bad_request"
	ErrorTypeUnavailable   ErrorType = "unavailable"
)

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying error for errors.Is and errors.As
func (e *AppError) Unwrap() error {
	return e.Cause
}

// Is implements error comparison for errors.Is
func (e *AppError) Is(target error) bool {
	var appErr *AppError
	if errors.As(target, &appErr) {
		return e.Type == appErr.Type && e.Code == appErr.Code
	}
	return false
}

// WithDetails adds additional details to the error
func (e *AppError) WithDetails(details any) *AppError {
	e.Details = details
	return e
}

// WithCause adds a cause to the error
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// New creates a new AppError
func New(errorType ErrorType, code, message string, httpStatus int) *AppError {
	return &AppError{
		Type:       errorType,
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Wrap wraps an existing error with AppError context
func Wrap(err error, errorType ErrorType, code, message string, httpStatus int) *AppError {
	return &AppError{
		Type:       errorType,
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Cause:      err,
	}
}

// GetHTTPStatus returns the HTTP status code for the error
func GetHTTPStatus(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}

// GetErrorType returns the error type
func GetErrorType(err error) ErrorType {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type
	}
	return ErrorTypeInternal
}

// GetErrorCode returns the error code
func GetErrorCode(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return "UNKNOWN_ERROR"
}