package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/dkucheru/Calendar/api"
	"github.com/dkucheru/Calendar/structs"
)

func TestAddingUserIntegration(t *testing.T) {
	//preset
	testCases := map[string]struct {
		jsonInput  []byte
		codeExpect int
		response   string
	}{
		"Ok new user": {
			[]byte(`{
				"username": "test",
				"password": "12345678",
				"location": "Local"
			 }`),
			200,
			`{"Status":1,"Data":"test Local"}`,
		},
		"Username already exists": {
			[]byte(`{
				"username": "userExists",
				"password": "12345678",
				"location": "Local"
			 }`),
			500,
			`{"Status":500,"Data":"postgres error : error occured when adding new user : pq: duplicate key value violates unique constraint \"users_pkey\""}`,
		},
		"No location": {
			[]byte(`{
				"username": "test",
				"password": "12345678"
			 }`),
			400,
			`{"Status":400,"Data":"validator : Invalid Data Format"}`,
		},
		"Bad request body": {
			[]byte(`{
			 }`),
			400,
			`{"Status":400,"Data":"validator : Invalid Data Format"}`,
		},
	}
	client := http.Client{}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			downMigrate := false
			app, err := New(downMigrate)
			if err != nil {
				t.Errorf("error launching app : " + err.Error())
				return
			}
			go func() {
				err = app.Run()
				if err != nil {
					log.Println(err.Error())
				}
				return
			}()
			log.Printf("Started server")
			CheckSignals := make(chan os.Signal, 1)
			signal.Notify(CheckSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				log.Println(fmt.Sprint(<-CheckSignals))
				log.Println("Stopping API server.")
				app.Stop()
				close(CheckSignals)
			}()
			err = app.EventsRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			err = app.UsersRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}

			app.UsersRepo.AddUser(structs.CreateUser{Username: "userExists", Password: "12345678", Location: "Local"})

			received, getUserErr := app.UsersRepo.GetUser("test")
			if getUserErr == nil && received.Username != "test" {
				t.Errorf("wanted user %s, received user %s", "test", received.Username)
			}
			//execute request
			response, err := client.Post("http://localhost:8080/users", "application/json", bytes.NewBuffer(test.jsonInput))
			if err != nil {
				t.Errorf("Error %s", err)
				CheckSignals <- os.Interrupt
			}
			body, err := ioutil.ReadAll(response.Body)
			if test.response != string(body) {
				t.Errorf("\n Wanted response : %s \n Got response : %s", test.response, string(body))
				CheckSignals <- os.Interrupt
			}
			if getUserErr == nil && response.StatusCode == http.StatusOK {
				t.Errorf("added user with the same username as another user in DB")
				CheckSignals <- os.Interrupt
			}
			if getUserErr != nil && errors.Is(getUserErr, structs.ErrNoMatch) && response.StatusCode != test.codeExpect {
				t.Errorf("unexpected response code ; response :" + string(body))
				CheckSignals <- os.Interrupt
			}
			response.Body.Close()
			CheckSignals <- os.Interrupt
		})
	}
}

