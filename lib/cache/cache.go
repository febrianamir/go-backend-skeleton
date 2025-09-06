package cache

import (
	"app/lib/logger"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Cache struct {
	*redis.Client
}

func (r *Cache) Get(ctx context.Context, key string) (data string, err error) {
	data, err = r.Client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.LogError(ctx, "error cache.Get", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"cache", "Get"}),
		}...)
		return "", err
	}

	return data, nil
}

func (r *Cache) GetWithTtl(ctx context.Context, key string) (string, time.Duration, error) {
	pipe := r.Client.Pipeline()

	get := pipe.Get(ctx, key)
	ttl := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.LogError(ctx, "error cache.Exec", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"cache", "GetWithTtl"}),
		}...)
		return "", 0, err
	}

	data, err := get.Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.LogError(ctx, "error cache.Get", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"cache", "GetWithTtl"}),
		}...)
		return "", 0, err
	}

	ttlDuration, err := ttl.Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.LogError(ctx, "error cache.TTL", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"cache", "GetWithTtl"}),
		}...)
		return "", 0, err
	}

	return data, ttlDuration, nil
}

// Set is function to set data into redis, expiration time must less than file storage expiration time.
func (r *Cache) Set(ctx context.Context, key string, data any, expiration time.Duration) (err error) {
	err = r.Client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		logger.LogError(ctx, "error cache.Set", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"cache", "Set"}),
		}...)
		return err
	}

	return nil
}

func (r *Cache) Del(ctx context.Context, keys ...string) (err error) {
	err = r.Client.Del(ctx, keys...).Err()
	if err != nil {
		logger.LogError(ctx, "error cache.Del", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"cache", "Del"}),
		}...)
		return err
	}

	return nil
}
