package dataservice

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type DataService interface {
	Load(ctx context.Context, id int) (*Person, error)
	LoadAll(ctx context.Context) ([]*Person, error)
	Save(ctx context.Context, fullName, phone, currency, price string) int
}
