package couchbase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/gocb/v2"
	"go.uber.org/zap"

	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
)

type VehicleRepository struct {
	cluster    *gocb.Cluster
	bucket     *gocb.Bucket
	collection *gocb.Collection
}

func NewVehicleRepository(couchbaseUrl string, username string, password string) *VehicleRepository {
	cluster, err := gocb.Connect(couchbaseUrl, gocb.ClusterOptions{
		TimeoutsConfig: gocb.TimeoutsConfig{
			ConnectTimeout: 10 * time.Second,
			KVTimeout:      5 * time.Second,
			QueryTimeout:   10 * time.Second,
		},
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
		Transcoder: gocb.NewJSONTranscoder(),
	})
	if err != nil {
		zap.L().Fatal("Failed to connect to couchbase", zap.Error(err))
	}

	bucket := cluster.Bucket("vehicles")
	bucket.WaitUntilReady(10*time.Second, &gocb.WaitUntilReadyOptions{})

	collection := bucket.DefaultCollection()

	return &VehicleRepository{
		cluster:    cluster,
		bucket:     bucket,
		collection: collection,
	}
}

// GetVehicle retrieves a vehicle by ID
func (r *VehicleRepository) GetVehicle(ctx context.Context, id string) (*domain.Vehicle, error) {
	if id == "" {
		return nil, apperrors.ErrInvalidID
	}

	data, err := r.collection.Get(id, &gocb.GetOptions{
		Timeout: 5 * time.Second,
		Context: ctx,
	})
	if err != nil {
		return nil, r.convertDBError("get_vehicle", err)
	}

	var vehicle domain.Vehicle
	if err := data.Content(&vehicle); err != nil {
		return nil, apperrors.NewDatabaseError("decode_vehicle", err)
	}

	return &vehicle, nil
}

// GetVehicleByVIN retrieves a vehicle by VIN using lookup operation
func (r *VehicleRepository) GetVehicleByVIN(ctx context.Context, vin string) (*domain.Vehicle, error) {
	vinKey := "vin::" + vin

	result, err := r.collection.Get(vinKey, &gocb.GetOptions{
		Timeout: 5 * time.Second,
		Context: ctx,
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, apperrors.NewNotFoundError("vehicle", vin)
		}
		return nil, r.convertDBError("get_vehicle_by_vin", err)
	}

	var vehicleRef struct {
		VehicleID string `json:"vehicle_id"`
	}
	if err := result.Content(&vehicleRef); err != nil {
		return nil, apperrors.NewDatabaseError("decode_vin_reference", err)
	}

	// Now get the actual vehicle document
	return r.GetVehicle(ctx, vehicleRef.VehicleID)
}

// CreateVehicle creates a new vehicle using atomic operations
func (r *VehicleRepository) CreateVehicle(ctx context.Context, vehicle *domain.Vehicle) error {
	now := time.Now()
	vehicle.CreatedAt = now
	vehicle.UpdatedAt = now

	vinKey := "vin::" + vehicle.VIN
	vinRef := map[string]string{"vehicle_id": vehicle.ID}

	_, err := r.cluster.Transactions().Run(func(attempt *gocb.TransactionAttemptContext) error {
		_, err := attempt.Insert(r.collection, vinKey, vinRef)
		if err != nil {
			return err
		}

		_, err = attempt.Insert(r.collection, vehicle.ID, vehicle)
		if err != nil {
			return err
		}

		return nil
	}, &gocb.TransactionOptions{
		Timeout:         10 * time.Second,
		DurabilityLevel: gocb.DurabilityLevelMajority,
	})

	if err != nil {
		if errors.Is(err, gocb.ErrDocumentExists) {
			return apperrors.NewConflictError("vehicle", fmt.Sprintf("Vehicle with VIN %s already exists", vehicle.VIN))
		}
		return r.convertDBError("create_vehicle", err)
	}
	return nil
}

