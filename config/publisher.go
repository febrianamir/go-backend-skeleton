package config

import (
	"app/lib/task"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

func (c *Config) NewPublisher() (*task.Publisher, error) {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", c.REDIS_HOST, c.REDIS_PORT),
		Password: c.REDIS_PASSWORD,
		DB:       1,
	})

	if err := client.Ping(); err != nil {
		return nil, err
	}

	log.Println("success connect to publisher")
	return &task.Publisher{Client: client}, nil
}
