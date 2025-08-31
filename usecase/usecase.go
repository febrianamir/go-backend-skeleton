package usecase

import (
	"app/lib"
	"app/lib/storage"
	"app/repository"
)

type Usecase struct {
	repo    *repository.Repository
	storage storage.Storage
	redis   *lib.Redis
}

func NewUsecase(repo *repository.Repository, storage storage.Storage, redis *lib.Redis) Usecase {
	return Usecase{
		repo:    repo,
		storage: storage,
		redis:   redis,
	}
}
