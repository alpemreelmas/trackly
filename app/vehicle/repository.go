package vehicle

import (
	"context"
	"microservicetest/domain"
)

// Repository defines the interface for vehicle data operations
type Repository interface {
	// Basic CRUD operations
	GetVehicle(ctx context.Context, id string) (*domain.Vehicle, error)
	GetVehicleByVIN(ctx context.Context, vin string) (*domain.Vehicle, error)
	GetVehiclesByOwner(ctx context.Context, ownerID string) ([]*domain.Vehicle, error)
	CreateVehicle(ctx context.Context, vehicle *domain.Vehicle) error
	UpdateVehicle(ctx context.Context, vehicle *domain.Vehicle) error
	DeleteVehicle(ctx context.Context, id string) error

	// Document operations
	AddDocument(ctx context.Context, vehicleID string, document domain.Document) error

	// Picture operations
	AddPicture(ctx context.Context, vehicleID string, picture domain.Picture) error
}
