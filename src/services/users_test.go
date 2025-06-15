package services

import (
	"go-fiber-unittest/domain/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockIUsersRepository struct {
	mock.Mock
}

func (m *MockIUsersRepository) FindByID(id string) (*entities.UserDataFormat, error) {
	return &entities.UserDataFormat{UserID: id, Username: "Mock User"}, nil
}

func (m *MockIUsersRepository) FindAll() ([]entities.UserDataFormat, error) {
	args := m.Called()
	return args.Get(0).([]entities.UserDataFormat), args.Error(1)
}

func (m *MockIUsersRepository) InsertNewUser(data *entities.NewUserBody) bool {
	args := m.Called(data)
	return args.Bool(0)
}

func TestGetAllUser(t *testing.T) {
	mockRepo := new(MockIUsersRepository)
	mockService := NewUsersService(mockRepo)

	mockUsers := []entities.UserDataFormat{
		{
			UserID:   "1",
			Username: "John",
			Email:    "j@j.com"},

		{
			UserID:   "2",
			Username: "Jane",
			Email:    "j@j.com"},
	}

	mockRepo.On("FindAll").Return(mockUsers, nil)

	users, err := mockService.GetAllUser()

	assert.Nil(t, err)
	assert.Equal(t, mockUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestInsertNewAccount(t *testing.T) {
	mockRepo := new(MockIUsersRepository)
	mockService := NewUsersService(mockRepo)

	newUser := &entities.NewUserBody{
		UserID:   "1",
		Username: "John",
		Email:    "j@j.com",
	}

	mockRepo.On("InsertNewUser", newUser).Return(true)

	status := mockService.InsertNewAccount(newUser)

	assert.True(t, status)
	mockRepo.AssertExpectations(t)
}
