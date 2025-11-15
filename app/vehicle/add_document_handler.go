package vehicle

import (
	"microservicetest/app"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AddDocumentRequest struct {
	VehicleID string `param:"id" validate:"required"`
	//	Type           string  `formData:"type" validate:"required,oneof=insurance_policy insurance_card registration title inspection emission_test purchase_agreement service_record warranty receipt accident_report other"`
	//	Name           string  `formData:"name" validate:"required,min=1,max=200"`
	//	Description    string  `formData:"description" validate:"omitempty,max=500"`
	//	FileName       string  `formData:"file_name" validate:"required,min=1,max=255"`
	//	FileSize       int64   `formData:"file_size" validate:"required,gt=0"`
	//	MimeType       string  `formData:"mime_type" validate:"required"`
	//	ExpiryDate     *string `formData:"expiry_date" validate:"omitempty"`
	//	IssuedDate     *string `formData:"issued_date" validate:"omitempty"`
	//	IssuedBy       string  `formData:"issued_by" validate:"omitempty,max=100"`
	//	DocumentNumber string  `formData:"document_number" validate:"omitempty,max=100"`
	//	UploadedBy     string  `formData:"uploaded_by" validate:"required"`
}

type AddDocumentResponse struct {
	DocumentID string    `json:"document_id"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type AddDocumentHandler struct {
	repository     Repository
	storageService app.Storage
}

func NewAddDocumentHandler(repository Repository, storageService app.Storage) *AddDocumentHandler {
	return &AddDocumentHandler{
		repository:     repository,
		storageService: storageService,
	}
}

func (h *AddDocumentHandler) Handle(ctx *fiber.Ctx, req *AddDocumentRequest) (*AddDocumentResponse, error) {
	vehicleID := ctx.Params("id") // param:"id" mapping
	docType := ctx.FormValue("type")
	name := ctx.FormValue("name")
	description := ctx.FormValue("description")
	fileName := ctx.FormValue("file_name")
	fileSizeStr := ctx.FormValue("file_size")
	mimeType := ctx.FormValue("mime_type")
	uploadedBy := ctx.FormValue("uploaded_by")
	expiryDateStr := ctx.FormValue("expiry_date")
	issuedDateStr := ctx.FormValue("issued_date")
	issuedBy := ctx.FormValue("issued_by")
	documentNumber := ctx.FormValue("document_number")

	_, err := h.repository.GetVehicle(ctx.UserContext(), vehicleID)
	if err != nil {
		return nil, err
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, apperrors.ErrInternalServer.WithCause(err)
	}
	file, err := fileHeader.Open()
	if err != nil {
		return nil, apperrors.ErrInternalServer.WithCause(err)
	}
	defer file.Close()

	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		fileSize = fileHeader.Size
	}

	filenameUUID, _ := uuid.NewUUID()

	fileURL, err := h.storageService.Upload(ctx.UserContext(), file, filenameUUID.String(), mimeType)
	if err != nil {
		return nil, apperrors.ErrInternalServer.WithCause(err)
	}

	var expiryDate, issuedDate *time.Time
	if expiryDateStr != "" {
		t, err := time.Parse(time.RFC3339, expiryDateStr)
		if err != nil {
			return nil, apperrors.ErrInvalidFormat.WithDetails(map[string]string{
				"field":   "expiry_date",
				"message": "must be in RFC3339 format",
			})
		}
		expiryDate = &t
	}
	if issuedDateStr != "" {
		t, err := time.Parse(time.RFC3339, issuedDateStr)
		if err != nil {
			return nil, apperrors.ErrInvalidFormat.WithDetails(map[string]string{
				"field":   "issued_date",
				"message": "must be in RFC3339 format",
			})
		}
		issuedDate = &t
	}

	now := time.Now()
	document := domain.Document{
		ID:             domain.GenerateDocumentID(),
		Type:           domain.DocumentType(docType),
		Name:           name,
		Description:    description,
		FileURL:        fileURL,
		FileName:       fileName,
		FileSize:       fileSize,
		MimeType:       mimeType,
		IssuedBy:       issuedBy,
		DocumentNumber: documentNumber,
		UploadedAt:     now,
		UploadedBy:     uploadedBy,
		ExpiryDate:     expiryDate,
		IssuedDate:     issuedDate,
		IsVerified:     false,
	}

	if err := h.repository.AddDocument(ctx.UserContext(), vehicleID, document); err != nil {
		return nil, apperrors.ErrDatabaseQuery.WithCause(err).WithDetails(map[string]string{
			"operation": "add_document",
		})
	}

	return &AddDocumentResponse{
		DocumentID: document.ID,
		UploadedAt: document.UploadedAt,
	}, nil
}
