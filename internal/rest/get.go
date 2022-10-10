package rest

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/logging"
	"github.com/gorilla/mux"
)

const (
	// default person id (returned on error)
	defaultPersonID = 0
)

// GetHandler is the HTTP handler for the "Get Person" endpoint
// In this simplified example we are assuming all possible errors are user errors and returning "bad request" HTTP 400
// or "not found" HTTP 404
// There are some programmer errors possible but hopefully these will be caught in testing.

// ServeHTTP implements http.Handler
func (s *Server) handlerGet(response http.ResponseWriter, request *http.Request) {
	// extract person id from request
	id, err := extractID(request)
	if err != nil {
		// output error
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	person, err := s.DataService.Load(id)

	if err != nil {
		// not need to log here as we can expect other layers to do so
		response.WriteHeader(http.StatusNotFound)
		return
	}

	// happy path
	err = writeJSON(response, person)
	if err != nil {
		// this error should not happen but if it does there is nothing we can do to recover
		response.WriteHeader(http.StatusInternalServerError)
	}
}

// extract the person ID from the request
func extractID(request *http.Request) (int, error) {
	// ID is part of the URL, so we extract it from there
	vars := mux.Vars(request)
	idAsString, exists := vars["id"]
	if !exists {
		// log and return error
		err := errors.New("[get] person id missing from request")
		logging.L.Warn(err.Error())
		return defaultPersonID, err
	}

	// convert ID to int
	id, err := strconv.Atoi(idAsString)
	if err != nil {
		// log and return error
		err = fmt.Errorf("[get] failed to convert person id into a number. err: %s", err)
		logging.L.Error(err.Error())
		return defaultPersonID, err
	}

	return id, nil
}

type writeReponseJson interface {
	WriteJson(writer io.Writer) error
}

// output the supplied person as JSON
func writeJSON(writer io.Writer, data writeReponseJson) error {
	return data.WriteJson(writer)
}
