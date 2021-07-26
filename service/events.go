package service

import (
	"errors"
	"sort"
	"time"

	"github.com/dkucheru/Calendar/db"
	"github.com/dkucheru/Calendar/structs"
)

type eventService struct {
	repository db.EventsRepository
}

func newEventsService(repository db.EventsRepository) *eventService {
	s := eventService{
		repository: repository,
	}
	return &s
}

func (s *eventService) AddEvent(loc time.Location, newEvent structs.Event) (structs.Event, error) {
	approved, err := s.checkData(newEvent)
	if !approved {
		return structs.Event{}, err
	}
	// log.Println("UTC ??? " + newEvent.Start.String())
	returnedEvent, err := s.repository.Add(newEvent)
	if err != nil {
		return structs.Event{}, err
	}

	returnedEvent.Start = newEvent.Start.In(&loc)
	returnedEvent.End = newEvent.End.In(&loc)
	if returnedEvent.Alert != (time.Time{}) {
		returnedEvent.Alert = newEvent.Alert.In(&loc)
	}

	return returnedEvent, nil
}

func (s *eventService) DeleteEvent(id int, user string) error {
	foundEvent, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}
	s.repository.Delete(foundEvent)
	return nil
}

func (s *eventService) GetById(id int, loc time.Location) (structs.Event, error) {

	returnedEvent, err := s.repository.GetByID(id)
	if err != nil {
		return structs.Event{}, err
	}
	// fmt.Println("returned start time before appliyng location : " + returnedEvent.Start.String())
	returnedEvent.Start = returnedEvent.Start.In(&loc)
	returnedEvent.End = returnedEvent.End.In(&loc)
	if returnedEvent.Alert != (time.Time{}) {
		returnedEvent.Alert = returnedEvent.Alert.In(&loc)
	}
	return returnedEvent, nil
}

func (s *eventService) UpdateEvent(id int, newEvent structs.Event, loc time.Location) (updated structs.Event, err error) {
	approved, err := s.checkData(newEvent)
	if !approved {
		return structs.Event{}, err
	}
	returnedEvent, err := s.repository.Update(id, newEvent)

	returnedEvent.Start = newEvent.Start.In(&loc)
	returnedEvent.End = newEvent.End.In(&loc)
	if returnedEvent.Alert != (time.Time{}) {
		returnedEvent.Alert = returnedEvent.Alert.In(time.Local)
	}
	return returnedEvent, err
}

func (s *eventService) GetEventsOfTheDay(p structs.EventParams, loc time.Location) ([]structs.Event, error) {
	var result []structs.Event
	if p.Day < 0 || p.Week < 0 || p.Month < 0 || p.Year < 0 {
		return nil, errors.New("bad date parameters")
	}
	receivedEvents, err := s.repository.Get(p)
	if err != nil {
		return []structs.Event{}, err
	}
	for _, event := range receivedEvents {
		event.Start = event.Start.In(&loc)
		event.End = event.End.In(&loc)
		if event.Alert != (time.Time{}) {
			event.Alert = event.Alert.In(&loc)
		}
		result = append(result, event)
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
