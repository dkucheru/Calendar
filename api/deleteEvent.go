package api

import (
	"net/http"
	"strconv"
)

func (rest *Rest) deleteEvent(w http.ResponseWriter, r *http.Request) {
	var id int
	var err error

	query := r.URL.Query()
	receivedId := query.Get("id")

	id, err = strconv.Atoi(receivedId)
	if err != nil {
		rest.sendError(w, http.StatusExpectationFailed, err)
		return
	}

	err = rest.service.Events.DeleteEvent(id)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}

	rest.sendData(w, "Deleted Event")
}
