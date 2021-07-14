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

func (a *ArrayRepository) GetAll() []*structs.Event {
	return a.ArrayRepo
}

func (a *ArrayRepository) GetByID(id int) (structs.Event, error) {
	for _, event := range a.ArrayRepo {
		if event.Id == id {
			return *event, nil
		}
	}
	message := "event with id [" + fmt.Sprint(id) + "] does not exist"
	return structs.Event{}, errors.New(message)
}

func (a *ArrayRepository) Add(e structs.Event) error {
	for _, event := range a.ArrayRepo {
		if event.Id == e.Id {
			message := "event with id [" + fmt.Sprint(e.Id) + "] already exists"
			return errors.New(message)
		}
	}
	a.ArrayRepo = append(a.ArrayRepo, &e)
	return nil
}

func (a *ArrayRepository) Delete(e structs.Event) {
	for i, event := range a.ArrayRepo {
		if event.Id == e.Id {
			a.ArrayRepo = append(a.ArrayRepo[:i], a.ArrayRepo[i+1:]...)
			return
		}
	}
}

func (a *ArrayRepository) GetNextId() int {
	a.ArrayId++
	return a.ArrayId - 1
}

func (a *ArrayRepository) GetLastUsedId() int {
	return a.ArrayId - 1
}

func (a *ArrayRepository) MatchParams(event structs.Event, p structs.EventParams) bool {
	_, weekI := event.Start.ISOWeek()

	if !(event.Start.Day() == p.Day || p.Day == 0) {
		return false
	}
	if !(p.Month == 0 || event.Start.Month() == time.Month(p.Month)) {
		return false
	}
	if !(p.Year == 0 || event.Start.Year() == p.Year) {
		return false
	}
	if !(p.Week == 0 || weekI == p.Week) {
		return false
	}
	if !(p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name)) {
		return false
	}
	if !(p.Start == (time.Time{}) || event.Start == p.Start) {
		return false
	}
	if !(p.End == (time.Time{}) || event.End == p.End) {
		return false
	}
	return true
}

// Implementation of Repository based on map.
type MapRepository struct {
	MapRepo map[int]*structs.Event
	MapId   int
}

func NewMapRepository() (*MapRepository, error) {
	events := make(map[int]*structs.Event)
	repo := &MapRepository{
		MapRepo: events,
		MapId:   1,
	}
	return repo, nil
}

func (m *MapRepository) GetAll() []*structs.Event {
	events := make([]*structs.Event, 0, len(m.MapRepo))
	for _, v := range m.MapRepo {
		events = append(events, v)
	}
	return events
}

func (m *MapRepository) GetByID(id int) (structs.Event, error) {
	_, ok := m.MapRepo[id]
	if !ok {
		message := "event with id [" + fmt.Sprint(id) + "] does not exist"
		return structs.Event{}, errors.New(message)
	}
	return *m.MapRepo[id], nil
}

func (m *MapRepository) Add(e structs.Event) error {
	_, ok := m.MapRepo[e.Id]
	if ok {
		message := "event with id [" + fmt.Sprint(e.Id) + "] already exists"
		return errors.New(message)
	}
	m.MapRepo[e.Id] = &e
	return nil
}

func (m *MapRepository) Delete(e structs.Event) {
	delete(m.MapRepo, e.Id)
}

func (m *MapRepository) GetNextId() int {
	m.MapId++
	return m.MapId - 1
}

func (m *MapRepository) GetLastUsedId() int {
	return m.MapId - 1
}

func (m *MapRepository) MatchParams(event structs.Event, p structs.EventParams) bool {
	_, weekI := event.Start.ISOWeek()

	if !(event.Start.Day() == p.Day || p.Day == 0) {
		return false
	}
	if !(p.Month == 0 || event.Start.Month() == time.Month(p.Month)) {
		return false
	}
	if !(p.Year == 0 || event.Start.Year() == p.Year) {
		return false
	}
	if !(p.Week == 0 || weekI == p.Week) {
		return false
	}
	if !(p.Name == "" || strings.ToLower(event.Name) == strings.ToLower(p.Name)) {
		return false
	}
	if !(p.Start == (time.Time{}) || event.Start == p.Start) {
		return false
	}
	if !(p.End == (time.Time{}) || event.End == p.End) {
		return false
	}
	return true
}
