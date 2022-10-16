package dataservice

import "errors"

type MockHappyDataService struct {
}

func InitHappyMockDataService() *MockHappyDataService {
	return &MockHappyDataService{}
}

func (d *MockHappyDataService) Load(id int) (*Person, error) {

	return &Person{ID: 1, FullName: "John", Phone: "0123456780", Currency: "USD", Price: 100}, nil
}

func (d *MockHappyDataService) LoadAll() ([]*Person, error) {
	return nil, nil
}

func (d *MockHappyDataService) Save(fullName, phone, currency, price string) int {
	return -1
}

type MockBadNotFoundDataService struct {
}

func InitMockBadNotFoundDataService() *MockBadNotFoundDataService {
	return &MockBadNotFoundDataService{}
}

func (d *MockBadNotFoundDataService) Load(id int) (*Person, error) {

	return nil, errors.New("error")
}

func (d *MockBadNotFoundDataService) LoadAll() ([]*Person, error) {
	return nil, nil
}

func (d *MockBadNotFoundDataService) Save(fullName, phone, currency, price string) int {
	return -1
}
