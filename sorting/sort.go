package sorting

//not needed package

import "github.com/dkucheru/Calendar/structs"

type ByStartTime []structs.Event

func (a ByStartTime) Len() int           { return len(a) }
func (a ByStartTime) Less(i, j int) bool { return a[i].Start.Unix() < a[j].Start.Unix() }
func (a ByStartTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
