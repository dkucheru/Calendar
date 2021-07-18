package service

import (
	"errors"
	"sort"
	"time"

	"github.com/dkucheru/Calendar/db"
	"github.com/dkucheru/Calendar/structs"
)

type eventService struct {
	repository db.Repository
}

func newEventsService(repository db.Repository) *eventService {
	s := eventService{
		repository: repository,
	}
	return &s
}

func (s *eventService) AddEvent(newEvent structs.Event) (structs.Event, error) {
	approved, err := s.checkData(newEvent)
	if !approved {
		return structs.Event{}, err
	}
	newEvent.Id = s.repository.GetNextId()
	s.repository.AddEvent(newEvent)
	return newEvent, nil
}

// FIXME : separate user service
func (s *eventService) AddUser(newUser structs.CreateUser) error {
	return s.repository.AddUser(newUser)
}

func (s *eventService) CheckPassword(user string, pass string) error {
	return s.repository.CheckCredentials(user, pass)
}

func (s *eventService) GetUserLocation(username string) (time.Location, error) {
	userInfo, err := s.repository.GetUser(username)
	if err != nil {
		return time.Location{}, err
	}
	return userInfo.Location, nil
}

func (s *eventService) UpdateLocation(user string, newLocation time.Location) (time.Location, error) {
	updatedUser, err := s.repository.GetUser(user)
	if err != nil {
		panic("Implement me")
	}

	updatedUser.Location = newLocation

	return updatedUser.Location, nil
}

func (s *eventService) DeleteEvent(id int) error {
	foundEvent, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}
	s.repository.Delete(foundEvent)
	return nil
}

func (s *eventService) UpdateEvent(id int, newEvent structs.Event) (updated structs.Event, err error) {
	approved, err := s.checkData(newEvent)
	if !approved {
		return structs.Event{}, err
	}

	for _, event := range s.repository.GetAll() {
		if event.Id == id {
			event.Name = newEvent.Name
			event.Start = newEvent.Start
			event.End = newEvent.End
			event.Alert = newEvent.Alert
			event.Description = newEvent.Description
			return *event, nil
		}
	}
	return structs.Event{}, errors.New("no event with such id")
}

func (s *eventService) GetEventsOfTheDay(p structs.EventParams) ([]structs.Event, error) {
	var result []structs.Event

	if p.Day < 0 || p.Week < 0 || p.Month < 0 || p.Year < 0 {
		return nil, errors.New("bad date parameters")
	}
	for _, event := range s.repository.GetAll() {
		if s.repository.MatchParams(*event, p) {
			result = append(result, *event)
		}
	}

	if p.Sorting {
		return s.sortResults(result), nil
	}

	return result, nil
}

func (s *eventService) sortResults(events []structs.Event) []structs.Event {
	sort.Sort(ByStartTime(events))
	return events
}

func (s *eventService) checkData(newEvent structs.Event) (bool, error) {
	if newEvent.Name == "" {
		return false, &structs.MandatoryFieldError{FieldName: "name"}
	}
	if newEvent.Start == (time.Time{}) {
		return false, &structs.MandatoryFieldError{FieldName: "start"}
	}
	if newEvent.End == (time.Time{}) {
		return false, &structs.MandatoryFieldError{FieldName: "end"}
	}
	if newEvent.Start.Unix() > newEvent.End.Unix() {
		return false, errors.New("end of the event is ahead of the start")
	}
	return true, nil
}
