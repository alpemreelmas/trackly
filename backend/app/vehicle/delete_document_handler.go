package vehicle

import (
	"microservicetest/app"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DeleteDocumentRequest struct {
	VehicleID  string `param:"id" validate:"required"`
	DocumentID string `param:"doc_id" validate:"required"`
}

type DeleteDocumentResponse struct {
	Message string `json:"message"`
}

type DeleteDocumentHandler struct {
	repository Repository
	storage    app.Storage
}

func NewDeleteDocumentHandler(repository Repository, storage app.Storage) *DeleteDocumentHandler {
	return &DeleteDocumentHandler{
		repository: repository,
		storage:    storage,
	}
}

func (h *DeleteDocumentHandler) Handle(ctx *fiber.Ctx, req *DeleteDocumentRequest) (*DeleteDocumentResponse, error) {
	vehicleID := ctx.Params("id")
	documentID := ctx.Params("doc_id")

	// Get vehicle to find document FileURL
	vehicle, err := h.repository.GetVehicle(ctx.UserContext(), vehicleID)
	if err != nil {
		return nil, err
	}

	// Find document and extract blob filename
	var blobFilename string
	for _, doc := range vehicle.Documents {
		if doc.ID == documentID {
			parts := strings.Split(doc.FileURL, "/")
			if len(parts) > 0 {
				blobFilename = parts[len(parts)-1]
			}
			break
		}
	}

	// Delete from database
	if err := h.repository.DeleteDocument(ctx.UserContext(), vehicleID, documentID); err != nil {
		return nil, err
	}

	// Delete from Azure Blob Storage if we found the filename
	if blobFilename != "" {
		if err := h.storage.Remove(ctx.UserContext(), blobFilename); err != nil {
			zap.L().Error("Failed to delete blob from storage",
				zap.String("filename", blobFilename),
				zap.Error(err))
		}
	}

	return &DeleteDocumentResponse{
		Message: "Document deleted successfully",
	}, nil
}
