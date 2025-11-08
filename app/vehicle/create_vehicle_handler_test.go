package vehicle

import (
	"context"
	"errors"
	"microservicetest/domain"
	apperrors "microservicetest/pkg/errors"
	"testing"
	"time"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	GetVehicleFunc          func(ctx context.Context, id string) (*domain.Vehicle, error)
	GetVehicleByVINFunc     func(ctx context.Context, vin string) (*domain.Vehicle, error)
	CreateVehicleFunc       func(ctx context.Context, vehicle *domain.Vehicle) error
	UpdateVehicleFunc       func(ctx context.Context, vehicle *domain.Vehicle) error
	DeleteVehicleFunc       func(ctx context.Context, id string) error
	GetVehiclesByOwnerFunc  func(ctx context.Context, ownerID string) ([]*domain.Vehicle, error)
	SearchVehiclesFunc      func(ctx context.Context, criteria map[string]interface{}) ([]*domain.Vehicle, error)
	GetVehiclesWithExpiredInsuranceFunc func(ctx context.Context) ([]*domain.Vehicle, error)
	GetVehiclesWithExpiringInsuranceFunc func(ctx context.Context, days int) ([]*domain.Vehicle, error)
	UpdateInsuranceFunc     func(ctx context.Context, vehicleID string, insurance domain.InsuranceInfo) error
	AddDocumentFunc         func(ctx context.Context, vehicleID string, document domain.Document) error
	AddPictureFunc          func(ctx context.Context, vehicleID string, picture domain.Picture) error
}

