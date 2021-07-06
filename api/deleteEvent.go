package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (rest *Rest) deleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	start := params["date"]

	var startTime time.Time
	err := json.Unmarshal([]byte(start), &startTime)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte("Invalid Date Format"))
		return
	}

	// fmt.Println(name)
	// fmt.Println(startTime)

	err = rest.service.Events.DeleteEvent(&name, &startTime)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write([]byte(err.Error()))
		return
	}

	// return after writing Body
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Deleted Event"))
}
