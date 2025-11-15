package app

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
	Download(ctx context.Context, filename string) ([]byte, string, error)
	Remove(ctx context.Context, filename string) error
}
