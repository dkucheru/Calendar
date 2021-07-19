package app

import (
	"github.com/dkucheru/Calendar/api"
	"github.com/dkucheru/Calendar/db"
	"github.com/dkucheru/Calendar/service"
)

type App struct {
	Repository db.Repository
	Service    *service.Service
	Api        *api.Rest
}

func New() (*App, error) {
	var err error
	app := &App{}

	app.Repository, err = db.NewArrayRepository()
	if err != nil {
		return nil, err
	}

	app.Service = service.NewService(&service.Config{Repository: app.Repository})

	app.Api = api.New(":8080", app.Service)
	return app, nil
}

func (a *App) Run() error {
	return a.Api.Listen()
}
