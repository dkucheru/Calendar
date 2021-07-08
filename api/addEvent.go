package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dkucheru/Calendar/structs"
)

func (rest *Rest) addEvent(w http.ResponseWriter, r *http.Request) {
	//Read incoming JSON from request body
	data, err := ioutil.ReadAll(r.Body)

	// If no body is associated return with StatusBadRequest
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, err)
		return
	}

	// Check if data is proper JSON (data validation)
	var event structs.Event
	err = json.Unmarshal(data, &event)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}

	//add event to memory
	err = rest.service.Events.AddEvent(&event)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}

	rest.sendData(w, "Added new event")
}
