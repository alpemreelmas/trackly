package domain

import (
	"fmt"
	"time"
)

// Vehicle represents a vehicle in the system
type Vehicle struct {
	ID          string    `json:"id" couchbase:"id"`
	VIN         string    `json:"vin" couchbase:"vin"`                     // Vehicle Identification Number
	Make        string    `json:"make" couchbase:"make"`                   // Toyota, BMW, etc.
	Model       string    `json:"model" couchbase:"model"`                 // Camry, X5, etc.
	Year        int       `json:"year" couchbase:"year"`                   // Manufacturing year
	Color       string    `json:"color" couchbase:"color"`                 // Vehicle color
	LicensePlate string   `json:"license_plate" couchbase:"license_plate"` // License plate number
	
	// Owner information
	OwnerID     string `json:"owner_id" couchbase:"owner_id"`
	OwnerName   string `json:"owner_name" couchbase:"owner_name"`
	OwnerEmail  string `json:"owner_email" couchbase:"owner_email"`
	OwnerPhone  string `json:"owner_phone" couchbase:"owner_phone"`
	
	// Vehicle specifications
	Engine      EngineInfo      `json:"engine" couchbase:"engine"`
	Transmission string         `json:"transmission" couchbase:"transmission"` // Manual, Automatic, CVT
	FuelType    FuelType       `json:"fuel_type" couchbase:"fuel_type"`
	Mileage     int            `json:"mileage" couchbase:"mileage"`           // Current mileage
	
	// Insurance information
	Insurance   InsuranceInfo  `json:"insurance" couchbase:"insurance"`
	
	// Documents and media
	Documents   []Document     `json:"documents" couchbase:"documents"`
	Pictures    []Picture      `json:"pictures" couchbase:"pictures"`
	
	// Status and metadata
	Status      VehicleStatus  `json:"status" couchbase:"status"`
	CreatedAt   time.Time      `json:"created_at" couchbase:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" couchbase:"updated_at"`
	CreatedBy   string         `json:"created_by" couchbase:"created_by"`
	UpdatedBy   string         `json:"updated_by" couchbase:"updated_by"`
}

// EngineInfo contains engine specifications
type EngineInfo struct {
	Displacement float64 `json:"displacement" couchbase:"displacement"` // Engine size in liters
	Cylinders    int     `json:"cylinders" couchbase:"cylinders"`       // Number of cylinders
	Horsepower   int     `json:"horsepower" couchbase:"horsepower"`     // Engine power
	Torque       int     `json:"torque" couchbase:"torque"`             // Engine torque
}

// InsuranceInfo contains insurance details
type InsuranceInfo struct {
	PolicyNumber    string            `json:"policy_number" couchbase:"policy_number"`
	Provider        string            `json:"provider" couchbase:"provider"`         // Insurance company name
	PolicyType      InsurancePolicyType `json:"policy_type" couchbase:"policy_type"`
	CoverageAmount  float64           `json:"coverage_amount" couchbase:"coverage_amount"`
	Deductible      float64           `json:"deductible" couchbase:"deductible"`
	PremiumAmount   float64           `json:"premium_amount" couchbase:"premium_amount"`
	StartDate       time.Time         `json:"start_date" couchbase:"start_date"`
	EndDate         time.Time         `json:"end_date" couchbase:"end_date"`
	IsActive        bool              `json:"is_active" couchbase:"is_active"`
	ContactInfo     InsuranceContact  `json:"contact_info" couchbase:"contact_info"`
}

// InsuranceContact contains insurance provider contact information
type InsuranceContact struct {
	Phone       string `json:"phone" couchbase:"phone"`
	Email       string `json:"email" couchbase:"email"`
	Address     string `json:"address" couchbase:"address"`
	ClaimsPhone string `json:"claims_phone" couchbase:"claims_phone"`
	Website     string `json:"website" couchbase:"website"`
}

