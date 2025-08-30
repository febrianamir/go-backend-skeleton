package app

import (
	"app/lib"
	"app/repository"
	"app/usecase"
)

type App struct {
	Usecase *usecase.Usecase
}

func NewApp(db *lib.Database) App {
	repository := repository.NewRepository(db)
	usecase := usecase.NewUsecase(&repository)

	return App{
		Usecase: &usecase,
	}
}
