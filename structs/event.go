package structs

import (
	"fmt"
	"time"
)

type Event struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Description string    `json:"description"`
	Alert       time.Time `json:"alert"`
}

type EventParams struct {
	Day     int
	Week    int
	Month   int
	Year    int
	Name    string
	Start   time.Time
	End     time.Time
	Sorting bool
}

var GlobalId int

type MandatoryFieldError struct {
	FieldName string
}

func (e *MandatoryFieldError) Error() string {
	return fmt.Sprintf("mandatory field *%v* is not filled. Please, add a %v to the event", e.FieldName, e.FieldName)
}
