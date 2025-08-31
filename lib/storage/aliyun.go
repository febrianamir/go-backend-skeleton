package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Aliyun struct {
	*oss.Client
	AliyunEndpoint          string
	StorageBucketName       string
	StoragePublicBucketName string
	StorageTmpUrlExpiration int
}

// Define the function that is used to handle progress change events.
// Documentation on: https://www.alibabacloud.com/help/en/oss/user-guide/upload-progress-bar?spm=a2c63.p38356.0.0.3adb6037PRDG5P#5f44b2d405h71
func (listener *Aliyun) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

func (aliyun *Aliyun) UploadFile(ctx context.Context, bucketName, fileName, contentType string, file io.Reader) error {
	if bucketName == "" {
		bucketName = aliyun.StorageBucketName
	}

	bucket, err := aliyun.Client.Bucket(bucketName)
	if err != nil {
		return err
	}

	return bucket.PutObject(fileName, file, oss.CacheControl("max-age=31536000"))
}

func (aliyun *Aliyun) GetFileTemporaryURL(ctx context.Context, bucketName, filename string) (string, error) {
	if bucketName == "" {
		bucketName = aliyun.StorageBucketName
	}

	bucket, err := aliyun.Client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	isExist, err := bucket.IsObjectExist(filename)
	if err != nil {
		return "", err
	}

	if !isExist {
		return "", errors.New("file not found")
	}

	expiredInSec := aliyun.getFileStorageExpiration()
	return bucket.SignURL(filename, oss.HTTPGet, expiredInSec, oss.ResponseCacheControl("max-age=31536000"))
}

func (aliyun *Aliyun) GetFilePublicURL(ctx context.Context, bucketName, filename string) (string, error) {
	if bucketName == "" {
		bucketName = aliyun.StoragePublicBucketName
	}

	return fmt.Sprintf("https://%s.%s/%s", bucketName, aliyun.AliyunEndpoint, filename), nil
}

func (aliyun *Aliyun) GetObject(ctx context.Context, bucketName, filename string) (io.Reader, error) {
	if bucketName == "" {
		bucketName = aliyun.StorageBucketName
	}

	bucket, err := aliyun.Client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return bucket.GetObject(filename)
}

func (aliyun *Aliyun) FGetObject(ctx context.Context, bucketName, filename, destination string) error {
	if bucketName == "" {
		bucketName = aliyun.StorageBucketName
	}

	bucket, err := aliyun.Client.Bucket(bucketName)
	if err != nil {
		return err
	}

	return bucket.GetObjectToFile(filename, destination)
}

func (aliyun *Aliyun) FPutObject(ctx context.Context, bucketName, filename, source string) error {
	if bucketName == "" {
		bucketName = aliyun.StorageBucketName
	}

	objectPath := fmt.Sprintf("%s/%s", source, filename)
	file, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		return err
	}

	bucket, err := aliyun.Client.Bucket(bucketName)
	if err != nil {
		return err
	}

	return bucket.PutObject(objectPath, bytes.NewReader(buf.Bytes()))
}

func (aliyun *Aliyun) RemoveFile(ctx context.Context, bucketName, pathFilename string) error {
	if bucketName == "" {
		bucketName = aliyun.StorageBucketName
	}

	bucket, err := aliyun.Client.Bucket(bucketName)
	if err != nil {
		return err
	}

	return bucket.DeleteObject(pathFilename)
}

func (aliyun *Aliyun) IsFileExist(ctx context.Context, bucketName, fileptah string) (bool, error) {
	if bucketName == "" {
		bucketName = aliyun.StorageBucketName
	}

	bucket, err := aliyun.Client.Bucket(bucketName)
	if err != nil {
		return false, err
	}

	return bucket.IsObjectExist(fileptah)
}

func (aliyun *Aliyun) FCopyObject(ctx context.Context, bucketName, src, dst string) error {
	if bucketName == "" {
		bucketName = aliyun.StorageBucketName
	}

	bucket, err := aliyun.Client.Bucket(bucketName)
	if err != nil {
		return err
	}

	_, err = bucket.CopyObject(src, dst)
	if err != nil {
		return err
	}

	return nil
}

// getFileStorageExpiration get temporary URL duration. The response is in Second(s).
func (aliyun *Aliyun) getFileStorageExpiration() int64 {
	expiredInSec := int64(aliyun.StorageTmpUrlExpiration)
	if expiredInSec == 0 {
		expiredInSec = 3600 // 1 Hour Default
	}
	return expiredInSec
}
