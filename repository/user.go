package repository

import (
	"app/lib"
	"app/lib/logger"
	"app/model"
	"app/request"
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (repo *Repository) GetUsers(ctx context.Context, req request.GetUsers) (res []model.User, total int64, err error) {
	tx, ok := ctx.Value(TrxKey{}).(*lib.Database)
	if !ok {
		tx = repo.db
	}

	stmt := tx.Model(&model.User{})
	if req.Search != "" {
		search := fmt.Sprintf("%s%s%s", "%", req.Search, "%")
		stmt = stmt.Where("name ILIKE ? OR email ILIKE ? OR phone_number ILIKE ?", search, search)
	}

	err = stmt.Count(&total).Error
	if err != nil {
		logger.LogError(ctx, "Error Count GetUsers", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"postgres", "user", "repo"}),
		}...)
		return res, total, err
	}

	if req.GetOrderQuery() != "" {
		stmt = stmt.Order(req.GetOrderQuery())
	}

	if req.Limit > 0 {
		stmt = stmt.Limit(int(req.Limit))
	}

	if req.GetOffset() > 0 {
		stmt = stmt.Offset(int(req.GetOffset()))
	}

	if len(req.Preloads) > 0 {
		for _, preload := range req.Preloads {
			stmt = stmt.Preload(preload)
		}
	}

	err = stmt.Find(&res).Error
	if err != nil {
		logger.LogError(ctx, "Error GetUsers", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"postgres", "user", "repo"}),
		}...)
		return res, total, err
	}

	return res, total, nil
}

func (repo *Repository) GetUser(ctx context.Context, req request.GetUser) (res model.User, err error) {
	tx, ok := ctx.Value(TrxKey{}).(*lib.Database)
	if !ok {
		tx = repo.db
	}

	stmt := tx.Model(&model.User{})
	if req.ID > 0 {
		stmt = stmt.Where("id = ?", req.ID)
	}

	if req.Name != "" {
		stmt = stmt.Where("name = ?", req.Name)
	}

	if req.Email != "" {
		stmt = stmt.Where("email = ?", req.Email)
	}

	if len(req.Preloads) > 0 {
		for _, preload := range req.Preloads {
			stmt = stmt.Preload(preload)
		}
	}

	err = stmt.First(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.LogError(ctx, "Error GetUser", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"postgres", "user", "repo"}),
		}...)
		return res, err
	}

	return res, nil
}

func (repo *Repository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	tx, ok := ctx.Value(TrxKey{}).(*lib.Database)
	if !ok {
		tx = repo.db
	}

	err := tx.Create(&user).Error
	if err != nil {
		logger.LogError(ctx, "Error CreateUser", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"postgres", "user", "repo"}),
		}...)
		return user, err
	}

	return user, nil
}

func (repo *Repository) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	tx, ok := ctx.Value(TrxKey{}).(*lib.Database)
	if !ok {
		tx = repo.db
	}

	err := tx.Save(&user).Error
	if err != nil {
		logger.LogError(ctx, "Error UpdateUser", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"postgres", "user", "repo"}),
		}...)
		return user, err
	}

	return user, nil
}

func (repo *Repository) DeleteUser(ctx context.Context, id uint) error {
	tx, ok := ctx.Value(TrxKey{}).(*lib.Database)
	if !ok {
		tx = repo.db
	}

	var user model.User
	err := tx.Delete(&user, id).Error
	if err != nil {
		logger.LogError(ctx, "Error DeleteUser", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"postgres", "user", "repo"}),
		}...)
		return err
	}

	return nil
}
