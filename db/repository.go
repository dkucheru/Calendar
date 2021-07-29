package db

import (
	"time"

	"github.com/dkucheru/Calendar/structs"
)

type EventsRepository interface {
	Add(structs.Event) (structs.Event, error)
	Get(structs.EventParams) ([]structs.Event, error)
	GetByID(id int) (structs.Event, error)
	Update(id int, newEvent structs.Event) (updated structs.Event, err error)
	Delete(structs.Event) error
	GetLastUsedId() int //this function currently is used only for testing purpuses
	ClearRepoData() error
}

type UserRepository interface {
	AddUser(structs.CreateUser) (structs.HashedInfo, error)
	GetUser(string) (structs.HashedInfo, error)
	UpdateLocation(user string, loc time.Location) (structs.HashedInfo, error)
	ClearRepoData() error
}
