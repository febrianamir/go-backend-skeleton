package app

import (
	"app/config"
	"app/lib"
	"app/lib/cache"
	"app/lib/mailer"
	"app/lib/storage"
	"app/lib/task"
	"app/lib/websocket"
	"app/repository"
	"app/usecase"
)

type App struct {
	Usecase *usecase.Usecase
}

func NewApp(config *config.Config, db *lib.Database, mailer *mailer.SMTP, storage storage.Storage, cache *cache.Cache, publisher *task.Publisher, wsPool *websocket.WebsocketPool) *App {
	repository := repository.NewRepository(config, db, mailer, publisher, cache, wsPool)
	usecase := usecase.NewUsecase(config, &repository, storage)

	return &App{
		Usecase: &usecase,
	}
}
