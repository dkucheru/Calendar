package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/dkucheru/Calendar/db"
	"github.com/dkucheru/Calendar/structs"
)

func TestAddOnMap(t *testing.T) {
	var testRepo, _ = db.NewMapRepository()
	var testService = newEventsService(testRepo)

	testCases := map[string]struct {
		event        structs.Event
		result       structs.Event
		errorMessage string
	}{
		"Ok Event": {
			structs.Event{
				Name:        "Ok Test Event",
				Description: "an ok event for testing",
				Start:       time.Now(),
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			structs.Event{
				Id:          testService.repository.GetLastUsedId(),
				Name:        "Ok Test Event",
				Description: "an ok event for testing",
				Start:       time.Now(),
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			"",
		},
		"No description Test Event": {
			structs.Event{
				Name:  "No description Test Event",
				Start: time.Now(),
				End:   time.Now().Add(time.Hour),
				Alert: time.Now(),
			},
			structs.Event{
				Id:    testService.repository.GetLastUsedId(),
				Name:  "No description Test Event",
				Start: time.Now(),
				End:   time.Now().Add(time.Hour),
				Alert: time.Now(),
			},
			"",
		},
		"Only mandatory fields filled Test Event": {
			structs.Event{
				Name:  "Only mandatory fields filled Test Event",
				Start: time.Now(),
				End:   time.Now().Add(time.Hour),
			},
			structs.Event{
				Id:    testService.repository.GetLastUsedId(),
				Name:  "Only mandatory fields filled Test Event",
				Start: time.Now(),
				End:   time.Now().Add(time.Hour),
			},
			"",
		},
		"Name field not filled Event": {
			structs.Event{
				Description: "name field not filled event for testing",
				Start:       time.Now(),
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "name"}).Error(),
		},
		"No start time Event": {
			structs.Event{
				Name:        "No start time Event",
				Description: "No start time event for testing",
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "start"}).Error(),
		},
		"No End time Event": {
			structs.Event{
				Name:        "No End time Event",
				Description: "No end time event for testing",
				Start:       time.Now(),
				Alert:       time.Now(),
			},
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "end"}).Error(),
		},
		"Wrong duration of an Event": {
			structs.Event{
				Name:        "Wrong duration of an Event",
				Description: "wrong duration event for testing",
				Start:       time.Now().Add(time.Hour),
				End:         time.Now(),
				Alert:       time.Now(),
			},
			structs.Event{},
			"end of the event is ahead of the start",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			newEvent, err := testService.AddEvent(*time.Local, test.event)

			if !ErrorContains(err, test.errorMessage) {
				t.Errorf("wrong error : got %q, wanted %q", err, test.errorMessage)
			}

			if (newEvent == (structs.Event{}) && err == nil) || (newEvent == test.event && err != nil) {
				t.Errorf("event was added incorrectly")
			}

			newEventFromRepo, err2 := testService.GetById(newEvent.Id, *time.Local)
			if err2 != nil && err == nil {
				t.Errorf("event with id [%v] was not found", newEvent.Id)
			}

			if newEventFromRepo != (structs.Event{}) && err == nil && !structs.CompareTwoEvents(newEvent, newEventFromRepo) {
				t.Errorf("event returned by add function is not equal to the test event")
			}
		})

	}
}

func TestUpdateEventOnMap(t *testing.T) {
	var testRepo, _ = db.NewMapRepository()
	var testService = newEventsService(testRepo)

	testService.AddEvent(*time.Local, structs.Event{
		Name:        "Ok Test Event",
		Description: "an ok event for testing",
		Start:       time.Now(),
		End:         time.Now().Add(time.Hour),
		Alert:       time.Now(),
	})

	testCases := map[string]struct {
		event structs.Event
		id    int

		result       structs.Event
		errorMessage string
	}{
		"Ok Updated event": {
			structs.Event{
				Name:        "Updated Ok Event",
				Description: "an ok event for testing",
				Start:       time.Now(),
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			testService.repository.GetLastUsedId(),
			structs.Event{
				Id:          testService.repository.GetLastUsedId(),
				Name:        "Updated Ok Event",
				Description: "an ok event for testing",
				Start:       time.Now(),
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			"",
		},
		"Bad Id Update Event": {
			structs.Event{
				Name:        "Updated Event With Negative ID",
				Description: "an event for testing",
				Start:       time.Now(),
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			-1,
			structs.Event{},
			"event with id [-1] does not exist",
		},
		"No Name Field Update Event": {
			structs.Event{
				Description: "No name field",
				Start:       time.Now(),
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			testService.repository.GetLastUsedId(),
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "name"}).Error(),
		},
		"No Start Field Update Event": {
			structs.Event{
				Name:        "No Start date Event",
				Description: "no start",
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			testService.repository.GetLastUsedId(),
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "start"}).Error(),
		},
		"No End Field Update Event": {
			structs.Event{
				Name:        "No End date Event",
				Description: "no end",
				Start:       time.Now(),
				Alert:       time.Now(),
			},
			testService.repository.GetLastUsedId(),
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "end"}).Error(),
		},
		"Only Mandatory Fields Update Event": {
			structs.Event{
				Name:  "Only Mandatory Fields",
				Start: time.Now(),
				End:   time.Now().Add(time.Hour),
			},
			testService.repository.GetLastUsedId(),
			structs.Event{
				Id:    testService.repository.GetLastUsedId(),
				Name:  "Only Mandatory Fields",
				Start: time.Now(),
				End:   time.Now().Add(time.Hour),
			},
			"",
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			updatedEvent, err := testService.UpdateEvent(test.id, test.event, *time.Local)
			if (updatedEvent == structs.Event{} && err == nil) || (updatedEvent == test.event && err != nil) {
				t.Errorf("result returned by update function is incorrect")
			}

			wasUpdated, err2 := testService.GetById(test.id, *time.Local)
			if err2 != nil && err == nil {
				t.Errorf("event with id [%v] was not found", test.id)
			}
			//check if event was indeed updated
			if err == nil && err2 == nil && !structs.CompareTwoEvents(updatedEvent, wasUpdated) {
				t.Errorf("event with id [%v] was not updated correctly", test.id)
			}

			test.event.Id = testService.repository.GetLastUsedId()
			updatedEvent.Start = updatedEvent.Start.In(time.Local)
			updatedEvent.End = updatedEvent.End.In(time.Local)
			if updatedEvent.Alert != (time.Time{}) {
				updatedEvent.Alert = updatedEvent.Alert.In(time.Local)
			}
			if err == nil && !structs.CompareTwoEvents(updatedEvent, test.event) {
				t.Errorf("result returned by update function is incorrect")
			}

			if !ErrorContains(err, test.errorMessage) {
				t.Errorf("wrong error : got %q, wanted %q", err, test.errorMessage)
			}
		})

	}
}

