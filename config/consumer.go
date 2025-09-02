package config

import (
	"fmt"

	"github.com/hibiken/asynq"
)

func (c *Config) NewConsumer() *asynq.Server {
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%s", c.REDIS_HOST, c.REDIS_PORT),
			Password: c.REDIS_PASSWORD,
			DB:       1,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	return server
}
