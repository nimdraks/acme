package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/dataservice"
)

// ListHandler is the HTTP handler for the "List Do people" endpoint
// In this simplified example we are assuming all possible errors are system errors (HTTP 500)

func NewListHandler(d dataservice.DataService) *ListHandler {
	return &ListHandler{dataService: d}
}

type ListHandler struct {
	dataService dataservice.DataService
}

// ServeHTTP implements http.Handler
func (h *ListHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// attempt loadAll
	people, err := h.dataService.LoadAll(context.TODO())
	if err != nil {
		// not need to log here as we can expect other layers to do so
		response.WriteHeader(http.StatusNotFound)
		return
	}

	// happy path
	err = h.writeJSON(response, people)
	if err != nil {
		// this error should not happen but if it does there is nothing we can do to recover
		response.WriteHeader(http.StatusInternalServerError)
	}
}

// output the result as JSON
func (h *ListHandler) writeJSON(writer io.Writer, people []*dataservice.Person) error {
	output := &listResponseFormat{
		People: make([]*listResponseItemFormat, len(people)),
	}

	for index, record := range people {
		output.People[index] = &listResponseItemFormat{
			ID:       record.ID,
			FullName: record.FullName,
			Phone:    record.Phone,
		}
	}

	// call to http.ResponseWriter.Write() will cause HTTP OK (200) to be output as well
	return json.NewEncoder(writer).Encode(output)
}

type listResponseFormat struct {
	People []*listResponseItemFormat `json:"people"`
}

type listResponseItemFormat struct {
	ID       int    `json:"id"`
	FullName string `json:"name"`
	Phone    string `json:"phone"`
}
