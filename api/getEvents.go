package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dkucheru/Calendar/structs"
)

type paramsCheck struct {
	intParams  map[string]*int
	timeParams map[string]*time.Time

	valuesFromURL map[string]string
}

const (
	day   = "Day"
	week  = "Week"
	month = "Month"
	year  = "Year"
	start = "Start"
	end   = "End"
)

func (rest *Rest) allEvents(w http.ResponseWriter, r *http.Request) {
	user, _, _ := r.BasicAuth()
	loc, err := rest.service.Users.GetUserLocation(user)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, err)
		return
	}

	query := r.URL.Query()
	params, err := LoadParameters(query)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, err)
		return
	}

	events, err := rest.service.Events.GetEventsOfTheDay(params, loc)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}
	eventsJSON, err := json.Marshal(events)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}
	_, err = w.Write(eventsJSON)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}
}

func LoadParameters(query url.Values) (structs.EventParams, error) {
	params := structs.EventParams{}
	params.Name = query.Get("name")
	if sort := query.Get("sorting"); sort != "" {
		params.Sorting = true
	}

	receivedParams := map[string]string{
		day:   query.Get("day"),
		week:  query.Get("week"),
		month: query.Get("month"),
		year:  query.Get("year"),
		start: query.Get("start"),
		end:   query.Get("end"),
	}

	isNum := map[string]bool{
		day:   true,
		week:  true,
		month: true,
		year:  true,
	}
	isTime := map[string]bool{
		start: true,
		end:   true,
	}

	for paramName, url := range receivedParams {
		if url != "" {
			if isNum[paramName] {
				newInt, err := strconv.Atoi(url)
				if err != nil {
					return structs.EventParams{}, errors.New("error parsing date part value")
				}
				switch paramName {
				case day:
					params.Day = newInt
				case week:
					params.Week = newInt
				case month:
					params.Month = newInt
				case year:
					params.Year = newInt
				}
			} else if isTime[paramName] {
				var dest *time.Time
				err := json.Unmarshal([]byte(url), dest)
				if err != nil {
					return structs.EventParams{}, err
				}
				switch paramName {
				case start:
					params.Start = *dest
				case end:
					params.End = *dest
				}
			}
		}
	}
	return params, nil
}
