package repository

import (
	"app/lib/logger"
	"app/lib/signoz"
	"app/model"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (repo *Repository) GetAppSettingByName(ctx context.Context, name string) (model.AppSetting, error) {
	ctx, span := signoz.StartSpan(ctx, "repository.GetAppSettingByName")
	defer span.Finish()

	var appSetting model.AppSetting
	err := repo.db.Where("name = ?", name).First(&appSetting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogError(ctx, "AppSetting not found", []zap.Field{
				zap.String("name", name),
				zap.Strings("tags", []string{"postgres", "app_setting", "repo"}),
			}...)
		} else {
			logger.LogError(ctx, "Error GetAppSettingByName", []zap.Field{
				zap.Error(err),
				zap.String("name", name),
				zap.Strings("tags", []string{"postgres", "app_setting", "repo"}),
			}...)
		}
		return appSetting, err
	}

	return appSetting, nil
}

func (repo *Repository) GetAppSettingBySlug(ctx context.Context, slug string) (model.AppSetting, error) {
	ctx, span := signoz.StartSpan(ctx, "repository.GetAppSettingBySlug")
	defer span.Finish()

	var appSetting model.AppSetting
	err := repo.db.Where("slug = ?", slug).First(&appSetting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogError(ctx, "AppSetting not found", []zap.Field{
				zap.String("slug", slug),
				zap.Strings("tags", []string{"postgres", "app_setting", "repo"}),
			}...)
		} else {
			logger.LogError(ctx, "Error GetAppSettingBySlug", []zap.Field{
				zap.Error(err),
				zap.String("slug", slug),
				zap.Strings("tags", []string{"postgres", "app_setting", "repo"}),
			}...)
		}
		return appSetting, err
	}

	return appSetting, nil
}
