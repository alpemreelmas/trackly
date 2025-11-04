package product

import (
	"context"
	"microservicetest/domain"
)

// Repository interface for products (keeping existing for backward compatibility)
type Repository interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
}

// VehicleRepository defines the interface for vehicle data operations
type VehicleRepository interface {
	// Basic CRUD operations
	GetVehicle(ctx context.Context, id string) (*domain.Vehicle, error)
	GetVehicleByVIN(ctx context.Context, vin string) (*domain.Vehicle, error)
	CreateVehicle(ctx context.Context, vehicle *domain.Vehicle) error
	UpdateVehicle(ctx context.Context, vehicle *domain.Vehicle) error
	DeleteVehicle(ctx context.Context, id string) error

	// Query operations
	GetVehiclesByOwner(ctx context.Context, ownerID string) ([]*domain.Vehicle, error)
	SearchVehicles(ctx context.Context, criteria map[string]interface{}) ([]*domain.Vehicle, error)

	// Insurance-related operations
	GetVehiclesWithExpiredInsurance(ctx context.Context) ([]*domain.Vehicle, error)
	GetVehiclesWithExpiringInsurance(ctx context.Context, days int) ([]*domain.Vehicle, error)
	UpdateInsurance(ctx context.Context, vehicleID string, insurance domain.InsuranceInfo) error

	// Document and picture operations
	AddDocument(ctx context.Context, vehicleID string, document domain.Document) error
	AddPicture(ctx context.Context, vehicleID string, picture domain.Picture) error
}
