package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dkucheru/Calendar/structs"
	"github.com/gorilla/mux"
)

func (rest *Rest) updateEvent(w http.ResponseWriter, r *http.Request) {
	mux := mux.Vars(r)
	receivedId := mux["id"]

	id, err := strconv.Atoi(receivedId)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}

	//Read incoming JSON from request body
	data, err := ioutil.ReadAll(r.Body)

	// If no body is associated return with StatusBadRequest
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
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
	updatedEvent, err := rest.service.Events.UpdateEvent(id, event)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}

	rest.sendData(w, updatedEvent)
}
