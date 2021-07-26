package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dkucheru/Calendar/api"
	"github.com/dkucheru/Calendar/structs"
)

func TestAddingUserIntegration(t *testing.T) {
	//preset
	app, err := New()
	if err != nil {
		t.Errorf("error launching app : " + err.Error())
	}
	err = app.EventsRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	err = app.UsersRepo.ClearRepoData()
	if err != nil {
		t.Errorf(err.Error())
	}
	testCases := map[string]struct {
		jsonInput  []byte
		codeExpect int
	}{
		"Ok new user": {
			[]byte(`{
				"username": "test",
				"password": "12345678",
				"location": "Local"
			 }`),
			200,
		},
		"Username already exists": {
			[]byte(`{
				"username": "test",
				"password": "12345678",
				"location": "Local"
			 }`),
			500,
		},
		"No location": {
			[]byte(`{
				"username": "test",
				"password": "12345678"
			 }`),
			400,
		},
		"Bad request body": {
			[]byte(`{
			 }`),
			400,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			if err != nil {
				t.Errorf("error occured when starting the app; error : " + err.Error())
			}
			received, getUserErr := app.UsersRepo.GetUser("test")
			if getUserErr == nil && received.Username != "test" {
				t.Errorf("wanted user %s, received user %s", "test", received.Username)
			}

			req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(test.jsonInput))
			req.Header.Set("Content-Type", "application/json")
			// req.SetBasicAuth("test", "testpass")
			if err != nil {
				t.Errorf("error occured when sending http request to add user : " + err.Error())
			}
			//execute request
			response := httptest.NewRecorder()
			app.Api.Mux.ServeHTTP(response, req)
			if getUserErr == nil && response.Code == http.StatusOK {
				t.Errorf("added user with the same username as another user in DB")
			}
			if getUserErr != nil && getUserErr.Error() == "no matching record" && response.Code != test.codeExpect {
				t.Errorf("unexpected response code ; response :" + response.Body.String())
			}
		})
	}
}

