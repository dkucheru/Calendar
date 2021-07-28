package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dkucheru/Calendar/structs"
	"github.com/gorilla/mux"
)

func (rest *Rest) deleteEvent(w http.ResponseWriter, r *http.Request) {
	var id int
	var err error

	mux := mux.Vars(r)
	receivedId := mux["id"]

	id, err = strconv.Atoi(receivedId)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}
	user, _, _ := r.BasicAuth()
	err = rest.service.Events.DeleteEvent(id, user)
	if err != nil {
		if errors.Is(err, structs.ErrNoMatch) {
			rest.sendError(w, http.StatusNotFound, err)
			return
		}
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}

	rest.sendData(w, "Deleted Event")
}
