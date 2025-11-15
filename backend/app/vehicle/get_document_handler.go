package vehicle

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type DocumentFilter struct {
	Type           string
	IsVerified     *bool
	IsExpired      *bool
	UploadedBy     string
	IssuedBy       string
	DocumentNumber string
}

type GetDocumentsRequest struct {
	VehicleID string `param:"id" validate:"required"`
	// Query filters
	Type           string `query:"type" validate:"omitempty,oneof=insurance_policy insurance_card registration title inspection emission_test purchase_agreement service_record warranty receipt accident_report other"`
	IsVerified     string `query:"is_verified"`     // "true", "false", or empty
	IsExpired      string `query:"is_expired"`      // "true", "false", or empty
	UploadedBy     string `query:"uploaded_by"`
	IssuedBy       string `query:"issued_by"`
	DocumentNumber string `query:"document_number"`
}

type DocumentResponse struct {
	ID             string     `json:"id"`
	Type           string     `json:"type"`
	Name           string     `json:"name"`
	Description    string     `json:"description,omitempty"`
	FileURL        string     `json:"file_url"`
	FileName       string     `json:"file_name"`
	FileSize       int64      `json:"file_size"`
	MimeType       string     `json:"mime_type"`
	IssuedBy       string     `json:"issued_by,omitempty"`
	DocumentNumber string     `json:"document_number,omitempty"`
	UploadedAt     time.Time  `json:"uploaded_at"`
	UploadedBy     string     `json:"uploaded_by,omitempty"`
	ExpiryDate     *time.Time `json:"expiry_date,omitempty"`
	IssuedDate     *time.Time `json:"issued_date,omitempty"`
	IsVerified     bool       `json:"is_verified"`
	IsExpired      bool       `json:"is_expired"`
}

type GetDocumentsResponse struct {
	Documents []DocumentResponse `json:"documents"`
	Total     int                `json:"total"`
}

type GetDocumentsHandler struct {
	repository Repository
}

func NewGetDocumentsHandler(repository Repository) *GetDocumentsHandler {
	return &GetDocumentsHandler{
		repository: repository,
	}
}

func (h *GetDocumentsHandler) Handle(ctx *fiber.Ctx, req *GetDocumentsRequest) (*GetDocumentsResponse, error) {
	vehicleID := ctx.Params("id")

	// Verify vehicle exists
	_, err := h.repository.GetVehicle(ctx.UserContext(), vehicleID)
	if err != nil {
		return nil, err
	}

	// Convert string booleans to *bool for filter
	var isVerified, isExpired *bool
	if req.IsVerified != "" {
		val := req.IsVerified == "true"
		isVerified = &val
	}
	if req.IsExpired != "" {
		val := req.IsExpired == "true"
		isExpired = &val
	}

	// Build filter from request
	filter := DocumentFilter{
		Type:           req.Type,
		IsVerified:     isVerified,
		IsExpired:      isExpired,
		UploadedBy:     req.UploadedBy,
		IssuedBy:       req.IssuedBy,
		DocumentNumber: req.DocumentNumber,
	}

	// Query documents with filters at DB level
	docs, err := h.repository.GetDocuments(ctx.UserContext(), vehicleID, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	documents := make([]DocumentResponse, 0, len(docs))
	now := time.Now()
	
	for _, doc := range docs {
		isExpired := doc.ExpiryDate != nil && doc.ExpiryDate.Before(now)
		documents = append(documents, DocumentResponse{
			ID:             doc.ID,
			Type:           string(doc.Type),
			Name:           doc.Name,
			Description:    doc.Description,
			FileURL:        doc.FileURL,
			FileName:       doc.FileName,
			FileSize:       doc.FileSize,
			MimeType:       doc.MimeType,
			IssuedBy:       doc.IssuedBy,
			DocumentNumber: doc.DocumentNumber,
			UploadedAt:     doc.UploadedAt,
			UploadedBy:     doc.UploadedBy,
			ExpiryDate:     doc.ExpiryDate,
			IssuedDate:     doc.IssuedDate,
			IsVerified:     doc.IsVerified,
			IsExpired:      isExpired,
		})
	}

	return &GetDocumentsResponse{
		Documents: documents,
		Total:     len(documents),
	}, nil
}
