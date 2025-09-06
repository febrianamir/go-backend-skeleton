package config

import (
	"app/lib/cache"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func (c *Config) NewCache() (*cache.Cache, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.REDIS_HOST, c.REDIS_PORT),
		Password: c.REDIS_PASSWORD,
		DB:       0,
	})

	ctx := context.Background()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	log.Println("success connect to redis")
	return &cache.Cache{Client: redisClient}, nil
}
