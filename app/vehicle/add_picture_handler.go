package vehicle

import (
	"context"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
	"microservicetest/pkg/validator"
	"strings"
	"time"
)

type AddPictureRequest struct {
	VehicleID    string  `json:"vehicle_id" param:"id" validate:"required"`
	Type         string  `json:"type" validate:"required,oneof=exterior_front exterior_back exterior_left exterior_right interior_front interior_back dashboard engine trunk wheels damage accident other"`
	Title        string  `json:"title" validate:"required,min=1,max=200"`
	Description  string  `json:"description" validate:"omitempty,max=500"`
	URL          string  `json:"url" validate:"required,url"`
	ThumbnailURL string  `json:"thumbnail_url" validate:"omitempty,url"`
	FileName     string  `json:"file_name" validate:"required,min=1,max=255"`
	FileSize     int64   `json:"file_size" validate:"required,gt=0"`
	Width        int     `json:"width" validate:"required,gt=0"`
	Height       int     `json:"height" validate:"required,gt=0"`
	MimeType     string  `json:"mime_type" validate:"required"`
	TakenAt      *string `json:"taken_at" validate:"omitempty"`
	UploadedBy   string  `json:"uploaded_by" validate:"required"`
	IsMain       bool    `json:"is_main"`
}

type AddPictureResponse struct {
	PictureID  string    `json:"picture_id"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type AddPictureHandler struct {
	repository Repository
}

func NewAddPictureHandler(repository Repository) *AddPictureHandler {
	return &AddPictureHandler{
		repository: repository,
	}
}

func (h *AddPictureHandler) Handle(ctx context.Context, req *AddPictureRequest) (*AddPictureResponse, error) {
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
	picture := domain.Picture{
		ID:           domain.GeneratePictureID(),
		Type:         domain.PictureType(req.Type),
		Title:        strings.TrimSpace(req.Title),
		Description:  strings.TrimSpace(req.Description),
		URL:          req.URL,
		ThumbnailURL: req.ThumbnailURL,
		FileName:     strings.TrimSpace(req.FileName),
		FileSize:     req.FileSize,
		Width:        req.Width,
		Height:       req.Height,
		MimeType:     req.MimeType,
		UploadedAt:   now,
		UploadedBy:   req.UploadedBy,
		IsMain:       req.IsMain,
		SortOrder:    0,
	}

	if req.TakenAt != nil && *req.TakenAt != "" {
		takenAt, err := time.Parse(time.RFC3339, *req.TakenAt)
		if err != nil {
			return nil, apperrors.ErrInvalidFormat.WithDetails(map[string]string{
				"field":   "taken_at",
				"message": "must be in RFC3339 format",
			})
		}
		picture.TakenAt = &takenAt
	}

	if err := h.repository.AddPicture(ctx, req.VehicleID, picture); err != nil {
		return nil, apperrors.ErrDatabaseQuery.WithCause(err).WithDetails(map[string]string{
			"operation": "add_picture",
		})
	}

	return &AddPictureResponse{
		PictureID:  picture.ID,
		UploadedAt: picture.UploadedAt,
	}, nil
}
