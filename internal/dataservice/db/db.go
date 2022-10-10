package db

import "database/sql"

const (
	// default person id (returned on error)
	defaultPersonID = 0
)

type DBService interface {
	Save(fullName, phone, currency, price string) (int, error)
	Load(ID int) *sql.Row
	LoadAll() (*sql.Rows, error)
}
