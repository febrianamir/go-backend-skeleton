package usecase

import (
	"app/lib/signoz"
	"context"
	"fmt"
)

func (usecase *Usecase) CronTest(ctx context.Context) error {
	ctx, span := signoz.StartSpan(ctx, "usecase.CronTest")
	defer span.Finish()

	fmt.Println("cron test")
	return nil
}
