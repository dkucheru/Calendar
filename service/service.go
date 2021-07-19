package service

import (
	"github.com/dkucheru/Calendar/db"
)

type Config struct {
	EventsRepo db.EventsRepository
	UsersRepo  db.UserRepository
}

type Service struct {
	eventsRepo db.EventsRepository
	usersRepo  db.UserRepository
	Events     *eventService
	Users      *usersService
}

func NewService(conf *Config) *Service {
	service := &Service{
		eventsRepo: conf.EventsRepo,
		usersRepo:  conf.UsersRepo,
	}

	service.Events = newEventsService(service.eventsRepo)
	service.Users = newUsersService(service.usersRepo)
	return service
}
