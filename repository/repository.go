package repository

import (
	"app/lib"
	"context"
)

type TrxKey struct{}

type Repository struct {
	db *lib.Database
}

func NewRepository(db *lib.Database) Repository {
	return Repository{
		db: db,
	}
}

func (repo *Repository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	trx := repo.db.Begin()

	ctx = context.WithValue(ctx, TrxKey{}, &lib.Database{DB: trx})
	if err := fn(ctx); err != nil {
		trx.Rollback()
		return err
	}

	return trx.Commit().Error
}
