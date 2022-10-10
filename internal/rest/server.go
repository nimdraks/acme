package rest

import (
	"net/http"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/config"
	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/dataservice"
	"github.com/gorilla/mux"
)

// New will create and initialize the server
func New(config *config.Config) *Server {
	dataService := dataservice.InitDataService(config.DSN)
	return &Server{
		config:          config,
		DataService:     dataService,
		handlerList:     &ListHandler{},
		handlerNotFound: notFoundHandler,
		handlerRegister: &RegisterHandler{},
	}
}

// Server is the HTTP REST server
type Server struct {
	server      *http.Server
	config      *config.Config
	DataService dataservice.DataService

	handlerList     http.Handler
	handlerNotFound http.HandlerFunc
	handlerRegister http.Handler
}

// Listen will start a HTTP rest for this service
func (s *Server) Listen(stop <-chan struct{}) {
	router := s.buildRouter()

	// create the HTTP server
	s.server = &http.Server{
		Handler: router,
		Addr:    s.config.Address,
	}

	// listen for shutdown
	go func() {
		// wait for shutdown signal
		<-stop

		_ = s.server.Close()
	}()

	// start the HTTP server
	_ = s.server.ListenAndServe()
}

// configure the endpoints to handlers
func (s *Server) buildRouter() http.Handler {
	router := mux.NewRouter()

	// map URL endpoints to HTTP handlers
	router.HandleFunc("/person/{id}/", s.handlerGet).Methods("GET")
	router.Handle("/person/list", s.handlerList).Methods("GET")
	router.Handle("/person/register", s.handlerRegister).Methods("POST")

	// convert a "catch all" not found handler
	router.NotFoundHandler = s.handlerNotFound

	return router
}
