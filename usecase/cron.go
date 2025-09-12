package usecase

import (
	"context"
	"fmt"
)

func (usecase *Usecase) CronTest(ctx context.Context) error {
	fmt.Println("cron test")
	return nil
}
