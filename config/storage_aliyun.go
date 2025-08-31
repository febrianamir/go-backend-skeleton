package config

import (
	"app/lib/storage"
	"fmt"
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func (c *Config) NewAliyunClient() (*storage.Aliyun, error) {
	client, err := oss.New(
		fmt.Sprintf("https://%s", c.ALIYUN_ENDPOINT),
		c.ALIYUN_ACCESS_KEY_ID,
		c.ALIYUN_ACCESS_KEY_SECRET,
	)
	if err != nil {
		log.Println("failed connect to aliyun: ", err)
		return nil, err
	}

	log.Println("successfully connected to aliyun")
	aliyun := storage.Aliyun{
		Client:                  client,
		StorageBucketName:       c.STORAGE_BUCKET_NAME,
		StoragePublicBucketName: c.STORAGE_PUBLIC_BUCKET_NAME,
	}
	return &aliyun, nil
}
