package couchbase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"

	"github.com/couchbase/gocb/v2"
)

// VehicleRepository implements the vehicle.Repository interface
type VehicleRepository struct {
	cluster    *gocb.Cluster
	bucket     *gocb.Bucket
	collection *gocb.Collection
}

// NewVehicleRepository creates a new vehicle repository
func NewVehicleRepository(cluster *gocb.Cluster, bucketName, scopeName, collectionName string) (*VehicleRepository, error) {
	bucket := cluster.Bucket(bucketName)
	
	err := bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		return nil, fmt.Errorf("bucket not ready: %w", err)
	}

	collection := bucket.Scope(scopeName).Collection(collectionName)

	return &VehicleRepository{
		cluster:    cluster,
		bucket:     bucket,
		collection: collection,
	}, nil
}

// GetVehicle retrieves a vehicle by ID
func (r *VehicleRepository) GetVehicle(ctx context.Context, id string) (*domain.Vehicle, error) {
	result, err := r.collection.Get(id, &gocb.GetOptions{
		Context: ctx,
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, apperrors.NewNotFoundError("vehicle", id)
		}
		return nil, apperrors.NewDatabaseError("get_vehicle", err)
	}

	var vehicle domain.Vehicle
	if err := result.Content(&vehicle); err != nil {
		return nil, apperrors.NewDatabaseError("decode_vehicle", err)
	}

	return &vehicle, nil
}

// GetVehicleByVIN retrieves a vehicle by VIN
func (r *VehicleRepository) GetVehicleByVIN(ctx context.Context, vin string) (*domain.Vehicle, error) {
	query := `
		SELECT META().id, v.*
		FROM vehicles v
		WHERE v.vin = $1
		LIMIT 1
	`

	result, err := r.cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{vin},
		Context:              ctx,
	})
	if err != nil {
		return nil, apperrors.NewDatabaseError("query_vehicle_by_vin", err)
	}
	defer result.Close()

	if !result.Next() {
		return nil, apperrors.NewNotFoundError("vehicle", vin)
	}

	var vehicle domain.Vehicle
	if err := result.Row(&vehicle); err != nil {
		return nil, apperrors.NewDatabaseError("decode_vehicle", err)
	}

	return &vehicle, nil
}

// CreateVehicle creates a new vehicle
func (r *VehicleRepository) CreateVehicle(ctx context.Context, vehicle *domain.Vehicle) error {
	_, err := r.collection.Insert(vehicle.ID, vehicle, &gocb.InsertOptions{
		Context: ctx,
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentExists) {
			return apperrors.NewConflictError("vehicle", "vehicle already exists")
		}
		return apperrors.NewDatabaseError("create_vehicle", err)
	}

	return nil
}

// UpdateVehicle updates an existing vehicle
func (r *VehicleRepository) UpdateVehicle(ctx context.Context, vehicle *domain.Vehicle) error {
	_, err := r.collection.Replace(vehicle.ID, vehicle, &gocb.ReplaceOptions{
		Context: ctx,
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return apperrors.NewNotFoundError("vehicle", vehicle.ID)
		}
		return apperrors.NewDatabaseError("update_vehicle", err)
	}

	return nil
}

// DeleteVehicle deletes a vehicle
func (r *VehicleRepository) DeleteVehicle(ctx context.Context, id string) error {
	_, err := r.collection.Remove(id, &gocb.RemoveOptions{
		Context: ctx,
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return apperrors.NewNotFoundError("vehicle", id)
		}
		return apperrors.NewDatabaseError("delete_vehicle", err)
	}

	return nil
}

// GetVehiclesByOwner retrieves all vehicles for an owner
func (r *VehicleRepository) GetVehiclesByOwner(ctx context.Context, ownerID string) ([]*domain.Vehicle, error) {
	query := `
		SELECT META().id, v.*
		FROM vehicles v
		WHERE v.owner_id = $1
		ORDER BY v.created_at DESC
	`

	result, err := r.cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{ownerID},
		Context:              ctx,
	})
	if err != nil {
		return nil, apperrors.NewDatabaseError("query_vehicles_by_owner", err)
	}
	defer result.Close()

	var vehicles []*domain.Vehicle
	for result.Next() {
		var vehicle domain.Vehicle
		if err := result.Row(&vehicle); err != nil {
			return nil, apperrors.NewDatabaseError("decode_vehicle", err)
		}
		vehicles = append(vehicles, &vehicle)
	}

	return vehicles, nil
}

