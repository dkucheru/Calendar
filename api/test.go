package api

import (
	"net/http"
)

func (rest *Rest) showSomething(w http.ResponseWriter, r *http.Request) {
	rest.sendData(w, "Successful function")
}
