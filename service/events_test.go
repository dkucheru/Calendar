package service

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dkucheru/Calendar/db"
	"github.com/dkucheru/Calendar/structs"
)

// var testRepo []*structs.Event
// FIXME: Move initialization to separate test cases.
var testRepo, _ = db.NewDBRepository()
var testService = newEventsService(testRepo)

type addEventTest struct {
	arg1      *structs.Event
	errorText string
}

type updateEventTest struct {
	arg1      int
	arg2      *structs.Event
	errorText string
}

var updateEventTests = []updateEventTest{
	updateEventTest{
		1,
		&structs.Event{
			Name:        "Ok Test Event",
			Description: "an ok event for testing",
			Start:       time.Now(),
			End:         time.Now().Add(time.Hour),
			Alert:       time.Now(),
		},
		"",
	},
	updateEventTest{
		-1,
		&structs.Event{
			Name:        "Ok Test Event",
			Description: "an ok event for testing",
			Start:       time.Now(),
			End:         time.Now().Add(time.Hour),
			Alert:       time.Now(),
		},
		"No event with such id",
	},
	updateEventTest{
		1,
		&structs.Event{
			Description: "No name field",
			Start:       time.Now(),
			End:         time.Now().Add(time.Hour),
			Alert:       time.Now(),
		},
		(&structs.MandatoryFieldError{FieldName: "name"}).Error(),
	},
	updateEventTest{
		1,
		&structs.Event{
			Name:        "No Start date Event",
			Description: "no start",
			End:         time.Now().Add(time.Hour),
			Alert:       time.Now(),
		},
		(&structs.MandatoryFieldError{FieldName: "start"}).Error(),
	},
	updateEventTest{
		1,
		&structs.Event{
			Name:        "No End date Event",
			Description: "no end",
			Start:       time.Now(),
			Alert:       time.Now(),
		},
		(&structs.MandatoryFieldError{FieldName: "end"}).Error(),
	},
	updateEventTest{
		-1,
		&structs.Event{
			Name:  "Only Mandatory Fields",
			Start: time.Now(),
			End:   time.Now().Add(time.Hour),
		},
		"No event with such id",
	},
}

type getEventsTest struct {
	arg1      *structs.EventParams
	errorText string
}

var getEventTests = []getEventsTest{
	{
		&structs.EventParams{
			Day:     time.Now().Day(),
			Month:   int(time.Now().Month()),
			Year:    time.Now().Year(),
			Name:    "Ok Test Event",
			Start:   time.Now(),
			End:     time.Now().Add(time.Hour),
			Sorting: true,
		},
		"",
	},
	{
		&structs.EventParams{
			Day: time.Now().Day(),
		},
		"",
	},
	{
		&structs.EventParams{
			Day: -1,
		},
		"Bad date parameters",
	},
	{
		&structs.EventParams{
			Month: int(time.Now().Month()),
		},
		"",
	},
}

type deleteEventTest struct {
	arg1      int
	errorText string
}

var deleteEventTests = []deleteEventTest{
	{1, ""},
	{-1, "No event with such id was found"},
}

