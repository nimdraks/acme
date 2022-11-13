package register

import (
	"context"
	"errors"
	"fmt"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/dataservice"
	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/logging"
	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/modules/data"
	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/modules/exchange"
)

const (
	// default person id (returned on error)
	defaultPersonID = 0
)

var (
	// validation errors
	errNameMissing     = errors.New("name is missing")
	errPhoneMissing    = errors.New("phone is missing")
	errCurrencyMissing = errors.New("currency is missing")
	errInvalidCurrency = errors.New("currency is invalid, supported types are AUD, CNY, EUR, GBP, JPY, MYR, SGD, USD")

	// a little trick to make checking for supported currencies easier
	supportedCurrencies = map[string]struct{}{
		"AUD": {},
		"CNY": {},
		"EUR": {},
		"GBP": {},
		"JPY": {},
		"MYR": {},
		"SGD": {},
		"USD": {},
	}
)

// Registerer validates the supplied person, calculates the price in the requested currency and saves the result.
// It will return an error when:
// -the person object does not include all the fields
// -the currency is invalid
// -the exchange rate cannot be loaded
// -the data layer throws an error.
type Config interface {
	GetBasePrice() float64
	GetExchangeRateBaseURL() string
	GetExchangeRateAPIKey() string
}

type Registerer struct {
	config      Config
	dataService dataservice.DataService
}

func NewRegisterer(config Config, d dataservice.DataService) *Registerer {
	return &Registerer{config: config, dataService: d}
}

// Do is API for this struct
func (r *Registerer) Do(ctx context.Context, in *data.Person) (int, error) {
	// validate the request
	err := r.validateInput(in)
	if err != nil {
		logging.L.Warn("input validation failed with err: %s", err)
		return defaultPersonID, err
	}

	// get price in the requested currency
	price, err := r.getPrice(ctx, in.Currency)
	if err != nil {
		return defaultPersonID, err
	}

	// save registration
	id, err := r.save(ctx, in, price)
	if err != nil {
		// no need to log here as we expect the data layer to do so
		return defaultPersonID, err
	}

	return id, nil
}

// validate input and return error on fail
func (r *Registerer) validateInput(in *data.Person) error {
	if in.FullName == "" {
		return errNameMissing
	}
	if in.Phone == "" {
		return errPhoneMissing
	}
	if in.Currency == "" {
		return errCurrencyMissing
	}

	if _, found := supportedCurrencies[in.Currency]; !found {
		return errInvalidCurrency
	}

	// happy path
	return nil
}

// get price in the requested currency
func (r *Registerer) getPrice(ctx context.Context, currency string) (float64, error) {
	price, err := exchange.NewConverter(r.config).Do(ctx, r.config.GetBasePrice(), currency)
	if err != nil {
		logging.L.Warn("failed to convert the price. err: %s", err)
		return defaultPersonID, err
	}

	return price, nil
}

// save the registration
func (r *Registerer) save(ctx context.Context, in *data.Person, price float64) (int, error) {
	result := r.dataService.Save(ctx, in.FullName, in.Phone, in.Currency, fmt.Sprintf("%f", price))
	if result < 0 {
		return result, fmt.Errorf("Error")
	}

	return result, nil
}
