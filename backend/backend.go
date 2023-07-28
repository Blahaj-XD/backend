package backend

import (
	"github.com/BlahajXD/backend/repo"
	"github.com/gojek/heimdall/v7/httpclient"
)

type Dependency struct {
	repo       *repo.Dependency
	httpclient *httpclient.Client
}

func New(repo *repo.Dependency, httpclient *httpclient.Client) *Dependency {
	return &Dependency{repo: repo, httpclient: httpclient}
}
