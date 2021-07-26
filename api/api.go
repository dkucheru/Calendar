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
	Mux      *mux.Router
	listener net.Listener
	service  *service.Service
}

func New(address string, service *service.Service) *Rest {
	rest := &Rest{
		address: address,
		service: service,
	}

	api := mux.NewRouter()
	api.HandleFunc("/users", rest.addUser).Methods("POST")
	api.Handle("/users/{username}", rest.BasicAuthMiddleware(http.HandlerFunc(rest.changeTimezone))).Methods("PUT")

	api.Handle("/events", rest.BasicAuthMiddleware(http.HandlerFunc(rest.addEvent))).Methods("POST")
	api.Handle("/events", rest.BasicAuthMiddleware(http.HandlerFunc(rest.allEvents))).Methods("GET")
	api.Handle("/events/{id}", rest.BasicAuthMiddleware(http.HandlerFunc(rest.deleteEvent))).Methods("DELETE")
	api.Handle("/events/{id}", rest.BasicAuthMiddleware(http.HandlerFunc(rest.updateEvent))).Methods("PUT")

	rest.Mux = api

	return rest
}

func (rest *Rest) BasicAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		err := rest.service.Users.CheckPassword(user, pass)
		if !ok || err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			log.Println(err.Error())
			return
		}
		handler(w, r)
	}
}

func (rest *Rest) Listen() (err error) {
	rest.listener, err = net.Listen("tcp", rest.address)
	if err != nil {
		return err
	}

	r := http.NewServeMux()
	r.Handle("/", rest.Mux)
	server := &http.Server{
		Handler: r,
	}

	rest.setupMiddleware()

	return server.Serve(rest.listener)
}

func (rest *Rest) setupMiddleware() {
	rest.Mux.Use(func(handler http.Handler) http.Handler {
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
