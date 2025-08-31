package usecase

import (
	"app/lib/storage"
	"app/repository"
)

type Usecase struct {
	repo    *repository.Repository
	storage storage.Storage
}

func NewUsecase(repo *repository.Repository, storage storage.Storage) Usecase {
	return Usecase{
		repo:    repo,
		storage: storage,
	}
}
