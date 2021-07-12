package db

import "github.com/dkucheru/Calendar/structs"

// FIXME: Use more generic names, not bound to specific implementation (e.g. Add, Delete etc.).
// The interface should be suitable for different implementation, for example for repository based on map.
// Add implementation of Repository based on map.
type Repository interface {
	Array() []*structs.Event
	AddToArray(*structs.Event)
	RemoveFromArray(int)
}
