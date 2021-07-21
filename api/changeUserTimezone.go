package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (rest *Rest) changeTimezone(w http.ResponseWriter, r *http.Request) {
	mux := mux.Vars(r)
	username := mux["username"]

	query := r.URL.Query()
	location := query.Get("location")

	if username == "" {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid username"))
		return
	}

	loc, err := time.LoadLocation(location)
	if err != nil || location == "" {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid location parameter"))
		return
	}

	newLocation, err := rest.service.Users.UpdateLocation(username, *loc)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}

	rest.sendData(w, newLocation)
}
