package structs

import (
	"encoding/json"
	"fmt"
	"time"
)

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type HashedInfo struct {
	Username   string
	Location   time.Location
	HashedPass string
}

type Event struct {
	Id          int       `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Start       time.Time `json:"start" validate:"required"`
	End         time.Time `json:"end" validate:"required"`
	Description string    `json:"description"`
	Alert       time.Time `json:"alert"`
}

// FIXME: Event Creation structure

// FIXME:
func (e *Event) ParseJSON(data []byte, loc time.Location) error {
	var err error
	if err = json.Unmarshal(data, &e); err != nil {
		panic(err)
	}

	s := e.Start
	end := e.End
	a := e.Alert

	t := time.Date(s.Year(), s.Month(), s.Day(), s.Hour(), s.Minute(), s.Second(), s.Nanosecond(), &loc)
	e.Start = t.UTC()
	t = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), &loc)
	e.End = t.UTC()
	if e.Alert != (time.Time{}) {
		t = time.Date(a.Year(), a.Month(), a.Day(), a.Hour(), a.Minute(), a.Second(), a.Nanosecond(), &loc)
		e.Alert = t.UTC()
	}

	return nil
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
