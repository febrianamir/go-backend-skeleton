package app

import (
	"app/config"
	"app/lib"
	"app/lib/cache"
	"app/lib/mailer"
	"app/lib/storage"
	"app/lib/task"
	"app/repository"
	"app/usecase"
)

type App struct {
	Usecase *usecase.Usecase
}

func NewApp(config *config.Config, db *lib.Database, mailer *mailer.SMTP, storage storage.Storage, cache *cache.Cache, publisher *task.Publisher) *App {
	repository := repository.NewRepository(config, db, mailer, publisher, cache)
	usecase := usecase.NewUsecase(&repository, storage)

	return &App{
		Usecase: &usecase,
	}
}
