package repositories

import (
	"go-fiber-unittest/domain/entities"

	"github.com/stretchr/testify/mock"
)

type MockIUserRepository struct {
	mock.Mock
}

func (m *MockIUserRepository) InsertNewUser(data *entities.NewUserBody) bool {
	args := m.Called(data)
	return args.Bool(0)
}

func (m *MockIUserRepository) FindByID(id string) (*entities.UserDataFormat, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.UserDataFormat), args.Error(1)
}

func (m *MockIUserRepository) FindAll() ([]entities.UserDataFormat, error) {
	args := m.Called()
	return args.Get(0).([]entities.UserDataFormat), args.Error(1)
}
