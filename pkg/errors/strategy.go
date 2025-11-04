package errors

import (
	"context"
	"fmt"
	"strings"
)

// Error handling strategy examples for different layers

// =============================================================================
// REPOSITORY LAYER - Domain-specific error conversion
// =============================================================================

type ProductRepository struct {
	// your db client
}

func (r *ProductRepository) GetProduct(ctx context.Context, id string) (*Product, error) {
	// Validate input at repository level
	if id == "" {
		return nil, ErrInvalidID
	}

	// Database operation
	result, err := r.dbClient.Get(ctx, id)
	if err != nil {
		// Convert database-specific errors to domain errors immediately
		return nil, r.convertDBError("get_product", err)
	}

	return result, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *Product) error {
	if product == nil {
		return ErrInvalidInput.WithDetails(map[string]string{
			"field": "product",
			"issue": "cannot be nil",
		})
	}

	err := r.dbClient.Insert(ctx, product)
	if err != nil {
		return r.convertDBError("create_product", err)
	}

	return nil
}

// Convert database-specific errors to application errors
func (r *ProductRepository) convertDBError(operation string, err error) error {
	errMsg := strings.ToLower(err.Error())
	
	// Database-specific error mappings
	switch {
	case strings.Contains(errMsg, "document not found") || 
		 strings.Contains(errMsg, "key not found"):
		return ErrResourceNotFound.WithCause(err)
		
	case strings.Contains(errMsg, "duplicate") || 
		 strings.Contains(errMsg, "already exists") ||
		 strings.Contains(errMsg, "unique constraint"):
		return ErrResourceExists.WithCause(err)
		
	case strings.Contains(errMsg, "timeout") ||
		 strings.Contains(errMsg, "deadline exceeded"):
		return ErrRequestTimeout.WithCause(err)
		
	case strings.Contains(errMsg, "connection") ||
		 strings.Contains(errMsg, "network"):
		return ErrDatabaseConnection.WithCause(err)
		
	default:
		// Generic database error
		return NewDatabaseError(operation, err)
	}
}

// =============================================================================
// SERVICE LAYER - Business logic error handling
// =============================================================================

type ProductService struct {
	repo ProductRepository
}

func (s *ProductService) GetProductWithValidation(ctx context.Context, id string) (*Product, error) {
	// Business validation
	if !isValidProductID(id) {
		return nil, NewValidationError("id", "Product ID must be a valid UUID")
	}

	// Call repository - errors are already converted
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		// Repository errors are already properly typed, just add context if needed
		var appErr *AppError
		if errors.As(err, &appErr) {
			// Add business context to existing error
			if appErr.Type == ErrorTypeNotFound {
				return nil, appErr.WithDetails(map[string]string{
					"resource": "product",
					"id":       id,
					"context":  "user_request",
				})
			}
		}
		return nil, err // Pass through as-is
	}

	// Business logic validation
	if product.IsDeleted {
		return nil, ErrResourceNotFound.WithDetails(map[string]string{
			"reason": "product is deleted",
			"id":     id,
		})
	}

	return product, nil
}

// =============================================================================
// HANDLER LAYER - HTTP-specific error handling
// =============================================================================

func (h *ProductHandler) GetProduct(ctx context.Context, req *GetProductRequest) (*GetProductResponse, error) {
	// Input validation (HTTP-specific)
	if req.ID == "" {
		return nil, NewValidationError("id", "Product ID is required in URL path")
	}

	// Call service - errors are already properly typed
	product, err := h.service.GetProductWithValidation(ctx, req.ID)
	if err != nil {
		// Errors are already properly typed, just return them
		// The error handler middleware will convert to HTTP response
		return nil, err
	}

	return &GetProductResponse{Product: product}, nil
}

// =============================================================================
// EXTERNAL SERVICE LAYER - Third-party error handling
// =============================================================================

type ExternalAPIClient struct {
	httpClient HTTPClient
}

func (c *ExternalAPIClient) CallExternalAPI(ctx context.Context, endpoint string) (*APIResponse, error) {
	resp, err := c.httpClient.Get(ctx, endpoint)
	if err != nil {
		// Convert HTTP client errors to application errors
		return nil, c.convertHTTPError("external_api", err)
	}

	if resp.StatusCode >= 500 {
		return nil, ErrExternalServiceUnavailable.WithDetails(map[string]any{
			"service":     "external_api",
			"endpoint":    endpoint,
			"status_code": resp.StatusCode,
		})
	}

	if resp.StatusCode >= 400 {
		return nil, ErrExternalService.WithDetails(map[string]any{
			"service":     "external_api", 
			"endpoint":    endpoint,
			"status_code": resp.StatusCode,
		})
	}

	return parseResponse(resp)
}

func (c *ExternalAPIClient) convertHTTPError(service string, err error) error {
	errMsg := strings.ToLower(err.Error())
	
	switch {
	case strings.Contains(errMsg, "timeout") ||
		 strings.Contains(errMsg, "deadline exceeded"):
		return ErrExternalServiceTimeout.WithCause(err).WithDetails(map[string]string{
			"service": service,
		})
		
	case strings.Contains(errMsg, "connection refused") ||
		 strings.Contains(errMsg, "no such host"):
		return ErrExternalServiceUnavailable.WithCause(err).WithDetails(map[string]string{
			"service": service,
		})
		
	default:
		return NewExternalServiceError(service, err)
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func isValidProductID(id string) bool {
	// UUID validation logic
	return len(id) == 36 // Simplified
}

// Dummy types for example
type Product struct {
	ID        string
	Name      string
	IsDeleted bool
}

type HTTPClient interface {
	Get(ctx context.Context, url string) (*HTTPResponse, error)
}

type HTTPResponse struct {
	StatusCode int
	Body       []byte
}

type APIResponse struct {
	Data any
}

func parseResponse(resp *HTTPResponse) (*APIResponse, error) {
	return &APIResponse{}, nil
}