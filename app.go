package app

import (
	"app/lib"
	"app/lib/mailer"
	"app/lib/storage"
	"app/repository"
	"app/usecase"
)

type App struct {
	Usecase *usecase.Usecase
}

func NewApp(db *lib.Database, mailer *mailer.SMTP, storage storage.Storage) *App {
	repository := repository.NewRepository(db, mailer)
	usecase := usecase.NewUsecase(&repository, storage)

	return &App{
		Usecase: &usecase,
	}
}
