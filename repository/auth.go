package repository

import (
	"app/model"
	"context"
)

func (repo *Repository) CreateAuth(ctx context.Context, auth model.UserAuth) (model.UserAuth, error) {
	ctx, tx := repo.prepareDBWithContext(ctx, "CreateAuth")

	err := tx.Create(&auth).Error
	if err != nil {
		return auth, err
	}

	return auth, nil
}
