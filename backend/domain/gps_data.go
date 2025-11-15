package domain

import "time"

// GPSData represents GPS location data from IoT devices
type GPSData struct {
	ID        string  `json:"id"`
	DeviceID  string  `json:"device_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp float64 `json:"timestamp"` // Unix timestamp as float64
}

// GetTimestamp converts the Unix timestamp to time.Time
func (g *GPSData) GetTimestamp() time.Time {
	return time.Unix(int64(g.Timestamp), 0)
}

// GPSDataResponse represents GPS data in API responses with formatted timestamp
type GPSDataResponse struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

// ToResponse converts GPSData to GPSDataResponse with proper time formatting
func (g *GPSData) ToResponse() GPSDataResponse {
	return GPSDataResponse{
		ID:        g.ID,
		DeviceID:  g.DeviceID,
		Latitude:  g.Latitude,
		Longitude: g.Longitude,
		Timestamp: g.GetTimestamp(),
	}
}
