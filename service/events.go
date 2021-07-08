package service

import (
	"errors"
	"fmt"
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

//return (int,error) mb event
func (s *eventService) AddEvent(newEvent *structs.Event) error {
	//check if event has mandatory field "Name" filled
	if newEvent.Name == "" {
		return &structs.MandatoryFieldError{FieldName: "name"}
	}
	if newEvent.Start == (time.Time{}) {
		return &structs.MandatoryFieldError{FieldName: "start"}
	}
	if newEvent.End == (time.Time{}) {
		return &structs.MandatoryFieldError{FieldName: "end"}
	}

	if newEvent.Start.Unix() > newEvent.End.Unix() {
		return errors.New("End of the event is ahead of the start")
	}

	//getNext function
	newEvent.Id = structs.GlobalId
	structs.GlobalId++

	// Add new Event to our calendar
	s.repository = append(s.repository, newEvent)

	//debuging
	fmt.Println(newEvent)

	return nil
}

func (s *eventService) DeleteEvent(id int) error {
	for i, event := range s.repository {
		if event.Id == id {
			s.repository = append(s.repository[:i], s.repository[i+1:]...)
			return nil
		}
	}
	return errors.New("No event with such name and start date was found")
}

func (s *eventService) UpdateEvent(id int, newEvent *structs.Event) error {
	//check if event has mandatory field "Name" filled
	if newEvent.Name == "" {
		return &structs.MandatoryFieldError{FieldName: "name"}
	}
	if newEvent.Start == (time.Time{}) {
		return &structs.MandatoryFieldError{FieldName: "start"}
	}
	if newEvent.End == (time.Time{}) {
		return &structs.MandatoryFieldError{FieldName: "end"}
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

func (s *eventService) GetAll() (error, []structs.Event) { //remove
	var result []structs.Event
	for _, event := range s.repository {
		result = append(result, *event)
	}
	return nil, result
}

//add separate method for sorting

func (s *eventService) GetEventsOfTheDay(p *structs.EventParams) ([]structs.Event, error) {
	var result []structs.Event

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

	return result, nil

}
