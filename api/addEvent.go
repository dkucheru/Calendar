package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dkucheru/Calendar/structs"
)

func (rest *Rest) addEvent(w http.ResponseWriter, r *http.Request) {
	//Read incoming JSON from request body
	data, err := ioutil.ReadAll(r.Body)

	// If no body is associated return with StatusBadRequest
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if data is proper JSON (data validation)
	var event structs.Event
	err = json.Unmarshal(data, &event)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("Invalid Data Format"))
		return
	}

	//add event to memory
	err = rest.service.Events.AddEvent(&event)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(err.Error()))
		return
	}

	// return after writing Body
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Added New Product"))
}
