package worker

import (
	"context"
	"errors"
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

func (s *Worker) RegisterWorker(mux *asynq.ServeMux, taskType, taskName string, skipRetry bool, fn func(ctx context.Context, t *asynq.Task) error) {
	mux.HandleFunc(taskType, func(ctx context.Context, t *asynq.Task) (err error) {
		ctx = context.WithValue(ctx, logger.CtxProcessID, t.ResultWriter().TaskID())
		defer func() {
			if r := recover(); r != nil {
				panicErr := handleWorkerPanic(ctx, taskName, r)
				if skipRetry {
					err = fmt.Errorf("%v: %w", panicErr, asynq.SkipRetry)
				} else {
					err = panicErr
				}
			}
		}()

		logger.LogInfo(ctx, "start process task", []zap.Field{
			zap.Any("payload", string(t.Payload())),
			zap.Strings("tags", []string{"worker", taskName}),
		}...)

		err = fn(ctx, t)
		if err != nil {
			logger.LogError(ctx, "process task error", []zap.Field{
				zap.Error(err),
				zap.Strings("tags", []string{"worker", taskName}),
			}...)
			if skipRetry {
				return fmt.Errorf("%v: %w", err, asynq.SkipRetry)
			}
			return err
		}

		logger.LogInfo(ctx, "success process task", []zap.Field{
			zap.Strings("tags", []string{"worker", taskName}),
		}...)
		return nil
	})
}

func handleWorkerPanic(ctx context.Context, taskName string, panicValue any) (err error) {
	var errorMsg string
	switch err := panicValue.(type) {
	case error:
		errorMsg = fmt.Sprintf("PANIC: %s", err.Error())
	default:
		errorMsg = fmt.Sprintf("PANIC: unknown error: %v", err)
	}

	logger.LogError(ctx, errorMsg, []zap.Field{
		zap.Strings("tags", []string{"worker", taskName}),
	}...)

	return errors.New(errorMsg)
}
