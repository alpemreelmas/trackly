package vehicle

import (
	"microservicetest/app"
	apperrors "microservicetest/pkg/errors"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type DownloadDocumentRequest struct {
	VehicleID  string `param:"id" validate:"required"`
	DocumentID string `param:"doc_id" validate:"required"`
}

type DownloadDocumentHandler struct {
	repository     Repository
	storageService app.Storage
}

func NewDownloadDocumentHandler(repository Repository, storageService app.Storage) *DownloadDocumentHandler {
	return &DownloadDocumentHandler{
		repository:     repository,
		storageService: storageService,
	}
}

func (h *DownloadDocumentHandler) Handle(ctx *fiber.Ctx, req *DownloadDocumentRequest) error {

	// Get vehicle
	vehicle, err := h.repository.GetVehicle(ctx.UserContext(), req.VehicleID)
	if err != nil {
		return err
	}

	// Find document
	var document *struct {
		FileURL  string
		FileName string
		MimeType string
	}

	for _, doc := range vehicle.Documents {
		if doc.ID == req.DocumentID {
			document = &struct {
				FileURL  string
				FileName string
				MimeType string
			}{
				FileURL:  doc.FileURL,
				FileName: doc.FileName,
				MimeType: doc.MimeType,
			}
			break
		}
	}

	if document == nil {
		return apperrors.ErrResourceNotFound.WithDetails(map[string]string{
			"resource": "document",
			"id":       req.DocumentID,
		})
	}

	// Extract filename from URL
	parsedURL, err := url.Parse(document.FileURL)
	if err != nil {
		return apperrors.ErrInternalServer.WithCause(err)
	}

	// Get the last part of the path (filename)
	pathParts := strings.Split(parsedURL.Path, "/")
	blobFilename := pathParts[len(pathParts)-1]

	// Download from Azure Blob
	data, contentType, err := h.storageService.Download(ctx.UserContext(), blobFilename)
	if err != nil {
		return apperrors.ErrInternalServer.WithCause(err).WithDetails(map[string]string{
			"operation": "download_blob",
		})
	}

	// Use stored content type if available, otherwise use downloaded one
	if document.MimeType != "" {
		contentType = document.MimeType
	}

	// Set headers
	ctx.Set("Content-Type", contentType)
	ctx.Set("Content-Disposition", "attachment; filename=\""+document.FileName+"\"")

	// Send file
	return ctx.Send(data)
}
