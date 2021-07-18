package api

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dkucheru/Calendar/structs"
	"github.com/go-playground/validator/v10"
)

func (rest *Rest) addEvent(w http.ResponseWriter, r *http.Request) {
	//Read incoming JSON from request body
	data, err := ioutil.ReadAll(r.Body)

	// If no body is associated return with StatusBadRequest
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}
	// FIXME: send location to service
	user, _, _ := r.BasicAuth()
	userLocation, err := rest.service.Events.GetUserLocation(user)

	var event structs.Event

	// FIXME: make ParseJSON a constructor
	err = event.ParseJSON(data, userLocation)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("invalid data format"))
		return
	}
	validate := validator.New()
	err = validate.Struct(event)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("validator : invalid data format"))
		return
	}

	newEvent, err := rest.service.Events.AddEvent(event)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}

	newEvent.Start = newEvent.Start.In(&userLocation)
	newEvent.End = newEvent.Start.In(&userLocation)
	newEvent.Alert = newEvent.Start.In(&userLocation)

	rest.sendData(w, newEvent)
}
