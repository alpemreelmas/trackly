package vehicle

import (
	"context"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
	"microservicetest/pkg/validator"
	"strings"
	"time"
)

type CreateVehicleRequest struct {
	VIN          string  `json:"vin" validate:"required,min=17,max=17"`
	Make         string  `json:"make" validate:"required,min=1,max=50"`
	Model        string  `json:"model" validate:"required,min=1,max=50"`
	Year         int     `json:"year" validate:"required,gte=1900,lte=2100"`
	Color        string  `json:"color" validate:"omitempty,max=30"`
	LicensePlate string  `json:"license_plate" validate:"omitempty,max=20"`
	OwnerID      string  `json:"owner_id" validate:"required"`
	OwnerName    string  `json:"owner_name" validate:"required,min=1,max=100"`
	OwnerEmail   string  `json:"owner_email" validate:"required,email"`
	OwnerPhone   string  `json:"owner_phone" validate:"omitempty,min=10,max=20"`
	Transmission string  `json:"transmission" validate:"omitempty,oneof=manual automatic cvt"`
	FuelType     string  `json:"fuel_type" validate:"required,oneof=gasoline diesel electric hybrid lpg cng"`
	Mileage      int     `json:"mileage" validate:"omitempty,gte=0"`
	CreatedBy    string  `json:"created_by" validate:"required"`
}

type CreateVehicleResponse struct {
	ID        string    `json:"id"`
	VIN       string    `json:"vin"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateVehicleHandler struct {
	repository Repository
}

func NewCreateVehicleHandler(repository Repository) *CreateVehicleHandler {
	return &CreateVehicleHandler{
		repository: repository,
	}
}

func (h *CreateVehicleHandler) Handle(ctx context.Context, req *CreateVehicleRequest) (*CreateVehicleResponse, error) {
	if err := validator.Validate(req); err != nil {
		return nil, apperrors.ErrInvalidInput.WithDetails(map[string]string{
			"validation": err.Error(),
		})
	}

	// Check if vehicle with VIN already exists
	existing, err := h.repository.GetVehicleByVIN(ctx, req.VIN)
	if err == nil && existing != nil {
		return nil, apperrors.ErrResourceExists.WithDetails(map[string]string{
			"resource": "vehicle",
			"vin":      req.VIN,
		})
	}

	now := time.Now()
	vehicle := &domain.Vehicle{
		ID:           domain.GenerateVehicleID(),
		VIN:          strings.ToUpper(strings.TrimSpace(req.VIN)),
		Make:         strings.TrimSpace(req.Make),
		Model:        strings.TrimSpace(req.Model),
		Year:         req.Year,
		Color:        strings.TrimSpace(req.Color),
		LicensePlate: strings.ToUpper(strings.TrimSpace(req.LicensePlate)),
		OwnerID:      req.OwnerID,
		OwnerName:    strings.TrimSpace(req.OwnerName),
		OwnerEmail:   strings.ToLower(strings.TrimSpace(req.OwnerEmail)),
		OwnerPhone:   strings.TrimSpace(req.OwnerPhone),
		Transmission: req.Transmission,
		FuelType:     domain.FuelType(req.FuelType),
		Mileage:      req.Mileage,
		Status:       domain.VehicleStatusActive,
		Documents:    make([]domain.Document, 0),
		Pictures:     make([]domain.Picture, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    req.CreatedBy,
		UpdatedBy:    req.CreatedBy,
	}

	if err := h.repository.CreateVehicle(ctx, vehicle); err != nil {
		return nil, apperrors.ErrDatabaseQuery.WithCause(err).WithDetails(map[string]string{
			"operation": "create_vehicle",
		})
	}

	return &CreateVehicleResponse{
		ID:        vehicle.ID,
		VIN:       vehicle.VIN,
		CreatedAt: vehicle.CreatedAt,
	}, nil
}