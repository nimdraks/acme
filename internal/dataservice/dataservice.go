package dataservice

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

type DataService interface {
	Load(id int) (*Person, error)
	LoadAll() ([]*Person, error)
	Save(fullName, phone, currency, price string) int
}
