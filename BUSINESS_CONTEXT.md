# Business Context Description - Trackly Project

## 📋 Project Overview
**Trackly** is a vehicle fleet management and tracking system that combines vehicle information management with real-time GPS location tracking. It enables vehicle owners and fleet managers to monitor vehicle status, manage documentation, and track vehicle locations in real-time.

---

## 👥 Stakeholders

| Stakeholder | Description |
|------------|-------------|
| **Vehicle Owners** | Individual or corporate owners managing one or multiple vehicles |
| **Fleet Managers** | Managers overseeing fleet operations and maintenance |
| **Insurance Companies** | Entities interested in vehicle insurance information and documentation |
| **IoT Device Operators** | Personnel managing GPS tracking devices |
| **System Administrators** | Team responsible for system maintenance and configuration |

---

## 🎯 Use Cases

### **1. Vehicle Management**

| Use Case | Description | Actors |
|----------|-------------|--------|
| **UC-001: Register Vehicle** | Owner registers a new vehicle with complete specifications, VIN, owner info, and insurance details | Vehicle Owner, System |
| **UC-002: View Vehicle Details** | Owner/Manager retrieves complete vehicle information including specs, insurance, and documents | Vehicle Owner, Manager |
| **UC-003: Update Vehicle Information** | Owner modifies vehicle specifications, insurance, or owner contact information | Vehicle Owner, System |
| **UC-004: Manage Vehicle Status** | Manager updates vehicle status (active, inactive, sold, stolen, accident, scrapped) | Fleet Manager, System |

### **2. Document Management**

| Use Case | Description | Actors |
|----------|-------------|--------|
| **UC-005: Upload Vehicle Documents** | Owner uploads critical documents (insurance, registration, title, inspection, warranty) | Vehicle Owner, System, Azure Storage |
| **UC-006: Retrieve Documents** | Owner/Manager views list of documents associated with a vehicle | Vehicle Owner, Manager |
| **UC-007: Download Document** | User downloads a specific vehicle document | Vehicle Owner, Manager |
| **UC-008: Delete Document** | Owner removes expired or outdated documents from the system | Vehicle Owner, System |


### **3. GPS Tracking**

| Use Case | Description | Actors |
|----------|-------------|--------|
| **UC-013: Stream GPS Data** | IoT device sends real-time location data to the system | GPS Device, IoT Hub, System |
| **UC-014: Query GPS Data** | Manager retrieves GPS location history for a vehicle | Fleet Manager, System |
| **UC-015: Track Vehicle Location** | User monitors current vehicle location in real-time | Vehicle Owner, Manager |

### **4. System Health**

| Use Case | Description | Actors |
|----------|-------------|--------|
| **UC-016: Health Check** | Monitor system availability and operational status | System, Monitoring Tools |

---

## 📖 User Stories

### **Vehicle Owner Persona**
*"As a vehicle owner, I want to..."*

1. **Register my vehicle** - ...register my new vehicle with complete information so that I can track it and manage its documentation.
2. **Upload insurance documents** - ...upload my insurance policy and documents so I have a centralized backup and can easily access them.
3. **Monitor insurance expiration** - ...receive alerts when my insurance is about to expire so I can renew it on time.
4. **Track vehicle location** - ...see my vehicle's real-time location to monitor its whereabouts and movement.
5. **Manage vehicle documents** - ...store and organize all vehicle documents (registration, title, inspection, warranty) in one place.
6. **Download documents** - ...easily download any stored document when needed for insurance claims or authorities.
7. **Update vehicle information** - ...modify vehicle specifications or owner contact details when circumstances change.

### **Fleet Manager Persona**
*"As a fleet manager, I want to..."*

1. **View all vehicle details** - ...access comprehensive information about each vehicle in my fleet including specifications and insurance.
2. **Monitor vehicle status** - ...track the status of each vehicle (active, sold, accident, stolen) to maintain fleet awareness.
3. **Query historical GPS data** - ...retrieve past location data for audit trails and accountability.
4. **Check document status** - ...identify vehicles with expired or expiring documents for compliance management.
5. **Manage fleet inventory** - ...update vehicle status across my fleet for better operational control.

### **System Administrator Persona**
*"As a system administrator, I want to..."*

