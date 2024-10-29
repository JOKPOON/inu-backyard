package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/team-inu/inu-backyard/entity"
	"github.com/team-inu/inu-backyard/internal/config"
	"github.com/team-inu/inu-backyard/mocks"
)

func setupTestApp() (*fiber.App, *AuthController, *mocks.MockAuthUseCase, *mocks.MockValidator, *mocks.MockUserUseCase, *mocks.Turnstile) {
	app := fiber.New()
	mockAuthUseCase := new(mocks.MockAuthUseCase)
	mockValidator := new(mocks.MockValidator)
	mockUserUseCase := new(mocks.MockUserUseCase)
	mockAuthConfig := config.AuthConfig{}
	mockTurnstile := new(mocks.Turnstile)

	authController := NewAuthController(mockValidator, mockAuthConfig, mockTurnstile, mockAuthUseCase, mockUserUseCase)

	// Setup routes
	app.Post("/auth/login", authController.SignIn)
	app.Get("/auth/logout", authController.SignOut)
	app.Get("/auth/me", authController.Me)

	return app, authController, mockAuthUseCase, mockValidator, mockUserUseCase, mockTurnstile
}

func TestAuthController_SignIn(t *testing.T) {
	app, _, mockAuthUseCase, mockValidator, _, mockTurnstile := setupTestApp()

	mockValidator.On("Validate", mock.Anything, mock.Anything).Return(true, nil)
	mockTurnstile.On("Validate", mock.Anything, mock.Anything).Return(true, nil)
	mockAuthUseCase.On("SignIn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&fiber.Cookie{Name: "session", Value: "testToken"}, nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check for a valid session cookie
	cookie := resp.Cookies()
	assert.Len(t, cookie, 1)
	assert.Equal(t, "session", cookie[0].Name)
	assert.Equal(t, "testToken", cookie[0].Value)
}

func TestAuthController_SignOut(t *testing.T) {
	app, _, mockAuthUseCase, mockValidator, _, _ := setupTestApp()

	mockValidator.On("Validate", mock.Anything, mock.Anything).Return(true, nil)
	mockAuthUseCase.On("SignOut", mock.Anything).Return(&fiber.Cookie{Name: "session", Value: "", MaxAge: -1}, nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/logout", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check if the session cookie has been cleared
	cookie := resp.Cookies()
	assert.Len(t, cookie, 1)
	assert.Equal(t, "session", cookie[0].Name)
	assert.Empty(t, cookie[0].Value)
}

func TestAuthController_Me(t *testing.T) {
	app, _, _, mockValidator, mockUserUseCase, _ := setupTestApp()

	mockValidator.On("Validate", mock.Anything, mock.Anything).Return(true, nil)
	mockUserUseCase.On("GetUserFromCtx", mock.Anything).Return(&entity.User{
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var res mocks.Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "John", res.Data.(map[string]interface{})["firstName"])
	assert.Equal(t, "Doe", res.Data.(map[string]interface{})["lastName"])
}
