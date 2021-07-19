package db

import "github.com/dkucheru/Calendar/structs"

type Repository interface {
	// GetAll() []*structs.Event
	Get(structs.EventParams) []*structs.Event
	GetByID(id int) (*structs.Event, error)
	Add(structs.Event) structs.Event
	Update(id int, newEvent structs.Event) (updated structs.Event, err error)
	Delete(structs.Event)
	GetLastUsedId() int //this function currently is used only for testing purpuses
	// MatchParams(structs.Event, structs.EventParams) bool
}