func TestAddingEventIntegration(t *testing.T) {
	testCases := map[string]struct {
		jsonInput  []byte
		codeExpect int
		response   string
	}{
		"No name event": {
			[]byte(`{
				"start": "2018-12-10T13:45:00.000Z",
				"end": "2018-12-10T14:00:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
			`{"Status":400,"Data":"validator : invalid data format"}`,
		},
		"No start time event": {
			[]byte(`{
				"name": "1 on 1 Meeting",
				"end": "2018-12-10T14:00:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
			`{"Status":400,"Data":"validator : invalid data format"}`,
		},
		"No end time event": {
			[]byte(`{
				"name": "1 on 1 Meeting",
				"start": "2018-12-10T13:45:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
			`{"Status":400,"Data":"validator : invalid data format"}`,
		},
		"No mandatory fields": {
			[]byte(`{
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
			`{"Status":400,"Data":"validator : invalid data format"}`,
		},
		"No alert time": {
			[]byte(`{
				"name": "Group Meeting",
				"start": "2021-12-10T16:30:00.000Z",
				"end": "2021-12-10T17:00:00.000Z",
				"description": "Group Meeting on Daily updates"
			 }`),
			200,
			`{"Status":1,"Data":{"id":1,"name":"Group Meeting","start":"2021-12-10T16:30:00+02:00","end":"2021-12-10T17:00:00+02:00","description":"Group Meeting on Daily updates","alert":"0001-01-01T00:00:00Z"}}`,
		},
		"Only mandatory fields filled": {
			[]byte(`{
				"name": "Group Meeting",
				"start": "2001-12-10T13:45:00.000Z",
				"end": "2022-12-10T14:00:00.000Z"
			 }`),
			200,
			`{"Status":1,"Data":{"id":1,"name":"Group Meeting","start":"2001-12-10T13:45:00+02:00","end":"2022-12-10T14:00:00+02:00","description":"","alert":"0001-01-01T00:00:00Z"}}`,
		},
		"Ok meeting": {
			[]byte(`{
				"name": "Event added when db is runnning",
				"start": "2021-07-15T14:35:00.000Z",
				"end": "2021-07-15T15:00:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			200,
			`{"Status":1,"Data":{"id":1,"name":"Event added when db is runnning","start":"2021-07-15T14:35:00+03:00","end":"2021-07-15T15:00:00+03:00","description":"1 on 1 Meeting on Onboarding","alert":"2019-12-10T14:00:00+02:00"}}`,
		},
	}
	client := http.Client{}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			//preset

			downMigrate := false
			app, err := New(downMigrate)
			if err != nil {
				t.Errorf("error launching app : " + err.Error())
				return
			}
			go func() {
				err = app.Run()
				if err != nil {
					log.Println(err.Error())
				}
				return
			}()
			log.Printf("Started server")
			CheckSignals := make(chan os.Signal, 1)
			signal.Notify(CheckSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				log.Println(fmt.Sprint(<-CheckSignals))
				log.Println("Stopping API server.")
				app.Stop()
				close(CheckSignals)
			}()
			err = app.EventsRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			err = app.UsersRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			_, err = app.UsersRepo.AddUser(structs.CreateUser{Username: "test", Password: "12345678", Location: "Local"})
			if err != nil {
				t.Errorf(err.Error())
			}

			//execute request
			req, err := http.NewRequest("POST", "http://localhost:8080/events", bytes.NewBuffer(test.jsonInput))
			req.Header.Set("Content-Type", "application/json")
			req.SetBasicAuth("test", "12345678")
			if err != nil {
				t.Errorf("error occured when sending http request to add event : " + err.Error())
				CheckSignals <- os.Interrupt
				return
			}

			response, err := client.Do(req)
			if err != nil {
				t.Errorf("error occured when sending http request to add event : " + err.Error())
				CheckSignals <- os.Interrupt
				return
			}
			defer response.Body.Close()

			//check
			if err != nil {
				t.Errorf("error occured when adding event")
				CheckSignals <- os.Interrupt
				return

			}
			if err == nil && response.StatusCode != test.codeExpect {
				t.Errorf("unexpected response code ; response code : %v", response.StatusCode)
				CheckSignals <- os.Interrupt
				return

			}

			body, err := ioutil.ReadAll(response.Body)
			if test.response != string(body) {
				t.Errorf("\n Wanted response : %s \n Got response : %s", test.response, string(body))
				CheckSignals <- os.Interrupt
				return

			}

			apiResp := api.Response{}
			if err = json.Unmarshal([]byte(body), &apiResp); err != nil {
				t.Errorf("error unmarshalling response body : " + err.Error() + "\n body : " + string(body))
				CheckSignals <- os.Interrupt
				return

			}
			var respEvent structs.Event
			if err == nil && response.StatusCode == http.StatusOK {
				bodyBytes, _ := json.Marshal(apiResp.Data)
				json.Unmarshal(bodyBytes, &respEvent)
				var e structs.EventCreation
				if err = json.Unmarshal(test.jsonInput, &e); err != nil {
					t.Errorf("error unmarshalling input body : " + err.Error())
				}
				event, err := structs.CreateEvent(*time.Local, e)
				if err != nil {
					t.Errorf("error creating event object : " + err.Error())
				}
				event.Start = event.Start.Local()
				event.End = event.End.Local()
				if event.Alert != (time.Time{}) {
					event.Alert = event.Alert.Local()
				}
				event.Id = app.EventsRepo.GetLastUsedId()

				if !structs.CompareTwoEvents(respEvent, event) {
					t.Errorf("returned event is different from input event :\n wanted : %v\n got : %v", event, respEvent)
				}
			}
			response.Body.Close()
			CheckSignals <- os.Interrupt
			return
		})
	}
}

func TestChangingUserLocationIntegration(t *testing.T) {
	testCases := map[string]struct {
		url        string
		codeExpect int
		response   string
	}{
		"Bad location var": {
			"/users/test?location=amEriGa", 400,
			`{"Status":400,"Data":"Invalid location parameter"}`,
		},
		"User does not exist": {
			"/users/fakeUser?location=America/New_York", 404,
			`{"Status":404,"Data":"no matching record : user with username [fakeUser] does not exist "}`,
		},
		"Ok user ok location": {
			"/users/test?location=America/New_York",
			200,
			`{"Status":1,"Data":"America/New_York"}`,
		},
	}
	client := http.Client{}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			//preset
			downMigrate := false
			app, err := New(downMigrate)
			if err != nil {
				t.Errorf("error launching app : " + err.Error())
			}
			go func() {
				err = app.Run()
				if err != nil {
					log.Println(err.Error())
				}
				return
			}()
			log.Printf("Started server")
			CheckSignals := make(chan os.Signal, 1)
			signal.Notify(CheckSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				log.Println(fmt.Sprint(<-CheckSignals))
				log.Println("Stopping API server.")
				app.Stop()
				close(CheckSignals)
			}()
			err = app.EventsRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			err = app.UsersRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			_, err = app.UsersRepo.AddUser(structs.CreateUser{Username: "test", Password: "12345678", Location: "Local"})
			if err != nil {
				t.Errorf(err.Error())
			}

			//execute request
			if err != nil {
				t.Errorf("error occured when starting the app; error : " + err.Error())
			}
			req, err := http.NewRequest("PUT", "http://localhost:8080"+test.url, nil)
			req.SetBasicAuth("test", "12345678")
			if err != nil {
				t.Errorf("error occured when creating http request : " + err.Error())
			}
			response, err := client.Do(req)
			if err != nil {
				t.Errorf("error occured when changinf user location via http request : " + err.Error())
				CheckSignals <- os.Interrupt
				return
			}
			if err == nil && response.StatusCode != test.codeExpect {
				t.Errorf("unexpected response code ; status : " + fmt.Sprint(response.StatusCode))
				CheckSignals <- os.Interrupt
				return
			}

			body, err := ioutil.ReadAll(response.Body)
			if string(body) != test.response {
				t.Errorf("\n Wanted response : %s \n Got response : %s", test.response, string(body))
				CheckSignals <- os.Interrupt
				return
			}
			response.Body.Close()
			CheckSignals <- os.Interrupt
			return
		})
	}
}

func TestDeleteEventIntegration(t *testing.T) {
	testCases := map[string]struct {
		url        string
		codeExpect int
		response   string
	}{
		"Bad event id": {
			"/events/1on1", 400, `{"Status":400,"Data":"Invalid Data Format"}`,
		},
		"Event does not exist": {
			"/events/-1",
			404,
			`{"Status":404,"Data":"no matching record : event with id [-1] does not exist "}`,
		},
		"Ok delete event": {
			"/events/1",
			200,
			`{"Status":1,"Data":"Deleted Event"}`,
		},
	}
	client := http.Client{}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			//preset
			downMigrate := false
			app, err := New(downMigrate)
			if err != nil {
				t.Errorf("error launching app : " + err.Error())
			}
			go func() {
				err = app.Run()
				if err != nil {
					log.Println(err.Error())
				}
				return
			}()
			log.Printf("Started server")
			CheckSignals := make(chan os.Signal, 1)
			signal.Notify(CheckSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				log.Println(fmt.Sprint(<-CheckSignals))
				log.Println("Stopping API server.")
				app.Stop()
				close(CheckSignals)
			}()
			err = app.EventsRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			err = app.UsersRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			_, err = app.UsersRepo.AddUser(structs.CreateUser{Username: "test", Password: "12345678", Location: "Local"})
			if err != nil {
				t.Errorf(err.Error())
			}
			_, err = app.EventsRepo.Add(structs.Event{Name: "Ok event1", Start: time.Now(), End: time.Now().Add(time.Hour)})
			if err != nil {
				t.Errorf("error occured when adding user; error : " + err.Error())
			}
			_, err = app.EventsRepo.Add(structs.Event{Name: "Ok event2", Start: time.Now(), End: time.Now().Add(time.Hour)})
			if err != nil {
				t.Errorf("error occured when adding user; error : " + err.Error())
			}
			_, err = app.EventsRepo.Add(structs.Event{Name: "Ok event3", Start: time.Now(), End: time.Now().Add(time.Hour)})
			if err != nil {
				t.Errorf("error occured when adding user; error : " + err.Error())
			}

			//execute request
			req, err := http.NewRequest("DELETE", "http://localhost:8080"+test.url, nil)
			req.SetBasicAuth("test", "12345678")
			if err != nil {
				t.Errorf("error occured when creating http request : " + err.Error())
			}
			response, err := client.Do(req)
			if err != nil {
				t.Errorf("error occured when deleting event via http request : " + err.Error())
				CheckSignals <- os.Interrupt
			}

			//check
			if err == nil && response.StatusCode != test.codeExpect {
				t.Errorf("unexpected response code ; status : " + fmt.Sprint(response.StatusCode))
				CheckSignals <- os.Interrupt
			}
			body, err := ioutil.ReadAll(response.Body)
			if test.response != string(body) {
				t.Errorf("\n Wanted response : %s \n Got response : %s", test.response, string(body))
				CheckSignals <- os.Interrupt
			}
			response.Body.Close()
			CheckSignals <- os.Interrupt
		})
	}
}

func TestGettingEventsIntegration(t *testing.T) {
	today := time.Date(2021, 7, 28, 14, 30, 0, 0, time.Local)
	yesterday := time.Date(2021, 7, 27, 14, 30, 0, 0, time.Local)

	testCases := map[string]struct {
		url        string
		codeExpect int
		response   string
	}{
		"Get All events": {
			"/events", 200,
			`[{"id":1,"name":"Ok event","start":"2021-07-28T17:30:00+03:00","end":"2021-07-28T18:30:00+03:00","description":"","alert":"0001-01-01T00:00:00Z"},{"id":2,"name":"was yesterday","start":"2021-07-27T17:30:00+03:00","end":"2021-07-27T18:30:00+03:00","description":"","alert":"0001-01-01T00:00:00Z"},{"id":3,"name":"will be in 2 hours","start":"2021-07-28T19:30:00+03:00","end":"2021-07-28T20:30:00+03:00","description":"","alert":"0001-01-01T00:00:00Z"}]`,
		},
		"Get All events Sorted": {
			"/events?sorting=yes",
			200,
			`[{"id":2,"name":"was yesterday","start":"2021-07-27T17:30:00+03:00","end":"2021-07-27T18:30:00+03:00","description":"","alert":"0001-01-01T00:00:00Z"},{"id":1,"name":"Ok event","start":"2021-07-28T17:30:00+03:00","end":"2021-07-28T18:30:00+03:00","description":"","alert":"0001-01-01T00:00:00Z"},{"id":3,"name":"will be in 2 hours","start":"2021-07-28T19:30:00+03:00","end":"2021-07-28T20:30:00+03:00","description":"","alert":"0001-01-01T00:00:00Z"}]`,
		},
		"Full date": {
			fmt.Sprintf("/events?day=%d&month=%d&year=%d", today.Day(), today.Month(), today.Year()),
			200,
			`[{"id":1,"name":"Ok event","start":"2021-07-28T17:30:00+03:00","end":"2021-07-28T18:30:00+03:00","description":"","alert":"0001-01-01T00:00:00Z"},{"id":3,"name":"will be in 2 hours","start":"2021-07-28T19:30:00+03:00","end":"2021-07-28T20:30:00+03:00","description":"","alert":"0001-01-01T00:00:00Z"}]`,
		},
		"No matching events": {
			fmt.Sprintf("/events?day=%d&month=%d&year=%d", today.Day()+1, today.Month(), today.Year()),
			200, `[]`,
		},
		"Bad date parameters": {
			"/events?day=today&month=thisMonth&year=current",
			400,
			`{"Status":400,"Data":"error parsing date part value"}`,
		},
	}
	client := http.Client{}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			//preset
			downMigrate := false
			app, err := New(downMigrate)
			if err != nil {
				t.Errorf("error launching app : " + err.Error())
			}
			go func() {
				err = app.Run()
				if err != nil {
					log.Println(err.Error())
				}
				return
			}()
			log.Printf("Started server")
			CheckSignals := make(chan os.Signal, 1)
			signal.Notify(CheckSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				log.Println(fmt.Sprint(<-CheckSignals))
				log.Println("Stopping API server.")
				app.Stop()
				close(CheckSignals)
			}()
			err = app.EventsRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			err = app.UsersRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			_, err = app.UsersRepo.AddUser(structs.CreateUser{Username: "test", Password: "12345678", Location: "Local"})
			if err != nil {
				t.Errorf(err.Error())
			}
			_, err = app.EventsRepo.Add(structs.Event{Name: "Ok event", Start: today, End: today.Add(time.Hour)})
			if err != nil {
				t.Errorf("error occured when starting the app; error : " + err.Error())
			}
			_, err = app.EventsRepo.Add(structs.Event{Name: "was yesterday", Start: yesterday, End: yesterday.Add(time.Hour)})
			if err != nil {
				t.Errorf("error occured when starting the app; error : " + err.Error())
			}
			_, err = app.EventsRepo.Add(structs.Event{Name: "will be in 2 hours", Start: today.Add(time.Hour * 2), End: today.Add(time.Hour * 3)})
			if err != nil {
				t.Errorf("error occured when starting the app; error : " + err.Error())
			}

			req, err := http.NewRequest("GET", "http://localhost:8080"+test.url, nil)
			req.SetBasicAuth("test", "12345678")
			if err != nil {
				t.Errorf("error occured when sending http request to get events : " + err.Error())
			}

			//execute request
			response, err := client.Do(req)
			if err != nil {
				t.Errorf("Error %s", err)
				CheckSignals <- os.Interrupt
				return
			}
			if err != nil && response.StatusCode == http.StatusOK {
				t.Errorf("error occured when adding event, but returned status is Ok(200)")
			}
			if err == nil && response.StatusCode != test.codeExpect {
				t.Errorf("unexpected response code ; status : " + fmt.Sprint(response.StatusCode))
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("error reading body : " + err.Error())
				CheckSignals <- os.Interrupt
			}
			if test.response != string(body) {
				t.Errorf("\n Wanted response : %s \n Got response : %s", test.response, string(body))
				CheckSignals <- os.Interrupt
			}

			if err == nil && response.StatusCode == http.StatusOK && string(body) != "No events found" {
				query := req.URL.Query()
				params, err := api.LoadParameters(query)
				if err != nil {
					t.Errorf("error parsing parameters : " + err.Error())
				}

				resps := make([]structs.Event, 0)
				if bodyErr := json.Unmarshal([]byte(body), &resps); bodyErr != nil && err == nil {
					t.Errorf(string(body))
					t.Errorf("error unmarshalling response body : " + bodyErr.Error())
				}

				for _, respEvent := range resps {
					if !structs.SuitsParams(params, respEvent) {
						t.Errorf("returned event does not correspond to input params :\n wanted : %v\n got : %v", params, respEvent)
					}
				}

			}
			response.Body.Close()
			CheckSignals <- os.Interrupt
		})
	}
}

func TestUpdatingEventIntegration(t *testing.T) {
	testCases := map[string]struct {
		url        string
		jsonInput  []byte
		codeExpect int
		response   string
	}{
		"Nonexistent id": {
			"/events/-1",
			[]byte(`{
				"name": "Updated Name",
				"start": "2008-12-10T13:45:00.000Z",
				"end": "2008-12-10T14:00:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			404,
			`{"Status":404,"Data":"no matching record : event with id [-1] does not exist "}`,
		},
		"Bad id": {
			"/events/first",
			[]byte(`{
				"name": "Updated Name",
				"start": "2008-12-10T13:45:00.000Z",
				"end": "2008-12-10T14:00:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
			`{"Status":400,"Data":"Invalid Data Format"}`,
		},
		"Bad input": {
			"/events/1",
			[]byte(`{
				,}`),
			400,
			`{"Status":400,"Data":"invalid character ',' looking for beginning of object key string"}`,
		},
		"No mandatory fields": {
			"/events/1",
			[]byte(`{
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
			`{"Status":400,"Data":"validator : invalid data format"}`,
		},
		"Ok updated event": {
			"/events/1",
			[]byte(`{
				"name": "Updated Name",
				"start": "2008-12-10T13:45:00.000Z",
				"end": "2008-12-10T14:00:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			200,
			`{"Status":1,"Data":{"id":1,"name":"Updated Name","start":"2008-12-10T13:45:00+02:00","end":"2008-12-10T14:00:00+02:00","description":"1 on 1 Meeting on Onboarding","alert":"2019-12-10T14:00:00+02:00"}}`,
		},
	}
	client := http.Client{}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			//preset
			downMigrate := false
			app, err := New(downMigrate)
			if err != nil {
				t.Errorf("error launching app : " + err.Error())
			}
			go func() {
				err = app.Run()
				if err != nil {
					log.Println(err.Error())
				}
				return
			}()
			log.Printf("Started server")
			CheckSignals := make(chan os.Signal, 1)
			signal.Notify(CheckSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				log.Println(fmt.Sprint(<-CheckSignals))
				log.Println("Stopping API server.")
				app.Stop()
				close(CheckSignals)
			}()
			err = app.EventsRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			err = app.UsersRepo.ClearRepoData()
			if err != nil {
				t.Errorf(err.Error())
			}
			_, err = app.UsersRepo.AddUser(structs.CreateUser{Username: "test", Password: "12345678", Location: "Local"})
			if err != nil {
				t.Errorf(err.Error())
			}
			_, err = app.EventsRepo.Add(structs.Event{Name: "Ok event", Start: time.Now(), End: time.Now().Add(time.Hour)})

			if err != nil {
				t.Errorf("error occured when adding event : " + err.Error())
			}

			//execute request

			req, err := http.NewRequest("PUT", "http://localhost:8080"+test.url, bytes.NewBuffer(test.jsonInput))
			req.Header.Set("Content-Type", "application/json")
			req.SetBasicAuth("test", "12345678")
			if err != nil {
				t.Errorf("error occured when sending http request to update event : " + err.Error())
			}
			response, err := client.Do(req)
			if err != nil {
				t.Errorf("error occured when sending http request")
			}

			//check

			if err == nil && response.StatusCode != test.codeExpect {
				t.Errorf("unexpected response code ; response code : %v", response.StatusCode)
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("error reading body : " + err.Error())
				CheckSignals <- os.Interrupt
			}
			if test.response != string(body) {
				t.Errorf("\n Wanted response : %s \n Got response : %s", test.response, string(body))
				CheckSignals <- os.Interrupt
			}
			apiResp := api.Response{}
			if err = json.Unmarshal(body, &apiResp); err != nil {
				t.Errorf("error unmarshalling response body : " + err.Error() + "\n" + string(body))
			}
			var respEvent structs.Event
			if err == nil && response.StatusCode == http.StatusOK {
				bodyBytes, _ := json.Marshal(apiResp.Data)
				json.Unmarshal(bodyBytes, &respEvent)
				var e structs.EventCreation
				if err = json.Unmarshal(test.jsonInput, &e); err != nil {
					t.Errorf("error unmarshalling input body : " + err.Error())
				}
				event, err := structs.CreateEvent(*time.Local, e)
				if err != nil {
					t.Errorf("error creating event object : " + err.Error())
				}
				event.Start = event.Start.Local()
				event.End = event.End.Local()
				if event.Alert != (time.Time{}) {
					event.Alert = event.Alert.Local()
				}
				event.Id = app.EventsRepo.GetLastUsedId()

				if !structs.CompareTwoEvents(respEvent, event) {
					t.Errorf("returned event is different from input event :\n wanted : %v\n got : %v", event, respEvent)
				}
			}
			response.Body.Close()
			CheckSignals <- os.Interrupt
		})
	}
}
