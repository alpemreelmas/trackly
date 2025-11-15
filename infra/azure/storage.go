package azure

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	blob2 "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
)

type Storage struct {
	account       string
	containerName string
	client        *azblob.Client
}

// NewStorage initializes Azure Blob service
//
// Requires env:
//
//	AZURE_STORAGE_CONNECTION_STRING
func NewStorage(connString string, containerName string) (*Storage, error) {
	//cred, err := azidentity.NewDefaultAzureCredential(nil)
	//if err != nil {
	//	return nil, err
	//}

	client, err := azblob.NewClientFromConnectionString(connString, nil)
	if err != nil {
		return nil, err
	}

	accountName, err := extractAccountName(connString)
	if err != nil {
		return nil, err
	}

	return &Storage{
		account:       accountName,
		client:        client,
		containerName: containerName,
	}, nil
}

// Upload file to Azure Blob Storage
func (s *Storage) Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	options := &azblob.UploadStreamOptions{
		HTTPHeaders: &blob2.HTTPHeaders{
			BlobContentType: &contentType,
		},
	}

	_, err := s.client.UploadStream(ctx, s.containerName, filename, file, options)
	if err != nil {
		return "", err
	}

	return s.URL(filename), nil
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
	const key = "AccountName="
	start := len(key)

	idx := -1
	for i := 0; i < len(conn)-start; i++ {
		if conn[i:i+start] == key {
			idx = i + start
			break
		}
	}

	if idx == -1 {
		return "", fmt.Errorf("AccountName not found in connection string")
	}

	end := idx
	for end < len(conn) && conn[end] != ';' {
		end++
	}

	return conn[idx:end], nil
}
