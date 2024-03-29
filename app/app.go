package app

import (
	"os"

	"github.com/dkucheru/Calendar/api"
	"github.com/dkucheru/Calendar/db"
	"github.com/dkucheru/Calendar/service"
)

type App struct {
	EventsRepo db.EventsRepository
	UsersRepo  db.UserRepository
	Service    *service.Service
	Api        *api.Rest
}

func New(downMigrateFlag bool) (*App, error) {
	var err error
	app := &App{}

	database, err := db.Initialize(os.Getenv("DSN"), downMigrateFlag)
	if err != nil {
		return nil, err
	}
	app.EventsRepo, err = db.NewDatabaseRepository(database)
	if err != nil {
		return nil, err
	}

	app.UsersRepo, err = db.NewUsersDBRepository(database)
	if err != nil {
		return nil, err
	}

	app.Service = service.NewService(&service.Config{EventsRepo: app.EventsRepo, UsersRepo: app.UsersRepo})

	app.Api = api.New(":8080", app.Service)
	return app, nil
}

func (a *App) Run() error {
	return a.Api.Listen()
}

func (a *App) Stop() {
	a.Api.Stop()
}
