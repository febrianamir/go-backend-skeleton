package storage

import (
	"context"
	"io"
)

type Storage interface {
	UploadFile(ctx context.Context, bucketName, fileName, contentType string, file io.Reader) error
	GetFileTemporaryURL(ctx context.Context, bucketName, filename string) (string, error)
	GetFilePublicURL(ctx context.Context, bucketName, filename string) (string, error)
	GetObject(ctx context.Context, bucketName, filename string) (io.Reader, error)
	FGetObject(ctx context.Context, bucketName, filename, destination string) error
	FPutObject(ctx context.Context, bucketName, filename, filepath string) error
	FCopyObject(ctx context.Context, bucketName, src, dst string) error
	RemoveFile(ctx context.Context, bucketName, pathFilename string) error
	IsFileExist(ctx context.Context, bucketName, fileptah string) (bool, error)
}
