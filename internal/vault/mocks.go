package vault

import "github.com/stretchr/testify/mock"

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Init() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) All() Records {
	args := m.Called()
	return args.Get(0).(Records)
}

func (m *MockStorage) Get(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *MockStorage) Exists(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

func (m *MockStorage) Upsert(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockStorage) Destroy(key string) error {
	args := m.Called(key)
	return args.Error(0)
}