// Document represents various vehicle documents
type Document struct {
	ID           string       `json:"id" couchbase:"id"`
	Type         DocumentType `json:"type" couchbase:"type"`
	Name         string       `json:"name" couchbase:"name"`
	Description  string       `json:"description" couchbase:"description"`
	FileURL      string       `json:"file_url" couchbase:"file_url"`
	FileName     string       `json:"file_name" couchbase:"file_name"`
	FileSize     int64        `json:"file_size" couchbase:"file_size"`     // Size in bytes
	MimeType     string       `json:"mime_type" couchbase:"mime_type"`     // application/pdf, image/jpeg, etc.
	ExpiryDate   *time.Time   `json:"expiry_date" couchbase:"expiry_date"` // For documents that expire
	IssuedDate   *time.Time   `json:"issued_date" couchbase:"issued_date"`
	IssuedBy     string       `json:"issued_by" couchbase:"issued_by"`     // Issuing authority
	DocumentNumber string     `json:"document_number" couchbase:"document_number"`
	UploadedAt   time.Time    `json:"uploaded_at" couchbase:"uploaded_at"`
	UploadedBy   string       `json:"uploaded_by" couchbase:"uploaded_by"`
	IsVerified   bool         `json:"is_verified" couchbase:"is_verified"`
	VerifiedAt   *time.Time   `json:"verified_at" couchbase:"verified_at"`
	VerifiedBy   string       `json:"verified_by" couchbase:"verified_by"`
}

// Picture represents vehicle images
type Picture struct {
	ID          string      `json:"id" couchbase:"id"`
	Type        PictureType `json:"type" couchbase:"type"`
	Title       string      `json:"title" couchbase:"title"`
	Description string      `json:"description" couchbase:"description"`
	URL         string      `json:"url" couchbase:"url"`
	ThumbnailURL string     `json:"thumbnail_url" couchbase:"thumbnail_url"`
	FileName    string      `json:"file_name" couchbase:"file_name"`
	FileSize    int64       `json:"file_size" couchbase:"file_size"`
	Width       int         `json:"width" couchbase:"width"`
	Height      int         `json:"height" couchbase:"height"`
	MimeType    string      `json:"mime_type" couchbase:"mime_type"`
	TakenAt     *time.Time  `json:"taken_at" couchbase:"taken_at"`
	UploadedAt  time.Time   `json:"uploaded_at" couchbase:"uploaded_at"`
	UploadedBy  string      `json:"uploaded_by" couchbase:"uploaded_by"`
	IsMain      bool        `json:"is_main" couchbase:"is_main"`      // Main/primary picture
	SortOrder   int         `json:"sort_order" couchbase:"sort_order"` // Display order
}

// Enums and constants

type VehicleStatus string

const (
	VehicleStatusActive    VehicleStatus = "active"
	VehicleStatusInactive  VehicleStatus = "inactive"
	VehicleStatusSold      VehicleStatus = "sold"
	VehicleStatusScrapped  VehicleStatus = "scrapped"
	VehicleStatusStolen    VehicleStatus = "stolen"
	VehicleStatusAccident  VehicleStatus = "accident"
)

type FuelType string

const (
	FuelTypeGasoline FuelType = "gasoline"
	FuelTypeDiesel   FuelType = "diesel"
	FuelTypeElectric FuelType = "electric"
	FuelTypeHybrid   FuelType = "hybrid"
	FuelTypeLPG      FuelType = "lpg"
	FuelTypeCNG      FuelType = "cng"
)

type InsurancePolicyType string

const (
	InsurancePolicyLiability    InsurancePolicyType = "liability"
	InsurancePolicyComprehensive InsurancePolicyType = "comprehensive"
	InsurancePolicyCollision    InsurancePolicyType = "collision"
	InsurancePolicyFullCoverage InsurancePolicyType = "full_coverage"
)

type DocumentType string

const (
	DocumentTypeInsurancePolicy    DocumentType = "insurance_policy"
	DocumentTypeInsuranceCard      DocumentType = "insurance_card"
	DocumentTypeRegistration       DocumentType = "registration"
	DocumentTypeTitle              DocumentType = "title"
	DocumentTypeInspection         DocumentType = "inspection"
	DocumentTypeEmissionTest       DocumentType = "emission_test"
	DocumentTypePurchaseAgreement  DocumentType = "purchase_agreement"
	DocumentTypeServiceRecord      DocumentType = "service_record"
	DocumentTypeWarranty           DocumentType = "warranty"
	DocumentTypeReceipt            DocumentType = "receipt"
	DocumentTypeAccidentReport     DocumentType = "accident_report"
	DocumentTypeOther              DocumentType = "other"
)

type PictureType string

