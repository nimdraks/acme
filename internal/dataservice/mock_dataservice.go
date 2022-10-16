package dataservice

type MockDataService struct {
}

func InitMockDataService() *MockDataService {
	return &MockDataService{}
}

func (d *MockDataService) Load(id int) (*Person, error) {

	return &Person{ID: 1, FullName: "John", Phone: "0123456780", Currency: "USD", Price: 100}, nil
}

func (d *MockDataService) LoadAll() ([]*Person, error) {
	return nil, nil
}

func (d *MockDataService) Save(fullName, phone, currency, price string) int {
	return -1
}
