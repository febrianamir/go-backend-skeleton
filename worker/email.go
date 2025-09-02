package worker

import (
	"app/lib/logger"
	"app/request"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func (w *Worker) WorkerSendEmail(ctx context.Context, t *asynq.Task) error {
	var p request.SendEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		logger.LogError(ctx, "json unmarshal error", []zap.Field{
			zap.Error(err),
			zap.String("task_id", t.ResultWriter().TaskID()),
			zap.Strings("tags", []string{"worker", "WorkerSendEmail"}),
		}...)
		return err
	}
	logger.LogInfo(ctx, "start process task", []zap.Field{
		zap.String("task_id", t.ResultWriter().TaskID()),
		zap.Any("payload", p),
		zap.Strings("tags", []string{"worker", "WorkerSendEmail"}),
	}...)

	err := w.App.Usecase.SendEmail(ctx, p)
	if err != nil {
		return err
	}

	logger.LogInfo(ctx, "success process task", []zap.Field{
		zap.String("task_id", t.ResultWriter().TaskID()),
		zap.Strings("tags", []string{"worker", "WorkerSendEmail"}),
	}...)
	return nil
}
