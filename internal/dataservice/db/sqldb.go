package db

import (
	"database/sql"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/logging"
)

const (
	// default person id (returned on error)
	defaultPersonID = 0
)

type SqlDB struct {
	sql *sql.DB
}

func InitDBService(DSN string) *SqlDB {
	if DSN == "" {
		return nil
	}

	db, err := sql.Open("mysql", DSN)
	if err != nil {
		panic(err.Error())
	}

	return &SqlDB{sql: db}
}

func (d *SqlDB) Load(ID int) *sql.Row {
	// perform DB select
	query := "SELECT id, fullname, phone, currency, price FROM person WHERE id = ? LIMIT 1"
	row := d.sql.QueryRow(query, ID)
	return row
}

// LoadAll will attempt to load all people in the database
// It will return ErrNotFound when there are not people in the database
// Any other errors returned are caused by the underlying database or our connection to it.
func (d *SqlDB) LoadAll() (*sql.Rows, error) {
	// perform DB select
	query := "SELECT id, fullname, phone, currency, price FROM person"
	rows, err := d.sql.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Save will save the supplied person and return the ID of the newly created person or an error.
// Errors returned are caused by the underlying database or our connection to it.
func (d *SqlDB) Save(fullName, phone, currency, price string) (int, error) {
	// perform DB insert
	query := "INSERT INTO person (fullname, phone, currency, price) VALUES (?, ?, ?, ?)"
	result, err := d.sql.Exec(query, fullName, phone, currency, price)
	if err != nil {
		logging.L.Error("failed to save person into DB. err: %s", err)
		return defaultPersonID, err
	}

	// retrieve and return the ID of the person created
	id, err := result.LastInsertId()
	if err != nil {
		logging.L.Error("failed to retrieve id of last saved person. err: %s", err)
		return defaultPersonID, err
	}
	return int(id), nil
}