1. **Monitor system health** - ...check the system's operational status through health check endpoints.
2. **Verify document authenticity** - ...mark documents as verified after manual review for compliance.
3. **Configure system settings** - ...manage system configuration through environment variables.

---

## 🔑 Key Features & Business Rules

### **Vehicle Management Features**
- Complete vehicle profile with VIN, make, model, year, color, license plate
- Engine specifications tracking (displacement, cylinders, horsepower, torque)
- Transmission and fuel type information
- Mileage tracking
- Owner information management (name, email, phone)

### **Insurance Management**
- Policy tracking (number, provider, type, coverage amount)
- Premium and deductible management
- Insurance expiration monitoring
- Active status tracking
- Insurance provider contact information

### **Document Management**
- Support for 12 document types: Insurance Policy, Insurance Card, Registration, Title, Inspection, Emission Test, Purchase Agreement, Service Record, Warranty, Receipt, Accident Report, Other
- Document tracking with expiration dates
- Document verification workflow
- Azure Blob Storage integration
- **Business Rules:**
  - Expiry tracking for relevant documents
  - Verification status tracking
  - Document status: no_documents, has_expired, has_expiring, up_to_date


### **GPS Tracking**
- Real-time location data collection from IoT devices
- Latitude/Longitude coordinate tracking
- Timestamp recording
- Cosmos DB integration for high-volume data storage
- **Business Rules:**
  - GPS data persisted with device ID and timestamp
  - Unix timestamp conversion for standardization

### **Vehicle Status Management**
- Status types: Active, Inactive, Sold, Scrapped, Stolen, Accident
- Status update audit trail (updated_by, updated_at)
- Creation audit trail (created_by, created_at)

---

## 🏗️ Technical Architecture Highlights

### **Technology Stack**
- **Backend:** Go with Fiber web framework
- **Databases:** 
  - Couchbase (Vehicle data)
  - Azure Cosmos DB (GPS data)
- **Storage:** Azure Blob Storage (Document/Media files)
- **IoT:** Azure IoT Hub with MQTT protocol
- **Logging:** Uber's Zap logger
- **Validation:** Go Playground validator

### **Data Flow**
```
IoT Device (GPS-IoT.py) 
    ↓ (MQTT via Azure IoT Hub)
Cosmos DB (GPS Data)
    ↓
Backend API (Go/Fiber)
    ↓
Couchbase (Vehicle Data) + Azure Storage (Documents/Media)
```

---

## 📊 Core Entities

```
Vehicle
├── Basic Info (VIN, Make, Model, Year, Color, License Plate)
├── Owner Info (ID, Name, Email, Phone)
├── Engine Specifications
├── Insurance Info
├── Documents[] (Type, Name, File Reference)
└── Status & Timestamps

GPSData
├── Device ID
├── Coordinates (Latitude, Longitude)
└── Timestamp
```

---

## 🔄 Business Processes

### **Vehicle Onboarding Process**
1. Owner creates vehicle record
2. Owner uploads vehicle documents
3. System initializes GPS tracking for vehicle
4. Fleet manager reviews and activates vehicle

### **Document Lifecycle**
1. Owner uploads document
2. Administrator verifies document
3. System monitors expiration dates
4. System alerts on upcoming expiration (30 days before)
5. Owner renews/uploads new document
6. System marks old document as expired

### **GPS Tracking Workflow**
1. IoT device initializes with device ID
2. GPS device sends location data via MQTT to Azure IoT Hub
3.Stores in Cosmos DB
4. Fleet manager queries historical or real-time data via API
5. System maintains location history for audit trails

---



## 🚀 Future Enhancements (Roadmap)

- **Maintenance Scheduling** - Predictive maintenance based on mileage
- **Driver Management** - Associate drivers with vehicles
- **Trip History Analysis** - Analytics on vehicle usage patterns
- **Geofencing Alerts** - Notifications when vehicles enter/exit defined zones
- **Cost Analytics** - Fuel consumption and maintenance cost tracking
- **Mobile App** - Native iOS/Android application for vehicle owners
- **Advanced Reporting** - Dashboard and report generation
- **Integration with External APIs** - Weather, traffic, and road condition services

---