const (
	PictureTypeExteriorFront  PictureType = "exterior_front"
	PictureTypeExteriorBack   PictureType = "exterior_back"
	PictureTypeExteriorLeft   PictureType = "exterior_left"
	PictureTypeExteriorRight  PictureType = "exterior_right"
	PictureTypeInteriorFront  PictureType = "interior_front"
	PictureTypeInteriorBack   PictureType = "interior_back"
	PictureTypeDashboard      PictureType = "dashboard"
	PictureTypeEngine         PictureType = "engine"
	PictureTypeTrunk          PictureType = "trunk"
	PictureTypeWheels         PictureType = "wheels"
	PictureTypeDamage         PictureType = "damage"
	PictureTypeAccident       PictureType = "accident"
	PictureTypeOther          PictureType = "other"
)

// Helper methods

// IsInsuranceExpired checks if the vehicle's insurance has expired
func (v *Vehicle) IsInsuranceExpired() bool {
	return time.Now().After(v.Insurance.EndDate)
}

// IsInsuranceExpiringSoon checks if insurance expires within the given days
func (v *Vehicle) IsInsuranceExpiringSoon(days int) bool {
	expiryThreshold := time.Now().AddDate(0, 0, days)
	return v.Insurance.EndDate.Before(expiryThreshold)
}

// GetMainPicture returns the main picture of the vehicle
func (v *Vehicle) GetMainPicture() *Picture {
	for _, picture := range v.Pictures {
		if picture.IsMain {
			return &picture
		}
	}
	return nil
}

// GetDocumentsByType returns documents of a specific type
func (v *Vehicle) GetDocumentsByType(docType DocumentType) []Document {
	var documents []Document
	for _, doc := range v.Documents {
		if doc.Type == docType {
			documents = append(documents, doc)
		}
	}
	return documents
}

// GetPicturesByType returns pictures of a specific type
func (v *Vehicle) GetPicturesByType(picType PictureType) []Picture {
	var pictures []Picture
	for _, pic := range v.Pictures {
		if pic.Type == picType {
			pictures = append(pictures, pic)
		}
	}
	return pictures
}

// HasExpiredDocuments checks if any documents have expired
func (v *Vehicle) HasExpiredDocuments() bool {
	now := time.Now()
	for _, doc := range v.Documents {
		if doc.ExpiryDate != nil && now.After(*doc.ExpiryDate) {
			return true
		}
	}
	return false
}

// GetExpiredDocuments returns all expired documents
func (v *Vehicle) GetExpiredDocuments() []Document {
	var expired []Document
	now := time.Now()
	for _, doc := range v.Documents {
		if doc.ExpiryDate != nil && now.After(*doc.ExpiryDate) {
			expired = append(expired, doc)
		}
	}
	return expired
}

// GetExpiringDocuments returns documents expiring within the given days
func (v *Vehicle) GetExpiringDocuments(days int) []Document {
	var expiring []Document
	threshold := time.Now().AddDate(0, 0, days)
	for _, doc := range v.Documents {
		if doc.ExpiryDate != nil && doc.ExpiryDate.Before(threshold) && doc.ExpiryDate.After(time.Now()) {
			expiring = append(expiring, doc)
		}
	}
	return expiring
}

// Business logic methods

// CalculateAge returns the age of the vehicle in years
func (v *Vehicle) CalculateAge() int {
	return time.Now().Year() - v.Year
}

// IsVintage checks if the vehicle is considered vintage (25+ years old)
func (v *Vehicle) IsVintage() bool {
	return v.CalculateAge() >= 25
}

// GetInsuranceStatus returns the current insurance status
func (v *Vehicle) GetInsuranceStatus() string {
	if !v.Insurance.IsActive {
		return "inactive"
	}
	
	if v.IsInsuranceExpired() {
		return "expired"
	}
	
	if v.IsInsuranceExpiringSoon(30) {
		return "expiring_soon"
	}
	
	return "active"
}

// GetDocumentStatus returns overall document status
func (v *Vehicle) GetDocumentStatus() string {
	if len(v.Documents) == 0 {
		return "no_documents"
	}
	
	if v.HasExpiredDocuments() {
		return "has_expired"
	}
	
	if len(v.GetExpiringDocuments(30)) > 0 {
		return "has_expiring"
	}
	
	return "up_to_date"
}

