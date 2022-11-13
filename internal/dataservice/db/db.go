package db

import (
	"context"
	"database/sql"
)

const (
	// default person id (returned on error)
	defaultPersonID = 0
)

type DBService interface {
	Save(ctx context.Context, fullName, phone, currency, price string) (int, error)
	Load(ctx context.Context, ID int) *sql.Row
	LoadAll(ctx context.Context) (*sql.Rows, error)
}
