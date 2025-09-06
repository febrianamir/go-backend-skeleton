package repository

import (
	"app/lib/constant"
	"context"
	"fmt"
	"time"
)

func (repo *Repository) SetVerificationDelayCache(ctx context.Context, userId uint, verificationType string) error {
	sendVerificationDelayKey := fmt.Sprintf(constant.SendVerificationDelayKeyPrefix, userId, verificationType)
	return repo.cache.Set(ctx, sendVerificationDelayKey, "default", time.Duration(repo.config.SEND_VERIFICATION_DELAY_TTL)*time.Second)
}

func (repo *Repository) GetVerificationDelayCacheWithTtl(ctx context.Context, userId uint, verificationType string) (string, time.Duration, error) {
	sendVerificationDelayKey := fmt.Sprintf(constant.SendVerificationDelayKeyPrefix, userId, verificationType)
	return repo.cache.GetWithTtl(ctx, sendVerificationDelayKey)
}