// UpdateTimestamp updates the UpdatedAt field and UpdatedBy
func (v *Vehicle) UpdateTimestamp(updatedBy string) {
	v.UpdatedAt = time.Now()
	v.UpdatedBy = updatedBy
}

// SetMainPicture sets a picture as the main picture and unsets others
func (v *Vehicle) SetMainPicture(pictureID string) error {
	found := false
	
	// First, unset all main pictures
	for i := range v.Pictures {
		if v.Pictures[i].ID == pictureID {
			v.Pictures[i].IsMain = true
			found = true
		} else {
			v.Pictures[i].IsMain = false
		}
	}
	
	if !found {
		return fmt.Errorf("picture with ID %s not found", pictureID)
	}
	
	return nil
}

// AddDocument adds a new document to the vehicle
func (v *Vehicle) AddDocument(doc Document) error {
	// Check for duplicate document IDs
	for _, existingDoc := range v.Documents {
		if existingDoc.ID == doc.ID {
			return fmt.Errorf("document with ID %s already exists", doc.ID)
		}
	}
	
	v.Documents = append(v.Documents, doc)
	return nil
}

// AddPicture adds a new picture to the vehicle
func (v *Vehicle) AddPicture(pic Picture) error {
	
	// Check for duplicate picture IDs
	for _, existingPic := range v.Pictures {
		if existingPic.ID == pic.ID {
			return fmt.Errorf("picture with ID %s already exists", pic.ID)
		}
	}
	
	// If this is the first picture, make it main
	if len(v.Pictures) == 0 {
		pic.IsMain = true
	}
	
	v.Pictures = append(v.Pictures, pic)
	return nil
}

// RemoveDocument removes a document by ID
func (v *Vehicle) RemoveDocument(documentID string) error {
	for i, doc := range v.Documents {
		if doc.ID == documentID {
			v.Documents = append(v.Documents[:i], v.Documents[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("document with ID %s not found", documentID)
}

// RemovePicture removes a picture by ID
func (v *Vehicle) RemovePicture(pictureID string) error {
	for i, pic := range v.Pictures {
		if pic.ID == pictureID {
			wasMain := pic.IsMain
			v.Pictures = append(v.Pictures[:i], v.Pictures[i+1:]...)
			
			// If we removed the main picture, set the first remaining picture as main
			if wasMain && len(v.Pictures) > 0 {
				v.Pictures[0].IsMain = true
			}
			
			return nil
		}
	}
	return fmt.Errorf("picture with ID %s not found", pictureID)
}

// Factory methods

// NewVehicle creates a new vehicle with default values
func NewVehicle(vin, make, model string, year int, ownerID string) *Vehicle {
	now := time.Now()
	
	return &Vehicle{
		ID:        GenerateVehicleID(),
		VIN:       vin,
		Make:      make,
		Model:     model,
		Year:      year,
		OwnerID:   ownerID,
		Status:    VehicleStatusActive,
		Documents: make([]Document, 0),
		Pictures:  make([]Picture, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewDocument creates a new document with default values
func NewDocument(docType DocumentType, name, fileURL, fileName string, fileSize int64, uploadedBy string) *Document {
	return &Document{
		ID:         GenerateDocumentID(),
		Type:       docType,
		Name:       name,
		FileURL:    fileURL,
		FileName:   fileName,
		FileSize:   fileSize,
		UploadedAt: time.Now(),
		UploadedBy: uploadedBy,
		IsVerified: false,
	}
}

// NewPicture creates a new picture with default values
func NewPicture(picType PictureType, title, url, fileName string, fileSize int64, width, height int, uploadedBy string) *Picture {
	return &Picture{
		ID:         GeneratePictureID(),
		Type:       picType,
		Title:      title,
		URL:        url,
		FileName:   fileName,
		FileSize:   fileSize,
		Width:      width,
		Height:     height,
		UploadedAt: time.Now(),
		UploadedBy: uploadedBy,
		IsMain:     false,
		SortOrder:  0,
	}
}

// Helper functions for ID generation (you can implement these based on your needs)
func GenerateVehicleID() string {
	return "VEH_" + time.Now().Format("20060102150405")
}

func GenerateDocumentID() string {
	return "DOC_" + time.Now().Format("20060102150405")
}

func GeneratePictureID() string {
	return "PIC_" + time.Now().Format("20060102150405")
}