package config

import (
	"app/lib"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func (c *Config) NewRedis() *lib.Redis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.REDIS_HOST, c.REDIS_PORT),
		Password: c.REDIS_PASSWORD,
		DB:       0,
	})

	ctx := context.Background()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Println("failed connect to redis: ", err)
	} else {
		log.Println("success connect to redis")
	}

	return &lib.Redis{Client: *redisClient}
}
