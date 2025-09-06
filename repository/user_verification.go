package repository

import (
	"app/model"
	"app/request"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (repo *Repository) CreateUserVerification(ctx context.Context, userVerification model.UserVerification) (model.UserVerification, error) {
	ctx, tx := repo.prepareDBWithContext(ctx, "CreateUserVerification")

	err := tx.Create(&userVerification).Error
	if err != nil {
		return userVerification, err
	}

	return userVerification, nil
}

func (repo *Repository) GetUserVerification(ctx context.Context, req request.GetUserVerification) (res model.UserVerification, err error) {
	ctx, tx := repo.prepareDBWithContext(ctx, "GetUserVerification")

	stmt := tx.Model(&model.UserVerification{})
	if req.UserID > 0 {
		stmt = stmt.Where("user_id = ?", req.UserID)
	}
	if len(req.Preloads) > 0 {
		for _, preload := range req.Preloads {
			stmt = stmt.Preload(preload)
		}
	}

	err = stmt.First(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return res, err
	}

	return res, nil
}
