package repository

import (
	"app/model"
	"app/request"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (repo *Repository) GetUsers(ctx context.Context, req request.GetUsers) (res []model.User, total int64, err error) {
	ctx, span, tx := repo.prepareRepoContext(ctx, "GetUsers")
	defer span.Finish()

	stmt := tx.Model(&model.User{})
	if req.Search != "" {
		search := fmt.Sprintf("%s%s%s", "%", req.Search, "%")
		stmt = stmt.Where("name ILIKE ? OR email ILIKE ? OR phone_number ILIKE ?", search, search)
	}

	err = stmt.Count(&total).Error
	if err != nil {
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
		return res, total, err
	}

	return res, total, nil
}

func (repo *Repository) GetUser(ctx context.Context, req request.GetUser) (res model.User, err error) {
	ctx, span, tx := repo.prepareRepoContext(ctx, "GetUser")
	defer span.Finish()

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
		return res, err
	}

	return res, nil
}

func (repo *Repository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	ctx, span, tx := repo.prepareRepoContext(ctx, "CreateUser")
	defer span.Finish()

	err := tx.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *Repository) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	ctx, span, tx := repo.prepareRepoContext(ctx, "UpdateUser")
	defer span.Finish()

	err := tx.Save(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *Repository) DeleteUser(ctx context.Context, id uint) error {
	ctx, span, tx := repo.prepareRepoContext(ctx, "DeleteUser")
	defer span.Finish()

	var user model.User
	err := tx.Delete(&user, id).Error
	if err != nil {
		return err
	}

	return nil
}
