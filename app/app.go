package app

import (
	"github.com/dkucheru/Calendar/api"
	"github.com/dkucheru/Calendar/service"
	"github.com/dkucheru/Calendar/structs"
)

type App struct {
	Repository []*structs.Event
	Service    *service.Service
	Api        *api.Rest
}

func New() (*App, error) {
	// var err error
	app := &App{}

	var events []*structs.Event
	app.Repository = events

	app.Service = service.NewService(&service.Config{Repository: app.Repository})

	app.Api = api.New(":8080", app.Service)
	return app, nil
}

func (a *App) Run() error {
	return a.Api.Listen()
}
