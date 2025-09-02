package app

import (
	"app/lib"
	"app/lib/mailer"
	"app/lib/storage"
	"app/lib/task"
	"app/repository"
	"app/usecase"
)

type App struct {
	Usecase *usecase.Usecase
}

func NewApp(db *lib.Database, mailer *mailer.SMTP, storage storage.Storage, redis *lib.Redis, publisher *task.Publisher) *App {
	repository := repository.NewRepository(db, mailer, publisher)
	usecase := usecase.NewUsecase(&repository, storage, redis)

	return &App{
		Usecase: &usecase,
	}
}
