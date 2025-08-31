package config

import (
	"app/lib/storage"
)

func (c *Config) NewStorage() storage.Storage {
	switch c.STORAGE_VENDOR {
	case "aliyun":
		aliyunClient, _ := c.NewAliyunClient()
		return aliyunClient
	case "minio":
		minioClient := c.NewMinioClient()
		return minioClient
	default:
		return NewLocalStorage()
	}
}
