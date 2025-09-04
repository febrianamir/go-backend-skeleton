package worker

import (
	"app/lib/logger"
	"app/request"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func (w *Worker) WorkerSendEmail(ctx context.Context, t *asynq.Task) error {
	ctx = context.WithValue(ctx, logger.CtxProcessID, t.ResultWriter().TaskID())
	defer recoverWorkerPanic(ctx)

	logger.LogInfo(ctx, "start process task", []zap.Field{
		zap.Any("payload", string(t.Payload())),
		zap.Strings("tags", []string{"worker", "WorkerSendEmail"}),
	}...)

	var p request.SendEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		logger.LogError(ctx, "json unmarshal error", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"worker", "WorkerSendEmail"}),
		}...)
		return err
	}

	err := w.App.Usecase.SendEmail(ctx, p)
	if err != nil {
		return err
	}

	logger.LogInfo(ctx, "success process task", []zap.Field{
		zap.Strings("tags", []string{"worker", "WorkerSendEmail"}),
	}...)
	return nil
}

func recoverWorkerPanic(ctx context.Context) {
	if r := recover(); r != nil {
		var errorMsg string
		switch err := r.(type) {
		case error:
			errorMsg = fmt.Sprintf("PANIC: %s", err.Error())
		default:
			errorMsg = fmt.Sprintf("PANIC: unknown error: %v", err)
		}
		logger.LogError(ctx, errorMsg, []zap.Field{
			zap.Strings("tags", []string{"worker", "WorkerSendEmail"}),
		}...)
	}
}