func TestGetEventOnMap(t *testing.T) {
	var testRepo, _ = db.NewMapRepository()
	var testService = newEventsService(testRepo)

	testService.AddEvent(*time.Local, structs.Event{
		Name:        "Ok Test Event",
		Description: "an ok event for testing",
		Start:       time.Now(),
		End:         time.Now().Add(time.Hour),
		Alert:       time.Now(),
	})

	testCases := map[string]struct {
		params       structs.EventParams
		result       []structs.Event
		errorMessage string
	}{
		"Normal Parameters Get Event Test": {
			structs.EventParams{
				Day:     time.Now().Day(),
				Month:   int(time.Now().Month()),
				Year:    time.Now().Year(),
				Name:    "Ok Test Event",
				Start:   time.Now(),
				End:     time.Now().Add(time.Hour),
				Sorting: true,
			},
			[]structs.Event{
				{
					Name:        "Ok Test Event",
					Description: "an ok event for testing",
					Start:       time.Now(),
					End:         time.Now().Add(time.Hour),
					Alert:       time.Now(),
				},
			},
			"",
		},
		"One Parameter Get Event Test ": {
			structs.EventParams{
				Day: time.Now().Day(),
			},
			[]structs.Event{
				{
					Name:        "Ok Test Event",
					Description: "an ok event for testing",
					Start:       time.Now(),
					End:         time.Now().Add(time.Hour),
					Alert:       time.Now(),
				},
			},
			"",
		},
		"Bad Day Parameter Get Event Test ": {
			structs.EventParams{
				Day: -1,
			},
			[]structs.Event{},
			"bad date parameters",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			events, err := testService.GetEventsOfTheDay(test.params)
			for _, v := range events {

				if (v == structs.Event{} && err == nil) {
					t.Errorf("result returned by get function is incorrect")
				}

				event, err2 := testService.GetById(v.Id, *time.Local)
				if v != (structs.Event{}) && err2 != nil {
					t.Errorf("event with id [%v] does not exist", v.Id)
				}

				resultMatchesInputParams := false
				for _, match := range testService.repository.Get(test.params) {
					if structs.CompareTwoEvents(event, *match) {
						resultMatchesInputParams = true
					}
				}

				if !resultMatchesInputParams {
					t.Errorf("result returned by get function does not correspond to input parameters")
				}
			}

			if !ErrorContains(err, test.errorMessage) {
				t.Errorf("wrong error : got %q, wanted %q", err, test.errorMessage)
			}
		})

	}
}

func TestDeleteEventOnMap(t *testing.T) {
	var testRepo, _ = db.NewMapRepository()
	var testService = newEventsService(testRepo)

	testService.AddEvent(*time.Local, structs.Event{
		Name:        "Ok Test Event",
		Description: "an ok event for testing",
		Start:       time.Now(),
		End:         time.Now().Add(time.Hour),
		Alert:       time.Now(),
	})

	testCases := map[string]struct {
		id           int
		errorMessage string
	}{
		"Ok Delete Event Test":     {1, ""},
		"Bad Id Delete Event Test": {-1, "event with id [-1] does not exist"},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			err := testService.DeleteEvent(test.id, "testUsername")
			if !ErrorContains(err, test.errorMessage) {
				t.Errorf("got %q, wanted %q", err, test.errorMessage)
			}

			_, err2 := testService.GetById(test.id, *time.Local)

			message := "event with id [" + fmt.Sprint(test.id) + "] does not exist"
			if err2.Error() != message && err == nil {
				t.Errorf("event with id [" + fmt.Sprint(test.id) + "] was not deleted")
			}
		})

	}
}
