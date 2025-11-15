package app

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
}
