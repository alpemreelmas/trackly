package product

import (
	"context"
	"errors"
	"io"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type GetProductRequest struct {
	ID string `json:"id" param:"id"`
}

type GetProductResponse struct {
	Product *domain.Product `json:"product"`
}

type GetProductHandler struct {
	repository Repository
	httpClient *retryablehttp.Client
	breaker    *gobreaker.CircuitBreaker
	httpServer string
}

func NewGetProductHandler(repository Repository, httpClient *retryablehttp.Client, httpServer string) *GetProductHandler {
	// Configure the circuit breaker
	breakerSettings := gobreaker.Settings{
		Name:        "http-client",
		MaxRequests: 3,                // Number of requests allowed in half-open state
		Interval:    5 * time.Second,  // Time window for counting failures
		Timeout:     10 * time.Second, // Time to wait before switching from open to half-open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			zap.L().Info("Circuit breaker state changed",
				zap.String("name", name),
				zap.String("from", from.String()),
				zap.String("to", to.String()))
			// You could add logging here
		},
	}

	return &GetProductHandler{
		repository: repository,
		httpClient: httpClient,
		breaker:    gobreaker.NewCircuitBreaker(breakerSettings),
		httpServer: httpServer,
	}
}

func (h *GetProductHandler) Handle(ctx context.Context, req *GetProductRequest) (*GetProductResponse, error) {
	// Validate input
	if req.ID == "" {
		return nil, apperrors.NewValidationError("id", "Product ID is required")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, h.httpServer+"/random-error", nil)
	if err != nil {
		return nil, apperrors.NewExternalServiceError("http-server", err)
	}

	retryableReq, err := retryablehttp.FromRequest(httpReq)
	if err != nil {
		return nil, apperrors.NewExternalServiceError("http-server", err)
	}

	// Execute the HTTP request through the circuit breaker
	resp, err := h.breaker.Execute(func() (interface{}, error) {
		return h.httpClient.Do(retryableReq)
	})
	if err != nil {
		// Check if it's a circuit breaker error
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, apperrors.ErrExternalServiceUnavailable.WithCause(err).WithDetails(map[string]string{
				"service": "http-server",
				"reason":  "circuit breaker open",
			})
		}
		return nil, apperrors.NewExternalServiceError("http-server", err)
	}

	httpResp := resp.(*http.Response)
	defer httpResp.Body.Close()
	
	if httpResp.StatusCode >= 500 {
		return nil, apperrors.ErrExternalServiceUnavailable.WithDetails(map[string]any{
			"service":     "http-server",
			"status_code": httpResp.StatusCode,
		})
	}
	
	if _, err = io.ReadAll(httpResp.Body); err != nil {
		return nil, apperrors.NewExternalServiceError("http-server", err)
	}

	product, err := h.repository.GetProduct(ctx, req.ID)
	if err != nil {
		// Check if it's a not found error from repository
		if isNotFoundError(err) {
			return nil, apperrors.NewNotFoundError("product", req.ID)
		}
		// Assume other repository errors are database related
		return nil, apperrors.NewDatabaseError("get_product", err)
	}

	return &GetProductResponse{Product: product}, nil
}

// Helper function to check if error indicates resource not found
// This would depend on your repository implementation
func isNotFoundError(err error) bool {
	// Add your specific logic here based on your repository implementation
	// For example, if using Couchbase, you might check for specific error types
	return err.Error() == "document not found" || err.Error() == "key not found"
}
