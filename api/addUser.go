package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dkucheru/Calendar/structs"
	"github.com/go-playground/validator/v10"
)

func (rest *Rest) addUser(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}

	var newUser structs.CreateUser
	err = json.Unmarshal(data, &newUser)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("Invalid Data Format"))
		return
	}
	validate := validator.New()
	err = validate.Struct(newUser)
	if err != nil {
		rest.sendError(w, http.StatusBadRequest, errors.New("validator : Invalid Data Format"))
		return
	}

	_, err = rest.service.Users.AddUser(newUser)
	if err != nil {
		rest.sendError(w, http.StatusInternalServerError, err)
		return
	}
	loc, err := time.LoadLocation(newUser.Location)
	rest.sendData(w, time.Now().In(loc).String())
}
