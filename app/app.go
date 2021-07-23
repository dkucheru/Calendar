package app

import (
	"fmt"
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

func New() (*App, error) {
	var err error
	app := &App{}

	host, port, dbUser, dbPassword, dbName :=
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, dbUser, dbPassword, dbName)
	fmt.Println(dsn)
	database, err := db.Initialize(dsn)

	// app.EventsRepo, err = db.NewArrayRepository()
	app.EventsRepo, err = db.NewDatabaseRepository(database)
	if err != nil {
		return nil, err
	}

	// app.UsersRepo, err = db.NewUsersInMemoryRepository()
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
