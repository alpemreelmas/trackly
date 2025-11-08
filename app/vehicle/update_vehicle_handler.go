package vehicle

import (
	"context"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
	"microservicetest/pkg/validator"
	"strings"
)

type UpdateVehicleRequest struct {
	ID           string  `json:"id" param:"id" validate:"required"`
	Color        *string `json:"color" validate:"omitempty,max=30"`
	LicensePlate *string `json:"license_plate" validate:"omitempty,max=20"`
	OwnerName    *string `json:"owner_name" validate:"omitempty,min=1,max=100"`
	OwnerEmail   *string `json:"owner_email" validate:"omitempty,email"`
	OwnerPhone   *string `json:"owner_phone" validate:"omitempty,min=10,max=20"`
	Transmission *string `json:"transmission" validate:"omitempty,oneof=manual automatic cvt"`
	Mileage      *int    `json:"mileage" validate:"omitempty,gte=0"`
	Status       *string `json:"status" validate:"omitempty,oneof=active inactive sold scrapped stolen accident"`
	UpdatedBy    string  `json:"updated_by" validate:"required"`
}

type UpdateVehicleResponse struct {
	Vehicle *domain.Vehicle `json:"vehicle"`
}

type UpdateVehicleHandler struct {
	repository Repository
}

func NewUpdateVehicleHandler(repository Repository) *UpdateVehicleHandler {
	return &UpdateVehicleHandler{
		repository: repository,
	}
}

func (h *UpdateVehicleHandler) Handle(ctx context.Context, req *UpdateVehicleRequest) (*UpdateVehicleResponse, error) {
	if err := validator.Validate(req); err != nil {
		return nil, apperrors.ErrInvalidInput.WithDetails(map[string]string{
			"validation": err.Error(),
		})
	}

	vehicle, err := h.repository.GetVehicle(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	if req.Color != nil {
		vehicle.Color = strings.TrimSpace(*req.Color)
	}
	if req.LicensePlate != nil {
		vehicle.LicensePlate = strings.ToUpper(strings.TrimSpace(*req.LicensePlate))
	}
	if req.OwnerName != nil {
		vehicle.OwnerName = strings.TrimSpace(*req.OwnerName)
	}
	if req.OwnerEmail != nil {
		vehicle.OwnerEmail = strings.ToLower(strings.TrimSpace(*req.OwnerEmail))
	}
	if req.OwnerPhone != nil {
		vehicle.OwnerPhone = strings.TrimSpace(*req.OwnerPhone)
	}
	if req.Transmission != nil {
		vehicle.Transmission = *req.Transmission
	}
	if req.Mileage != nil {
		vehicle.Mileage = *req.Mileage
	}
	if req.Status != nil {
		vehicle.Status = domain.VehicleStatus(*req.Status)
	}

	vehicle.UpdateTimestamp(req.UpdatedBy)

	if err := h.repository.UpdateVehicle(ctx, vehicle); err != nil {
		return nil, apperrors.ErrDatabaseQuery.WithCause(err).WithDetails(map[string]string{
			"operation": "update_vehicle",
		})
	}

	return &UpdateVehicleResponse{Vehicle: vehicle}, nil
}