func TestAdd(t *testing.T) {
	// testCases := map[string]struct {
	// 	event  structs.Event
	// 	result structs.Event
	// }{
	// 	"": {},
	// }

	var addEventTests = []addEventTest{
		{
			&structs.Event{
				Name:        "Ok Test Event",
				Description: "an ok event for testing",
				Start:       time.Now(),
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			"",
		},
		{
			&structs.Event{
				Name:  "No description Test Event",
				Start: time.Now(),
				End:   time.Now().Add(time.Hour),
				Alert: time.Now(),
			},
			"",
		},
		{
			&structs.Event{
				Name:  "Only mandatory fields filled Test Event",
				Start: time.Now(),
				End:   time.Now().Add(time.Hour),
			},
			"",
		},
		{
			&structs.Event{
				Description: "name field not filled event for testing",
				Start:       time.Now(),
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			(&structs.MandatoryFieldError{FieldName: "name"}).Error(),
		},
		{
			&structs.Event{
				Name:        "No start time Event",
				Description: "No start time event for testing",
				End:         time.Now().Add(time.Hour),
				Alert:       time.Now(),
			},
			(&structs.MandatoryFieldError{FieldName: "start"}).Error(),
		},
		{
			&structs.Event{
				Name:        "No End time Event",
				Description: "No end time event for testing",
				Start:       time.Now(),
				Alert:       time.Now(),
			},
			(&structs.MandatoryFieldError{FieldName: "end"}).Error(),
		},
		{
			&structs.Event{
				Name:        "Wrong duration of an Event",
				Description: "wrong duration event for testing",
				Start:       time.Now().Add(time.Hour),
				End:         time.Now(),
				Alert:       time.Now(),
			},
			"End of the event is ahead of the start",
		},
	}

	// FIXME: Use t.Run() and add names for test cases.
	for _, test := range addEventTests {
		// FIXME: Check the result
		if _, output := testService.AddEvent(test.arg1); !ErrorContains(output, test.errorText) {
			t.Errorf("got %q, wanted %q", output, test.errorText)
		}
	}
}
func TestUpdateEvent(t *testing.T) {
	for _, test := range updateEventTests {
		// FIXME: Check the result, check that event was indeed updated.
		if _, output := testService.UpdateEvent(test.arg1, test.arg2); !ErrorContains(output, test.errorText) {
			t.Errorf("got %q, wanted %q", output, test.errorText)
		}
	}
}

func TestGetEvent(t *testing.T) {
	for _, test := range getEventTests {
		// FIXME: Check the result
		if _, output := testService.GetEventsOfTheDay(test.arg1); !ErrorContains(output, test.errorText) {
			t.Errorf("got %q, wanted %q", output, test.errorText)
		}
	}
}

func TestDeleteEvent(t *testing.T) {
	for _, test := range deleteEventTests {
		// FIXME: Check that event was indeed deleted.
		if output := testService.DeleteEvent(test.arg1); !ErrorContains(output, test.errorText) {
			t.Errorf("got %q, wanted %q", output, test.errorText)
		}
	}
}

func BenchmarkAddEvent(b *testing.B) {
	testService := new(eventService)
	var testEvent *structs.Event
	testEvent = &structs.Event{
		Name:        "Ok Test Event",
		Description: "an ok event for testing",
		Start:       time.Now(),
		End:         time.Now().Add(time.Hour),
		Alert:       time.Now()}

	for i := 0; i < b.N; i++ {
		testService.AddEvent(testEvent)
	}
}

func BenchmarkUpdateEvent(b *testing.B) {
	testService := new(eventService)
	var testEvent *structs.Event
	testEvent = &structs.Event{
		Name:        "Ok Test Event",
		Description: "an ok event for testing",
		Start:       time.Now(),
		End:         time.Now().Add(time.Hour),
		Alert:       time.Now()}

	for i := 0; i < b.N; i++ {
		testService.UpdateEvent(1, testEvent)
	}
}

func BenchmarkGetEvent(b *testing.B) {
	testService := new(eventService)
	var testEventParams *structs.EventParams
	testEventParams = &structs.EventParams{
		Day:     time.Now().Day(),
		Month:   int(time.Now().Month()),
		Year:    time.Now().Year(),
		Name:    "Ok Test Event",
		Start:   time.Now(),
		End:     time.Now().Add(time.Hour),
		Sorting: true}

	for i := 0; i < b.N; i++ {
		testService.GetEventsOfTheDay(testEventParams)
	}
}

func ExampleAddEvent() {
	testService := new(eventService)
	fmt.Println(testService.AddEvent(&structs.Event{
		Name:        "Ok Test Event",
		Description: "an ok event for testing",
		Start:       time.Now(),
		End:         time.Now().Add(time.Hour),
		Alert:       time.Now(),
	}))
}
func ExampleUpdateEvent() {
	testService := new(eventService)
	fmt.Println(testService.GetEventsOfTheDay(&structs.EventParams{
		Day:     time.Now().Day(),
		Month:   int(time.Now().Month()),
		Year:    time.Now().Year(),
		Name:    "Ok Test Event",
		Start:   time.Now(),
		End:     time.Now().Add(time.Hour),
		Sorting: true}))
}

func ExampleGetEvent() {
	testService := new(eventService)
	fmt.Println(testService.UpdateEvent(1, &structs.Event{
		Name:        "Ok Test Event",
		Description: "an ok event for testing",
		Start:       time.Now(),
		End:         time.Now().Add(time.Hour),
		Alert:       time.Now(),
	}))
}

func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}