func (m *MockRepository) GetVehicle(ctx context.Context, id string) (*domain.Vehicle, error) {
	if m.GetVehicleFunc != nil {
		return m.GetVehicleFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *MockRepository) GetVehicleByVIN(ctx context.Context, vin string) (*domain.Vehicle, error) {
	if m.GetVehicleByVINFunc != nil {
		return m.GetVehicleByVINFunc(ctx, vin)
	}
	return nil, apperrors.ErrResourceNotFound
}

func (m *MockRepository) CreateVehicle(ctx context.Context, vehicle *domain.Vehicle) error {
	if m.CreateVehicleFunc != nil {
		return m.CreateVehicleFunc(ctx, vehicle)
	}
	return nil
}

func (m *MockRepository) UpdateVehicle(ctx context.Context, vehicle *domain.Vehicle) error {
	if m.UpdateVehicleFunc != nil {
		return m.UpdateVehicleFunc(ctx, vehicle)
	}
	return nil
}

func (m *MockRepository) DeleteVehicle(ctx context.Context, id string) error {
	if m.DeleteVehicleFunc != nil {
		return m.DeleteVehicleFunc(ctx, id)
	}
	return nil
}

func (m *MockRepository) GetVehiclesByOwner(ctx context.Context, ownerID string) ([]*domain.Vehicle, error) {
	if m.GetVehiclesByOwnerFunc != nil {
		return m.GetVehiclesByOwnerFunc(ctx, ownerID)
	}
	return nil, nil
}

func (m *MockRepository) SearchVehicles(ctx context.Context, criteria map[string]interface{}) ([]*domain.Vehicle, error) {
	if m.SearchVehiclesFunc != nil {
		return m.SearchVehiclesFunc(ctx, criteria)
	}
	return nil, nil
}

func (m *MockRepository) GetVehiclesWithExpiredInsurance(ctx context.Context) ([]*domain.Vehicle, error) {
	if m.GetVehiclesWithExpiredInsuranceFunc != nil {
		return m.GetVehiclesWithExpiredInsuranceFunc(ctx)
	}
	return nil, nil
}

func (m *MockRepository) GetVehiclesWithExpiringInsurance(ctx context.Context, days int) ([]*domain.Vehicle, error) {
	if m.GetVehiclesWithExpiringInsuranceFunc != nil {
		return m.GetVehiclesWithExpiringInsuranceFunc(ctx, days)
	}
	return nil, nil
}

func (m *MockRepository) UpdateInsurance(ctx context.Context, vehicleID string, insurance domain.InsuranceInfo) error {
	if m.UpdateInsuranceFunc != nil {
		return m.UpdateInsuranceFunc(ctx, vehicleID, insurance)
	}
	return nil
}

func (m *MockRepository) AddDocument(ctx context.Context, vehicleID string, document domain.Document) error {
	if m.AddDocumentFunc != nil {
		return m.AddDocumentFunc(ctx, vehicleID, document)
	}
	return nil
}

func (m *MockRepository) AddPicture(ctx context.Context, vehicleID string, picture domain.Picture) error {
	if m.AddPictureFunc != nil {
		return m.AddPictureFunc(ctx, vehicleID, picture)
	}
	return nil
}

func TestCreateVehicleHandler_Success(t *testing.T) {
	mockRepo := &MockRepository{
		GetVehicleByVINFunc: func(ctx context.Context, vin string) (*domain.Vehicle, error) {
			return nil, apperrors.ErrResourceNotFound
		},
		CreateVehicleFunc: func(ctx context.Context, vehicle *domain.Vehicle) error {
			return nil
		},
	}

	handler := NewCreateVehicleHandler(mockRepo)

	req := &CreateVehicleRequest{
		VIN:          "1HGBH41JXMN109186",
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		Color:        "Silver",
		LicensePlate: "ABC123",
		OwnerID:      "owner-123",
		OwnerName:    "John Doe",
		OwnerEmail:   "john@example.com",
		OwnerPhone:   "+1234567890",
		Transmission: "automatic",
		FuelType:     "gasoline",
		Mileage:      15000,
		CreatedBy:    "admin-user",
	}

	resp, err := handler.Handle(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	if resp.ID == "" {
		t.Error("Expected vehicle ID to be set")
	}

	if resp.VIN != "1HGBH41JXMN109186" {
		t.Errorf("Expected VIN to be 1HGBH41JXMN109186, got %s", resp.VIN)
	}

	if resp.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestCreateVehicleHandler_ValidationError_MissingVIN(t *testing.T) {
	mockRepo := &MockRepository{}
	handler := NewCreateVehicleHandler(mockRepo)

	req := &CreateVehicleRequest{
		Make:       "Toyota",
		Model:      "Camry",
		Year:       2023,
		OwnerID:    "owner-123",
		OwnerName:  "John Doe",
		OwnerEmail: "john@example.com",
		FuelType:   "gasoline",
		CreatedBy:  "admin-user",
	}

	_, err := handler.Handle(context.Background(), req)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}

	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("Expected AppError, got %T", err)
	}

	if appErr.Type != apperrors.ErrorTypeValidation {
		t.Errorf("Expected validation error type, got %s", appErr.Type)
	}
}

func TestCreateVehicleHandler_ValidationError_InvalidVINLength(t *testing.T) {
	mockRepo := &MockRepository{}
	handler := NewCreateVehicleHandler(mockRepo)

	req := &CreateVehicleRequest{
		VIN:        "SHORT",
		Make:       "Toyota",
		Model:      "Camry",
		Year:       2023,
		OwnerID:    "owner-123",
		OwnerName:  "John Doe",
		OwnerEmail: "john@example.com",
		FuelType:   "gasoline",
		CreatedBy:  "admin-user",
	}

	_, err := handler.Handle(context.Background(), req)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestCreateVehicleHandler_ValidationError_InvalidEmail(t *testing.T) {
	mockRepo := &MockRepository{}
	handler := NewCreateVehicleHandler(mockRepo)

	req := &CreateVehicleRequest{
		VIN:        "1HGBH41JXMN109186",
		Make:       "Toyota",
		Model:      "Camry",
		Year:       2023,
		OwnerID:    "owner-123",
		OwnerName:  "John Doe",
		OwnerEmail: "not-an-email",
		FuelType:   "gasoline",
		CreatedBy:  "admin-user",
	}

	_, err := handler.Handle(context.Background(), req)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestCreateVehicleHandler_DuplicateVIN(t *testing.T) {
	existingVehicle := &domain.Vehicle{
		ID:  "VEH_123",
		VIN: "1HGBH41JXMN109186",
	}

	mockRepo := &MockRepository{
		GetVehicleByVINFunc: func(ctx context.Context, vin string) (*domain.Vehicle, error) {
			return existingVehicle, nil
		},
	}

	handler := NewCreateVehicleHandler(mockRepo)

	req := &CreateVehicleRequest{
		VIN:        "1HGBH41JXMN109186",
		Make:       "Toyota",
		Model:      "Camry",
		Year:       2023,
		OwnerID:    "owner-123",
		OwnerName:  "John Doe",
		OwnerEmail: "john@example.com",
		FuelType:   "gasoline",
		CreatedBy:  "admin-user",
	}

	_, err := handler.Handle(context.Background(), req)

	if err == nil {
		t.Fatal("Expected conflict error, got nil")
	}

	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("Expected AppError, got %T", err)
	}

	if appErr.Type != apperrors.ErrorTypeConflict {
		t.Errorf("Expected conflict error type, got %s", appErr.Type)
	}
}

func TestCreateVehicleHandler_DatabaseError(t *testing.T) {
	mockRepo := &MockRepository{
		GetVehicleByVINFunc: func(ctx context.Context, vin string) (*domain.Vehicle, error) {
			return nil, apperrors.ErrResourceNotFound
		},
		CreateVehicleFunc: func(ctx context.Context, vehicle *domain.Vehicle) error {
			return errors.New("database connection failed")
		},
	}

	handler := NewCreateVehicleHandler(mockRepo)

	req := &CreateVehicleRequest{
		VIN:        "1HGBH41JXMN109186",
		Make:       "Toyota",
		Model:      "Camry",
		Year:       2023,
		OwnerID:    "owner-123",
		OwnerName:  "John Doe",
		OwnerEmail: "john@example.com",
		FuelType:   "gasoline",
		CreatedBy:  "admin-user",
	}

	_, err := handler.Handle(context.Background(), req)

	if err == nil {
		t.Fatal("Expected database error, got nil")
	}
}

func TestCreateVehicleHandler_DataNormalization(t *testing.T) {
	var capturedVehicle *domain.Vehicle

	mockRepo := &MockRepository{
		GetVehicleByVINFunc: func(ctx context.Context, vin string) (*domain.Vehicle, error) {
			return nil, apperrors.ErrResourceNotFound
		},
		CreateVehicleFunc: func(ctx context.Context, vehicle *domain.Vehicle) error {
			capturedVehicle = vehicle
			return nil
		},
	}

	handler := NewCreateVehicleHandler(mockRepo)

	req := &CreateVehicleRequest{
		VIN:          "  1hgbh41jxmn109186  ",
		Make:         "  Toyota  ",
		Model:        "  Camry  ",
		Year:         2023,
		LicensePlate: "  abc123  ",
		OwnerID:      "owner-123",
		OwnerName:    "  John Doe  ",
		OwnerEmail:   "  JOHN@EXAMPLE.COM  ",
		FuelType:     "gasoline",
		CreatedBy:    "admin-user",
	}

	_, err := handler.Handle(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if capturedVehicle.VIN != "1HGBH41JXMN109186" {
		t.Errorf("Expected VIN to be uppercase and trimmed, got %s", capturedVehicle.VIN)
	}

	if capturedVehicle.Make != "Toyota" {
		t.Errorf("Expected Make to be trimmed, got %s", capturedVehicle.Make)
	}

	if capturedVehicle.LicensePlate != "ABC123" {
		t.Errorf("Expected LicensePlate to be uppercase and trimmed, got %s", capturedVehicle.LicensePlate)
	}

	if capturedVehicle.OwnerEmail != "john@example.com" {
		t.Errorf("Expected OwnerEmail to be lowercase and trimmed, got %s", capturedVehicle.OwnerEmail)
	}
}
