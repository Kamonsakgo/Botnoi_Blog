package gateways_test

import (
	"bytes"
	"encoding/json"
	"go-fiber-unittest/domain/entities"
	"go-fiber-unittest/src/gateways"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock UserService
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

func TestGetAllUserData(t *testing.T) {
	app := fiber.New()

	mockUserService := new(MockUserService)
	mockUserService.On("GetAllUser").Return([]entities.UserDataFormat{}, nil)

	gateway := &gateways.HTTPGateway{UserService: mockUserService}
	app.Get("/users", gateway.GetAllUserData)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockUserService.AssertExpectations(t)
}

func TestCreateNewUserAccount(t *testing.T) {
	app := fiber.New()

	mockUserService := new(MockUserService)
	newUser := &entities.NewUserBody{
		UserID:   "1",
		Username: "john.doe",
		Email:    "john.doe@example.com",
	}
	mockUserService.On("InsertNewAccount", newUser).Return(true)

	gateway := &gateways.HTTPGateway{UserService: mockUserService}
	app.Post("/add_user", gateway.CreateNewUserAccount)

	body, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/add_user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockUserService.AssertExpectations(t)
}

func TestCreateNewUserAccount_BadRequest(t *testing.T) {
	app := fiber.New()

	mockUserService := new(MockUserService)
	gateway := &gateways.HTTPGateway{UserService: mockUserService}
	app.Post("/add_user", gateway.CreateNewUserAccount)

	req := httptest.NewRequest(http.MethodPost, "/add_user", bytes.NewBufferString(`{"name": "John Doe"`)) // malformed JSON
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
