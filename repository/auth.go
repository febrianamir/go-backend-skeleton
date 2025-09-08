package repository

import (
	"app/model"
	"app/request"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (repo *Repository) CreateAuth(ctx context.Context, auth model.UserAuth) (model.UserAuth, error) {
	ctx, tx := repo.prepareDBWithContext(ctx, "CreateAuth")

	err := tx.Create(&auth).Error
	if err != nil {
		return auth, err
	}

	return auth, nil
}

func (repo *Repository) GetAuth(ctx context.Context, req request.GetAuth) (res model.UserAuth, err error) {
	ctx, tx := repo.prepareDBWithContext(ctx, "GetAuth")

	stmt := tx.Model(&model.UserAuth{})
	if req.ID > 0 {
		stmt = stmt.Where("id = ?", req.ID)
	}

	if req.RefreshToken != "" {
		stmt = stmt.Where("refresh_token = ?", req.RefreshToken)
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

func (repo *Repository) UpdateAuth(ctx context.Context, auth model.UserAuth) (model.UserAuth, error) {
	ctx, tx := repo.prepareDBWithContext(ctx, "UpdateAuth")

	err := tx.Save(&auth).Error
	if err != nil {
		return auth, err
	}

	return auth, nil
}
