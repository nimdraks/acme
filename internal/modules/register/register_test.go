package register

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/dataservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockOutConfig struct{}

func (m *mockOutConfig) GetBasePrice() float64 {
	return 100
}
func (m *mockOutConfig) GetExchangeRateBaseURL() string {
	return "https://api.apilayer.com/currency_data"
}
func (m *mockOutConfig) GetExchangeRateAPIKey() string {
	return "KqLgYfBfNgCgGfHG6UFdJWp3qOdaoGYc"
}

type mockInConfig struct {
	testUrl string
}

func (m *mockInConfig) GetBasePrice() float64 {
	return 100
}

func (m *mockInConfig) GetExchangeRateBaseURL() string {
	return m.testUrl
}

func (m *mockInConfig) GetExchangeRateAPIKey() string {
	return "test"
}
func TestRegisterer_Do(t *testing.T) {
	// inputs
	in := &dataservice.Person{
		FullName: "Chang",
		Phone:    "11122233345",
		Currency: "CNY",
	}

	// call method
	registerer := NewRegisterer(&mockOutConfig{}, dataservice.InitHappyMockDataService())
	ID, err := registerer.Do(context.TODO(), in)

	// validate expectations
	require.NoError(t, err)
	assert.True(t, ID > 0)
}

func TestRegisterer_Do_Refactored(t *testing.T) {
	// inputs
	in := &dataservice.Person{
		FullName: "Chang",
		Phone:    "11122233345",
		Currency: "AUD",
	}

	for _, test := range []struct {
		name    string
		handler http.HandlerFunc
	}{
		{
			name: "Happy",
			handler: func(w http.ResponseWriter, r *http.Request) {
				payload := []byte(`
				{
				   "success":true,
				   "historical":true,
				   "date":"2010-11-09",
				   "timestamp":1289347199,
				   "source":"USD",
				   "quotes":{
					  "USDAUD":0.989981
				   }
				}`)
				w.Write(payload)
			},
		},
		{
			name: "Not Happy1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				payload := []byte(`
				{
				   "success":true,
				   "historical":true,
				   "date":"2010-11-09",
				   "timestamp":1289347199,
				   "source":"USD",
				   "quotes":{
					  "USDCNY":0.989981
				   }
				}`)
				w.Write(payload)
			},
		},
		{
			name: "Not Happy2",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(500)
			},
		},
	} {
		// set Servers
		server := httptest.NewServer(test.handler)
		defer server.Close()

		// call method
		registerer := NewRegisterer(&mockInConfig{server.URL}, dataservice.InitHappyMockDataService())
		ID, err := registerer.Do(context.TODO(), in)

		// validate expectations
		if test.name == "Happy" {
			require.NoError(t, err)
			assert.True(t, ID > 0)
		} else {
			require.Error(t, err)
		}
	}
}
