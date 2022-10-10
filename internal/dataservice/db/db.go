package db

import "database/sql"

type DBService interface {
	Save(fullName, phone, currency, price string) (int, error)
	Load(ID int) *sql.Row
	LoadAll() (*sql.Rows, error)
}
