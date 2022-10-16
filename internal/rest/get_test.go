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
			d:              dataservice.InitMockDataService(),
			expectedStatus: http.StatusOK,
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
