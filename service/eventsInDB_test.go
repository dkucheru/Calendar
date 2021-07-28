package service

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dkucheru/Calendar/db"
	"github.com/dkucheru/Calendar/structs"
)

func TestAddToDB(t *testing.T) {
	repo, err := db.Initialize(os.Getenv("DSN"))
	if err != nil {
		t.Errorf(err.Error())
	}
	var testRepo, _ = db.NewDatabaseRepository(repo)
	var testService = newEventsService(testRepo)
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	testCases := map[string]struct {
		event        structs.Event
		result       structs.Event
		errorMessage string
	}{
		"Ok Event": {
			structs.Event{
				Name:        "Ok Test Event",
				Description: "an ok event for testing",
				Start:       time.Now().In(time.UTC),
				End:         time.Now().Add(time.Hour).In(time.UTC),
				Alert:       time.Now().In(time.UTC),
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
				Start: time.Now().In(time.UTC),
				End:   time.Now().Add(time.Hour).In(time.UTC),
				Alert: time.Now().In(time.UTC),
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
				Start: time.Now().In(time.UTC),
				End:   time.Now().Add(time.Hour).In(time.UTC),
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
				Start:       time.Now().In(time.UTC),
				End:         time.Now().Add(time.Hour).In(time.UTC),
				Alert:       time.Now().In(time.UTC),
			},
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "name"}).Error(),
		},
		"No start time Event": {
			structs.Event{
				Name:        "No start time Event",
				Description: "No start time event for testing",
				End:         time.Now().Add(time.Hour).In(time.UTC),
				Alert:       time.Now().In(time.UTC),
			},
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "start"}).Error(),
		},
		"No End time Event": {
			structs.Event{
				Name:        "No End time Event",
				Description: "No end time event for testing",
				Start:       time.Now().In(time.UTC),
				Alert:       time.Now().In(time.UTC),
			},
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "end"}).Error(),
		},
		"Wrong duration of an Event": {
			structs.Event{
				Name:        "Wrong duration of an Event",
				Description: "wrong duration event for testing",
				Start:       time.Now().Add(time.Hour).In(time.UTC),
				End:         time.Now().In(time.UTC),
				Alert:       time.Now().In(time.UTC),
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
			test.event.Id = testRepo.GetLastUsedId()
			if (newEvent == (structs.Event{}) && err == nil) || (structs.CompareTwoEvents(newEvent, test.event) && err != nil) {
				t.Errorf("event was added incorrectly:\n wanted %v\n got %v\n", test.event, newEvent)
			}

			newEventFromRepo, err2 := testService.GetById(newEvent.Id, *time.Local)
			if err2 != nil && err == nil {
				t.Errorf("freshy added event was not found in the db")
			}

			if err2 == nil && err == nil && newEventFromRepo != (structs.Event{}) && !structs.CompareTwoEvents(newEvent, newEventFromRepo) {
				t.Errorf("event was added incorrectly:\n wanted %v\n got %v\n", newEvent, newEventFromRepo)
			}
		})

	}

}

