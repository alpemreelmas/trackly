package errors

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ErrorResponse represents the JSON structure for error responses
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains the error information
type ErrorDetail struct {
	Type    ErrorType `json:"type"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
}

// HandleError converts an error to an appropriate HTTP response
func HandleError(c *fiber.Ctx, err error) error {
	requestID := c.Locals("requestID")
	if requestID == nil {
		requestID = "unknown"
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		// Log the error with context
		logError(requestID.(string), c, appErr)

		// Return structured error response
		return c.Status(appErr.HTTPStatus).JSON(ErrorResponse{
			Error: ErrorDetail{
				Type:    appErr.Type,
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		})
	}

	// Handle unknown errors
	logError(requestID.(string), c, &AppError{
		Type:       ErrorTypeInternal,
		Code:       "UNKNOWN_ERROR",
		Message:    "An unexpected error occurred",
		HTTPStatus: 500,
		Cause:      err,
	})

	return c.Status(500).JSON(ErrorResponse{
		Error: ErrorDetail{
			Type:    ErrorTypeInternal,
			Code:    "UNKNOWN_ERROR",
			Message: "An unexpected error occurred",
		},
	})
}

// logError logs the error with appropriate level based on error type
func logError(requestID string, c *fiber.Ctx, appErr *AppError) {
	fields := []zap.Field{
		zap.String("request_id", requestID),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("error_type", string(appErr.Type)),
		zap.String("error_code", appErr.Code),
		zap.Int("http_status", appErr.HTTPStatus),
	}

	if appErr.Details != nil {
		fields = append(fields, zap.Any("error_details", appErr.Details))
	}

	if appErr.Cause != nil {
		fields = append(fields, zap.Error(appErr.Cause))
	}

	// Log with appropriate level based on error type
	switch appErr.Type {
	case ErrorTypeValidation, ErrorTypeNotFound, ErrorTypeUnauthorized, ErrorTypeForbidden, ErrorTypeConflict:
		// Client errors - log as info/warn
		zap.L().Warn("Client error", fields...)
	case ErrorTypeInternal, ErrorTypeDatabaseConnection, ErrorTypeDatabaseQuery:
		// Server errors - log as error
		zap.L().Error("Server error", fields...)
	case ErrorTypeExternal, ErrorTypeTimeout, ErrorTypeUnavailable:
		// External/infrastructure errors - log as warn
		zap.L().Warn("External service error", fields...)
	case ErrorTypeRateLimit:
		// Rate limiting - log as info
		zap.L().Info("Rate limit exceeded", fields...)
	default:
		// Unknown error type - log as error
		zap.L().Error("Unknown error type", fields...)
	}
}

// IsClientError checks if the error is a client error (4xx)
func IsClientError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatus >= 400 && appErr.HTTPStatus < 500
	}
	return false
}

// IsServerError checks if the error is a server error (5xx)
func IsServerError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatus >= 500
	}
	return false
}

// IsRetryable checks if the error indicates a retryable condition
func IsRetryable(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		switch appErr.Type {
		case ErrorTypeTimeout, ErrorTypeUnavailable, ErrorTypeExternal:
			return true
		case ErrorTypeRateLimit:
			return true
		default:
			return false
		}
	}
	return false
}