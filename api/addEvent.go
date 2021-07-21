package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dkucheru/Calendar/structs"
)

func (rest *Rest) addEvent(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}
	var e structs.EventCreation
	if err = json.Unmarshal(data, &e); err != nil {
		rest.sendError(w, http.StatusBadRequest, err)
		return
	}
	user, _, _ := r.BasicAuth()
	loc, err := rest.service.Users.GetUserLocation(user)
	event, err := structs.CreateEvent(loc, e)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, err)
		return
	}
	newEvent, err := rest.service.Events.AddEvent(loc, event)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}

	rest.sendData(w, newEvent)
}
