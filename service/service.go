package service

import "github.com/dkucheru/Calendar/structs"

type Config struct {
	Repository []*structs.Event
}

type Service struct {
	repository []*structs.Event
	Events     *eventService
}

func NewService(conf *Config) *Service {
	service := &Service{
		repository: conf.Repository,
	}

	service.Events = newEventsService(service, service.repository)
	return service
}