func TestAddingEventIntegration(t *testing.T) {
	//preset
	app, err := New()
	if err != nil {
		t.Errorf("error launching app : " + err.Error())
	}
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
	testCases := map[string]struct {
		jsonInput  []byte
		codeExpect int
	}{
		"No name event": {
			[]byte(`{
				"start": "2018-12-10T13:45:00.000Z",
				"end": "2018-12-10T14:00:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
		},
		"No start time event": {
			[]byte(`{
				"name": "1 on 1 Meeting",
				"end": "2018-12-10T14:00:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
		},
		"No end time event": {
			[]byte(`{{
				"name": "1 on 1 Meeting",
				"start": "2018-12-10T13:45:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
		},
		"No mandatory fields": {
			[]byte(`{
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
		},
		"No alert time": {
			[]byte(`{
				"name": "Group Meeting",
				"start": "2021-12-10T16:30:00.000Z",
				"end": "2021-12-10T17:00:00.000Z",
				"description": "Group Meeting on Daily updates"
			 }`),
			200,
		},
		"Only mandatory fields filled": {
			[]byte(`{
				"name": "Group Meeting",
				"start": "2001-12-10T13:45:00.000Z",
				"end": "2022-12-10T14:00:00.000Z"
			 }`),
			200,
		},
		"Ok meeting": {
			[]byte(`{
				"name": "Event added when docker is runnning",
				"start": "2021-07-15T14:35:00.000Z",
				"end": "2021-07-15T15:00:00.000Z",
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			200,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			if err != nil {
				t.Errorf("error occured when starting the app; error : " + err.Error())
			}
			req, err := http.NewRequest("POST", "/events", bytes.NewBuffer(test.jsonInput))
			req.Header.Set("Content-Type", "application/json")
			req.SetBasicAuth("test", "12345678")
			if err != nil {
				t.Errorf("error occured when sending http request to add user : " + err.Error())
			}
			//execute request
			response := httptest.NewRecorder()
			app.Api.Mux.ServeHTTP(response, req)

			//check
			if err != nil && response.Code == http.StatusOK {
				t.Errorf("error occured when adding event, but returned status is Ok(200)")
			}
			if err == nil && response.Code != test.codeExpect {
				t.Errorf("unexpected response code ; response :" + response.Body.String())
			}
			body, bodyErr := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("error reading response body : " + bodyErr.Error())
			}
			apiResp := api.Response{}
			if bodyErr = json.Unmarshal(body, &apiResp); bodyErr != nil {
				t.Errorf("error unmarshalling response body : " + bodyErr.Error())
			}
			var respEvent structs.Event
			if bodyErr == nil && err == nil && response.Code == http.StatusOK {
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

		})
	}
}

func TestChangingUserLocationIntegration(t *testing.T) {
	//preset
	app, err := New()
	if err != nil {
		t.Errorf("error launching app : " + err.Error())
	}
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

	testCases := map[string]struct {
		url        string
		codeExpect int
	}{
		"Bad location var": {
			"/users/test?location=amEriGa", 400,
		},
		"User does not exist": {
			"/users/fakeUser?location=America/New_York",
			500,
		},
		"Ok user ok location": {
			"/users/test?location=America/New_York",
			200,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			if err != nil {
				t.Errorf("error occured when starting the app; error : " + err.Error())
			}
			req, err := http.NewRequest("PUT", test.url, nil)
			req.SetBasicAuth("test", "12345678")
			if err != nil {
				t.Errorf("error occured when sending http request to add user : " + err.Error())
			}
			//execute request
			response := httptest.NewRecorder()
			app.Api.Mux.ServeHTTP(response, req)
			if err != nil && response.Code == http.StatusOK {
				t.Errorf("error occured when adding event, but returned status is Ok(200)")
			}
			if err == nil && response.Code != test.codeExpect {
				t.Errorf("unexpected response code ; response :" + response.Body.String() + "status : " + fmt.Sprint(response.Code))
			}
		})
	}
}
func TestDeleteEventIntegration(t *testing.T) {
	//preset
	app, err := New()
	if err != nil {
		t.Errorf("error launching app : " + err.Error())
	}
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
	testCases := map[string]struct {
		url        string
		codeExpect int
	}{
		"Bad event id": {
			"/events/1on1", 400,
		},
		"Event does not exist": {
			"/events/-1",
			500,
		},
		"Ok delete event": {
			"/events/1",
			200,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			if err != nil {
				t.Errorf("error occured when starting the app; error : " + err.Error())
			}
			req, err := http.NewRequest("DELETE", test.url, nil)
			req.SetBasicAuth("test", "12345678")
			if err != nil {
				t.Errorf("error occured when sending http request to add user : " + err.Error())
			}
			//execute request
			response := httptest.NewRecorder()
			app.Api.Mux.ServeHTTP(response, req)
			if err != nil && response.Code == http.StatusOK {
				t.Errorf("error occured when adding event, but returned status is Ok(200)")
			}
			if err == nil && response.Code != test.codeExpect {
				t.Errorf("unexpected response code ; response :" + response.Body.String() + "status : " + fmt.Sprint(response.Code))
			}
		})
	}
}

func TestGettingEventsIntegration(t *testing.T) {
	//preset
	app, err := New()
	if err != nil {
		t.Errorf("error launching app : " + err.Error())
	}
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
	testCases := map[string]struct {
		url        string
		codeExpect int
	}{
		"Get All events": {
			"/events", 200,
		},
		"Get All events Sorted": {
			"/events?sorting=yes",
			200,
		},
		"Full date": {
			fmt.Sprintf("/events?day=%d&month=%d&year=%d", time.Now().Day(), time.Now().Month(), time.Now().Year()),
			200,
		},
		"No matching events": {
			fmt.Sprintf("/events?day=%d&month=%d&year=%d", time.Now().Day()+1, time.Now().Month(), time.Now().Year()),
			200,
		},
		"Bad date parameters": {
			"/events?day=today&month=thisMonth&year=current",
			400,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			if err != nil {
				t.Errorf("error occured when starting the app; error : " + err.Error())
			}
			req, err := http.NewRequest("GET", test.url, nil)
			req.SetBasicAuth("test", "12345678")
			if err != nil {
				t.Errorf("error occured when sending http request to add user : " + err.Error())
			}

			//execute request
			response := httptest.NewRecorder()
			app.Api.Mux.ServeHTTP(response, req)
			if err != nil && response.Code == http.StatusOK {
				t.Errorf("error occured when adding event, but returned status is Ok(200)")
			}
			if err == nil && response.Code != test.codeExpect {
				t.Errorf("unexpected response code ; response :" + response.Body.String() + "status : " + fmt.Sprint(response.Code))
			}
			body, bodyErr := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("error reading response body : " + bodyErr.Error())
			}
			type ApiResponses struct {
				events []structs.Event
			}
			if bodyErr == nil && err == nil && response.Code == http.StatusOK {
				query := req.URL.Query()
				params, err := api.LoadParameters(query)
				if err != nil {
					t.Errorf("error parsing parameters : " + err.Error())
				}
				if string(body) == "" {
					//check if really no event corresponds to input parameters
					allEvents, err := app.EventsRepo.Get(structs.EventParams{})
					if err != nil {
						t.Errorf("error getting all events from Db : " + err.Error())
					}
					for _, e := range allEvents {
						if structs.SuitsParams(params, e) {
							t.Errorf("an event which suits input parameter was not in the json output")
						}
					}
				} else {
					//check if events in output json actually suit the parameters
					resps := ApiResponses{}
					if bodyErr = json.Unmarshal(body, &resps); bodyErr != nil && err == nil {
						t.Errorf(string(body))
						t.Errorf("error unmarshalling response body : " + bodyErr.Error())
					}

					for _, respEvent := range resps.events {
						if !structs.SuitsParams(params, respEvent) {
							t.Errorf("returned event does not correspond to input params :\n wanted : %v\n got : %v", params, respEvent)
						}
					}

				}

			}
		})
	}
}

func TestUpdatingEventIntegration(t *testing.T) {
	//preset
	app, err := New()
	if err != nil {
		t.Errorf("error launching app : " + err.Error())
	}
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

	testCases := map[string]struct {
		url        string
		jsonInput  []byte
		codeExpect int
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
			500,
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
		},
		"Bad input": {
			"/events/1",
			[]byte(`{
				,}`),
			400,
		},
		"No mandatory fields": {
			"/events/1",
			[]byte(`{
				"description": "1 on 1 Meeting on Onboarding",
				"alert": "2019-12-10T14:00:00.000Z"
			 }`),
			400,
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
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			if err != nil {
				t.Errorf("error occured when starting the app; error : " + err.Error())
			}
			req, err := http.NewRequest("PUT", test.url, bytes.NewBuffer(test.jsonInput))
			req.Header.Set("Content-Type", "application/json")
			req.SetBasicAuth("test", "12345678")
			if err != nil {
				t.Errorf("error occured when sending http request to add user : " + err.Error())
			}
			//execute request
			response := httptest.NewRecorder()
			app.Api.Mux.ServeHTTP(response, req)
			if err != nil && response.Code == http.StatusOK {
				t.Errorf("error occured when adding event, but returned status is Ok(200)")
			}
			if err == nil && response.Code != test.codeExpect {
				t.Errorf("unexpected response code ; response :" + response.Body.String())
			}
			body, bodyErr := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("error reading response body : " + bodyErr.Error())
			}
			apiResp := api.Response{}
			if bodyErr = json.Unmarshal(body, &apiResp); bodyErr != nil {
				t.Errorf("error unmarshalling response body : " + bodyErr.Error())
			}
			var respEvent structs.Event
			if bodyErr == nil && err == nil && response.Code == http.StatusOK {
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
		})
	}
}
