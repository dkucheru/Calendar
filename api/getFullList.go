package api

import (
	"encoding/json"
	"net/http"
)

func (rest *Rest) allEvents(w http.ResponseWriter, r *http.Request) {

	err, events := rest.service.Events.GetAll()

	if err != nil {
		rest.sendError(w, err)
		return
	}
	for _, event := range events {
		json.NewEncoder(w).Encode(event)

	}
}
