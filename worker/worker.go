package worker

import (
	"context"
	"fmt"

	"app"
	"app/lib/logger"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type Worker struct {
	App *app.App
}

func NewWorker(a *app.App) Worker {
	return Worker{App: a}
}

func (s *Worker) RegisterWorker(mux *asynq.ServeMux, taskType, taskName string, fn func(ctx context.Context, t *asynq.Task) error) {
	mux.HandleFunc(taskType, func(ctx context.Context, t *asynq.Task) error {
		ctx = context.WithValue(ctx, logger.CtxProcessID, t.ResultWriter().TaskID())
		defer recoverWorkerPanic(ctx, taskName)

		logger.LogInfo(ctx, "start process task", []zap.Field{
			zap.Any("payload", string(t.Payload())),
			zap.Strings("tags", []string{"worker", taskName}),
		}...)

		err := fn(ctx, t)
		if err != nil {
			logger.LogError(ctx, "process task error", []zap.Field{
				zap.Error(err),
				zap.Strings("tags", []string{"worker", taskName}),
			}...)
			return err
		}

		logger.LogInfo(ctx, "success process task", []zap.Field{
			zap.Strings("tags", []string{"worker", taskName}),
		}...)
		return nil
	})
}

func recoverWorkerPanic(ctx context.Context, taskName string) {
	if r := recover(); r != nil {
		var errorMsg string
		switch err := r.(type) {
		case error:
			errorMsg = fmt.Sprintf("PANIC: %s", err.Error())
		default:
			errorMsg = fmt.Sprintf("PANIC: unknown error: %v", err)
		}
		logger.LogError(ctx, errorMsg, []zap.Field{
			zap.Strings("tags", []string{"worker", taskName}),
		}...)
	}
}
