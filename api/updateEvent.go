package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dkucheru/Calendar/structs"
	"github.com/gorilla/mux"
)

func (rest *Rest) updateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	receivedId := params["id"]

	id, err := strconv.Atoi(receivedId)
	if err != nil {
		rest.sendError(w, err)
		return
	}

	//Read incoming JSON from request body
	data, err := ioutil.ReadAll(r.Body)

	// If no body is associated return with StatusBadRequest
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if data is proper JSON (data validation)
	var event structs.Event
	err = json.Unmarshal(data, &event)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("Invalid Data Format"))
		return
	}

	//add event to memory
	err = rest.service.Events.UpdateEvent(id, &event)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(err.Error()))
		return
	}

	// return after writing Body
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Upsated Eb=vet"))
}
