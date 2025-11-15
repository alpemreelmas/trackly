package azure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

type Storage struct {
	account       string
	containerName string
	accountKey    string
	client        *azblob.Client
}

// NewStorage initializes Azure Blob service
//
// Requires env:
//
//	AZURE_STORAGE_CONNECTION_STRING
func NewStorage(connString string, containerName string) (*Storage, error) {
	client, err := azblob.NewClientFromConnectionString(connString, nil)
	if err != nil {
		return nil, err
	}

	accountName, err := extractAccountName(connString)
	if err != nil {
		return nil, err
	}

	accountKey, err := extractAccountKey(connString)
	if err != nil {
		return nil, err
	}

	return &Storage{
		account:       accountName,
		accountKey:    accountKey,
		client:        client,
		containerName: containerName,
	}, nil
}

// Upload file to Azure Blob Storage with SAS token
func (s *Storage) Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// Read file into buffer
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate SAS token for upload
	sasURL, err := s.generateUploadSAS(filename)
	if err != nil {
		return "", fmt.Errorf("failed to generate SAS token: %w", err)
	}

	// Create client with SAS URL
	blobClient, err := blockblob.NewClientWithNoCredential(sasURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create blob client: %w", err)
	}

	// Create a ReadSeekCloser from bytes
	reader := bytes.NewReader(data)

	// Upload with options
	options := &blockblob.UploadOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: &contentType,
		},
	}

	_, err = blobClient.Upload(ctx, &readSeekNopCloser{reader}, options)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}

	return s.URL(filename), nil
}

// generateUploadSAS creates a SAS token for uploading a blob
func (s *Storage) generateUploadSAS(filename string) (string, error) {
	// Create shared key credential
	credential, err := azblob.NewSharedKeyCredential(s.account, s.accountKey)
	if err != nil {
		return "", fmt.Errorf("failed to create credential: %w", err)
	}

	// Set SAS token permissions and expiry
	now := time.Now().UTC()
	expiry := now.Add(15 * time.Minute) // Token valid for 15 minutes

	// Create SAS query parameters
	permissions := sas.BlobPermissions{Write: true, Create: true}
	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     now.Add(-5 * time.Minute), // Start 5 minutes ago to handle clock skew
		ExpiryTime:    expiry,
		Permissions:   permissions.String(),
		ContainerName: s.containerName,
		BlobName:      filename,
	}.SignWithSharedKey(credential)

	if err != nil {
		return "", fmt.Errorf("failed to sign SAS: %w", err)
	}

	// Build full URL with SAS token
	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
		s.account,
		s.containerName,
		filename,
		sasQueryParams.Encode(),
	)

	return sasURL, nil
}

// // Delete file
//
//	func (s *Storage) Delete(ctx context.Context, filename string) error {
//		blob := s.container.NewBlockBlobClient(filename)
//
//		_, err := blob.Delete(ctx, nil)
//		return err
//	}
//
// Get public URL (works if blob is public OR SAS token appended)
func (s *Storage) URL(filename string) string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s",
		s.account,
		s.containerName,
		filename,
	)
}

// Parse account name from the connection string
func extractAccountName(conn string) (string, error) {
	return extractValue(conn, "AccountName=")
}

// Parse account key from the connection string
func extractAccountKey(conn string) (string, error) {
	return extractValue(conn, "AccountKey=")
}

// Generic helper to extract values from connection string
func extractValue(conn, key string) (string, error) {
	start := len(key)

	idx := -1
	for i := 0; i < len(conn)-start; i++ {
		if conn[i:i+start] == key {
			idx = i + start
			break
		}
	}

	if idx == -1 {
		return "", fmt.Errorf("%s not found in connection string", key)
	}

	end := idx
	for end < len(conn) && conn[end] != ';' {
		end++
	}

	return conn[idx:end], nil
}

// readSeekNopCloser wraps a bytes.Reader to implement io.ReadSeekCloser
type readSeekNopCloser struct {
	*bytes.Reader
}

func (r *readSeekNopCloser) Close() error {
	return nil
}
