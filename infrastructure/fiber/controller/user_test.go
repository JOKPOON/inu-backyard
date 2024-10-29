package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/infrastructure/fiber/request"
	"github.com/team-inu/inu-backyard/mocks"
)

func setupUserControllerTest() (*fiber.App, *mocks.MockUserUseCase, *mocks.MockAuthUseCase, *mocks.MockValidator) {
	app := fiber.New()
	mockUserUseCase := new(mocks.MockUserUseCase)
	mockAuthUseCase := new(mocks.MockAuthUseCase)
	mockValidator := new(mocks.MockValidator)

	userController := NewUserController(mockValidator, mockUserUseCase, mockAuthUseCase)

	// Setting up routes
	app.Get("/users", userController.GetAll)
	app.Get("/users/:userId", userController.GetById)
	app.Post("/users", userController.Create)
	app.Post("/users/bulk", userController.CreateMany)
	app.Put("/users/:userId", userController.Update)
	app.Delete("/users/:userId", userController.Delete)
	app.Put("/users/:userId/password", userController.ChangePassword)

	return app, mockUserUseCase, mockAuthUseCase, mockValidator
}

// Test cases for GetAll route
func TestGetAllUsers(t *testing.T) {
	app, mockUserUseCase, _, _ := setupUserControllerTest()

	// Successful retrieval
	t.Run("GetAll Users - Success", func(t *testing.T) {
		mockUserUseCase.On("GetAll").Return([]entity.User{
			{Id: "1", FirstName: "John", LastName: "Doe"},
			{Id: "2", FirstName: "Jane", LastName: "Doe"},
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res mocks.Response
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res.Data.([]interface{}), 2)
	})

	// Empty user list
	t.Run("GetAll Users - No Users Found", func(t *testing.T) {
		mockUserUseCase.On("GetAll").Return([]entity.User{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res mocks.Response
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Len(t, res.Data.([]interface{}), 0)
	})
}

// Test cases for GetById route
func TestGetUserById(t *testing.T) {
	app, mockUserUseCase, _, _ := setupUserControllerTest()

	// Valid user ID
	t.Run("GetById - Valid User", func(t *testing.T) {
		mockUserUseCase.On("GetById", "1").Return(&entity.User{Id: "1", FirstName: "John", LastName: "Doe"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res mocks.Response
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, "John", res.Data.(map[string]interface{})["firstName"])
	})

	// Invalid user ID
	t.Run("GetById - User Not Found", func(t *testing.T) {
		mockUserUseCase.On("GetById", "999").Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// Test cases for Create route
func TestCreateUser(t *testing.T) {
	app, mockUserUseCase, _, mockValidator := setupUserControllerTest()

	// Successful creation
	t.Run("Create User - Success", func(t *testing.T) {
		mockValidator.On("Validate", mock.Anything, mock.Anything).Return(true, nil)
		mockUserUseCase.On("Create", mock.Anything).Return(nil)

		payload := request.CreateUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  "password",
			Role:      entity.UserRoleLecturer,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// Validation failure
	t.Run("Create User - Validation Failure", func(t *testing.T) {
		mockValidator.On("Validate", mock.Anything, mock.Anything).Return(false, fiber.ErrBadRequest)

		payload := request.CreateUserPayload{
			FirstName: "",
			LastName:  "Doe",
			Email:     "invalid-email",
			Password:  "pwd",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// Test cases for Update route
func TestUpdateUser(t *testing.T) {
	app, mockUserUseCase, _, mockValidator := setupUserControllerTest()

	// Successful update
	t.Run("Update User - Success", func(t *testing.T) {
		mockValidator.On("Validate", mock.Anything, mock.Anything).Return(true, nil)
		mockUserUseCase.On("Update", mock.Anything, mock.Anything).Return(nil)

		payload := request.UpdateUserPayload{
			FirstName: "John",
			LastName:  "Smith",
			Email:     "john.smith@example.com",
			Role:      entity.UserRoleLecturer,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Update with non-existent user ID
	t.Run("Update User - User Not Found", func(t *testing.T) {
		mockUserUseCase.On("Update", mock.Anything, mock.Anything).Return(fiber.ErrNotFound)

		payload := request.UpdateUserPayload{
			FirstName: "Non-existent",
			LastName:  "User",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// Additional test cases for Delete and ChangePassword routes would follow similar structure
// Ensure coverage of successful actions and error handling for each endpoint.
