package usecase

import (
	"app/config"
	"app/lib/storage"
	"app/repository"
)

type Usecase struct {
	config  *config.Config
	repo    *repository.Repository
	storage storage.Storage
}

func NewUsecase(config *config.Config, repo *repository.Repository, storage storage.Storage) Usecase {
	return Usecase{
		config:  config,
		repo:    repo,
		storage: storage,
	}
}
