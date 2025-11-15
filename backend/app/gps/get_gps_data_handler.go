package gps

import (
	"context"
	"microservicetest/domain"
	cosmosdb "microservicetest/infra/cosmos"
	"time"

	"go.uber.org/zap"
)

type GetGPSDataRequest struct {
	DeviceID  string `query:"device_id" validate:"required"`
	StartDate string `query:"start_date"` // Format: 2006-01-02
	EndDate   string `query:"end_date"`   // Format: 2006-01-02
}

type GetGPSDataResponse struct {
	Data  []domain.GPSDataResponse `json:"data"`
	Count int                      `json:"count"`
}

type GetGPSDataHandler struct {
	repository *cosmosdb.GPSRepository
}

func NewGetGPSDataHandler(repository *cosmosdb.GPSRepository) *GetGPSDataHandler {
	return &GetGPSDataHandler{
		repository: repository,
	}
}

func (h *GetGPSDataHandler) Handle(ctx context.Context, req *GetGPSDataRequest) (*GetGPSDataResponse, error) {
	// Parse dates or use defaults
	var startDate, endDate time.Time
	var err error

	if req.StartDate == "" {
		// Default to today at 00:00:00
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			zap.L().Error("Failed to parse start_date", zap.Error(err))
			startDate = time.Now().Truncate(24 * time.Hour)
		}
	}

	if req.EndDate == "" {
		// Default to today at 23:59:59
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	} else {
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			zap.L().Error("Failed to parse end_date", zap.Error(err))
			endDate = time.Now()
		} else {
			// Set to end of day
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
		}
	}
	zap.L().Info("Fetching GPS data",
		zap.String("device_id", req.DeviceID),
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate),
	)

	gpsData, err := h.repository.GetGPSDataByDateRange(ctx, req.DeviceID, startDate, endDate)
	if err != nil {
		zap.L().Error("Failed to fetch GPS data", zap.Error(err))
		return nil, err
	}

	// Convert to response format with proper timestamp formatting
	responseData := make([]domain.GPSDataResponse, len(gpsData))
	for i, data := range gpsData {
		responseData[i] = data.ToResponse()
	}

	return &GetGPSDataResponse{
		Data:  responseData,
		Count: len(responseData),
	}, nil
}
