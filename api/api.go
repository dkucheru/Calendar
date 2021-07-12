package api

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/dkucheru/Calendar/service"
	"github.com/gorilla/mux"
)

type Rest struct {
	address  string
	mux      *mux.Router
	listener net.Listener
	service  *service.Service
}

func New(address string, service *service.Service) *Rest {
	rest := &Rest{
		address: address,
		service: service,
	}

	api := mux.NewRouter()

	api.HandleFunc("/events", rest.addEvent).Methods("POST")
	api.HandleFunc("/events/{id}", rest.deleteEvent).Methods("DELETE")
	api.HandleFunc("/events/{id}", rest.updateEvent).Methods("PUT")
	api.HandleFunc("/events", rest.allEvents).Methods("GET")
	rest.mux = api

	return rest
}

func (rest *Rest) Listen() (err error) {
	rest.listener, err = net.Listen("tcp", rest.address)
	if err != nil {
		return err
	}

	r := http.NewServeMux()
	r.Handle("/", rest.mux)

	server := &http.Server{Handler: r}

	rest.setupMiddleware()

	return server.Serve(rest.listener)
}

func (rest *Rest) setupMiddleware() {
	rest.mux.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("new request", r.RequestURI)
			handler.ServeHTTP(w, r)
		})
	})
}

type Response struct {
	Status int
	Data   interface{}
}

func (rest *Rest) sendError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)

	bytes, err := json.Marshal(Response{
		Status: statusCode,
		Data:   err.Error(),
	})

	if err != nil {
		log.Println(err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Println(err)
	}
}

func (rest *Rest) sendData(w http.ResponseWriter, data interface{}) {
	bytes, err := json.Marshal(Response{
		Status: 1,
		Data:   data,
	})

	if err != nil {
		log.Println(err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Println(err)
	}
}
