package service

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dkucheru/Calendar/structs"
)

type addEventTest struct {
	arg1      *structs.Event
	errorText string
}

var addEventTests = []addEventTest{
	addEventTest{
		&structs.Event{
			Name:        "Ok Test Event",
			Description: "an ok event for testing",
			Start:       time.Now(),
			End:         time.Now().Add(time.Hour),
			Alert:       time.Now(),
		},
		"",
	},
	addEventTest{
		&structs.Event{
			Name:  "No description Test Event",
			Start: time.Now(),
			End:   time.Now().Add(time.Hour),
			Alert: time.Now(),
		},
		"",
	},
	addEventTest{
		&structs.Event{
			Name:  "Only mandatory fields filled Test Event",
			Start: time.Now(),
			End:   time.Now().Add(time.Hour),
		},
		"",
	},
	addEventTest{
		&structs.Event{
			Description: "name field not filled event for testing",
			Start:       time.Now(),
			End:         time.Now().Add(time.Hour),
			Alert:       time.Now(),
		},
		(&structs.MandatoryFieldError{FieldName: "name"}).Error(),
	},
	addEventTest{
		&structs.Event{
			Name:        "No start time Event",
			Description: "No start time event for testing",
			End:         time.Now().Add(time.Hour),
			Alert:       time.Now(),
		},
		(&structs.MandatoryFieldError{FieldName: "start"}).Error(),
	},
	addEventTest{
		&structs.Event{
			Name:        "No End time Event",
			Description: "No end time event for testing",
			Start:       time.Now(),
			Alert:       time.Now(),
		},
		(&structs.MandatoryFieldError{FieldName: "end"}).Error(),
	},
	addEventTest{
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

func TestAdd(t *testing.T) {
	testService := new(eventService)

	for _, test := range addEventTests {
		if output := testService.AddEvent(test.arg1); !ErrorContains(output, test.errorText) {
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

func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}
