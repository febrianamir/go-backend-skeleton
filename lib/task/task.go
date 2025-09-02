package task

import "github.com/hibiken/asynq"

type Publisher struct {
	*asynq.Client
}
