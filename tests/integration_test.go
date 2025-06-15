package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-fiber-unittest/domain/entities"
	gw "go-fiber-unittest/src/gateways"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// SetupApp initializes the app with the necessary configurations and dependencies.
func SetupApp() *fiber.App {
	app := fiber.New()

	// Setup mock service
	mockUserService := new(MockUserService)

	// Mocking the responses
	mockUserService.On("InsertNewAccount", mock.Anything).Return(true)
	mockUserService.On("GetAllUser").Return([]entities.UserDataFormat{
		{UserID: "1", Username: "John Doe", Email: "john.doe@example.com"},
	}, nil)
	// Setup gateway with mock service
	gw.NewHTTPGateway(app, mockUserService)

	return app
}

func TestCreateNewUserAccount(t *testing.T) {
	app := SetupApp()

	newUser := &entities.NewUserBody{
		UserID:   "1",
		Username: "John Doe",
		Email:    "john.doe@example.com",
	}
	body, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/add_user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAllUserData(t *testing.T) {
	app := SetupApp()

	// First, create a new user to ensure there is data to retrieve
	newUser := &entities.NewUserBody{
		UserID:   "1",
		Username: "John Doe",
		Email:    "john.doe@example.com",
	}
	body, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Now test retrieving all users
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/users", nil)
	resp, _ = app.Test(req, -1)

	var response entities.ResponseModel
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", response.Message)
	assert.NotEmpty(t, response.Data)

	users := response.Data.([]interface{})
	assert.Equal(t, 1, len(users))

	user := users[0].(map[string]interface{})
	assert.Equal(t, "John Doe", user["username"])
	assert.Equal(t, "john.doe@example.com", user["email"])
}

func TestCreateNewUserAccount_BadRequest(t *testing.T) {
	app := SetupApp()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/add_user", bytes.NewBufferString(`{"name": "John Doe"`)) // malformed JSON
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	fmt.Println("status code", resp.StatusCode)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetAllUserData_NoUsers(t *testing.T) {
	app := fiber.New()

	// Setup mock service with no users
	mockUserService := new(MockUserService)
	mockUserService.On("GetAllUser").Return([]entities.UserDataFormat{}, nil)

	gw.NewHTTPGateway(app, mockUserService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/users", nil)
	resp, _ := app.Test(req, -1)

	var response entities.ResponseModel
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", response.Message)
	assert.Empty(t, response.Data)
}