// SearchVehicles searches vehicles based on criteria
func (r *VehicleRepository) SearchVehicles(ctx context.Context, criteria map[string]interface{}) ([]*domain.Vehicle, error) {
	// Build dynamic query based on criteria
	// This is a simplified example
	query := "SELECT META().id, v.* FROM vehicles v WHERE 1=1"
	params := make([]interface{}, 0)

	// Add criteria to query (simplified)
	if make, ok := criteria["make"].(string); ok && make != "" {
		query += " AND v.make = $" + fmt.Sprintf("%d", len(params)+1)
		params = append(params, make)
	}

	result, err := r.cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: params,
		Context:              ctx,
	})
	if err != nil {
		return nil, apperrors.NewDatabaseError("search_vehicles", err)
	}
	defer result.Close()

	var vehicles []*domain.Vehicle
	for result.Next() {
		var vehicle domain.Vehicle
		if err := result.Row(&vehicle); err != nil {
			return nil, apperrors.NewDatabaseError("decode_vehicle", err)
		}
		vehicles = append(vehicles, &vehicle)
	}

	return vehicles, nil
}

// GetVehiclesWithExpiredInsurance retrieves vehicles with expired insurance
func (r *VehicleRepository) GetVehiclesWithExpiredInsurance(ctx context.Context) ([]*domain.Vehicle, error) {
	query := `
		SELECT META().id, v.*
		FROM vehicles v
		WHERE v.insurance.end_date < NOW_STR()
		AND v.status = 'active'
	`

	result, err := r.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
	})
	if err != nil {
		return nil, apperrors.NewDatabaseError("query_expired_insurance", err)
	}
	defer result.Close()

	var vehicles []*domain.Vehicle
	for result.Next() {
		var vehicle domain.Vehicle
		if err := result.Row(&vehicle); err != nil {
			return nil, apperrors.NewDatabaseError("decode_vehicle", err)
		}
		vehicles = append(vehicles, &vehicle)
	}

	return vehicles, nil
}

// GetVehiclesWithExpiringInsurance retrieves vehicles with insurance expiring within days
func (r *VehicleRepository) GetVehiclesWithExpiringInsurance(ctx context.Context, days int) ([]*domain.Vehicle, error) {
	query := `
		SELECT META().id, v.*
		FROM vehicles v
		WHERE v.insurance.end_date BETWEEN NOW_STR() AND DATE_ADD_STR(NOW_STR(), $1, 'day')
		AND v.status = 'active'
	`

	result, err := r.cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{days},
		Context:              ctx,
	})
	if err != nil {
		return nil, apperrors.NewDatabaseError("query_expiring_insurance", err)
	}
	defer result.Close()

	var vehicles []*domain.Vehicle
	for result.Next() {
		var vehicle domain.Vehicle
		if err := result.Row(&vehicle); err != nil {
			return nil, apperrors.NewDatabaseError("decode_vehicle", err)
		}
		vehicles = append(vehicles, &vehicle)
	}

	return vehicles, nil
}

// UpdateInsurance updates vehicle insurance information
func (r *VehicleRepository) UpdateInsurance(ctx context.Context, vehicleID string, insurance domain.InsuranceInfo) error {
	_, err := r.collection.MutateIn(vehicleID, []gocb.MutateInSpec{
		gocb.ReplaceSpec("insurance", insurance, nil),
		gocb.ReplaceSpec("updated_at", time.Now(), nil),
	}, &gocb.MutateInOptions{
		Context: ctx,
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return apperrors.NewNotFoundError("vehicle", vehicleID)
		}
		return apperrors.NewDatabaseError("update_insurance", err)
	}

	return nil
}

// AddDocument adds a document to a vehicle
func (r *VehicleRepository) AddDocument(ctx context.Context, vehicleID string, document domain.Document) error {
	_, err := r.collection.MutateIn(vehicleID, []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("documents", []interface{}{document}, nil),
		gocb.ReplaceSpec("updated_at", time.Now(), nil),
	}, &gocb.MutateInOptions{
		Context: ctx,
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return apperrors.NewNotFoundError("vehicle", vehicleID)
		}
		return apperrors.NewDatabaseError("add_document", err)
	}

	return nil
}

// AddPicture adds a picture to a vehicle
func (r *VehicleRepository) AddPicture(ctx context.Context, vehicleID string, picture domain.Picture) error {
	_, err := r.collection.MutateIn(vehicleID, []gocb.MutateInSpec{
		gocb.ArrayAppendSpec("pictures", []interface{}{picture}, nil),
		gocb.ReplaceSpec("updated_at", time.Now(), nil),
	}, &gocb.MutateInOptions{
		Context: ctx,
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return apperrors.NewNotFoundError("vehicle", vehicleID)
		}
		return apperrors.NewDatabaseError("add_picture", err)
	}

	return nil
}
