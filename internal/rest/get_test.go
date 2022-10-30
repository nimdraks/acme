package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/dataservice"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestGetHandler_Refactored(t *testing.T) {

	for _, s := range []struct {
		desc           string
		inRequest      *http.Request
		d              dataservice.DataService
		expectedStatus int
	}{
		{
			desc:           "Happy path",
			inRequest:      buildValidGetRequest(),
			d:              dataservice.InitHappyMockDataService(),
			expectedStatus: http.StatusOK,
		}, {
			desc:           "Invalid id case1",
			inRequest:      buildInValidGetRequest(),
			d:              dataservice.InitHappyMockDataService(),
			expectedStatus: http.StatusBadRequest,
		}, {
			desc:           "Invalid id case2",
			inRequest:      buildInValidGetRequest2(),
			d:              dataservice.InitHappyMockDataService(),
			expectedStatus: http.StatusBadRequest,
		}, {
			desc:           "Error from data service",
			inRequest:      buildValidGetRequest(),
			d:              dataservice.InitMockBadNotFoundDataService(),
			expectedStatus: http.StatusNotFound,
		},
	} {
		response := httptest.NewRecorder()
		handler := NewGetHandle(s.d)
		handler.ServeHTTP(response, s.inRequest)
		require.Equal(t, s.expectedStatus, response.Code)
	}
}

func buildValidGetRequest() *http.Request {
	r, _ := http.NewRequest("GET", "/person/1/", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "1"})
	return r
}

func buildInValidGetRequest() *http.Request {
	r, _ := http.NewRequest("GET", "/person//", nil)
	r = mux.SetURLVars(r, map[string]string{"id": ""})
	return r
}

func buildInValidGetRequest2() *http.Request {
	r, _ := http.NewRequest("GET", "/person/a/", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "a"})
	return r
}
