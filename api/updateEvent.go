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

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}

	user, _, _ := r.BasicAuth()
	loc, err := rest.service.Users.GetUserLocation(user)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, err)
		return
	}
	var e structs.EventCreation
	if err = json.Unmarshal(data, &e); err != nil {
		rest.sendError(w, http.StatusBadRequest, err)
		return
	}
	event, err := structs.CreateEvent(loc, e)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, err)
		return
	}

	updatedEvent, err := rest.service.Events.UpdateEvent(id, event, loc)
	if err != nil {
		if errors.Is(err, structs.ErrNoMatch) {
			rest.sendError(w, http.StatusNotFound, err)
			return
		}
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}

	rest.sendData(w, updatedEvent)
}
