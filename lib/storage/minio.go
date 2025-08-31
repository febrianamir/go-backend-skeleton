package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type Minio struct {
	*minio.Client
	MinioCdnBaseDns string
	MinioCdnBaseUrl string
}

func (m *Minio) UploadFile(ctx context.Context, bucketName, fileName, contentType string, file io.Reader) error {
	_, err := m.Client.PutObject(ctx, bucketName, fileName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println("failed upload file to minio: ", err)
	}

	return err
}

func (m *Minio) GetFileTemporaryURL(ctx context.Context, bucketName, filename string) (string, error) {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)

	// Generates a presigned url.
	presignedURL, err := m.Client.PresignedGetObject(ctx, bucketName, filename, time.Second*24*60*60, reqParams)
	if err != nil {
		log.Println("failed get file temporary url: ", err)
		return "", err
	}

	baseURL := m.MinioCdnBaseDns
	if baseURL == "" {
		baseURL = m.MinioCdnBaseUrl
	}

	return fmt.Sprintf("%s%s?%s", baseURL, presignedURL.Path, presignedURL.RawQuery), nil
}

func (m *Minio) GetFilePublicURL(ctx context.Context, bucketName, filename string) (string, error) {
	return "", nil
}

func (m *Minio) GetObject(ctx context.Context, bucketName, filename string) (io.Reader, error) {
	return m.Client.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
}

func (m *Minio) FGetObject(ctx context.Context, bucketName, filename, destination string) error {
	return m.Client.FGetObject(ctx, bucketName, filename, destination, minio.GetObjectOptions{})
}

func (m *Minio) FPutObject(ctx context.Context, bucketName, filename, filepath string) error {
	_, err := m.Client.FPutObject(ctx, bucketName, filename, filepath, minio.PutObjectOptions{})
	return err
}

func (m *Minio) RemoveFile(ctx context.Context, bucketName, pathFilename string) error {
	opts := minio.RemoveObjectOptions{GovernanceBypass: true}
	err := m.Client.RemoveObject(ctx, bucketName, pathFilename, opts)
	if err != nil {
		log.Println("failed remove file minio: ", err)
	}

	return err
}

func (m *Minio) IsFileExist(ctx context.Context, bucketName, fileptah string) (bool, error) {
	opts := minio.GetObjectOptions{}
	_, err := m.Client.StatObject(ctx, bucketName, fileptah, opts)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (m *Minio) FCopyObject(ctx context.Context, bucketName, src, dst string) error {

	_, err := m.Client.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: bucketName,
		Object: dst,
	}, minio.CopySrcOptions{
		Bucket: bucketName,
		Object: src,
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Minio) DeleteDirectoryTmp(ctx context.Context, bucketName string) error {
	return nil
}
