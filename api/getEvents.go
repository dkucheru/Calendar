package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/dkucheru/Calendar/sorting"
	"github.com/dkucheru/Calendar/structs"
)

func (rest *Rest) allEvents(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	receivedDay := query.Get("day")
	receivedWeek := query.Get("week")
	receivedMonth := query.Get("month")
	receivedYear := query.Get("year")
	receivedSorting := query.Get("sorting")

	receivedName := query.Get("name")
	receivedStart := query.Get("start")
	receivedEnd := query.Get("end")

	var events []structs.Event
	var err error
	var startTime, endTime time.Time
	var day, month, year, week int
	var params structs.EventParams
	// startTime := (time.Time{})
	// endTime := (time.Time{})
	// day := 0
	// month := 0
	// year := 0
	// week := 0

	if receivedStart != "" {
		err = json.Unmarshal([]byte(receivedStart), &params.Start)
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("Invalid Date Format"))
			return
		}
	}

	if receivedEnd != "" {
		err = json.Unmarshal([]byte(receivedEnd), &params.End)
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write([]byte("Invalid Date Format"))
			return
		}
	}

	if receivedDay == "" && receivedMonth == "" &&
		receivedYear == "" && receivedWeek == "" {
		//calling a service function to receive an araay of all events in current calendar
		err, events = rest.service.Events.GetAll()
		if err != nil {
			rest.sendError(w, http.StatusBadRequest, err)
			return
		}
	} else {
		if receivedDay != "" {
			params.Day, err = strconv.Atoi(receivedDay)
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

		if receivedWeek != "" {
			params.Week, err = strconv.Atoi(receivedWeek)
			if err != nil {
				rest.sendError(w, http.StatusBadRequest, err)
				return
			}
		}

		events, err = rest.service.Events.GetEventsOfTheDay(&params)
	}

	//if optional parameter {sorting} was added we sort resulting list of events by time of the start
	if receivedSorting != "" {
		sort.Sort(sorting.ByStartTime(events))
	}

	for _, event := range events {
		json.NewEncoder(w).Encode(event)

	}
}
