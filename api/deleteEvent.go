package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (rest *Rest) deleteEvent(w http.ResponseWriter, r *http.Request) {
	var id int
	var err error

	mux := mux.Vars(r)
	// query := r.URL.Query()
	receivedId := mux["id"]

	id, err = strconv.Atoi(receivedId)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}

	err = rest.service.Events.DeleteEvent(id)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}

	rest.sendData(w, "Deleted Event")
}
