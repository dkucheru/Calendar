package structs

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/go-playground/validator/v10"
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

func CompareTwoEvents(f Event, s Event) bool {
	if f.Id != s.Id {
		return false
	}
	if f.Name != s.Name {
		return false
	}
	if f.Description != s.Description {
		return false
	}
	if math.Abs(float64(f.Start.Unix()-s.Start.Unix())) >= time.Minute.Seconds() {
		// fmt.Printf("diff between two starts is over a minute : %v\n", math.Abs(float64(f.Start.Unix()-s.Start.Unix())))
		// fmt.Println("first start time : " + f.Start.String())
		// fmt.Println("second start time : " + s.Start.String())
		return false
	}
	if math.Abs(float64(f.End.Unix()-s.End.Unix())) >= time.Minute.Seconds() {
		return false
	}
	if math.Abs(float64(f.Alert.Unix()-s.Alert.Unix())) >= time.Minute.Seconds() {
		return false
	}
	return true
}

type EventCreation struct {
	Name        string    `json:"name" validate:"required"`
	Start       time.Time `json:"start" validate:"required"`
	End         time.Time `json:"end" validate:"required"`
	Description string    `json:"description"`
	Alert       time.Time `json:"alert"`
}

func SuitsParams(p EventParams, e Event) bool {
	if p.Day != 0 {
		if e.Start.Day() != p.Day {
			return false
		}
	}
	if p.Week != 0 {
		_, week := e.Start.ISOWeek()
		if week != p.Week {
			return false
		}
	}
	if p.Month != 0 {
		if int(e.Start.Month()) != p.Month {
			return false
		}
	}
	if p.Year != 0 {
		if e.Start.Year() != p.Year {
			return false
		}
	}
	if p.Name != "" {
		if e.Name != p.Name {
			return false
		}
	}
	if p.Start != (time.Time{}) {
		if e.Start != p.Start {
			return false
		}
	}
	if p.End != (time.Time{}) {
		if e.End != p.End {
			return false
		}
	}
	return true
}

func CreateEvent(loc time.Location, newEvent EventCreation) (Event, error) {
	validate := validator.New()
	err := validate.Struct(newEvent)
	if err != nil {
		return Event{}, errors.New("validator : invalid data format")
	}
	if err != nil {
		return Event{}, err
	}

	st := newEvent.Start
	end := newEvent.End
	a := newEvent.Alert

	t := time.Date(st.Year(), st.Month(), st.Day(), st.Hour(), st.Minute(), st.Second(), st.Nanosecond(), &loc)
	newEvent.Start = t.In(time.UTC)
	t = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), &loc)
	newEvent.End = t.In(time.UTC)
	if newEvent.Alert != (time.Time{}) {
		t = time.Date(a.Year(), a.Month(), a.Day(), a.Hour(), a.Minute(), a.Second(), a.Nanosecond(), &loc)
		newEvent.Alert = t.In(time.UTC)
	}

	return Event{
		Name:        newEvent.Name,
		Start:       newEvent.Start,
		End:         newEvent.End,
		Alert:       newEvent.Alert,
		Description: newEvent.Description,
	}, nil
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

type URLParams struct {
	Day     string
	Week    string
	Month   string
	Year    string
	Name    string
	Start   string
	End     string
	Sorting string
}

var GlobalId int

type MandatoryFieldError struct {
	FieldName string
}

func (e *MandatoryFieldError) Error() string {
	return fmt.Sprintf("mandatory field *%v* is not filled. Please, add a %v to the event", e.FieldName, e.FieldName)
}
