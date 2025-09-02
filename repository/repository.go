package repository

import (
	"app/lib"
	"app/lib/mailer"
	"app/lib/task"
	"context"
)

type TrxKey struct{}

type Repository struct {
	db        *lib.Database
	mailer    *mailer.SMTP
	publisher *task.Publisher
}

func NewRepository(db *lib.Database, mailer *mailer.SMTP, publisher *task.Publisher) Repository {
	return Repository{
		db:        db,
		mailer:    mailer,
		publisher: publisher,
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
