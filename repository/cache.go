package repository

import (
	"app/lib/auth"
	"app/lib/constant"
	"app/lib/logger"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func (repo *Repository) SetVerificationDelayCache(ctx context.Context, userId uint, verificationType string) error {
	sendVerificationDelayKey := fmt.Sprintf(constant.SendVerificationDelayKeyPrefix, userId, verificationType)
	return repo.cache.Set(ctx, sendVerificationDelayKey, "default", time.Duration(repo.config.SEND_VERIFICATION_DELAY_TTL)*time.Second)
}

func (repo *Repository) GetVerificationDelayCacheWithTtl(ctx context.Context, userId uint, verificationType string) (string, time.Duration, error) {
	sendVerificationDelayKey := fmt.Sprintf(constant.SendVerificationDelayKeyPrefix, userId, verificationType)
	return repo.cache.GetWithTtl(ctx, sendVerificationDelayKey)
}

func (repo *Repository) SetMfaFlag(ctx context.Context, userId uint) error {
	mfaFlagKey := fmt.Sprintf(constant.MfaFlagKeyPrefix, userId)
	return repo.cache.Set(ctx, mfaFlagKey, "default", time.Duration(repo.config.MFA_FLAG_TTL)*time.Second)
}

func (repo *Repository) GetMfaFlag(ctx context.Context, userId uint) (string, error) {
	mfaFlagKey := fmt.Sprintf(constant.MfaFlagKeyPrefix, userId)
	return repo.cache.Get(ctx, mfaFlagKey)
}

func (repo *Repository) SetAccessToken(ctx context.Context, accessToken string, claims auth.AccessTokenClaims) error {
	accessTokenKey := fmt.Sprintf(constant.AccessTokenKeyPrefix, accessToken)

	data, err := json.Marshal(claims)
	if err != nil {
		logger.LogError(ctx, "error json.Marshal", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"repository", "SetAccessToken"}),
		}...)
		return err
	}

	return repo.cache.Set(ctx, accessTokenKey, string(data), time.Duration(repo.config.ACCESS_TOKEN_TTL)*time.Second)
}

func (repo *Repository) GetAccessToken(ctx context.Context, accessToken string) (auth.AccessTokenClaims, error) {
	accessTokenKey := fmt.Sprintf(constant.AccessTokenKeyPrefix, accessToken)
	accessTokenBytes, err := repo.cache.GetBytes(ctx, accessTokenKey)
	if err != nil {
		return auth.AccessTokenClaims{}, err
	}

	var claims auth.AccessTokenClaims
	err = json.Unmarshal(accessTokenBytes, &claims)
	if err != nil {
		logger.LogError(ctx, "error json.Unmarshal", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"repository", "GetAccessToken"}),
		}...)
		return auth.AccessTokenClaims{}, err
	}

	return claims, nil
}
