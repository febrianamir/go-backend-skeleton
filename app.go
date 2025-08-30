package app

import (
	"app/lib"
	"app/lib/mailer"
	"app/repository"
	"app/usecase"
)

type App struct {
	Usecase *usecase.Usecase
}

func NewApp(db *lib.Database, mailer *mailer.SMTP) *App {
	repository := repository.NewRepository(db, mailer)
	usecase := usecase.NewUsecase(&repository)

	return &App{
		Usecase: &usecase,
	}
}