func TestUpdateEventInDB(t *testing.T) {
	repo, err := db.Initialize(os.Getenv("DSN"))
	if err != nil {
		t.Errorf(err.Error())
	}
	var testRepo, _ = db.NewDatabaseRepository(repo)
	var testService = newEventsService(testRepo)
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	testService.AddEvent(*time.Local, structs.Event{
		Name:        "Ok Test Event",
		Description: "an ok event for testing",
		Start:       time.Now().In(time.UTC),
		End:         time.Now().Add(time.Hour).In(time.UTC),
		Alert:       time.Now().In(time.UTC),
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
				Start:       time.Now().In(time.UTC),
				End:         time.Now().Add(time.Hour).In(time.UTC),
				Alert:       time.Now().In(time.UTC),
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
				Start:       time.Now().In(time.UTC),
				End:         time.Now().Add(time.Hour).In(time.UTC),
				Alert:       time.Now().In(time.UTC),
			},
			-1,
			structs.Event{},
			"event with id [-1] does not exist",
		},
		"No Name Field Update Event": {
			structs.Event{
				Description: "No name field",
				Start:       time.Now().In(time.UTC),
				End:         time.Now().Add(time.Hour).In(time.UTC),
				Alert:       time.Now().In(time.UTC),
			},
			testService.repository.GetLastUsedId(),
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "name"}).Error(),
		},
		"No Start Field Update Event": {
			structs.Event{
				Name:        "No Start date Event",
				Description: "no start",
				End:         time.Now().Add(time.Hour).In(time.UTC),
				Alert:       time.Now().In(time.UTC),
			},
			testService.repository.GetLastUsedId(),
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "start"}).Error(),
		},
		"No End Field Update Event": {
			structs.Event{
				Name:        "No End date Event",
				Description: "no end",
				Start:       time.Now().In(time.UTC),
				Alert:       time.Now().In(time.UTC),
			},
			testService.repository.GetLastUsedId(),
			structs.Event{},
			(&structs.MandatoryFieldError{FieldName: "end"}).Error(),
		},
		"Only Mandatory Fields Update Event": {
			structs.Event{
				Name:  "Only Mandatory Fields",
				Start: time.Now().In(time.UTC),
				End:   time.Now().Add(time.Hour).In(time.UTC),
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
			if updatedEvent == (structs.Event{}) && err == nil {
				t.Errorf("no errors in update function occured, but returned result is an empty struct")
			}

			test.event.Id = testService.repository.GetLastUsedId()

			wasUpdated, err2 := testService.GetById(test.id, *time.Local)
			// check if event was indeed updated
			if err2 == nil && err == nil && !structs.CompareTwoEvents(updatedEvent, wasUpdated) {
				t.Errorf("event with id [%v] was not updated correctly", test.id)
			}

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

func TestGetEventFromDB(t *testing.T) {
	repo, err := db.Initialize(os.Getenv("DSN"))
	if err != nil {
		t.Errorf(err.Error())
	}
	var testRepo, _ = db.NewDatabaseRepository(repo)
	var testService = newEventsService(testRepo)
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	testService.AddEvent(*time.Local, structs.Event{
		Name:        "Ok Test Event",
		Description: "an ok event for testing",
		Start:       time.Now().In(time.UTC),
		End:         time.Now().Add(time.Hour).In(time.UTC),
		Alert:       time.Now().In(time.UTC),
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
				Start:   time.Now().In(time.UTC),
				End:     time.Now().Add(time.Hour).In(time.UTC),
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
			events, err := testService.GetEventsOfTheDay(test.params, *time.Local)
			for _, v := range events {

				if (v == structs.Event{} && err == nil) {
					t.Errorf("result returned by get function is incorrect")
				}

				event, err2 := testService.GetById(v.Id, *time.Local)
				if v != (structs.Event{}) && err2 != nil {
					t.Errorf("event with id [%v] does not exist", v.Id)
				}

				resultMatchesInputParams := false
				events, err := testService.repository.Get(test.params)
				if err != nil {
					t.Errorf(err.Error())
				}
				for _, match := range events {
					if structs.CompareTwoEvents(event, match) {
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

func TestDeleteEventFromDB(t *testing.T) {
	repo, err := db.Initialize(os.Getenv("DSN"))
	if err != nil {
		t.Errorf(err.Error())
	}
	var testRepo, _ = db.NewDatabaseRepository(repo)
	var testService = newEventsService(testRepo)
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = testRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	testService.AddEvent(*time.Local, structs.Event{
		Name:        "Ok Test Event",
		Description: "an ok event for testing",
		Start:       time.Now().In(time.UTC),
		End:         time.Now().Add(time.Hour).In(time.UTC),
		Alert:       time.Now().In(time.UTC),
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

			if !errors.Is(err2, structs.ErrNoMatch) && err == nil {
				t.Errorf("event with id [" + fmt.Sprint(test.id) + "] was not deleted")
			}
		})
	}
}
