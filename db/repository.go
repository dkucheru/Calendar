package db

import "github.com/dkucheru/Calendar/structs"

type Repository interface {
	Array() []*structs.Event
	AddToArray(*structs.Event)
	RemoveFromArray(int)
}
