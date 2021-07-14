package db

import "github.com/dkucheru/Calendar/structs"

type Repository interface {
	GetAll() []*structs.Event
	GetByID(id int) (structs.Event, error)
	Add(structs.Event) error
	Delete(structs.Event)
	GetNextId() int
	GetLastUsedId() int //this function currently is used only for testing purpuses
	MatchParams(structs.Event, structs.EventParams) bool
}
