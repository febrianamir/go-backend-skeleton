package scheduler

import (
	"context"
	"fmt"
	"log"

	"app"
	"app/lib"
	"app/lib/logger"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type Scheduler struct {
	gocron.Scheduler
	App *app.App
}

func NewScheduler(app *app.App) (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		Scheduler: s,
		App:       app,
	}, nil
}

func (s *Scheduler) RegisterJob(jobDefinition gocron.JobDefinition, cronName string, fn func(ctx context.Context) error) {
	_, err := s.NewJob(jobDefinition, gocron.NewTask(func(ctx context.Context) error {
		ctx = context.WithValue(ctx, logger.CtxProcessID, lib.GenerateUUID())
		defer recoverCronPanic(ctx, cronName)

		logger.LogInfo(ctx, "start process cron", []zap.Field{
			zap.Strings("tags", []string{"cron", "CronTest"}),
		}...)

		if err := fn(ctx); err != nil {
			logger.LogError(ctx, "error process cron", []zap.Field{
				zap.Error(err),
				zap.Strings("tags", []string{"cron", cronName}),
			}...)
			return err
		}

		logger.LogInfo(ctx, "success process cron", []zap.Field{
			zap.Strings("tags", []string{"cron", "CronTest"}),
		}...)
		return nil
	}))
	if err != nil {
		log.Fatal("failed to register job: ", err)
	}
}

func recoverCronPanic(ctx context.Context, cronName string) {
	if r := recover(); r != nil {
		var errorMsg string
		switch err := r.(type) {
		case error:
			errorMsg = fmt.Sprintf("PANIC: %s", err.Error())
		default:
			errorMsg = fmt.Sprintf("PANIC: unknown error: %v", err)
		}
		logger.LogError(ctx, errorMsg, []zap.Field{
			zap.Strings("tags", []string{"cron", cronName}),
		}...)
	}
}
