package db

import "github.com/dkucheru/Calendar/structs"

// FIXME: add one more repository for users

type Repository interface {
	GetAll() []*structs.Event
	GetByID(id int) (structs.Event, error)
	AddEvent(structs.Event) error
	AddUser(structs.CreateUser) error
	GetUser(string) (*structs.HashedInfo, error)
	CheckCredentials(username string, pass string) error
	Delete(structs.Event)
	GetNextId() int
	GetLastUsedId() int //this function currently is used only for testing purpuses
	MatchParams(structs.Event, structs.EventParams) bool
}
