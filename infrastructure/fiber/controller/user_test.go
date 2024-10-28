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
)

// Mocking UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) GetAll() ([]entity.User, error) {
	args := m.Called()
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetById(userId string) (*entity.User, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetBySessionId(sessionId string) (*entity.User, error) {
	args := m.Called(sessionId)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetByParams(params *entity.User, limit, offset int) ([]entity.User, error) {
	args := m.Called(params, limit, offset)
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *MockUserUseCase) Create(firstName, lastName, email, password string, role entity.UserRole) error {
	args := m.Called(firstName, lastName, email, password, role)
	return args.Error(0)
}

func (m *MockUserUseCase) CreateMany(users []entity.User) error {
	args := m.Called(users)
	return args.Error(0)
}

func (m *MockUserUseCase) Update(userId string, user *entity.User) error {
	args := m.Called(userId, user)
	return args.Error(0)
}

func (m *MockUserUseCase) Delete(userId string) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *MockUserUseCase) CheckUserRole(ctx *fiber.Ctx, userId string, role entity.UserRole) error {
	args := m.Called(ctx, userId, role)
	return args.Error(0)
}

// Mocking AuthUseCase
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) ChangePassword(userId, oldPassword, newPassword string) error {
	args := m.Called(userId, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockAuthUseCase) Authenticate(sessionId string) (*entity.User, error) {
	args := m.Called(sessionId)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthUseCase) SignIn(email string, password string, ipAddress string, userAgent string) (*fiber.Cookie, error) {
	args := m.Called(email, password, ipAddress, userAgent)
	return args.Get(0).(*fiber.Cookie), args.Error(1)
}

func (m *MockAuthUseCase) SignOut(header string) (*fiber.Cookie, error) {
	args := m.Called(header)
	return args.Get(0).(*fiber.Cookie), args.Error(1)
}

// Mocking PayloadValidator
type MockValidator struct {
	mock.Mock
}

// ValidateAuth implements validator.PayloadValidator.
func (m *MockValidator) ValidateAuth(ctx *fiber.Ctx) (string, error) {
	panic("unimplemented")
}

func (m *MockValidator) Validate(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	args := m.Called(payload, ctx)
	return args.Bool(0), args.Error(1)
}

type res struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func TestUserController(t *testing.T) {
	app := fiber.New()
	mockUserUseCase := new(MockUserUseCase)
	mockAuthUseCase := new(MockAuthUseCase)
	mockValidator := new(MockValidator)

	userController := NewUserController(mockValidator, mockUserUseCase, mockAuthUseCase)

	// Setting up routes
	app.Get("/users", userController.GetAll)
	app.Get("/users/:userId", userController.GetById)
	app.Post("/users", userController.Create)
	app.Post("/users/bulk", userController.CreateMany)
	app.Put("/users/:userId", userController.Update)
	app.Delete("/users/:userId", userController.Delete)
	app.Put("/users/:userId/password", userController.ChangePassword)

	// Test GetAll
	t.Run("GetAll", func(t *testing.T) {
		mockUserUseCase.On("GetAll").Return([]entity.User{
			{Id: "1", FirstName: "John", LastName: "Doe"},
			{Id: "2", FirstName: "Jane", LastName: "Doe"}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res res
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Len(t, res.Data.([]interface{}), 2)
		assert.Equal(t, "John", res.Data.([]interface{})[0].(map[string]interface{})["firstName"])
		assert.Equal(t, "Jane", res.Data.([]interface{})[1].(map[string]interface{})["firstName"])
	})

	// Test GetById
	t.Run("GetById Valid User", func(t *testing.T) {
		mockUserUseCase.On("GetById", "1").Return(&entity.User{Id: "1", FirstName: "John", LastName: "Doe"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res res
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "John", res.Data.(map[string]interface{})["firstName"])
		assert.Equal(t, "Doe", res.Data.(map[string]interface{})["lastName"])
	})

	t.Run("GetById Invalid User", func(t *testing.T) {
		mockUserUseCase.On("GetById", "2").Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/2", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	// Test Create
	t.Run("Create User", func(t *testing.T) {
		mockValidator.On("Validate", mock.Anything, mock.Anything).Return(true, nil)
		mockUserUseCase.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

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

	// Test CreateMany
	t.Run("CreateMany Users", func(t *testing.T) {
		mockValidator.On("Validate", mock.Anything, mock.Anything).Return(true, nil)
		mockUserUseCase.On("CreateMany", mock.Anything).Return(nil)

		payload := request.CreateBulkUserPayload{
			Users: []request.CreateUserPayload{
				{FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", Password: "password", Role: entity.UserRoleLecturer},
			},
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/users/bulk", bytes.NewBuffer(body))
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	// Test Update
	t.Run("Update User", func(t *testing.T) {
		mockValidator.On("Validate", mock.Anything, mock.Anything).Return(true, nil)
		mockUserUseCase.On("Update", mock.Anything, mock.Anything).Return(nil)
		mockUserUseCase.On("CheckUserRole", mock.Anything, mock.Anything, mock.Anything).Return(nil)

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

	// Test Delete
	t.Run("Delete User", func(t *testing.T) {
		mockUserUseCase.On("Delete", "1").Return(nil)
		mockUserUseCase.On("CheckUserRole", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Test ChangePassword
	t.Run("ChangePassword", func(t *testing.T) {
		mockValidator.On("Validate", mock.Anything, mock.Anything).Return(true, nil)
		mockAuthUseCase.On("ChangePassword", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockUserUseCase.On("CheckUserRole", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		payload := request.ChangePasswordPayload{
			OldPassword: "oldPassword",
			NewPassword: "newPassword",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/users/1/password", bytes.NewBuffer(body))
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
