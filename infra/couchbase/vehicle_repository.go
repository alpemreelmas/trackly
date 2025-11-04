package couchbase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
	"go.uber.org/zap"

	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
)

type VehicleRepository struct {
	cluster *gocb.Cluster
	bucket  *gocb.Bucket
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
		Tracer:     tracer,
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
	// Set timestamps
	now := time.Now()
	vehicle.CreatedAt = now
	vehicle.UpdatedAt = now

	// Use atomic operations to ensure VIN uniqueness and update indexes
	vinKey := "vin::" + vehicle.VIN
	vinRef := map[string]string{"vehicle_id": vehicle.ID}

	// Create both documents atomically using transactions
	err := r.cluster.Transactions().Run(func(attempt *gocb.TransactionAttempt) error {
		// Try to insert VIN reference first (this will fail if VIN exists)
		_, err := attempt.Insert(r.collection, vinKey, vinRef)
		if err != nil {
			return err
		}

		// Insert the vehicle document
		_, err = attempt.Insert(r.collection, vehicle.ID, vehicle)
		if err != nil {
			return err
		}

		return nil
	}, &gocb.TransactionOptions{
		Timeout: 10 * time.Second,
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

// convertDBError converts Couchbase errors to application errors
func (r *VehicleRepository) convertDBError(operation string, err error) error {
	errMsg := strings.ToLower(err.Error())
	
	switch {
	case errors.Is(err, gocb.ErrDocumentNotFound):
		return apperrors.ErrResourceNotFound.WithCause(err)
		
	case errors.Is(err, gocb.ErrDocumentExists):
		return apperrors.ErrResourceExists.WithCause(err)
		
	case strings.Contains(errMsg, "timeout") ||
		 strings.Contains(errMsg, "deadline exceeded"):
		return apperrors.ErrRequestTimeout.WithCause(err)
		
	case strings.Contains(errMsg, "connection") ||
		 strings.Contains(errMsg, "network") ||
		 strings.Contains(errMsg, "cluster"):
		return apperrors.ErrDatabaseConnection.WithCause(err)
		
	case strings.Contains(errMsg, "authentication") ||
		 strings.Contains(errMsg, "unauthorized"):
		return apperrors.ErrUnauthorized.WithCause(err)
		
	default:
		return apperrors.NewDatabaseError(operation, err)
	}
}
