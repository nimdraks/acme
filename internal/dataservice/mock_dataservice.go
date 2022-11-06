package dataservice

import (
	"context"
	"errors"
)

type MockHappyDataService struct {
}

func InitHappyMockDataService() *MockHappyDataService {
	return &MockHappyDataService{}
}

func (d *MockHappyDataService) Load(ctx context.Context, id int) (*Person, error) {

	return &Person{ID: 1, FullName: "John", Phone: "0123456780", Currency: "USD", Price: 100}, nil
}

func (d *MockHappyDataService) LoadAll(ctx context.Context) ([]*Person, error) {
	return []*Person{{ID: 1, FullName: "John", Phone: "0123456780", Currency: "USD", Price: 100}}, nil
}

func (d *MockHappyDataService) Save(ctx context.Context, fullName, phone, currency, price string) int {
	return -1
}

type MockBadNotFoundDataService struct {
}

func InitMockBadNotFoundDataService() *MockBadNotFoundDataService {
	return &MockBadNotFoundDataService{}
}

func (d *MockBadNotFoundDataService) Load(ctx context.Context, id int) (*Person, error) {

	return nil, errors.New("error")
}

func (d *MockBadNotFoundDataService) LoadAll(ctx context.Context) ([]*Person, error) {
	return nil, errors.New("error")
}

func (d *MockBadNotFoundDataService) Save(ctx context.Context, fullName, phone, currency, price string) int {
	return -1
}
