package db

import "github.com/dkucheru/Calendar/structs"

// FIXME: add one more repository for users

type Repository interface {
	Add(structs.Event) structs.Event
	Get(structs.EventParams) []*structs.Event
	GetByID(id int) (*structs.Event, error)
	Update(id int, newEvent structs.Event) (updated structs.Event, err error)
	Delete(structs.Event)
	GetLastUsedId() int //this function currently is used only for testing purpuses

	AddUser(structs.CreateUser) error
	GetUser(string) (*structs.HashedInfo, error)
	CheckCredentials(username string, pass string) error
}
