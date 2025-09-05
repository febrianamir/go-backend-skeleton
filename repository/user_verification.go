package repository

import (
	"app/model"
	"context"
)

func (repo *Repository) CreateUserVerification(ctx context.Context, userVerification model.UserVerification) (model.UserVerification, error) {
	ctx, tx := repo.prepareDBWithContext(ctx, "CreateUserVerification")

	err := tx.Create(&userVerification).Error
	if err != nil {
		return userVerification, err
	}

	return userVerification, nil
}
