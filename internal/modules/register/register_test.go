package register

import (
	"context"
	"testing"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/modules/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConfig struct{}

func (m *mockConfig) GetBasePrice() float64 {
	return 100
}
func (m *mockConfig) GetExchangeRateBaseURL() string {
	return "https://api.apilayer.com/currency_data"
}
func (m *mockConfig) GetExchangeRateAPIKey() string {
	return "KqLgYfBfNgCgGfHG6UFdJWp3qOdaoGYc"
}

func TestRegisterer_Do(t *testing.T) {
	// inputs
	in := &data.Person{
		FullName: "Chang",
		Phone:    "11122233345",
		Currency: "CNY",
	}

	// call method
	registerer := NewRegisterer(&mockConfig{})
	ID, err := registerer.Do(context.TODO(), in)

	// validate expectations
	require.NoError(t, err)
	assert.True(t, ID > 0)
}
