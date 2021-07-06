package structs

import "time"

type Event struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Description string    `json:"description"`
	Alert       time.Time `json:"alert"`
}

var GlobalId int
