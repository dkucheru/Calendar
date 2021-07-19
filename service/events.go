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
	returnedEvent := s.repository.Add(newEvent)
	return returnedEvent, nil
}

func (s *eventService) DeleteEvent(id int) error {
	foundEvent, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}
	s.repository.Delete(*foundEvent)
	return nil
}

func (s *eventService) UpdateEvent(id int, newEvent structs.Event) (updated structs.Event, err error) {
	approved, err := s.checkData(newEvent)
	if !approved {
		return structs.Event{}, err
	}
	return s.repository.Update(id, newEvent)
}

func (s *eventService) GetEventsOfTheDay(p structs.EventParams) ([]structs.Event, error) {
	var result []structs.Event

	if p.Day < 0 || p.Week < 0 || p.Month < 0 || p.Year < 0 {
		return nil, errors.New("bad date parameters")
	}
	for _, event := range s.repository.Get(p) {
		result = append(result, *event)
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
