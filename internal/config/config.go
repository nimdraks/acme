package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/logging"
)

// DefaultEnvVar is the default environment variable the points to the config file
const DefaultEnvVar = "ACME_CONFIG"

// App is the application config
var App *Config

// Config defines the JSON format for the config file
type Config struct {
	// DSN is the data source name (format: https://github.com/go-sql-driver/mysql/#dsn-data-source-name)
	DSN string

	// Address is the IP address and port to bind this rest to
	Address string

	// BasePrice is the price of registration
	BasePrice float64

	// ExchangeRateBaseURL is the server and protocol part of the URL from which to load the exchange rate
	ExchangeRateBaseURL string

	// ExchangeRateAPIKey is the API for the exchange rate API
	ExchangeRateAPIKey string
}

func (c *Config) GetDSN() string {
	return c.DSN
}

func (c *Config) GetAddress() string {
	return c.Address
}

func (c *Config) GetBasePrice() float64 {
	return c.BasePrice
}

func (c *Config) GetExchangeRateBaseURL() string {
	return c.ExchangeRateBaseURL
}

func (c *Config) GetExchangeRateAPIKey() string {
	return c.ExchangeRateAPIKey
}

// Load returns the config loaded from environment
func init() {
	filename, found := os.LookupEnv(DefaultEnvVar)
	if !found {
		logging.L.Error("failed to locate file specified by %s", DefaultEnvVar)
		return
	}

	_ = load(filename)
}

func load(filename string) error {
	App = &Config{}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logging.L.Error("failed to read config file. err: %s", err)
		return err
	}

	err = json.Unmarshal(bytes, App)
	if err != nil {
		logging.L.Error("failed to parse config file. err : %s", err)
		return err
	}

	return nil
}

func Load(filename string) *Config {
	if filename == "" {
		logging.L.Error("Empty filename")
		return nil
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logging.L.Error("failed to read config file. err: %s", err)
		return nil
	}

	config := &Config{}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		logging.L.Error("failed to parse config file. err : %s", err)
		return nil
	}

	return config
}
