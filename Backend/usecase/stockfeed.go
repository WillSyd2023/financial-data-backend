package usecase

import "Backend/repo"

type UsecaseItf interface{}

type Usecase struct {
	rp repo.RepoItf
}

func NewUsecase(rp repo.Repo) *Usecase {
	return &Usecase{
		rp: rp,
	}
}
