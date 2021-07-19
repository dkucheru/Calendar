package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dkucheru/Calendar/structs"
)

type ArrayRepository struct {
	ArrayRepo []*structs.Event
	ArrayId   int
}

func NewArrayRepository() (*ArrayRepository, error) {
	var events []*structs.Event
	repo := &ArrayRepository{
		ArrayRepo: events,
		ArrayId:   1,
	}
	return repo, nil
}

// func (a *ArrayRepository) GetAll() []*structs.Event {
// 	return a.ArrayRepo
// }

func (a *ArrayRepository) Get(p structs.EventParams) []*structs.Event {
	var matchedEvents []*structs.Event

	for _, event := range a.ArrayRepo {
		_, weekI := event.Start.ISOWeek()

		if !(event.Start.Day() == p.Day || p.Day == 0) {
			matchedEvents = append(matchedEvents, event)
		} else if !(p.Month == 0 || event.Start.Month() == time.Month(p.Month)) {
			matchedEvents = append(matchedEvents, event)
		} else if !(p.Year == 0 || event.Start.Year() == p.Year) {
			matchedEvents = append(matchedEvents, event)
		} else if !(p.Week == 0 || weekI == p.Week) {
			matchedEvents = append(matchedEvents, event)
		} else if !(p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name)) {
			matchedEvents = append(matchedEvents, event)
		} else if !(p.Start == (time.Time{}) || event.Start == p.Start) {
			matchedEvents = append(matchedEvents, event)
		} else if !(p.End == (time.Time{}) || event.End == p.End) {
			matchedEvents = append(matchedEvents, event)
		}

	}
	return matchedEvents
}

func (a *ArrayRepository) Update(id int, newEvent structs.Event) (updated structs.Event, err error) {
	foundEvent, err := a.GetByID(id)
	if err != nil {
		return structs.Event{}, err
	}
	foundEvent.Name = newEvent.Name
	foundEvent.Start = newEvent.Start
	foundEvent.End = newEvent.End
	foundEvent.Alert = newEvent.Alert
	foundEvent.Description = newEvent.Description
	return *foundEvent, nil
}

func (a *ArrayRepository) GetByID(id int) (*structs.Event, error) {
	for _, event := range a.ArrayRepo {
		if event.Id == id {
			return event, nil
		}
	}
	message := "event with id [" + fmt.Sprint(id) + "] does not exist"
	return nil, errors.New(message)
}

func (a *ArrayRepository) Add(e structs.Event) structs.Event {
	e.Id = a.ArrayId
	a.ArrayId++
	a.ArrayRepo = append(a.ArrayRepo, &e)
	return e
}

func (a *ArrayRepository) Delete(e structs.Event) {
	for i, event := range a.ArrayRepo {
		if event.Id == e.Id {
			a.ArrayRepo = append(a.ArrayRepo[:i], a.ArrayRepo[i+1:]...)
			return
		}
	}
}

func (a *ArrayRepository) GetLastUsedId() int {
	return a.ArrayId - 1
}

// func (a *ArrayRepository) MatchParams(event structs.Event, p structs.EventParams) bool {
// 	_, weekI := event.Start.ISOWeek()

// 	if !(event.Start.Day() == p.Day || p.Day == 0) {
// 		return false
// 	}
// 	if !(p.Month == 0 || event.Start.Month() == time.Month(p.Month)) {
// 		return false
// 	}
// 	if !(p.Year == 0 || event.Start.Year() == p.Year) {
// 		return false
// 	}
// 	if !(p.Week == 0 || weekI == p.Week) {
// 		return false
// 	}
// 	if !(p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name)) {
// 		return false
// 	}
// 	if !(p.Start == (time.Time{}) || event.Start == p.Start) {
// 		return false
// 	}
// 	if !(p.End == (time.Time{}) || event.End == p.End) {
// 		return false
// 	}
// 	return true
// }

// Implementation of Repository based on map.
type MapRepository struct {
	MapRepo map[int]structs.Event
	MapId   int
}

func NewMapRepository() (*MapRepository, error) {
	events := make(map[int]structs.Event)
	repo := &MapRepository{
		MapRepo: events,
		MapId:   1,
	}
	return repo, nil
}

// func (m *MapRepository) GetAll() []*structs.Event {
// 	events := make([]*structs.Event, 0, len(m.MapRepo))
// 	for _, v := range m.MapRepo {
// 		events = append(events, &v)
// 	}
// 	return events
// }

func (m *MapRepository) Get(p structs.EventParams) []*structs.Event {
	var matchedEvents []*structs.Event

	for _, event := range m.MapRepo {
		_, weekI := event.Start.ISOWeek()

		if !(event.Start.Day() == p.Day || p.Day == 0) {
			matchedEvents = append(matchedEvents, &event)
		} else if !(p.Month == 0 || event.Start.Month() == time.Month(p.Month)) {
			matchedEvents = append(matchedEvents, &event)
		} else if !(p.Year == 0 || event.Start.Year() == p.Year) {
			matchedEvents = append(matchedEvents, &event)
		} else if !(p.Week == 0 || weekI == p.Week) {
			matchedEvents = append(matchedEvents, &event)
		} else if !(p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name)) {
			matchedEvents = append(matchedEvents, &event)
		} else if !(p.Start == (time.Time{}) || event.Start == p.Start) {
			matchedEvents = append(matchedEvents, &event)
		} else if !(p.End == (time.Time{}) || event.End == p.End) {
			matchedEvents = append(matchedEvents, &event)
		}

	}
	return matchedEvents
}

func (m *MapRepository) Update(id int, newEvent structs.Event) (updated structs.Event, err error) {
	foundEvent, ok := m.MapRepo[id]
	if !ok {
		message := "event with id [" + fmt.Sprint(id) + "] does not exist"
		return structs.Event{}, errors.New(message)
	}

	foundEvent.Name = newEvent.Name
	foundEvent.Start = newEvent.Start
	foundEvent.End = newEvent.End
	foundEvent.Alert = newEvent.Alert
	foundEvent.Description = newEvent.Description

	m.MapRepo[id] = foundEvent

	return foundEvent, nil
}

func (m *MapRepository) GetByID(id int) (*structs.Event, error) {
	event, ok := m.MapRepo[id]
	if !ok {
		message := "event with id [" + fmt.Sprint(id) + "] does not exist"
		return nil, errors.New(message)
	}
	return &event, nil
}

func (m *MapRepository) Add(e structs.Event) structs.Event {
	e.Id = m.MapId
	m.MapId++

	m.MapRepo[e.Id] = e
	return e
}

func (m *MapRepository) Delete(e structs.Event) {
	delete(m.MapRepo, e.Id)
}

func (m *MapRepository) GetLastUsedId() int {
	return m.MapId - 1
}

// func (m *MapRepository) MatchParams(event structs.Event, p structs.EventParams) bool {
// 	_, weekI := event.Start.ISOWeek()

// 	if !(event.Start.Day() == p.Day || p.Day == 0) {
// 		return false
// 	}
// 	if !(p.Month == 0 || event.Start.Month() == time.Month(p.Month)) {
// 		return false
// 	}
// 	if !(p.Year == 0 || event.Start.Year() == p.Year) {
// 		return false
// 	}
// 	if !(p.Week == 0 || weekI == p.Week) {
// 		return false
// 	}
// 	if !(p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name)) {
// 		return false
// 	}
// 	if !(p.Start == (time.Time{}) || event.Start == p.Start) {
// 		return false
// 	}
// 	if !(p.End == (time.Time{}) || event.End == p.End) {
// 		return false
// 	}
// 	return true
// }
