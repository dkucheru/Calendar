package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dkucheru/Calendar/structs"
	"github.com/go-playground/validator/v10"
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

	user, _, _ := r.BasicAuth()
	userLocation, err := rest.service.Events.GetUserLocation(user)

	//Read incoming JSON from request body
	data, err := ioutil.ReadAll(r.Body)

	// If no body is associated return with StatusBadRequest
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}

	// Check if data is proper JSON (data validation)
	var event structs.Event
	// err = json.Unmarshal(data, &event)
	err = event.ParseJSON(data, userLocation)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}
	validate := validator.New()
	err = validate.Struct(event)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("validator : Invalid Data Format"))
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
