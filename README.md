# Trackly - Vehicle Fleet Management System

Trackly is a modern vehicle fleet management and tracking system with complete administrative capabilities.

## 🚀 Project Setup Guide

This guide provides step-by-step installation instructions for developers starting from scratch.

---

## 📋 Prerequisites

Install the following software:

### 1. **Go Programming Language**
   - **Download**: [golang.org](https://golang.org/dl)
   - **Version**: 1.24.0 or higher
   - Verify after installation:
     ```powershell
     go version
     ```

### 2. **Python**
   - **Download**: [python.org](https://www.python.org/downloads)
   - **Version**: 3.8 or higher
   - Verify after installation:
     ```powershell
     py --version
     ```

### 3. **Git**
   - **Download**: [git-scm.com](https://git-scm.com)
   - Verify after installation:
     ```powershell
     git --version
     ```

### 4. **Couchbase Server** (Database)
   - **Download**: [couchbase.com/downloads](https://www.couchbase.com/downloads)
   - **Version**: 7.2 or higher
   - Download and run Windows MSI installer
   - Set admin password during installation: `password`

## 📁 Project Structure

```
trackly/
├── backend/                 # Go Backend API Server
│   ├── app/                # Application handlers
│   ├── domain/             # Domain models
│   ├── infra/              # Infrastructure (Couchbase, Cosmos, Azure)
│   ├── pkg/                # Utility packages
│   ├── config/             # Configuration files
│   ├── main.go             # Application entry point
│   ├── go.mod              # Go modules
│   └── Dockerfile          # Docker configuration
├── iot/                    # Python IoT GPS Simulator
│   ├── gps-iot.py          # GPS simulator script
│   └── requirements.txt     # Python dependencies
├── BUSINESS_CONTEXT.md     # Project description
└── README.md               # This file
```

---

## ⚙️ Installation Steps

### Step 1: Download Project

```powershell
# Clone or navigate to the project directory
cd c:\Users\{username}\trackly
```

### Step 2: Start Couchbase

**Option A: Using Installed Couchbase Server**
1. Open "Couchbase Server" from Windows Start Menu
2. Navigate to `http://localhost:8091`
3. Login with: `Administrator` / `password`
4. Create a bucket:
   - Name: `vehicles`
   - Type: `Couchbase`
   - RAM Quota: 256 MB


### Step 3: Start Backend API

**Terminal 1 - Backend Server:**
```powershell
cd c:\Users\{username}\trackly\backend

# Download Go dependencies
go mod download

# Start the server
air
```

Expected output:
```
Server started on port 8080
┌───────────────────────────────────┐
│       Fiber v2.52.6               │
│   http://127.0.0.1:8080           │
│   Handlers ............ 16         │
└───────────────────────────────────┘
```

### Step 4: Test the API

**Terminal 2 - Health Check:**
```powershell
# Check if API is healthy
curl http://localhost:8080/healthcheck

# Expected response:
# {"status":"OK"}
```

### Step 5: Python IoT Simulator (Optional)

**Terminal 3 - GPS Simulator:**
```powershell
cd c:\Users\{username}\trackly\iot

# Install Python dependencies
py -m pip install -r requirements.txt

# Start simulator (for testing without Azure IoT Hub)
py gps-iot.py
```

---

## 📡 API Endpoints

All endpoints are available at `http://localhost:8080`

### Health Check
```
GET /healthcheck
Response: {"status":"OK"}
```

### Vehicle Management
```
POST   /vehicles              → Create new vehicle
GET    /vehicles/:id          → Get vehicle details
PUT    /vehicles/:id          → Update vehicle information
```

### Document Management
```
POST   /vehicles/:id/documents                    → Add document
GET    /vehicles/:id/documents                    → List documents
GET    /vehicles/:id/documents/:doc_id/download   → Download document
DELETE /vehicles/:id/documents/:doc_id            → Delete document
```

### GPS Data
```
GET /gps/data → Query GPS data
```

---

## 🧪 Example API Calls

### Create Vehicle

**PowerShell:**
```powershell
$body = @{
    id = "vehicle-001"
    vin = "WBADT43452G296706"
    make = "BMW"
    model = "3 Series"
    year = 2023
    color = "Black"
    licenseplate = "ABC-1234"
    ownerId = "owner-001"
    status = "active"
} | ConvertTo-Json

Invoke-WebRequest `
    -Uri "http://localhost:8080/vehicles" `
    -Method POST `
    -Body $body `
    -ContentType "application/json"
```

**cURL:**
```bash
curl -X POST http://localhost:8080/vehicles \
  -H "Content-Type: application/json" \
  -d '{
    "id": "vehicle-001",
    "vin": "WBADT43452G296706",
    "make": "BMW",
    "model": "3 Series",
    "year": 2023,
    "color": "Black",
    "licenseplate": "ABC-1234",
    "ownerId": "owner-001",
    "status": "active"
  }'
```

### Get Vehicle Details

```powershell
Invoke-WebRequest `
    -Uri "http://localhost:8080/vehicles/vehicle-001" `
    -Method GET
```

---
## 📊 System Architecture

```
IoT Device (GPS)
     ↓ MQTT
Azure IoT Hub
     ↓
Cosmos DB (GPS Data)
     ↓
Go/Fiber Backend API
     ↓
Couchbase (Vehicle Data)
Azure Blob Storage (Documents)
```

---

## 🔒 Configuration

Backend configuration is in `backend/config/config.yaml`:

```yaml
port: "8080"
couchbase_url: "couchbase://localhost"
couchbase_username: "Administrator"
couchbase_password: "password"
azure_connection_string: "DefaultEndpointsProtocol=https;..."
cosmosdb_endpoint: "https://localhost:8081/"
cosmosdb_key: "fake-key"
cosmosdb_database: "trackly"
cosmosdb_container: "gpsdata"
```

---

## 📚 Technologies Used

- **Backend**: Go 1.24 + Fiber v2.52
- **Databases**: Couchbase, Azure Cosmos DB
- **Storage**: Azure Blob Storage
- **IoT**: Python, MQTT, Azure IoT Hub
- **Logging**: Uber Zap
- **Validation**: Go Playground Validator

---

## 🎯 Next Steps

1. ✅ Start backend server
2. ✅ Setup Couchbase and create bucket
3. ✅ Test the API
4. ⏳ Configure IoT simulator (with Azure IoT Hub credentials)

**Happy development!** 🚀
