package repository

import (
	"app/lib/auth"
	"app/lib/constant"
	"app/lib/logger"
	"app/lib/signoz"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func (repo *Repository) SetVerificationDelayCache(ctx context.Context, userId uint, verificationType string) error {
	ctx, span := signoz.StartSpan(ctx, "repository.SetVerificationDelayCache")
	defer span.Finish()

	sendVerificationDelayKey := fmt.Sprintf(constant.SendVerificationDelayKeyPrefix, userId, verificationType)
	return repo.cache.Set(ctx, sendVerificationDelayKey, "default", time.Duration(repo.config.SEND_VERIFICATION_DELAY_TTL)*time.Second)
}

func (repo *Repository) GetVerificationDelayCacheWithTtl(ctx context.Context, userId uint, verificationType string) (string, time.Duration, error) {
	ctx, span := signoz.StartSpan(ctx, "repository.GetVerificationDelayCacheWithTtl")
	defer span.Finish()

	sendVerificationDelayKey := fmt.Sprintf(constant.SendVerificationDelayKeyPrefix, userId, verificationType)
	return repo.cache.GetWithTtl(ctx, sendVerificationDelayKey)
}

func (repo *Repository) SetMfaFlag(ctx context.Context, userId uint) error {
	ctx, span := signoz.StartSpan(ctx, "repository.SetMfaFlag")
	defer span.Finish()

	mfaFlagKey := fmt.Sprintf(constant.MfaFlagKeyPrefix, userId)
	return repo.cache.Set(ctx, mfaFlagKey, "default", time.Duration(repo.config.MFA_FLAG_TTL)*time.Second)
}

func (repo *Repository) GetMfaFlag(ctx context.Context, userId uint) (string, error) {
	ctx, span := signoz.StartSpan(ctx, "repository.GetMfaFlag")
	defer span.Finish()

	mfaFlagKey := fmt.Sprintf(constant.MfaFlagKeyPrefix, userId)
	return repo.cache.Get(ctx, mfaFlagKey)
}

func (repo *Repository) SetAccessToken(ctx context.Context, accessToken string, claims auth.AccessTokenClaims) error {
	ctx, span := signoz.StartSpan(ctx, "repository.SetAccessToken")
	defer span.Finish()

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
	ctx, span := signoz.StartSpan(ctx, "repository.GetAccessToken")
	defer span.Finish()

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

func (repo *Repository) GetSendOtpRateLimitCtrWithTtl(ctx context.Context, identifier uint, otpType string) (int, time.Duration, error) {
	ctx, span := signoz.StartSpan(ctx, "repository.GetSendOtpRateLimitCtrWithTtl")
	defer span.Finish()

	otpRateLimitCtrKey := fmt.Sprintf(constant.SendOtpCtrKeyPrefix, identifier, otpType)
	ctrString, duration, err := repo.cache.GetWithTtl(ctx, otpRateLimitCtrKey)
	if err != nil {
		return 0, 0, err
	}
	if ctrString == "" {
		return 0, 0, nil
	}

	ctr, err := strconv.Atoi(ctrString)
	if err != nil {
		logger.LogError(ctx, "error strconv.Atoi", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"repository", "GetSendOtpRateLimitCtrWithTtl"}),
		}...)
		return 0, 0, err
	}

	return ctr, duration, nil
}

func (repo *Repository) IncrSendOtpRateLimitCtr(ctx context.Context, identifier uint, otpType string) (int64, error) {
	ctx, span := signoz.StartSpan(ctx, "repository.IncrSendOtpRateLimitCtr")
	defer span.Finish()

	otpRateLimitCtrKey := fmt.Sprintf(constant.SendOtpCtrKeyPrefix, identifier, otpType)
	return repo.cache.Incr(ctx, otpRateLimitCtrKey)
}

func (repo *Repository) ExpSendOtpRateLimitCtr(ctx context.Context, identifier uint, otpType string) error {
	ctx, span := signoz.StartSpan(ctx, "repository.ExpSendOtpRateLimitCtr")
	defer span.Finish()

	otpRateLimitCtrKey := fmt.Sprintf(constant.SendOtpCtrKeyPrefix, identifier, otpType)
	return repo.cache.Expire(ctx, otpRateLimitCtrKey, time.Duration(repo.config.SEND_OTP_MAX_RATE_LIMIT_TTL)*time.Second)
}

func (repo *Repository) TtlSendOtpDelay(ctx context.Context, identifier uint, otpType string) (time.Duration, error) {
	ctx, span := signoz.StartSpan(ctx, "repository.TtlSendOtpDelay")
	defer span.Finish()

	sendOtpDelayKey := fmt.Sprintf(constant.SendOtpDelayKeyPrefix, identifier, otpType)
	return repo.cache.TTL(ctx, sendOtpDelayKey)
}

func (repo *Repository) SetSendOtpDelay(ctx context.Context, identifier uint, otpType string) error {
	ctx, span := signoz.StartSpan(ctx, "repository.SetSendOtpDelay")
	defer span.Finish()

	sendOtpDelayKey := fmt.Sprintf(constant.SendOtpDelayKeyPrefix, identifier, otpType)
	return repo.cache.Set(ctx, sendOtpDelayKey, "default", time.Duration(repo.config.SEND_OTP_DELAY_TTL)*time.Second)
}
