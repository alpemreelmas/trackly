package vehicle

import (
	"context"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
	"microservicetest/pkg/validator"
	"strings"
	"time"
)

type AddDocumentRequest struct {
	VehicleID      string  `json:"vehicle_id" param:"id" validate:"required"`
	Type           string  `json:"type" validate:"required,oneof=insurance_policy insurance_card registration title inspection emission_test purchase_agreement service_record warranty receipt accident_report other"`
	Name           string  `json:"name" validate:"required,min=1,max=200"`
	Description    string  `json:"description" validate:"omitempty,max=500"`
	FileURL        string  `json:"file_url" validate:"required,url"`
	FileName       string  `json:"file_name" validate:"required,min=1,max=255"`
	FileSize       int64   `json:"file_size" validate:"required,gt=0"`
	MimeType       string  `json:"mime_type" validate:"required"`
	ExpiryDate     *string `json:"expiry_date" validate:"omitempty"`
	IssuedDate     *string `json:"issued_date" validate:"omitempty"`
	IssuedBy       string  `json:"issued_by" validate:"omitempty,max=100"`
	DocumentNumber string  `json:"document_number" validate:"omitempty,max=100"`
	UploadedBy     string  `json:"uploaded_by" validate:"required"`
}

type AddDocumentResponse struct {
	DocumentID string    `json:"document_id"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type AddDocumentHandler struct {
	repository Repository
}

func NewAddDocumentHandler(repository Repository) *AddDocumentHandler {
	return &AddDocumentHandler{
		repository: repository,
	}
}

func (h *AddDocumentHandler) Handle(ctx context.Context, req *AddDocumentRequest) (*AddDocumentResponse, error) {
	if err := validator.Validate(req); err != nil {
		return nil, apperrors.ErrInvalidInput.WithDetails(map[string]string{
			"validation": err.Error(),
		})
	}

	// Verify vehicle exists
	_, err := h.repository.GetVehicle(ctx, req.VehicleID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	document := domain.Document{
		ID:             domain.GenerateDocumentID(),
		Type:           domain.DocumentType(req.Type),
		Name:           strings.TrimSpace(req.Name),
		Description:    strings.TrimSpace(req.Description),
		FileURL:        req.FileURL,
		FileName:       strings.TrimSpace(req.FileName),
		FileSize:       req.FileSize,
		MimeType:       req.MimeType,
		IssuedBy:       strings.TrimSpace(req.IssuedBy),
		DocumentNumber: strings.TrimSpace(req.DocumentNumber),
		UploadedAt:     now,
		UploadedBy:     req.UploadedBy,
		IsVerified:     false,
	}

	if req.ExpiryDate != nil && *req.ExpiryDate != "" {
		expiryDate, err := time.Parse(time.RFC3339, *req.ExpiryDate)
		if err != nil {
			return nil, apperrors.ErrInvalidFormat.WithDetails(map[string]string{
				"field":   "expiry_date",
				"message": "must be in RFC3339 format",
			})
		}
		document.ExpiryDate = &expiryDate
	}

	if req.IssuedDate != nil && *req.IssuedDate != "" {
		issuedDate, err := time.Parse(time.RFC3339, *req.IssuedDate)
		if err != nil {
			return nil, apperrors.ErrInvalidFormat.WithDetails(map[string]string{
				"field":   "issued_date",
				"message": "must be in RFC3339 format",
			})
		}
		document.IssuedDate = &issuedDate
	}

	if err := h.repository.AddDocument(ctx, req.VehicleID, document); err != nil {
		return nil, apperrors.ErrDatabaseQuery.WithCause(err).WithDetails(map[string]string{
			"operation": "add_document",
		})
	}

	return &AddDocumentResponse{
		DocumentID: document.ID,
		UploadedAt: document.UploadedAt,
	}, nil
}
