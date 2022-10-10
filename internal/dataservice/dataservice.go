package dataservice

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"

	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/dataservice/sqldb"
	"github.com/PacktPublishing/Hands-On-Dependency-Injection-in-Go/ch04/acme/internal/logging"
)

var ErrNotFound = errors.New("not found")

type DataServiceInterface interface {
	Load(id int) (*Person, error)
	LoadAll() ([]*Person, error)
	Save(fullName, phone, currency, price string) int
}

type DataService struct {
	db sqldb.DBService
}

func InitDataService(DSN string) *DataService {
	sqlDB := sqldb.InitDBService(DSN)
	return &DataService{db: sqlDB}
}

type Person struct {
	// ID is the unique ID for this person
	ID int
	// FullName is the name of this person
	FullName string
	// Phone is the phone for this person
	Phone string
	// Currency is the currency this person has paid in
	Currency string
	// Price is the amount (in the above currency) paid by this person
	Price float64
}

func (p *Person) WriteJson(writer io.Writer) error {
	// the JSON response format
	type getResponseFormat struct {
		ID       int     `json:"id"`
		FullName string  `json:"name"`
		Phone    string  `json:"phone"`
		Currency string  `json:"currency"`
		Price    float64 `json:"price"`
	}

	output := &getResponseFormat{
		ID:       p.ID,
		FullName: p.FullName,
		Phone:    p.Phone,
		Currency: p.Currency,
		Price:    p.Price,
	}
	return json.NewEncoder(writer).Encode(output)
}

// custom type so we can convert sql results to easily
type scanner func(dest ...interface{}) error

// reduce the duplication (and maintenance) between sql.Row and sql.Rows usage
func populatePerson(scanner scanner) (*Person, error) {
	out := &Person{}
	err := scanner(&out.ID, &out.FullName, &out.Phone, &out.Currency, &out.Price)
	return out, err
}

func (d *DataService) Load(id int) (*Person, error) {
	row := d.db.Load(id)

	// retrieve columns and populate the person object
	out, err := populatePerson(row.Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			logging.L.Warn("failed to load requested person '%d'. err: %s", id, err)
			return nil, ErrNotFound
		}

		logging.L.Error("failed to convert query result. err: %s", err)
		return nil, err
	}
	return out, nil
}

func (d *DataService) LoadAll() ([]*Person, error) {
	rows, err := d.db.LoadAll()
	if err != nil {
		return nil, err

	}

	defer func() {
		_ = rows.Close()
	}()

	var out []*Person

	for rows.Next() {
		// retrieve columns and populate the person object
		record, err := populatePerson(rows.Scan)
		if err != nil {
			logging.L.Error("failed to convert query result. err: %s", err)
			return nil, err
		}

		out = append(out, record)
	}

	if len(out) == 0 {
		logging.L.Warn("no people found in the database.")
		return nil, ErrNotFound
	}

	return out, nil

}

func (d *DataService) Save(fullName, phone, currency, price string) int {
	id, _ := d.db.Save(fullName, phone, currency, price)
	return id
}