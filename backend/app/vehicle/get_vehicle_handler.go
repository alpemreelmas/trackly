package vehicle

import (
	"context"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
	"microservicetest/pkg/validator"
)

type GetVehicleRequest struct {
	ID string `json:"id" param:"id" validate:"required"`
}

type GetVehicleResponse struct {
	Vehicle *domain.Vehicle `json:"vehicle"`
}

type GetVehicleHandler struct {
	repository Repository
}

func NewGetVehicleHandler(repository Repository) *GetVehicleHandler {
	return &GetVehicleHandler{
		repository: repository,
	}
}

func (h *GetVehicleHandler) Handle(ctx context.Context, req *GetVehicleRequest) (*GetVehicleResponse, error) {
	if err := validator.Validate(req); err != nil {
		return nil, apperrors.ErrInvalidInput.WithDetails(map[string]string{
			"validation": err.Error(),
		})
	}

	vehicle, err := h.repository.GetVehicle(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetVehicleResponse{Vehicle: vehicle}, nil
}