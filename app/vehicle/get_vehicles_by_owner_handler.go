package vehicle

import (
	"context"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
	"microservicetest/pkg/validator"
)

type GetVehiclesByOwnerRequest struct {
	OwnerID string `json:"owner_id" param:"owner_id" validate:"required"`
}

type GetVehiclesByOwnerResponse struct {
	Vehicles []*domain.Vehicle `json:"vehicles"`
	Count    int               `json:"count"`
}

type GetVehiclesByOwnerHandler struct {
	repository Repository
}

func NewGetVehiclesByOwnerHandler(repository Repository) *GetVehiclesByOwnerHandler {
	return &GetVehiclesByOwnerHandler{
		repository: repository,
	}
}

func (h *GetVehiclesByOwnerHandler) Handle(ctx context.Context, req *GetVehiclesByOwnerRequest) (*GetVehiclesByOwnerResponse, error) {
	if err := validator.Validate(req); err != nil {
		return nil, apperrors.ErrInvalidInput.WithDetails(map[string]string{
			"validation": err.Error(),
		})
	}

	vehicles, err := h.repository.GetVehiclesByOwner(ctx, req.OwnerID)
	if err != nil {
		return nil, apperrors.ErrDatabaseQuery.WithCause(err).WithDetails(map[string]string{
			"operation": "get_vehicles_by_owner",
		})
	}

	if vehicles == nil {
		vehicles = make([]*domain.Vehicle, 0)
	}

	return &GetVehiclesByOwnerResponse{
		Vehicles: vehicles,
		Count:    len(vehicles),
	}, nil
}
