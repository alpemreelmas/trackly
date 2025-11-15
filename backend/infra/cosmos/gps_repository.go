package cosmosdb

import (
	"context"
	"encoding/json"
	"fmt"
	"microservicetest/domain"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type GPSRepository struct {
	client        *azcosmos.Client
	database      *azcosmos.DatabaseClient
	container     *azcosmos.ContainerClient
	databaseName  string
	containerName string
}

func NewGPSRepository(endpoint, key, databaseName, containerName string) (*GPSRepository, error) {
	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cosmos client: %w", err)
	}

	database, err := client.NewDatabase(databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	container, err := database.NewContainer(containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get container: %w", err)
	}

	return &GPSRepository{
		client:        client,
		database:      database,
		container:     container,
		databaseName:  databaseName,
		containerName: containerName,
	}, nil
}

// GetGPSDataByDateRange retrieves GPS data within a date range
func (r *GPSRepository) GetGPSDataByDateRange(ctx context.Context, deviceID string, startDate, endDate time.Time) ([]domain.GPSData, error) {
	query := `SELECT * FROM c`

	// queryOptions := azcosmos.QueryOptions{
	// 	QueryParameters: []azcosmos.QueryParameter{
	// 		{Name: "@deviceID", Value: deviceID},
	// 		{Name: "@startDate", Value: startDate.Unix()},
	// 		{Name: "@endDate", Value: endDate.Unix()},
	// 	},
	// }

	// Create partition key with the device_id value
	pk := azcosmos.NewPartitionKeyString(deviceID)
	queryPager := r.container.NewQueryItemsPager(query, pk, nil)

	var gpsDataList []domain.GPSData

	for queryPager.More() {
		response, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to query items: %w", err)
		}

		for _, item := range response.Items {
			var gpsData domain.GPSData
			if err := json.Unmarshal(item, &gpsData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal item: %w", err)
			}
			gpsDataList = append(gpsDataList, gpsData)
		}
	}

	return gpsDataList, nil
}

// GetGPSDataByDevice retrieves all GPS data for a specific device
func (r *GPSRepository) GetGPSDataByDevice(ctx context.Context, deviceID string, limit int) ([]domain.GPSData, error) {
	query := fmt.Sprintf(`SELECT TOP %d * FROM c WHERE c.device_id = @deviceID ORDER BY c.timestamp DESC`, limit)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@deviceID", Value: deviceID},
		},
	}

	// Create partition key with the device_id value
	pk := azcosmos.NewPartitionKeyString(deviceID)
	queryPager := r.container.NewQueryItemsPager(query, pk, &queryOptions)

	var gpsDataList []domain.GPSData

	for queryPager.More() {
		response, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to query items: %w", err)
		}

		for _, item := range response.Items {
			var gpsData domain.GPSData
			if err := json.Unmarshal(item, &gpsData); err != nil {
				return nil, fmt.Errorf("failed to unmarshal item: %w", err)
			}
			gpsDataList = append(gpsDataList, gpsData)
		}
	}

	return gpsDataList, nil
}
