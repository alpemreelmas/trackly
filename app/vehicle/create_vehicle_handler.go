package product

import (
	"context"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
	"strings"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name string `json:"name"`
}

type CreateProductResponse struct {
	ID string `json:"id"`
}

type CreateProductHandler struct {
	repository Repository
}

func NewCreateProductHandler(repository Repository) *CreateProductHandler {
	return &CreateProductHandler{
		repository: repository,
	}
}

func (h *CreateProductHandler) Handle(ctx context.Context, req *CreateProductRequest) (*CreateProductResponse, error) {
	// Validate input
	if strings.TrimSpace(req.Name) == "" {
		return nil, apperrors.NewValidationError("name", "Product name is required")
	}

	if len(req.Name) > 100 {
		return nil, apperrors.NewValidationError("name", "Product name must be less than 100 characters")
	}

	productID := uuid.New().String()

	product := &domain.Product{
		ID:   productID,
		Name: strings.TrimSpace(req.Name),
	}

	err := h.repository.CreateProduct(ctx, product)
	if err != nil {
		// Check if it's a duplicate key error
		if isDuplicateError(err) {
			return nil, apperrors.NewConflictError("product", "Product with this name already exists")
		}
		// Assume other repository errors are database related
		return nil, apperrors.NewDatabaseError("create_product", err)
	}

	return &CreateProductResponse{ID: product.ID}, nil
}

// Helper function to check if error indicates duplicate resource
// This would depend on your repository implementation
func isDuplicateError(err error) bool {
	// Add your specific logic here based on your repository implementation
	// For example, if using Couchbase or SQL, you might check for specific error codes
	errorMsg := strings.ToLower(err.Error())
	return strings.Contains(errorMsg, "duplicate") || 
		   strings.Contains(errorMsg, "already exists") ||
		   strings.Contains(errorMsg, "unique constraint")
}
