package service

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/dkucheru/Calendar/structs"
)

type eventService struct {
	repository []*structs.Event
}

func newEventsService(repository []*structs.Event) *eventService {
	s := eventService{
		repository: repository,
	}
	return &s
}

func (s *eventService) AddEvent(newEvent *structs.Event) (*structs.Event, error) {
	approved, err := s.checkData(newEvent)
	if !approved {
		return nil, err
	}

	newEvent.Id = s.incrementId()
	s.repository = append(s.repository, newEvent)

	return newEvent, nil
}

func (s *eventService) DeleteEvent(id int) error {
	for i, event := range s.repository {
		if event.Id == id {
			s.repository = append(s.repository[:i], s.repository[i+1:]...)
			return nil
		}
	}
	return errors.New("No event with such id was found")
}

func (s *eventService) UpdateEvent(id int, newEvent *structs.Event) (updated *structs.Event, err error) {
	approved, err := s.checkData(newEvent)
	if !approved {
		return nil, err
	}

	for _, event := range s.repository {
		if event.Id == id {
			event.Name = newEvent.Name
			event.Start = newEvent.Start
			event.End = newEvent.End
			event.Alert = newEvent.Alert
			event.Description = newEvent.Description
			return event, nil
		}
	}
	return nil, errors.New("No event with such id")
}

func (s *eventService) GetEventsOfTheDay(p *structs.EventParams) ([]structs.Event, error) {
	var result []structs.Event

	if p.Day < 0 || p.Week < 0 || p.Month < 0 || p.Year < 0 {
		return nil, errors.New("Bad date parameters")
	}
	for _, event := range s.repository {
		_, weekI := event.Start.ISOWeek()

		if (event.Start.Day() == p.Day || p.Day == 0) &&
			(p.Month == 0 || event.Start.Month() == time.Month(p.Month)) &&
			(p.Year == 0 || event.Start.Year() == p.Year) &&
			(p.Week == 0 || weekI == p.Week) &&
			(p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name)) &&
			(p.Start == (time.Time{}) || event.Start == p.Start) &&
			(p.End == (time.Time{}) || event.End == p.End) {
			result = append(result, *event)
		}
	}

	if p.Sorting {
		return *s.sortResults(&result), nil
	}

	return result, nil
}

func (s *eventService) sortResults(events *[]structs.Event) *[]structs.Event {
	sort.Sort(ByStartTime(*events))
	return events
}

func (s *eventService) incrementId() int {
	structs.GlobalId++
	return structs.GlobalId - 1
}

func (s *eventService) checkData(newEvent *structs.Event) (bool, error) {
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
		return false, errors.New("End of the event is ahead of the start")
	}
	return true, nil
}
