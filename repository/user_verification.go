package repository

import (
	"app/model"
	"app/request"
	"context"
	"errors"

	"gorm.io/gorm"
)

func (repo *Repository) CreateUserVerification(ctx context.Context, userVerification model.UserVerification) (model.UserVerification, error) {
	ctx, span, tx := repo.prepareRepoContext(ctx, "CreateUserVerification")
	defer span.Finish()

	err := tx.Create(&userVerification).Error
	if err != nil {
		return userVerification, err
	}

	return userVerification, nil
}

func (repo *Repository) GetUserVerification(ctx context.Context, req request.GetUserVerification) (res model.UserVerification, err error) {
	ctx, span, tx := repo.prepareRepoContext(ctx, "GetUserVerification")
	defer span.Finish()

	stmt := tx.Model(&model.UserVerification{})
	if req.Type != "" {
		stmt = stmt.Where("type = ?", req.Type)
	}

	if req.UserID > 0 {
		stmt = stmt.Where("user_id = ?", req.UserID)
	}

	if req.Code != "" {
		stmt = stmt.Where("code = ?", req.Code)
	}

	if req.IsUsed != nil {
		if *req.IsUsed {
			stmt = stmt.Where("used_at IS NOT NULL")
		} else {
			stmt = stmt.Where("used_at IS NULL")
		}
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

func (repo *Repository) UpdateUserVerification(ctx context.Context, userVerification model.UserVerification) (model.UserVerification, error) {
	ctx, span, tx := repo.prepareRepoContext(ctx, "UpdateUserVerification")
	defer span.Finish()

	err := tx.Save(&userVerification).Error
	if err != nil {
		return userVerification, err
	}

	return userVerification, nil
}
