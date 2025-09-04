package repository

import (
	"app/lib/logger"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// PublishTask create new asynq task and publish it
func (repo *Repository) PublishTask(ctx context.Context, taskType string, payload any) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logger.LogError(ctx, "failed marshal payload", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"repository", "PublishTask"}),
		}...)
		return err
	}

	task := asynq.NewTask(taskType, jsonPayload)
	taskInfo, err := repo.publisher.Enqueue(task)
	if err != nil {
		logger.LogError(ctx, "failed enqueue task", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"repository", "PublishTask"}),
		}...)
		return err
	}

	logger.LogInfo(ctx, "success publish task", []zap.Field{
		zap.String("process_id", taskInfo.ID),
		zap.String("task_queue", taskInfo.Queue),
		zap.Any("payload", payload),
		zap.Strings("tags", []string{"repository", "PublishTask"}),
	}...)
	return nil
}
