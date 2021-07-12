package service

import (
	"github.com/dkucheru/Calendar/db"
)

type Config struct {
	// Repository []*structs.Event
	Repository db.Repository
}

type Service struct {
	// repository []*structs.Event
	repository db.Repository
	Events     *eventService
}

func NewService(conf *Config) *Service {
	service := &Service{
		repository: conf.Repository,
	}

	service.Events = newEventsService(service.repository)
	return service
}