// UpdateVehicle updates an existing vehicle
func (r *VehicleRepository) UpdateVehicle(ctx context.Context, vehicle *domain.Vehicle) error {
	vehicle.UpdatedAt = time.Now()

	_, err := r.collection.Replace(vehicle.ID, vehicle, &gocb.ReplaceOptions{
		Timeout: 5 * time.Second,
		Context: ctx,
	})
	if err != nil {
		return r.convertDBError("update_vehicle", err)
	}

	return nil
}

// DeleteVehicle soft deletes a vehicle by setting status to inactive
func (r *VehicleRepository) DeleteVehicle(ctx context.Context, id string) error {

	// Get the vehicle first
	vehicle, err := r.GetVehicle(ctx, id)
	if err != nil {
		return err
	}

	// Set status to inactive and update timestamp
	vehicle.Status = domain.VehicleStatusInactive
	vehicle.UpdatedAt = time.Now()

	return r.UpdateVehicle(ctx, vehicle)
}

// GetVehiclesByOwner retrieves all vehicles for a specific owner
func (r *VehicleRepository) GetVehiclesByOwner(ctx context.Context, ownerID string) ([]*domain.Vehicle, error) {
	if ownerID == "" {
		return nil, apperrors.ErrInvalidID
	}

	query := `
		SELECT v.* 
		FROM vehicles v 
		WHERE v.owner_id = $1 
		AND v.status != 'inactive'
		ORDER BY v.created_at DESC
	`

	result, err := r.cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{ownerID},
		Timeout:              10 * time.Second,
		Context:              ctx,
	})
	if err != nil {
		return nil, r.convertDBError("get_vehicles_by_owner", err)
	}
	defer result.Close()

	var vehicles []*domain.Vehicle
	for result.Next() {
		var vehicle domain.Vehicle
		if err := result.Row(&vehicle); err != nil {
			zap.L().Error("Failed to decode vehicle row", zap.Error(err))
			continue
		}
		vehicles = append(vehicles, &vehicle)
	}

	if err := result.Err(); err != nil {
		return nil, r.convertDBError("get_vehicles_by_owner_iteration", err)
	}

	return vehicles, nil
}

// AddDocument adds a document to a vehicle
func (r *VehicleRepository) AddDocument(ctx context.Context, vehicleID string, document domain.Document) error {
	vehicle, err := r.GetVehicle(ctx, vehicleID)
	if err != nil {
		return err
	}

	// Add the document to the vehicle
	if err := vehicle.AddDocument(document); err != nil {
		return apperrors.ErrInvalidInput.WithDetails(map[string]string{
			"error": err.Error(),
		})
	}

	// Update the vehicle
	return r.UpdateVehicle(ctx, vehicle)
}

// AddPicture adds a picture to a vehicle
func (r *VehicleRepository) AddPicture(ctx context.Context, vehicleID string, picture domain.Picture) error {
	vehicle, err := r.GetVehicle(ctx, vehicleID)
	if err != nil {
		return err
	}

	// Add the picture to the vehicle
	if err := vehicle.AddPicture(picture); err != nil {
		return apperrors.ErrInvalidInput.WithDetails(map[string]string{
			"error": err.Error(),
		})
	}

	// Update the vehicle
	return r.UpdateVehicle(ctx, vehicle)
}

// convertDBError converts Couchbase errors to application errors
func (r *VehicleRepository) convertDBError(operation string, err error) error {
	var timeoutErr *gocb.TimeoutError

	switch {
	case errors.Is(err, gocb.ErrDocumentNotFound):
		return apperrors.ErrResourceNotFound.WithCause(err)

	case errors.Is(err, gocb.ErrDocumentExists):
		return apperrors.ErrResourceExists.WithCause(err)

	case errors.As(err, &timeoutErr):
		return apperrors.ErrRequestTimeout.WithCause(timeoutErr)
	default:
		// If we can’t categorize it, just wrap it.
		return apperrors.NewDatabaseError(operation, err)
	}
}
