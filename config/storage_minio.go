package config

import (
	"app/lib/storage"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (c *Config) NewMinioClient() *storage.Minio {
	client, err := minio.New(c.MINIO_ENDPOINT, &minio.Options{
		Creds:  credentials.NewStaticV4(c.MINIO_USERNAME, c.MINIO_PASSWORD, ""),
		Secure: c.MINIO_SSL,
	})
	if err != nil {
		log.Fatal("failed connect to minio: ", err)
	}

	log.Println("success connect to minio")
	return &storage.Minio{Client: client}
}
