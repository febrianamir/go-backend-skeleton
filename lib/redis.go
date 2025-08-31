package lib

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	redis.Client
}

func (r *Redis) Get(ctx context.Context, key string) (data string, err error) {
	data, err = r.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Println("error redis.Get:", err)
		return "", err
	}

	return data, nil
}

// Set is function to set data into redis, expiration time must less than file storage expiration time.
func (r *Redis) Set(ctx context.Context, key string, data any, expiration time.Duration) (err error) {
	err = r.Client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		log.Println("error redis.Set:", err)
		return err
	}

	return nil
}
func (r *Redis) Del(ctx context.Context, keys ...string) (err error) {
	err = r.Client.Del(ctx, keys...).Err()
	if err != nil {
		log.Println("error redis.Del:", err)
		return err
	}

	return nil
}
