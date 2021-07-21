package api

import (
	"encoding/json"
	"errors"
	"net/http"
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
	day   = "day"
	week  = "week"
	month = "month"
	year  = "year"
	start = "start"
	end   = "end"
)

func (rest *Rest) allEvents(w http.ResponseWriter, r *http.Request) {
	var params structs.EventParams

	user, _, _ := r.BasicAuth()
	loc, err := rest.service.Users.GetUserLocation(user)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, err)
		return
	}

	query := r.URL.Query()
	receivedSorting := query.Get("sorting")
	params.Name = query.Get("name")

	receivedParams := map[string]string{
		day:   query.Get("day"),
		week:  query.Get("week"),
		month: query.Get("month"),
		year:  query.Get("year"),
		start: query.Get("start"),
		end:   query.Get("end"),
	}
	intParams := map[string]*int{
		day:   &params.Day,
		week:  &params.Week,
		month: &params.Month,
		year:  &params.Year,
	}
	timeParams := map[string]*time.Time{
		start: &params.Start,
		end:   &params.End,
	}

	var newStruct paramsCheck
	newStruct = paramsCheck{
		intParams:     intParams,
		timeParams:    timeParams,
		valuesFromURL: receivedParams,
	}

	if !checkAllParameters(newStruct) {
		rest.sendError(w, http.StatusBadRequest, errors.New("Not valid data input"))
		return
	}

	if receivedSorting != "" {
		params.Sorting = true
	}

	events, err := rest.service.Events.GetEventsOfTheDay(params, loc)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}
	for _, event := range events {
		json.NewEncoder(w).Encode(event)
	}

	rest.sendData(w, "Everything is fine")
}

func checkAllParameters(r paramsCheck) bool {
	for paramName := range r.intParams {
		if r.valuesFromURL[paramName] != "" {
			newInt, err := strconv.Atoi(r.valuesFromURL[paramName])
			if err != nil {
				return false
			}
			r.intParams[paramName] = &newInt
		}
	}

	for paramName, destination := range r.timeParams {
		if r.valuesFromURL[paramName] != "" {
			err := json.Unmarshal([]byte(r.valuesFromURL[paramName]), destination)
			if err != nil {
				return false
			}
		}
	}
	return true
}
