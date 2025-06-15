package tests

import (
	"go-fiber-unittest/domain/entities"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAllUser() ([]entities.UserDataFormat, error) {
	args := m.Called()
	return args.Get(0).([]entities.UserDataFormat), args.Error(1)
}

func (m *MockUserService) InsertNewAccount(newUser *entities.NewUserBody) bool {
	args := m.Called(newUser)
	return args.Bool(0)
}
