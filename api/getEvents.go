package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dkucheru/Calendar/structs"
)

func (rest *Rest) allEvents(w http.ResponseWriter, r *http.Request) {

	var events []structs.Event
	var err error
	var params structs.EventParams

	query := r.URL.Query()
	receivedDay := query.Get("day")
	receivedWeek := query.Get("week")
	receivedMonth := query.Get("month")
	receivedYear := query.Get("year")
	receivedSorting := query.Get("sorting")

	params.Name = query.Get("name")
	receivedStart := query.Get("start")
	receivedEnd := query.Get("end")

	if receivedDay != "" {
		params.Day, err = strconv.Atoi(receivedDay)
		if err != nil {
			rest.sendError(w, http.StatusBadRequest, err)
			return
		}
	}

	if receivedWeek != "" {
		params.Week, err = strconv.Atoi(receivedWeek)
		if err != nil {
			rest.sendError(w, http.StatusBadRequest, err)
			return
		}
	}

	if receivedMonth != "" {
		params.Month, err = strconv.Atoi(receivedMonth)
		if err != nil {
			rest.sendError(w, http.StatusBadRequest, err)
			return
		}
	}

	if receivedYear != "" {
		params.Year, err = strconv.Atoi(receivedYear)
		if err != nil {
			rest.sendError(w, http.StatusBadRequest, err)
			return
		}
	}

	if receivedStart != "" {
		err = json.Unmarshal([]byte(receivedStart), &params.Start)
		if err != nil {
			rest.sendError(w, http.StatusBadRequest, err)
			return
		}
	}

	if receivedEnd != "" {
		err = json.Unmarshal([]byte(receivedEnd), &params.End)
		if err != nil {
			rest.sendError(w, http.StatusBadRequest, err)
			return
		}
	}

	if receivedSorting != "" {
		params.Sorting = true
	}

	events, err = rest.service.Events.GetEventsOfTheDay(&params)

	for _, event := range events {
		json.NewEncoder(w).Encode(event)

	}
}
