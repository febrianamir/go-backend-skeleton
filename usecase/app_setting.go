package usecase

import (
	"app/lib/logger"
	"context"
	"strconv"

	"go.uber.org/zap"
)

// getAppSettingUint get data with uint data type from app_settings table
func (usecase *Usecase) getAppSettingUint(ctx context.Context, slug string) (uint, error) {
	appSetting, err := usecase.repo.GetAppSettingBySlug(ctx, slug)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseUint(appSetting.Value, 10, 64)
	if err != nil {
		logger.LogError(ctx, "Error getAppSettingUint", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"app_setting", "slug", slug}),
		}...)
		return 0, err
	}

	return uint(value), nil
}

// getAppSettingFloat get data with uint data type from app_settings table
func (usecase *Usecase) getAppSettingFloat(ctx context.Context, slug string) (float64, error) {
	appSetting, err := usecase.repo.GetAppSettingBySlug(ctx, slug)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(appSetting.Value, 64)
	if err != nil {
		logger.LogError(ctx, "Error getAppSettingFloat", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"app_setting", "slug", slug}),
		}...)
		return 0, err
	}

	return value, nil
}

// getAppSettingString get data with string data type from app_settings table
func (usecase *Usecase) getAppSettingString(ctx context.Context, slug string) (string, error) {
	appSetting, err := usecase.repo.GetAppSettingBySlug(ctx, slug)
	if err != nil {
		return "", err
	}
	return appSetting.Value, nil
}
