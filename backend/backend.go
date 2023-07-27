package backend

import (
	"github.com/BlahajXD/backend/repo"
)

type Dependency struct {
	repo *repo.Dependency
}

func New(repo *repo.Dependency) *Dependency {
	return &Dependency{repo: repo}
}
