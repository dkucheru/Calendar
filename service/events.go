package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dkucheru/Calendar/structs"
)

type eventService struct {
	service    *Service
	repository []*structs.Event
}

func newEventsService(service *Service, repository []*structs.Event) *eventService {
	s := eventService{
		service:    service,
		repository: repository,
	}
	return &s
}

func (s *eventService) AddEvent(newEvent *structs.Event) error {
	//check if event has mandatory field "Name" filled
	if newEvent.Name == "" {
		return errors.New("Mandatory field *NAME* is not filled. Please, add a name to the event")
	}
	if newEvent.Start == (time.Time{}) {
		return errors.New("Mandatory field *START* is not filled. Please, add a name to the event")
	}
	if newEvent.End == (time.Time{}) {
		return errors.New("Mandatory field *END* is not filled. Please, add a name to the event")
	}

	newEvent.Id = structs.GlobalId
	structs.GlobalId++

	// Add new Event to our calendar
	s.repository = append(s.repository, newEvent)

	//debuging
	fmt.Println(newEvent)

	return nil
}

func (s *eventService) DeleteEvent(name *string, startTime *time.Time) error {
	for i, event := range s.repository {
		if event.Name == *name && event.Start == *startTime {
			s.repository = append(s.repository[:i], s.repository[i+1:]...)
			return nil
		}
	}
	return errors.New("No event with such name and start date was found")
}

func (s *eventService) UpdateEvent(id int, newEvent *structs.Event) error {
	//check if event has mandatory field "Name" filled
	if newEvent.Name == "" {
		return errors.New("Mandatory field *NAME* is not filled. Please, add a name to the event")
	}
	if newEvent.Start == (time.Time{}) {
		return errors.New("Mandatory field *START* is not filled. Please, add a name to the event")
	}
	if newEvent.End == (time.Time{}) {
		return errors.New("Mandatory field *END* is not filled. Please, add a name to the event")
	}

	for _, event := range s.repository {
		if event.Id == id {
			event.Name = newEvent.Name
			event.Start = newEvent.Start
			event.End = newEvent.End
			event.Alert = newEvent.Alert
			event.Description = newEvent.Description
			return nil
		}
	}
	return errors.New("No event with such id")
}

func (s *eventService) GetAll() (error, []*structs.Event) {
	return nil, s.repository
}
